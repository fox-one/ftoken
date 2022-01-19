package order

import (
	"net/http"

	"github.com/fox-one/ftoken/core"
	"github.com/fox-one/ftoken/handler/render"
	"github.com/fox-one/ftoken/handler/render/views"
	"github.com/fox-one/pkg/httputil/param"
	"github.com/fox-one/pkg/uuid"
	"github.com/go-chi/chi"
	"github.com/shopspring/decimal"
	"github.com/twitchtv/twirp"
)

func HandleCreateOrder(system core.System, walletz core.WalletService, orders core.OrderStore, factories []core.Factory) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		var body struct {
			TraceID         string        `json:"trace_id,omitempty"`
			Platform        string        `json:"platform,omitempty"`
			Tokens          core.Tokens   `json:"tokens,omitempty"`
			UserID          string        `json:"user_id,omitempty"`
			ReceiverAddress *core.Address `json:"receiver_address,omitempty"`
		}

		if err := param.Binding(r, &body); err != nil {
			render.Error(w, twirp.InternalErrorWith(err))
			return
		}

		if body.TraceID == "" {
			body.TraceID = uuid.New()
		}

		if len(body.Tokens) == 0 {
			render.Error(w, twirp.RequiredArgumentError("tokens"))
			return
		}

		var factory core.Factory
		for _, f := range factories {
			if f.Platform() != body.Platform {
				continue
			}
			factory = f
		}

		if factory == nil {
			render.Error(w, twirp.RequiredArgumentError("platform"))
			return
		}

		order, err := orders.Find(ctx, body.TraceID)
		if err != nil {
			render.Error(w, twirp.InternalErrorWith(err))
			return
		}

		if order.ID == 0 {
			order = &core.Order{
				Version:  1,
				TraceID:  body.TraceID,
				State:    core.OrderStatePending,
				UserID:   body.UserID,
				FeeAsset: factory.GasAsset(),
				Platform: body.Platform,
				Tokens:   body.Tokens,
			}
			if body.ReceiverAddress != nil && body.ReceiverAddress.Destination != "" {
				order.Receiver = body.ReceiverAddress
			}
			if err := orders.Create(ctx, order); err != nil {
				render.Error(w, twirp.InternalErrorWith(err))
				return
			}
		} else if order.UserID != "" && order.UserID != body.UserID {
			render.Error(w, twirp.NewErrorf(twirp.AlreadyExists, "order with trace already exists"))
			return
		}

		tx, err := factory.CreateTransaction(ctx, body.Tokens, system.Addresses[factory.GasAsset()])
		if err != nil {
			render.Error(w, twirp.InternalErrorWith(err))
			return
		}

		order.FeeAmount = tx.Gas.Mul(decimal.New(5, 0))
		render.JSON(w, order)
	}
}

func HandleFetchOrder(orders core.OrderStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		traceID := chi.URLParam(r, "trace_id")
		order, err := orders.Find(ctx, traceID)
		if err != nil {
			render.Error(w, twirp.InternalErrorWith(err))
			return
		} else if order.ID == 0 {
			render.Error(w, twirp.NotFoundError("order not found"))
			return
		}

		render.JSON(w, views.OrderView(*order, false))
	}
}
