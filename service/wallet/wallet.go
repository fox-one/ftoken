package wallet

import (
	"context"
	"time"

	"github.com/fox-one/ftoken/core"
	"github.com/fox-one/mixin-sdk-go/v2"
	"github.com/fox-one/mixin-sdk-go/v2/mixinnet"
	"github.com/fox-one/pkg/logger"
	"github.com/shopspring/decimal"
	"github.com/sirupsen/logrus"
)

type (
	Config struct {
		Spend mixinnet.Key
	}

	mixinBot struct {
		client *mixin.Client
		spend  mixinnet.Key
	}
)

func New(cfg Config, client *mixin.Client) core.WalletService {
	if !cfg.Spend.HasValue() {
		panic("spend key is empty")
	}

	return &mixinBot{
		client: client,
		spend:  cfg.Spend,
	}
}

func (m *mixinBot) ListSnapshots(ctx context.Context, offset time.Time, limit int) ([]*core.Snapshot, error) {
	items, err := m.client.ReadSafeSnapshots(ctx, "", offset, "ASC", limit)
	if err != nil {
		return nil, err
	}
	return m.toSnapshots(ctx, items)
}

func (m *mixinBot) Transfer(ctx context.Context, req *core.Transfer) error {
	logger := logger.FromContext(ctx).WithFields(logrus.Fields{
		"trace_id": req.TraceID,
		"asset_id": req.AssetID,
		"amount":   req.Amount,
	})
	utxos, err := m.client.SafeListUtxos(ctx, mixin.SafeListUtxoOption{
		Members:   []string{m.client.ClientID},
		Threshold: 1,
		Limit:     255,
		Asset:     req.AssetID,
		State:     mixin.SafeUtxoStateUnspent,
	})
	if err != nil {
		logger.WithError(err).Errorln("SafeListUtxos")
		return err
	}

	balance := decimal.Zero
	for i, utxo := range utxos {
		balance = balance.Add(utxo.Amount)
		if balance.GreaterThanOrEqual(req.Amount) {
			utxos = utxos[:i+1]
			break
		}
	}

	if balance.LessThan(req.Amount) {
		req, err := m.client.SafeReadTransactionRequest(ctx, req.TraceID)
		if err != nil {
			logger.WithError(err).Errorln("SafeReadTransactionRequest")
			return err
		}

		if req.SnapshotHash != "" {
			return nil
		}
		logger.Errorln("insufficient balance")
		return &mixin.Error{Code: mixin.InsufficientBalance, Description: "insufficient balance"}
	}

	builder := mixin.NewSafeTransactionBuilder(utxos)
	tx, err := m.client.MakeTransaction(ctx, builder, []*mixin.TransactionOutput{
		{
			Amount:  req.Amount,
			Address: mixin.RequireNewMixAddress(req.Opponents, req.Threshold),
		},
	})
	if err != nil {
		logger.WithError(err).Errorln("MakeTransaction")
		return err
	}

	raw, err := tx.Dump()
	if err != nil {
		logger.WithError(err).Errorln("Dump")
		return err
	}

	request, err := m.client.SafeCreateTransactionRequest(ctx, &mixin.SafeTransactionRequestInput{
		RequestID:      req.TraceID,
		RawTransaction: raw,
	})
	if err != nil {
		logger.WithError(err).Errorln("SafeCreateTransactionRequest")
		return err
	}

	if err := mixin.SafeSignTransaction(
		tx,
		m.spend,
		request.Views,
		0,
	); err != nil {
		logger.WithError(err).Errorln("SafeSignTransaction")
		return err
	}

	if raw, err = tx.Dump(); err != nil {
		logger.WithError(err).Errorln("Dump")
		return err
	}

	if _, err := m.client.SafeSubmitTransactionRequest(ctx, &mixin.SafeTransactionRequestInput{
		RequestID:      request.RequestID,
		RawTransaction: raw,
	}); err != nil {
		logger.WithError(err).Errorln("SafeSubmitTransactionRequest")
		return err
	}

	return nil
}

func (m *mixinBot) toSnapshots(ctx context.Context, items []*mixin.SafeSnapshot) ([]*core.Snapshot, error) {
	var snapshots = make([]*core.Snapshot, len(items))
	for i, s := range items {
		snapshot := &core.Snapshot{
			CreatedAt:  s.CreatedAt,
			SnapshotID: s.SnapshotID,
			UserID:     s.UserID,
			OpponentID: s.OpponentID,
			TraceID:    s.RequestID,
			AssetID:    s.AssetID,
			Amount:     s.Amount,
			Memo:       s.Memo,
		}

		if s.Deposit == nil && s.Deposit.DepositHash != "" {
			snapshot.TransactionHash = s.Deposit.DepositHash
		}

		snapshots[i] = snapshot
	}
	return snapshots, nil
}
