package system

import (
	"net/http"
	"time"

	"github.com/fox-one/ftoken/handler/render"
)

func HandleTime() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		t := time.Now()
		render.JSON(w, render.H{
			"iso":   t.Format(time.RFC3339),
			"epoch": t.Unix(),
		})
	}
}
