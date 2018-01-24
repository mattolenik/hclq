package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// PrependCmd cobra command
var PrependCmd = &cobra.Command{
	Use:   "prepend <query> <value>",
	Short: "prepend something to query result",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(args)
	},
}
