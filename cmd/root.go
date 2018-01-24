package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// RootCmd command
var RootCmd = &cobra.Command{
	Use:   "hclq [flags] [command]",
	Short: "Query and modify HashiCorp HCL files",
	Long: `hclq is a tool for querying the values of HCL files, reminiscent of jq.

Queries can return either single or multiple values, which means that hclq commands work over ALL results of a query.
This means that commands such as set can work over many keys at once.

hclq outputs JSON by default. A tool such as jq is recommended for further processing.`,
}

// RootFlags flags
var RootFlags = RootCmd.PersistentFlags()

// Execute - cobra entry point
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	var outFile string
	var inFile string
	var modify bool
	var raw bool
	RootFlags.StringVarP(&outFile, "out", "o", "", "write output to this file, otherwise use stdout")
	RootFlags.StringVarP(&inFile, "in", "i", "", "read input from this file, otherwise use stdin")
	RootFlags.BoolVarP(&modify, "modify", "m", false, "modify the input file rather than printing output, conflicts with --out")
	RootFlags.BoolVarP(&raw, "raw", "r", false, "output raw format instead of JSON")
	RootCmd.AddCommand(GetCmd)
	RootCmd.AddCommand(SetCmd)
	RootCmd.AddCommand(StrReplaceCmd)
	RootCmd.AddCommand(PrependCmd)
	RootCmd.AddCommand(AppendCmd)
}
