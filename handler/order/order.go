package order

import (
	"bytes"
	"net/http"

	"github.com/fox-one/ftoken/core"
	"github.com/fox-one/ftoken/handler/render"
	"github.com/fox-one/ftoken/handler/render/views"
	"github.com/fox-one/ftoken/pkg/token"
	"github.com/fox-one/pkg/httputil/param"
	"github.com/fox-one/pkg/uuid"
	"github.com/go-chi/chi/v5"
	"github.com/shopspring/decimal"
	"github.com/twitchtv/twirp"
)

func HandleEstimateGas(system core.System, factories []core.Factory) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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

		var tokens = make(core.TokenItems, body.Count)
		for i := 0; i < body.Count; i++ {
			tokens[i] = &core.TokenItem{
				TotalSupply: 1000000,
			}
		}

		if factory == nil {
			render.Error(w, twirp.RequiredArgumentError("platform"))
			return
		}

		fee := system.Fees[factory.Platform()]
		render.JSON(w, render.H{
			"fee_asset":  fee.FeeAssetID,
			"fee_amount": fee.FeeAmount.Mul(decimal.NewFromInt(int64(body.Count))),
		})
	}
}

func HandleCreateOrder(
	system core.System,
	assets core.AssetStore,
	walletz core.WalletService,
	orders core.OrderStore,
	factories []core.Factory,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		var body struct {
			TraceID  string             `json:"trace_id,omitempty"`
			Platform string             `json:"platform,omitempty"`
			Tokens   core.TokenRequests `json:"tokens,omitempty"`
		}

		if err := param.Binding(r, &body); err != nil {
			render.Error(w, twirp.InternalErrorWith(err))
			return
		}

		if body.TraceID == "" {
			body.TraceID = uuid.New()
		}

		tokens, err := token.ExportTokenItems(ctx, assets, body.Tokens)
		if err != nil {
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
				Version:       1,
				TraceID:       body.TraceID,
				State:         core.OrderStateNew,
				FeeAsset:      factory.GasAsset(),
				Platform:      body.Platform,
				TokenRequests: tokens,
			}
			if err := orders.Create(ctx, order); err != nil {
				render.Error(w, twirp.InternalErrorWith(err))
				return
			}
		} else {
			t1, _ := core.EncodeTokens(order.TokenRequests)
			t2, _ := core.EncodeTokens(tokens)
			if !bytes.Equal(t1, t2) {
				render.Error(w, twirp.NewErrorf(twirp.AlreadyExists, "order with trace already exists"))
				return
			}
		}

		fee := system.Fees[factory.Platform()]
		order.FeeAsset = fee.FeeAssetID
		order.FeeAmount = fee.FeeAmount.Mul(decimal.NewFromInt(int64(len(body.Tokens))))
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

		render.JSON(w, views.OrderView(*order))
	}
}
