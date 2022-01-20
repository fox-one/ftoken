package payee

import (
	"context"
	"encoding/base64"

	"github.com/fox-one/ftoken/core"
	"github.com/fox-one/ftoken/pkg/mtg"
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
	}

	if order.ID == 0 {
		order = &core.Order{
			CreatedAt: snapshot.CreatedAt,
			Version:   1,
			TraceID:   snapshot.TraceID,
			State:     core.OrderStateNew,
			UserID:    snapshot.OpponentID,
			FeeAsset:  snapshot.AssetID,
			FeeAmount: snapshot.Amount,
			Platform:  factory.Platform(),
		}

		data := []byte(snapshot.Memo)
		if d, err := base64.StdEncoding.DecodeString(snapshot.Memo); err == nil {
			data = d
		}

		if data, err := mtg.Scan(data, &order.Tokens); err != nil || len(order.Tokens) == 0 {
			log.Infoln("refund: scan tokens failed")
			return w.refundOrder(ctx, order)
		} else {
			var receiver core.Address
			if _, err := mtg.Scan(data, &receiver); err == nil && receiver.Destination != "" {
				order.Receiver = &receiver
			}
		}

		if err := w.orders.Create(ctx, order); err != nil {
			log.WithError(err).Errorln("orders.Create")
			return err
		}
	} else {
		if order.FeeAsset != snapshot.AssetID {
			log.WithField("order_asset", order.FeeAsset).Infoln("skip: asset not matched")
			return nil
		}
		if order.UserID == "" {
			order.UserID = snapshot.OpponentID
		}
	}
	if order.Receiver == nil && snapshot.OpponentID == "" {
		log.Infoln("skip: empty reciever / address")
		return nil
	}

	if order.State == core.OrderStateNew {
		order.FeeAmount = snapshot.Amount
		receiver := order.Receiver
		if receiver == nil {
			receiver = w.system.Addresses[order.FeeAsset]
		}

		tx, err := factory.CreateTransaction(ctx, order.Tokens, receiver)
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
		log.Infoln("refund: scan tokens failed")
		return w.refundOrder(ctx, order)
	}

	return nil
}
