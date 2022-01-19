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

	"github.com/fox-one/ftoken/core"
	"github.com/spf13/cobra"
)

// deployCmd represents the deploy tokens command
var deployCmd = &cobra.Command{
	Use: "deploy",
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()
		exec, _ := cmd.Flags().GetBool("e")
		receiver, _ := cmd.Flags().GetString("receiver")
		tokenStr, _ := cmd.Flags().GetString("tokens")

		var tokens core.Tokens
		if err := json.Unmarshal([]byte(tokenStr), &tokens); err != nil {
			cmd.PrintErr("unmarshal tokens failed: ", tokens)
			return
		}

		factory := provideQuorumFactory()
		tx, err := factory.CreateTransaction(ctx, tokens, &core.Address{Destination: receiver})
		if err != nil {
			cmd.PrintErr("CreateTransaction failed:", err)
			return
		}

		if exec {
			if err := factory.SendTransaction(ctx, tx); err != nil {
				cmd.PrintErr("SendTransaction failed:", err)
				return
			}
		}

		bts, _ := json.MarshalIndent(tx, "", "  ")
		cmd.Println(string(bts))
	},
}

func init() {
	rootCmd.AddCommand(deployCmd)

	deployCmd.Flags().String("tokens", "[]", "tokens in json")
	deployCmd.Flags().String("receiver", "", "receiver address to receive tokens")
	deployCmd.Flags().Bool("e", false, "execute transaction directly")
}
