package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// AppendCmd cobra command
var AppendCmd = &cobra.Command{
	Use:   "append <query> <value>",
	Short: "append something to matching value(s)",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(args)
	},
}
