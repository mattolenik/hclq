package cmd

import (
	"github.com/spf13/cobra"
	"fmt"
)

var AppendCmd = &cobra.Command{
	Use: "append <query> <value>",
	Short: "append something to matching value(s)",
	Args: cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(args)
	},
}
