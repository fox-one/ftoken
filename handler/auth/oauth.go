package auth

import (
	"net/http"

	"github.com/fox-one/ftoken/core"
	"github.com/fox-one/ftoken/handler/render"
	"github.com/fox-one/mixin-sdk-go"
	"github.com/fox-one/pkg/httputil/param"
	"github.com/twitchtv/twirp"
)

func HandleOauth(system *core.System) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var body struct {
			Code string `json:"code,omitempty" valid:"minstringlength(6),required"`
		}

		if err := param.Binding(r, &body); err != nil {
			render.Error(w, err)
			return
		}

		ctx := r.Context()

		token, scope, err := mixin.AuthorizeToken(ctx, system.ClientID, system.ClientSecret, body.Code, "")
		if err != nil {
			render.Error(w, twirp.InvalidArgumentError("code", err.Error()))
			return
		}

		render.JSON(w, render.H{
			"token": token,
			"scope": scope,
		})
	}
}
