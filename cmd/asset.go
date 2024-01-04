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
	"github.com/spf13/cobra"
)

// assetCmd represents the asset command
var assetCmd = &cobra.Command{
	Use:   "asset",
	Short: "manager asset store",
	Run: func(cmd *cobra.Command, args []string) {
		ctx := cmd.Context()
		database, err := provideDatabase()
		if err != nil {
			cmd.PrintErrf("provideDatabase failed: %v", err)
			return
		}
		defer database.Close()

		client := provideMixinClient()
		assetz := provideAssetService(client)
		assets := provideAssetStore(database)

		assetID, ok := getArg(args, 0)
		if !ok {
			cmd.PrintErr("args[0]: asset id is empty")
			return
		}

		asset, err := assetz.Find(ctx, assetID)
		if err != nil {
			cmd.PrintErr(err)
			return
		}

		if err := assets.Save(ctx, asset); err != nil {
			cmd.PrintErr(err)
			return
		}

		cmd.Println(asset.Symbol, "saved")
	},
}

func init() {
	rootCmd.AddCommand(assetCmd)
}
