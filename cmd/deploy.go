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
