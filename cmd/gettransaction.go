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
