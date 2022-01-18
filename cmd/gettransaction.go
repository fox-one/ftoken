/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

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
	"context"
	"encoding/json"

	"github.com/spf13/cobra"
)

// gettransactionCmd represents the gettransaction command
var gettransactionCmd = &cobra.Command{
	Use: "gettransaction",
	Run: func(cmd *cobra.Command, args []string) {
		hash, err := cmd.Flags().GetString("hash")
		if err != nil {
			cmd.PrintErr("invalid hash", err)
			return
		}

		factory := provideQuorumFactory()
		tx, err := factory.ReadTransaction(context.Background(), hash)
		if err != nil {
			cmd.PrintErr("ReadTransaction failed:", err)
			return
		}

		bts, _ := json.MarshalIndent(tx, "", "  ")
		cmd.Println(string(bts))
	},
}

func init() {
	rootCmd.AddCommand(gettransactionCmd)

	gettransactionCmd.Flags().String("hash", "", "transaction hash")
}
