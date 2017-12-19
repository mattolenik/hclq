package cmd

import (
	"github.com/spf13/cobra"
	"fmt"
)

var SetCmd = &cobra.Command{
	Use: "set <query> <value>",
	Short: "set matching value(s)",
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(args)
	},
}
