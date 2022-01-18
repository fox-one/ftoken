package wallet

import (
	"context"
	"time"

	"github.com/fox-one/ftoken/core"
	"github.com/fox-one/mixin-sdk-go"
)

type (
	Config struct {
		Pin string
	}

	mixinBot struct {
		client *mixin.Client
		pin    string
	}
)

func New(cfg Config, client *mixin.Client) core.WalletService {
	return &mixinBot{
		client: client,
		pin:    cfg.Pin,
	}
}

func (m *mixinBot) ListSnapshots(ctx context.Context, offset time.Time, limit int) ([]*core.Snapshot, error) {
	items, err := m.client.ReadNetworkSnapshots(ctx, "", offset, "ASC", limit)
	if err != nil {
		return nil, err
	}
	return convertSnapshots(items), nil
}

func (m *mixinBot) Transfer(ctx context.Context, req *core.Transfer) error {
	input := &mixin.TransferInput{
		AssetID: req.AssetID,
		Amount:  req.Amount,
		TraceID: req.TraceID,
		Memo:    req.Memo,
	}

	var err error
	if len(req.Opponents) == 1 {
		input.OpponentID = req.Opponents[0]
		_, err = m.client.Transfer(ctx, input, m.pin)
	} else {
		input.OpponentMultisig.Threshold = req.Threshold
		input.OpponentMultisig.Receivers = req.Opponents
		_, err = m.client.Transaction(ctx, input, m.pin)
	}

	if err != nil {
		if e, ok := err.(*mixin.Error); ok && e.Code == mixin.InvalidTraceID {
			return core.ErrInvalidTrace
		}
	}
	return err
}

func convertSnapshots(items []*mixin.Snapshot) []*core.Snapshot {
	var snapshots = make([]*core.Snapshot, len(items))
	for i, s := range items {
		snapshots[i] = &core.Snapshot{
			CreatedAt:  s.CreatedAt,
			SnapshotID: s.SnapshotID,
			UserID:     s.UserID,
			OpponentID: s.OpponentID,
			TraceID:    s.TraceID,
			AssetID:    s.AssetID,
			Amount:     s.Amount,
			Memo:       s.Memo,
		}
	}
	return snapshots
}
