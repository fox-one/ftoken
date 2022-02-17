package payee

import (
	"context"

	"github.com/fox-one/ftoken/core"
	"github.com/fox-one/pkg/logger"
	"github.com/sirupsen/logrus"
)

func (w *Worker) handleSnapshot(ctx context.Context, snapshot *core.Snapshot) error {
	factory, ok := w.factories[snapshot.AssetID]
	if !ok {
		return nil
	}

	log := logger.FromContext(ctx).WithFields(logrus.Fields{
		"trace":     snapshot.TraceID,
		"pay_asset": snapshot.AssetID,
		"amount":    snapshot.Amount,
		"memo":      snapshot.Memo,
		"platform":  factory.Platform(),
	})
	ctx = logger.WithContext(ctx, log)

	order, err := w.orders.Find(ctx, snapshot.TraceID)
	if err != nil {
		log.WithError(err).Errorln("orders.Find")
		return err
	} else if order.ID == 0 {
		log.WithField("order_id", order.ID).WithField("state", order.State).Infoln("skip: order not exist")
		return nil
	}

	if order.FeeAsset != snapshot.AssetID {
		log.WithField("order_asset", order.FeeAsset).Infoln("skip: asset not matched")
		return nil
	}

	if order.State == core.OrderStateNew {
		order.FeeAmount = snapshot.Amount

		tx, err := factory.CreateTransaction(ctx, order.TokenRequests, order.TraceID)
		if err != nil {
			log.WithError(err).Errorln("factory.CreateTransaction")
			return err
		}

		if snapshot.Amount.LessThan(tx.Gas.Mul(w.system.Gas.StrictMultiplier)) {
			order.State = core.OrderStateFailed
		} else if min, ok := w.system.Gas.Mins[order.Platform]; ok && snapshot.Amount.LessThan(min) {
			order.State = core.OrderStateFailed
		} else {
			order.State = core.OrderStatePaid
		}

		if err := w.orders.Update(ctx, order); err != nil {
			log.WithError(err).Errorln("orders.Update")
			return err
		}
	}

	if order.State == core.OrderStateFailed {
		log.Infoln("refund: order rejected")
		return w.refundOrder(ctx, order)
	}

	return nil
}
