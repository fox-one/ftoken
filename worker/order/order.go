package order

import (
	"context"

	"github.com/fox-one/ftoken/core"
	"github.com/fox-one/pkg/logger"
	"golang.org/x/sync/errgroup"
)

type (
	Worker struct {
		system       core.System
		wallets      core.WalletStore
		orders       core.OrderStore
		transactions core.TransactionStore
		factories    map[string]core.Factory
	}
)

func New(
	system core.System,
	orders core.OrderStore,
	transactions core.TransactionStore,
	wallets core.WalletStore,
	factories []core.Factory,
) *Worker {
	factoryM := make(map[string]core.Factory, len(factories))
	for _, factory := range factories {
		factoryM[factory.GasAsset()] = factory
	}

	return &Worker{
		system:       system,
		wallets:      wallets,
		orders:       orders,
		transactions: transactions,
		factories:    factoryM,
	}
}

func (w *Worker) Run(ctx context.Context) error {
	log := logger.FromContext(ctx).WithField("worker", "payee")
	ctx = logger.WithContext(ctx, log)

	g, ctx := errgroup.WithContext(ctx)
	g.Go(func() error {
		return w.loopPaidOrders(ctx)
	})
	g.Go(func() error {
		return w.loopProcessingOrders(ctx)
	})
	return g.Wait()
}
