package payee

import (
	"context"
	"encoding/base64"
	"encoding/json"

	"github.com/fox-one/ftoken/core"
	"github.com/fox-one/pkg/logger"
	"github.com/shopspring/decimal"
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

		if err := json.Unmarshal(data, &order.Tokens); err != nil || len(order.Tokens) == 0 {
			log.Infoln("skip: scan tokens failed")
			return nil
		}

		if err := w.orders.Create(ctx, order); err != nil {
			log.WithError(err).Errorln("orders.Create")
			return err
		}
	}

	fee := w.system.Fees[factory.Platform()]
	if snapshot.Amount.LessThan(fee.FeeAmount.Mul(decimal.NewFromInt(int64(len(order.Tokens))))) {
		log.WithField("tokens:count", len(order.Tokens)).Infoln("skip: not enough fee")
		return nil
	}

	if order.State == core.OrderStateNew {
		order.UserID = snapshot.OpponentID
		order.FeeAmount = snapshot.Amount
		order.State = core.OrderStatePaid

		if err := w.orders.Update(ctx, order); err != nil {
			log.WithError(err).Errorln("orders.Update")
			return err
		}
	}

	return nil
}
