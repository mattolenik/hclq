package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// StrReplaceCmd cobra command
var StrReplaceCmd = &cobra.Command{
	Use:     "str:replace <query> <str> <newStr>",
	Short:   "perform a string replace on query result",
	Args:    cobra.ExactArgs(3),
	Example: "str:replace <query> foo bar",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(args)
	},
}
