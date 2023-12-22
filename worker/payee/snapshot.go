package payee

import (
	"context"
	"encoding/base64"
	"encoding/hex"
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

	memo := snapshot.Memo
	if m, err := hex.DecodeString(snapshot.Memo); err == nil {
		memo = string(m)
	}

	log := logger.FromContext(ctx).WithFields(logrus.Fields{
		"trace":     snapshot.TraceID,
		"pay_asset": snapshot.AssetID,
		"amount":    snapshot.Amount,
		"memo":      memo,
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

		data := []byte(memo)
		if d, err := base64.StdEncoding.DecodeString(memo); err == nil {
			data = d
		}

		if err := json.Unmarshal(data, &order.TokenRequests); err != nil || len(order.TokenRequests) == 0 {
			log.Infoln("skip: scan tokens failed")
			return nil
		}

		if err := w.orders.Create(ctx, order); err != nil {
			log.WithError(err).Errorln("orders.Create")
			return err
		}
	}

	fee := w.system.Fees[factory.Platform()]
	if snapshot.Amount.LessThan(fee.FeeAmount.Mul(decimal.NewFromInt(int64(len(order.TokenRequests))))) {
		log.WithField("tokens:count", len(order.TokenRequests)).Infoln("skip: not enough fee")
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
