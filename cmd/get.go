package cmd

import (
	"fmt"

	"github.com/mattolenik/hclq/hcl"
	"github.com/spf13/cobra"
)

var useRawOutput bool

// GetCmd command
var GetCmd = &cobra.Command{
	Use:   "get <query>",
	Short: "retrieve values matching <query>",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		reader, err := getInputReader()
		if err != nil {
			failExit(err)
		}
		result, err := hcl.Get(reader, args[0])
		if err != nil {
			failExit(err)
		}
		output, err := getOutput(result, useRawOutput)
		if err != nil {
			failExit(err)
		}
		fmt.Println(output)
	},
}

func init() {
	GetCmd.PersistentFlags().BoolVarP(&useRawOutput, "raw", "r", false, "output raw format instead of JSON")
	RootCmd.AddCommand(GetCmd)
}
