package order

import (
	"context"
	"time"

	"github.com/fox-one/ftoken/core"
	"github.com/fox-one/pkg/logger"
	"github.com/sirupsen/logrus"
)

func (w *Worker) loopPaidOrders(ctx context.Context) error {
	log := logger.FromContext(ctx).WithField("loop", "paid-orders")
	ctx = logger.WithContext(ctx, log)

	dur := time.Millisecond
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(dur):
			if err := w.loopPaidOrdersOnce(ctx); err == nil {
				dur = 100 * time.Millisecond
			} else {
				dur = time.Second
			}
		}
	}
}

func (w *Worker) loopPaidOrdersOnce(ctx context.Context) error {
	const LIMIT = 10
	log := logger.FromContext(ctx)

	orders, err := w.orders.List(ctx, core.OrderStatePaid, LIMIT)
	if err != nil {
		log.WithError(err).Errorln("list orders")
		return err
	}

	for _, order := range orders {
		if err := w.handlePaidOrder(ctx, order); err != nil {
			return err
		}
	}

	return nil
}

func (w *Worker) handlePaidOrder(ctx context.Context, order *core.Order) error {
	log := logger.FromContext(ctx).WithFields(logrus.Fields{
		"order":     order.TraceID,
		"fee_asset": order.FeeAsset,
	})
	ctx = logger.WithContext(ctx, log)

	factory, ok := w.factories[order.FeeAsset]
	if !ok {
		log.Errorln("factory not found")
		return nil
	}

	tx, err := w.latestTransaction(ctx, order)
	if err != nil {
		return err
	}

	if tx.ID == 0 {
		tx, err = factory.CreateTransaction(ctx, order.TokenRequests, order.TraceID)
		if err != nil {
			log.WithError(err).Errorln("factory.CreateTransaction failed")
			return err
		}

		tx.TraceID = order.TraceID
		if err := w.transactions.Create(ctx, tx); err != nil {
			log.WithError(err).Errorln("transactions.Create failed")
			return err
		}
	}

	switch tx.State {
	case core.TransactionStateNew:
		if err := factory.SendTransaction(ctx, tx); err != nil {
			log.WithError(err).Errorln("factory.SendTransaction failed")
			return err
		}
		tx.State = core.TransactionStatePending
		if err := w.transactions.Update(ctx, tx); err != nil {
			log.WithError(err).Errorln("transactions.Update failed")
			return err
		}
		order.State = core.OrderStateProcessing
		order.Transaction = tx.Hash

	case core.TransactionStatePending:
		order.State = core.OrderStateProcessing
		order.Transaction = tx.Hash

	case core.TransactionStateFailed, core.TransactionStateSuccess:
		panic("should never happen")
	}

	if err := w.orders.Update(ctx, order); err != nil {
		log.WithError(err).Errorln("orders.Update failed")
		return err
	}

	return nil
}

func (w *Worker) latestTransaction(ctx context.Context, order *core.Order) (*core.Transaction, error) {
	txs, err := w.transactions.FindTrace(ctx, order.TraceID)
	if err != nil {
		logger.FromContext(ctx).WithError(err).Errorln("transactions.Find failed")
		return nil, err
	}
	if len(txs) > 0 {
		if tx := txs[len(txs)-1]; tx.Hash != order.Transaction {
			return tx, nil
		}
	}
	return &core.Transaction{}, nil
}
