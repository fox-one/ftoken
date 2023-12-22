package cmd

import (
	"fmt"
	"net/http"

	"github.com/fox-one/ftoken/handler/hc"
	"github.com/fox-one/ftoken/worker"
	"github.com/fox-one/ftoken/worker/cashier"
	"github.com/fox-one/ftoken/worker/order"
	"github.com/fox-one/ftoken/worker/payee"
	"github.com/fox-one/pkg/logger"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/rs/cors"
	"github.com/spf13/cobra"
	"golang.org/x/sync/errgroup"
)

// workerCmd represents the worker command
var workerCmd = &cobra.Command{
	Use:   "worker",
	Short: "run dirtoracle worker",
	Run: func(cmd *cobra.Command, args []string) {
		ctx := cmd.Context()
		cfg.DB.ReadHost = ""
		database, err := provideDatabase()
		if err != nil {
			cmd.PrintErrf("provideDatabase failed: %v", err)
			return
		}

		defer database.Close()
		client := provideMixinClient()

		assets := provideAssetStore(database)
		wallets := provideWalletStore(database)
		walletz := provideWalletService(client)
		orders := provideOrderStore(database)
		transactions := provideTransactionStore(database)
		properties := providePropertyStore(database)
		factories := provideAllFactories()
		system, err := provideSystem(ctx, client, factories)
		if err != nil {
			cmd.PrintErrf("provideSystem failed: %v", err)
			return
		}

		workers := []worker.Worker{
			cashier.New(wallets, walletz),
			payee.New(
				payee.Config{ClientID: cfg.Dapp.ClientID},
				system,
				properties,
				assets,
				orders,
				transactions,
				wallets,
				walletz,
				factories,
			),
			order.New(
				system,
				orders,
				transactions,
				wallets,
				factories,
			),
		}

		cmd.Printf("ftoken worker with version %q launched!\n", rootCmd.Version)

		// worker api
		{
			mux := chi.NewMux()
			mux.Use(middleware.Recoverer)
			mux.Use(middleware.StripSlashes)
			mux.Use(cors.AllowAll().Handler)
			mux.Use(logger.WithRequestID)
			mux.Use(middleware.Logger)

			// hc
			{
				mux.Mount("/hc", hc.Handle(rootCmd.Version))
			}

			// launch server
			port, _ := cmd.Flags().GetInt("port")
			addr := fmt.Sprintf(":%d", port)

			go http.ListenAndServe(addr, mux)
		}

		g, ctx := errgroup.WithContext(ctx)
		for idx := range workers {
			w := workers[idx]
			g.Go(func() error {
				return w.Run(ctx)
			})
		}

		if err := g.Wait(); err != nil {
			cmd.PrintErrln("run worker", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(workerCmd)

	workerCmd.Flags().Int("port", 9301, "worker api port")
}
