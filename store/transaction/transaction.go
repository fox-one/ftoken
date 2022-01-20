package transaction

import (
	"context"

	"github.com/fox-one/ftoken/core"
	"github.com/fox-one/pkg/store"
	"github.com/fox-one/pkg/store/db"
)

func init() {
	db.RegisterMigrate(func(db *db.DB) error {
		tx := db.Update().Model(core.Transaction{})
		if err := tx.AutoMigrate(core.Transaction{}).Error; err != nil {
			return err
		}

		if err := tx.AddIndex("idx_transactions_trace", "trace_id").Error; err != nil {
			return err
		}

		if err := tx.AddUniqueIndex("idx_transactions_hash", "hash").Error; err != nil {
			return err
		}

		return nil
	})
}

func New(db *db.DB) core.TransactionStore {
	return &transactionStore{db: db}
}

type transactionStore struct {
	db *db.DB
}

func (s *transactionStore) Create(ctx context.Context, tx *core.Transaction) error {
	return s.db.Update().Model(tx).Where("trace_id = ?", tx.TraceID).FirstOrCreate(tx).Error
}

func (s *transactionStore) Update(ctx context.Context, tx *core.Transaction) error {
	params := toUpdateParams(tx)
	if query := s.db.Update().Model(tx).Where("version = ?", tx.Version).Updates(params); query.Error != nil {
		return query.Error
	} else if query.RowsAffected == 0 {
		return db.ErrOptimisticLock
	}
	return nil
}

func toUpdateParams(tx *core.Transaction) map[string]interface{} {
	return map[string]interface{}{
		"version": tx.Version + 1,
		"state":   tx.State,
		"gas":     tx.Gas,
		"tokens":  tx.Tokens,
	}
}

func (s *transactionStore) Find(ctx context.Context, hash string) (*core.Transaction, error) {
	var tx core.Transaction
	if db := s.db.View().Where("hash = ?", hash).First(&tx); db.Error != nil {
		if store.IsErrNotFound(db.Error) {
			return &core.Transaction{}, nil
		}
		return nil, db.Error
	}
	return &tx, nil
}

func (s *transactionStore) FindTrace(ctx context.Context, traceID string) ([]*core.Transaction, error) {
	var txs []*core.Transaction
	if err := s.db.View().Where("trace_id = ?", traceID).Find(&txs).Error; err != nil {
		return nil, err
	}
	return txs, nil
}
