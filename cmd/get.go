package cmd

import (
	"fmt"

	"github.com/mattolenik/hclq/config"
	"github.com/mattolenik/hclq/hclq"
	"github.com/spf13/cobra"
)

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
		doc := hclq.FromReader(reader)
		result, err := doc.Get(args[0])
		if err != nil {
			failExit(err)
		}
		output, err := getOutput(result, config.UseRawOutput)
		if err != nil {
			failExit(err)
		}
		fmt.Println(output)
	},
}

func init() {
	GetCmd.PersistentFlags().BoolVarP(&config.UseRawOutput, "raw", "r", false, "output raw format instead of JSON")
	RootCmd.AddCommand(GetCmd)
}
