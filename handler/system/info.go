package system

import (
	"net/http"

	"github.com/fox-one/ftoken/core"
	"github.com/fox-one/ftoken/handler/render"
)

func HandleInfo(system core.System, factories []core.Factory) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var _factories = make([]interface{}, len(factories))
		for i, factory := range factories {
			_factories[i] = render.H{
				"platform": factory.Platform(),
				"asset_id": factory.GasAsset(),
			}
		}

		render.JSON(w, render.H{
			"version":   system.Version,
			"client_id": system.ClientID,
			"factories": _factories,
		})
	}
}
