package payee

import (
	"context"

	"github.com/fox-one/ftoken/core"
	"github.com/fox-one/pkg/logger"
	"github.com/fox-one/pkg/uuid"
	"github.com/lib/pq"
)

func (w *Worker) refundOrder(ctx context.Context, order *core.Order) error {
	if err := w.wallets.CreateTransfers(ctx, []*core.Transfer{
		{
			TraceID:   uuid.Modify(order.TraceID, "refund"),
			AssetID:   order.FeeAsset,
			Amount:    order.FeeAmount.Sub(order.GasUsage),
			Memo:      "refund " + order.TraceID,
			Threshold: 1,
			Opponents: pq.StringArray{order.UserID},
		},
	}); err != nil {
		logger.FromContext(ctx).Errorln("CreateTransfers failed")
		return err
	}
	return nil
}
