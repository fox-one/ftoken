package handler

import (
	"net/http"

	"github.com/fox-one/ftoken/core"
	"github.com/fox-one/ftoken/handler/action"
	"github.com/fox-one/ftoken/handler/auth"
	"github.com/fox-one/ftoken/handler/ip"
	"github.com/fox-one/ftoken/handler/order"
	"github.com/fox-one/ftoken/handler/render"
	"github.com/go-chi/chi"
	"github.com/twitchtv/twirp"
)

type (
	Server struct {
		system    core.System
		orders    core.OrderStore
		txStore   core.TransactionStore
		walletz   core.WalletService
		factories []core.Factory
	}
)

func New(
	system core.System,
	orders core.OrderStore,
	txStore core.TransactionStore,
	walletz core.WalletService,
	factories []core.Factory,
) Server {
	return Server{
		system:    system,
		orders:    orders,
		txStore:   txStore,
		walletz:   walletz,
		factories: factories,
	}
}

func (s Server) Handle() http.Handler {
	r := chi.NewRouter()
	r.Use(ip.WithClientIP)

	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		render.Error(w, twirp.NotFoundError("not found"))
	})

	r.Post("/oauth", auth.HandleOauth(&s.system))

	r.Post("/actions", action.HandleCreateAction(s.system, s.walletz, s.factories))

	r.Route("/orders", func(r chi.Router) {
		r.Post("/", order.HandleCreateOrder(s.system, s.walletz, s.orders, s.factories))
		r.Get("/{trace_id}", order.HandleFetchOrder(s.orders))
	})

	return r
}
