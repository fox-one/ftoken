package order

import (
	"context"
	"time"

	"github.com/fox-one/ftoken/core"
	"github.com/fox-one/pkg/logger"
	"github.com/sirupsen/logrus"
)

func (w *Worker) loopProcessingOrders(ctx context.Context) error {
	log := logger.FromContext(ctx).WithField("loop", "processing-orders")
	ctx = logger.WithContext(ctx, log)

	dur := time.Millisecond
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(dur):
			if err := w.loopProcessingOrdersOnce(ctx); err == nil {
				dur = 100 * time.Millisecond
			} else {
				dur = time.Second
			}
		}
	}
}

func (w *Worker) loopProcessingOrdersOnce(ctx context.Context) error {
	const LIMIT = 10
	log := logger.FromContext(ctx)

	orders, err := w.orders.List(ctx, core.OrderStateProcessing, LIMIT)
	if err != nil {
		log.WithError(err).Errorln("list orders")
		return err
	}

	for _, order := range orders {
		if err := w.handleProcessingOrder(ctx, order); err != nil {
			return err
		}
	}

	return nil
}

func (w *Worker) handleProcessingOrder(ctx context.Context, order *core.Order) error {
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

	tx, err := w.transactions.Find(ctx, order.Transaction)
	if err != nil {
		return err
	}

	if tx.ID == 0 {
		panic("should never happen")
	}

	switch tx.State {
	case core.TransactionStateNew:
		panic("should never happen")

	case core.TransactionStatePending:
		tx1, err := factory.ReadTransaction(ctx, tx.Hash)
		if err != nil {
			log.WithError(err).Errorln("factory.ReadTransaction failed")
			return err
		}
		switch tx1.State {
		case core.TransactionStateNew:
			panic("should never happen")
		case core.TransactionStatePending:
			return nil
		case core.TransactionStateSuccess:
			order.State = core.OrderStateDone
			tx.Tokens = tx1.Tokens
			order.GasUsage = order.GasUsage.Add(tx.Gas)
		case core.TransactionStateFailed:
			order.State = core.OrderStateFailed
			order.GasUsage = order.GasUsage.Add(tx.Gas)
		}
		tx.Gas = tx1.Gas
		tx.State = tx1.State
		if err := w.transactions.Update(ctx, tx); err != nil {
			log.WithError(err).Errorln("transactions.Update failed")
			return err
		}

	case core.TransactionStateFailed:
		order.GasUsage = order.GasUsage.Add(tx.Gas)
		order.State = core.OrderStateFailed

	case core.TransactionStateSuccess:
		order.GasUsage = order.GasUsage.Add(tx.Gas)
		order.State = core.OrderStateDone
	}

	switch order.State {
	case core.OrderStateNew, core.OrderStatePaid:
		panic("should never happen")
	case core.OrderStateProcessing:
		return nil
	case core.OrderStateFailed:
		if err := w.refundOrder(ctx, order); err != nil {
			return err
		}
	case core.OrderStateDone:
		order.Tokens = tx.Tokens
	}

	if err := w.orders.Update(ctx, order); err != nil {
		log.WithError(err).Errorln("orders.Update failed")
		return err
	}
	return nil
}
