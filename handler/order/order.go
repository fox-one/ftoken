package order

import (
	"bytes"
	"net/http"

	"github.com/fox-one/ftoken/core"
	"github.com/fox-one/ftoken/handler/render"
	"github.com/fox-one/ftoken/handler/render/views"
	"github.com/fox-one/pkg/httputil/param"
	"github.com/fox-one/pkg/uuid"
	"github.com/go-chi/chi"
	"github.com/twitchtv/twirp"
)

func HandleEstimateGas(system core.System, factories []core.Factory) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		var body struct {
			Platform string `json:"platform,omitempty"`
			Count    int    `json:"count"`
		}

		if err := param.Binding(r, &body); err != nil {
			render.Error(w, twirp.InternalErrorWith(err))
			return
		}

		if body.Count < 1 {
			body.Count = 1
		}

		var factory core.Factory
		for _, f := range factories {
			if f.Platform() != body.Platform {
				continue
			}
			factory = f
		}

		var tokens = make(core.Tokens, body.Count)
		for i := 0; i < body.Count; i++ {
			tokens[i] = &core.Token{
				TotalSupply: 1000000,
			}
		}

		if factory == nil {
			render.Error(w, twirp.RequiredArgumentError("platform"))
			return
		}

		tx, err := factory.CreateTransaction(ctx, tokens, system.Addresses[factory.GasAsset()])
		if err != nil {
			render.Error(w, twirp.InternalErrorWith(err))
			return
		}

		feeAmount := tx.Gas.Mul(system.Gas.Multiplier)
		if min, ok := system.Gas.Mins[factory.Platform()]; ok && feeAmount.LessThan(min) {
			feeAmount = min
		}
		render.JSON(w, render.H{
			"fee_asset":  factory.GasAsset(),
			"fee_amount": feeAmount,
		})
	}
}

func HandleCreateOrder(system core.System, walletz core.WalletService, orders core.OrderStore, factories []core.Factory) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		var body struct {
			TraceID         string        `json:"trace_id,omitempty"`
			Platform        string        `json:"platform,omitempty"`
			Tokens          core.Tokens   `json:"tokens,omitempty"`
			ReceiverAddress *core.Address `json:"receiver_address,omitempty"`
		}

		if err := param.Binding(r, &body); err != nil {
			render.Error(w, twirp.InternalErrorWith(err))
			return
		}

		if body.TraceID == "" {
			body.TraceID = uuid.New()
		}

		var tokens core.Tokens
		for _, token := range body.Tokens {
			if token.Name != "" && token.Symbol != "" && token.TotalSupply > 0 {
				tokens = append(tokens, token)
			}
		}
		if len(tokens) == 0 {
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
				State:    core.OrderStateNew,
				FeeAsset: factory.GasAsset(),
				Platform: body.Platform,
				Tokens:   tokens,
			}
			if body.ReceiverAddress != nil && body.ReceiverAddress.Destination != "" {
				order.Receiver = body.ReceiverAddress
			}
			if err := orders.Create(ctx, order); err != nil {
				render.Error(w, twirp.InternalErrorWith(err))
				return
			}
		} else {
			t1, _ := core.EncodeTokens(order.Tokens)
			t2, _ := core.EncodeTokens(tokens)
			if order.UserID != "" && order.Receiver.Destination != body.ReceiverAddress.Destination || !bytes.Equal(t1, t2) {
				render.Error(w, twirp.NewErrorf(twirp.AlreadyExists, "order with trace already exists"))
				return
			}
		}

		tx, err := factory.CreateTransaction(ctx, tokens, system.Addresses[factory.GasAsset()])
		if err != nil {
			render.Error(w, twirp.InternalErrorWith(err))
			return
		}

		order.FeeAmount = tx.Gas.Mul(system.Gas.Multiplier)
		if min, ok := system.Gas.Mins[factory.Platform()]; ok && order.FeeAmount.LessThan(min) {
			order.FeeAmount = min
		}
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
