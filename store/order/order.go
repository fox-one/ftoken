package order

import (
	"context"

	"github.com/fox-one/ftoken/core"
	"github.com/fox-one/pkg/store"
	"github.com/fox-one/pkg/store/db"
)

func init() {
	db.RegisterMigrate(func(db *db.DB) error {
		tx := db.Update().Model(core.Order{})
		if err := tx.AutoMigrate(core.Order{}).Error; err != nil {
			return err
		}

		if err := tx.AddUniqueIndex("idx_orders_trace", "trace_id").Error; err != nil {
			return err
		}

		return nil
	})
}

func New(db *db.DB) core.OrderStore {
	return &orderStore{db: db}
}

type orderStore struct {
	db *db.DB
}

func (s *orderStore) Create(ctx context.Context, order *core.Order) error {
	return s.db.Update().Model(order).Where("trace_id = ?", order.TraceID).FirstOrCreate(order).Error
}

func (s *orderStore) Update(ctx context.Context, order *core.Order) error {
	params := toUpdateParams(order)
	if query := s.db.Update().Model(order).Where("version = ?", order.Version).Updates(params); query.Error != nil {
		return query.Error
	} else if query.RowsAffected == 0 {
		return db.ErrOptimisticLock
	}
	return nil
}

func toUpdateParams(order *core.Order) map[string]interface{} {
	return map[string]interface{}{
		"version":     order.Version + 1,
		"user_id":     order.UserID,
		"state":       order.State,
		"result":      order.Result,
		"gas_usage":   order.GasUsage,
		"transaction": order.Transaction,
	}
}

func (s *orderStore) Find(ctx context.Context, traceID string) (*core.Order, error) {
	var order core.Order
	if db := s.db.View().Where("trace_id = ?", traceID).First(&order); db.Error != nil {
		if store.IsErrNotFound(db.Error) {
			return &core.Order{}, nil
		}
		return nil, db.Error
	}
	return &order, nil
}

func (s *orderStore) List(ctx context.Context, state core.OrderState, limit int) ([]*core.Order, error) {
	var orders []*core.Order
	query := s.db.View().Limit(limit)

	switch state {
	case core.OrderStateNew, core.OrderStateProcessing, core.OrderStatePaid, core.OrderStateFailed, core.OrderStateDone:
		query = query.Where("state = ?", state)
	}

	if err := query.Find(&orders).Error; err != nil {
		return nil, err
	}
	return orders, nil
}
