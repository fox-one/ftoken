package system

import (
	"net/http"

	"github.com/fox-one/ftoken/handler/render"
	"github.com/fox-one/pkg/property"
)

func HandleProperty(properties property.Store) http.HandlerFunc {
	keys := map[string]string{
		"bwatch_outputs_checkpoint": "outputs",
		"bwatch_sync_checkpoint":    "sync_utxo",
	}

	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		values, err := properties.List(ctx)
		if err != nil {
			render.Error(w, err)
			return
		}

		view := render.H{}
		for k, v := range values {
			if key, ok := keys[k]; ok {
				view[key] = v.String()
			}
		}

		render.JSON(w, view)
	}
}
