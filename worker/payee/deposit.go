package payee

import (
	"context"

	"github.com/fox-one/ftoken/core"
	"github.com/fox-one/pkg/logger"
	"github.com/fox-one/pkg/uuid"
	"github.com/lib/pq"
	"github.com/sirupsen/logrus"
)

func (w *Worker) handleDepositSnapshot(ctx context.Context, snapshot *core.Snapshot) error {
	log := logger.FromContext(ctx).WithFields(logrus.Fields{
		"snapshot_id": snapshot.SnapshotID,
		"pay_asset":   snapshot.AssetID,
		"amount":      snapshot.Amount,
		"transaction": snapshot.TransactionHash,
	})
	ctx = logger.WithContext(ctx, log)

	tx, err := w.transactions.Find(ctx, snapshot.TransactionHash)
	if err != nil {
		log.WithError(err).Errorln("transactions.Find")
		return err
	}

	if tx.ID == 0 {
		return nil
	}

	order, err := w.orders.Find(ctx, tx.TraceID)
	if err != nil {
		log.WithError(err).Errorln("orders.Find")
		return err
	}

	if order.ID > 0 && order.UserID != "" {
		if err := w.wallets.CreateTransfers(ctx, []*core.Transfer{
			{
				TraceID:   uuid.Modify(snapshot.SnapshotID, "forward"),
				AssetID:   snapshot.AssetID,
				Amount:    snapshot.Amount,
				Threshold: 1,
				Opponents: pq.StringArray{order.UserID},
			},
		}); err != nil {
			logger.FromContext(ctx).Errorln("CreateTransfers failed")
			return err
		}
		return nil
	}

	return nil
}
