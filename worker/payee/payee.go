package payee

import (
	"context"
	"time"

	"github.com/asaskevich/govalidator"
	"github.com/fox-one/ftoken/core"
	"github.com/fox-one/pkg/logger"
	"github.com/fox-one/pkg/property"
	"github.com/patrickmn/go-cache"
	"golang.org/x/sync/errgroup"
)

const (
	snapshotCheckpoint = "ftoken:payee:checkpoint:snapshot"
)

type (
	Config struct {
		ClientID string `valid:"required"`
	}

	Worker struct {
		clientID string

		system       core.System
		properties   property.Store
		wallets      core.WalletStore
		walletz      core.WalletService
		orders       core.OrderStore
		transactions core.TransactionStore
		factories    map[string]core.Factory
		cache        *cache.Cache
	}
)

func New(
	cfg Config,
	system core.System,
	properties property.Store,
	orders core.OrderStore,
	transactions core.TransactionStore,
	wallets core.WalletStore,
	walletz core.WalletService,
	factories []core.Factory,
) *Worker {
	if _, err := govalidator.ValidateStruct(cfg); err != nil {
		panic(err)
	}

	factoryM := make(map[string]core.Factory, len(factories))
	for _, factory := range factories {
		factoryM[factory.GasAsset()] = factory
	}

	return &Worker{
		clientID:     cfg.ClientID,
		system:       system,
		properties:   properties,
		wallets:      wallets,
		walletz:      walletz,
		cache:        cache.New(time.Hour, time.Hour),
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
		return w.loopSnapshots(ctx)
	})
	g.Go(func() error {
		return w.loopPaidOrders(ctx)
	})
	return g.Wait()
}
