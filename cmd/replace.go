package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// ReplaceCmd cobra command
var ReplaceCmd = &cobra.Command{
	Use:   "replace <query> <seq> <newSeq>",
	Short: "replace a substring with a new string, or a sublist with a new list",
	Args:  cobra.ExactArgs(3),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(args)
	},
}
