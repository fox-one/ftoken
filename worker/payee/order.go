package payee

import (
	"context"
	"time"

	"github.com/fox-one/ftoken/core"
	"github.com/fox-one/pkg/logger"
	"github.com/sirupsen/logrus"
)

func (w *Worker) loopPaidOrders(ctx context.Context) error {
	log := logger.FromContext(ctx).WithField("loop", "snapshots")
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

	tx, err := w.transactions.Find(ctx, order.TraceID)
	if err != nil {
		log.WithError(err).Errorln("transactions.Find failed")
		return err
	}

	if tx.ID == 0 {
		tx, err = factory.CreateTransaction(ctx, order.Tokens, order.Receiver)
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
			tx.Gas = tx1.Gas
			tx.Tokens = tx1.Tokens
		case core.TransactionStateFailed:
			order.State = core.OrderStateFailed
			tx.Gas = tx1.Gas
		}
		if err := w.transactions.Update(ctx, tx); err != nil {
			log.WithError(err).Errorln("transactions.Update failed")
			return err
		}

	case core.TransactionStateFailed:
		order.State = core.OrderStateFailed

	case core.TransactionStateSuccess:
		order.State = core.OrderStateDone
	}

	switch order.State {
	case core.OrderStatePending:
		panic("should never happen")
	case core.OrderStatePaid:
		return nil
	case core.OrderStateFailed:
		if err := w.refundSnapshot(ctx, order.TraceID, order.UserID, order.FeeAsset, order.FeeAmount.Sub(tx.Gas), "deploy failed"); err != nil {
			return err
		}

	case core.OrderStateDone:
		order.Result = tx.Tokens
	}

	if order.State != core.OrderStatePaid {
		if err := w.orders.Update(ctx, order); err != nil {
			log.WithError(err).Errorln("orders.Update failed")
			return err
		}
	}
	return nil
}
