package wallet

import (
	"context"
	"net/http"
	"time"

	"github.com/fox-one/ftoken/core"
	"github.com/fox-one/mixin-sdk-go/v2"
	"github.com/pandodao/safe-wallet/handler/rpc/safewallet"
)

type (
	Config struct {
		SafeWalletHost string
	}

	mixinBot struct {
		client      *mixin.Client
		safeWallets safewallet.SafeWalletService
	}
)

func New(cfg Config, client *mixin.Client) core.WalletService {
	return &mixinBot{
		client:      client,
		safeWallets: safewallet.NewSafeWalletServiceProtobufClient(cfg.SafeWalletHost, http.DefaultClient),
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
	if _, err := m.safeWallets.CreateTransfer(ctx, &safewallet.CreateTransferRequest{
		AssetId:   req.AssetID,
		TraceId:   req.TraceID,
		Amount:    req.Amount.String(),
		Memo:      req.Memo,
		Opponents: req.Opponents,
		Threshold: uint32(req.Threshold),
	}); err != nil {
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
