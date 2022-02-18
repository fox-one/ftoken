/*
Copyright © 2020 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/shopspring/decimal"
	"github.com/spf13/cobra"
)

// ringsAssetCmd represents the load rings assets command
var ringsAssetCmd = &cobra.Command{
	Use:   "rings",
	Short: "load rings assets",
	Run: func(cmd *cobra.Command, args []string) {
		ctx := cmd.Context()
		database, err := provideDatabase()
		if err != nil {
			cmd.PrintErrf("provideDatabase failed: %v", err)
			return
		}
		defer database.Close()

		assets := provideAssetStore(database)
		client := provideMixinClient()
		assetz := provideAssetService(client)

		resp, err := http.Get("https://rings-api.pando.im/api/v1/markets/all")
		if err != nil {
			cmd.PrintErrf("fetch rings markets failed: %v", err)
			return
		}

		var body struct {
			Data []struct {
				AssetID        string          `json:"asset_id"`
				CtokenAssetID  string          `json:"ctoken_asset_id"`
				Price          decimal.Decimal `json:"price"`
				PriceUpdatedAt *time.Time      `json:"price_updated_at"`
			} `json:"data"`
		}

		if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
			cmd.PrintErrf("decode rings markets resp failed: %v", err)
			return
		}

		for _, item := range body.Data {
			{
				asset, err := assetz.Find(ctx, item.AssetID)
				if err != nil {
					cmd.PrintErr(err)
					return
				}

				asset.Price = item.Price
				asset.PriceUpdatedAt = item.PriceUpdatedAt
				if err := assets.Save(ctx, asset); err != nil {
					cmd.PrintErr(err)
					return
				}
				cmd.Println(asset.Symbol, "saved")
			}

			{
				asset, err := assetz.Find(ctx, item.CtokenAssetID)
				if err != nil {
					cmd.PrintErr(err)
					return
				}

				if err := assets.Save(ctx, asset); err != nil {
					cmd.PrintErr(err)
					return
				}
				cmd.Println(asset.Symbol, "saved")
			}
		}
	},
}

func init() {
	assetCmd.AddCommand(ringsAssetCmd)
}
