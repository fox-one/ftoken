package cmd

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/drone/signal"
	"github.com/fox-one/ftoken/handler"
	"github.com/fox-one/ftoken/handler/hc"
	"github.com/fox-one/ftoken/handler/ip"
	"github.com/fox-one/ftoken/handler/render"
	"github.com/fox-one/ftoken/handler/system"
	"github.com/fox-one/pkg/logger"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/rs/cors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// serverCmd represents the server command
var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "run bwatch server",
	Run: func(cmd *cobra.Command, args []string) {
		ctx := cmd.Context()
		database, err := provideDatabase()
		if err != nil {
			cmd.PrintErrf("provideDatabase failed: %v", err)
			return
		}
		defer database.Close()

		client := provideMixinClient()

		mux := chi.NewMux()
		mux.Use(middleware.Recoverer)
		mux.Use(middleware.StripSlashes)
		mux.Use(cors.AllowAll().Handler)
		mux.Use(logger.WithRequestID)
		mux.Use(middleware.Logger)
		mux.Use(middleware.NewCompressor(5).Handler)

		// debug
		if debugMode {
			mux.Mount("/debug", middleware.Profiler())
			render.ResponseErrorMessageAsHint = true
		}

		// hc
		{
			mux.Mount("/hc", hc.Handle(rootCmd.Version))
			mux.Get("/ip", ip.Handle())
			mux.Get("/time", system.HandleTime())
		}

		{
			factories := provideAllFactories()
			system, err := provideSystem(ctx, client, factories)
			if err != nil {
				cmd.PrintErrf("provideSystem failed: %v", err)
				return
			}

			svr := handler.New(
				system,
				provideAssetStore(database),
				provideOrderStore(database),
				provideTransactionStore(database),
				provideWalletService(client),
				factories,
			)
			mux.Mount("/api", svr.Handle())
		}

		// launch server
		port, _ := cmd.Flags().GetInt("port")
		addr := fmt.Sprintf(":%d", port)

		svr := &http.Server{
			Addr:    addr,
			Handler: mux,
		}

		done := make(chan struct{}, 1)
		ctx = signal.WithContextFunc(ctx, func() {
			logrus.Debug("shutdown server...")

			// create context with timeout
			ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
			defer cancel()

			if err := svr.Shutdown(ctx); err != nil {
				logrus.WithError(err).Error("graceful shutdown server failed")
			}

			close(done)
		})

		logrus.Infoln("serve at", addr)
		if err := svr.ListenAndServe(); err != http.ErrServerClosed {
			logrus.WithError(err).Fatal("server aborted")
		}

		<-done
	},
}

func init() {
	rootCmd.AddCommand(serverCmd)
	serverCmd.Flags().IntP("port", "p", 9302, "server port")
}
