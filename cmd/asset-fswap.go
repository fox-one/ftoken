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

	"github.com/fox-one/ftoken/core"
	"github.com/shopspring/decimal"
	"github.com/spf13/cobra"
)

// fswapAssetCmd represents the load fswap assets command
var fswapAssetCmd = &cobra.Command{
	Use:   "fswap",
	Short: "load fswap assets",
	Run: func(cmd *cobra.Command, args []string) {
		ctx := cmd.Context()
		database, err := provideDatabase()
		if err != nil {
			cmd.PrintErrf("provideDatabase failed: %v", err)
			return
		}
		defer database.Close()

		assets := provideAssetStore(database)

		resp, err := http.Get("https://api.4swap.org/api/assets")
		if err != nil {
			cmd.PrintErrf("fetch 4swap assets failed: %v", err)
			return
		}

		var body struct {
			Timestamp int64 `json:"ts"`
			Data      struct {
				Assets []struct {
					ID      string          `json:"id"`
					Name    string          `json:"name"`
					Symbol  string          `json:"symbol"`
					Logo    string          `json:"logo"`
					ChainID string          `json:"chain_id"`
					Price   decimal.Decimal `json:"price"`
				} `json:"assets"`
				Timestamp           int64           `json:"ts"`
				TransactionCount24H int             `json:"transaction_count_24h"`
				Volume24H           decimal.Decimal `json:"volume_24h"`
			} `json:"data"`
		}

		if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
			cmd.PrintErrf("decode 4swap assets resp failed: %v", err)
			return
		}

		for _, item := range body.Data.Assets {
			asset := core.Asset{
				AssetID: item.ID,
				Name:    item.Name,
				Symbol:  item.Symbol,
				Logo:    item.Logo,
				ChainID: item.ChainID,
			}
			if err := assets.Save(ctx, &asset); err != nil {
				cmd.PrintErr(err)
				return
			}

			cmd.Println(asset.Symbol, "saved")
		}
	},
}

func init() {
	assetCmd.AddCommand(fswapAssetCmd)
}
