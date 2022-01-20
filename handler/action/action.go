package action

import (
	"encoding/base64"
	"net/http"

	"github.com/fox-one/ftoken/core"
	"github.com/fox-one/ftoken/handler/render"
	"github.com/fox-one/ftoken/pkg/mtg"
	"github.com/fox-one/mixin-sdk-go"
	"github.com/fox-one/pkg/httputil/param"
	"github.com/fox-one/pkg/uuid"
	"github.com/lib/pq"
	"github.com/twitchtv/twirp"
)

func HandleCreateAction(system core.System, walletz core.WalletService, factories []core.Factory) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		var body struct {
			TraceID  string       `json:"trace_id,omitempty"`
			Platform string       `json:"platform,omitempty"`
			Tokens   core.Tokens  `json:"tokens,omitempty"`
			Receiver core.Address `json:"receiver,omitempty"`
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

		memoBts, err := mtg.Encode(body.Tokens, body.Receiver)
		if err != nil {
			render.Error(w, twirp.InternalErrorWith(err))
			return
		}
		memo := base64.StdEncoding.EncodeToString(memoBts)

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

		if body.Receiver.Destination == "" {
			body.Receiver = *system.Addresses[factory.GasAsset()]
		}

		tx, err := factory.CreateTransaction(ctx, body.Tokens, &body.Receiver)
		if err != nil {
			render.Error(w, twirp.InternalErrorWith(err))
			return
		}
		gas := tx.Gas.Mul(system.Gas.Multiplier)
		if min, ok := system.Gas.Mins[factory.Platform()]; ok && gas.LessThan(min) {
			gas = min
		}

		transfer := &core.Transfer{
			TraceID: body.TraceID,
			AssetID: factory.GasAsset(),
			Amount:  gas,
			Memo:    memo,

			Opponents: pq.StringArray{system.ClientID},
		}

		code, err := walletz.ReqTransfer(ctx, transfer)
		if err != nil {
			render.Error(w, twirp.InternalErrorWith(err))
			return
		}

		render.JSON(w, render.H{
			"asset_id": transfer.AssetID,
			"amount":   transfer.Amount,
			"memo":     memo,
			"code":     code,
			"code_url": mixin.URL.Codes(code),
		})
	}
}
