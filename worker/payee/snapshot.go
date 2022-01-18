package payee

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/fox-one/ftoken/core"
	"github.com/fox-one/pkg/logger"
	"github.com/fox-one/pkg/uuid"
	"github.com/lib/pq"
	"github.com/shopspring/decimal"
	"github.com/sirupsen/logrus"
)

func (w *Worker) loopSnapshots(ctx context.Context) error {
	log := logger.FromContext(ctx).WithField("loop", "snapshots")
	ctx = logger.WithContext(ctx, log)

	dur := time.Millisecond
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(dur):
			if err := w.loopSnapshotOnce(ctx); err == nil {
				dur = 100 * time.Millisecond
			} else {
				dur = time.Second
			}
		}
	}
}

func (w *Worker) loopSnapshotOnce(ctx context.Context) error {
	const LIMIT = 500
	log := logger.FromContext(ctx)

	v, err := w.properties.Get(ctx, snapshotCheckpoint)
	if err != nil {
		log.WithError(err).Errorln("properties.Get")
		return err
	}

	var (
		offset    = v.Time()
		newOffset = offset
	)

	snapshots, err := w.walletz.ListSnapshots(ctx, offset, LIMIT)
	if err != nil {
		log.WithError(err).Errorln("list snapshots")
		return err
	}

	for _, snapshot := range snapshots {
		newOffset = snapshot.CreatedAt
		if snapshot.UserID != w.clientID || snapshot.Amount.IsNegative() || snapshot.OpponentID == "" {
			continue
		}

		snapshotKey := fmt.Sprintf("snapshot:%s", snapshot.SnapshotID)
		if _, ok := w.cache.Get(snapshotKey); !ok {
			if err := w.handleSnapshot(ctx, snapshot); err != nil {
				return err
			}
			w.cache.SetDefault(snapshotKey, true)
		}
	}

	if newOffset.Equal(offset) {
		return errors.New("empty list")
	}

	if err := w.properties.Save(ctx, snapshotCheckpoint, newOffset); err != nil {
		log.WithError(err).Errorln("properties.Save", snapshotCheckpoint)
		return err
	}
	return nil
}

func (w *Worker) handleSnapshot(ctx context.Context, snapshot *core.Snapshot) error {
	factory, ok := w.factories[snapshot.AssetID]
	if !ok {
		return nil
	}

	log := logger.FromContext(ctx).WithFields(logrus.Fields{
		"trace":     snapshot.TraceID,
		"pay_asset": snapshot.AssetID,
		"asset":     snapshot.AssetID,
		"amount":    snapshot.Amount,
		"memo":      snapshot.Memo,
	})
	ctx = logger.WithContext(ctx, log)

	data := []byte(snapshot.Memo)
	if d, err := base64.StdEncoding.DecodeString(snapshot.Memo); err == nil {
		data = d
	}

	var input struct {
		Tokens   []byte `json:"t"`
		Receiver string `json:"r"`
	}

	if err := json.Unmarshal(data, &input); err != nil || len(input.Tokens) == 0 || input.Receiver == "" {
		return w.refundSnapshot(ctx, snapshot.TraceID, snapshot.OpponentID, snapshot.AssetID, snapshot.Amount)
	}

	tokens := core.DecodeTokens(input.Tokens)
	if len(tokens) == 0 {
		return w.refundSnapshot(ctx, snapshot.TraceID, snapshot.OpponentID, snapshot.AssetID, snapshot.Amount)
	}

	order, err := w.orders.Find(ctx, snapshot.TraceID)
	if err != nil {
		log.WithError(err).Errorln("orders.Find")
		return err
	}

	if order.ID == 0 {
		order = &core.Order{
			CreatedAt: snapshot.CreatedAt,
			TraceID:   snapshot.TraceID,
			State:     core.OrderStatePaid,
			UserID:    snapshot.OpponentID,
			FeeAsset:  snapshot.AssetID,
			FeeAmount: snapshot.Amount,
			Platform:  factory.Platform(),
			Tokens:    tokens,
		}

		tx, err := factory.CreateTransaction(ctx, tokens, input.Receiver)
		if err != nil {
			log.WithError(err).Errorln("factory.CreateTransaction")
			return err
		}

		if tx.Gas.GreaterThan(snapshot.Amount.Div(decimal.New(4, 0))) {
			return w.refundSnapshot(ctx, snapshot.TraceID, snapshot.OpponentID, snapshot.AssetID, snapshot.Amount, "payment too low to cover the gas fee")
		}

		if err := w.orders.Create(ctx, order); err != nil {
			log.WithError(err).Errorln("orders.Create")
			return err
		}
	}

	return nil
}

func (w *Worker) refundSnapshot(ctx context.Context, trace, opponent, asset string, amount decimal.Decimal, msg ...string) error {
	memo := "refund " + trace
	if len(msg) > 0 && msg[0] != "" {
		memo = memo + ", " + msg[0]
	}

	if err := w.wallets.CreateTransfers(ctx, []*core.Transfer{
		{
			TraceID:   uuid.Modify(trace, "refund"),
			AssetID:   asset,
			Amount:    amount,
			Memo:      memo,
			Threshold: 1,
			Opponents: pq.StringArray{opponent},
		},
	}); err != nil {
		logger.FromContext(ctx).Errorln("CreateTransfers failed")
		return err
	}
	return nil
}
