package payee

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/asaskevich/govalidator"
	"github.com/fox-one/ftoken/core"
	"github.com/fox-one/pkg/logger"
	"github.com/fox-one/pkg/property"
	"github.com/patrickmn/go-cache"
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
		assets       core.AssetStore
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
	assets core.AssetStore,
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
		assets:       assets,
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

	dur := time.Millisecond
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(dur):
			if err := w.run(ctx); err == nil {
				dur = 100 * time.Millisecond
			} else {
				dur = time.Second
			}
		}
	}
}

func (w *Worker) run(ctx context.Context) error {
	const LIMIT = 500
	log := logger.FromContext(ctx)

	v, err := w.properties.Get(ctx, snapshotCheckpoint)
	if err != nil {
		log.WithError(err).Errorln("properties.Get")
		return err
	}

	var (
		offset    = v.Time()
		newOffset = offset
	)

	snapshots, err := w.walletz.ListSnapshots(ctx, offset, LIMIT)
	if err != nil {
		log.WithError(err).Errorln("list snapshots")
		return err
	}

	for _, snapshot := range snapshots {
		newOffset = snapshot.CreatedAt
		if snapshot.UserID != w.clientID || snapshot.Amount.IsNegative() {
			continue
		}

		snapshotKey := fmt.Sprintf("snapshot:%s", snapshot.SnapshotID)
		if _, ok := w.cache.Get(snapshotKey); !ok {
			if err := w.handleSnapshot(ctx, snapshot); err != nil {
				return err
			}
			w.cache.SetDefault(snapshotKey, true)
		}
	}

	if newOffset.Equal(offset) {
		return errors.New("empty list")
	}

	if err := w.properties.Save(ctx, snapshotCheckpoint, newOffset); err != nil {
		log.WithError(err).Errorln("properties.Save", snapshotCheckpoint)
		return err
	}
	return nil
}
