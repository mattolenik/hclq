package cmd

import (
	"fmt"
	"os"
	"github.com/spf13/cobra"
)

var RootCmd = &cobra.Command{
	Use:   "hclq [flags] [command]",
	Short: "Query and modify HashiCorp HCL files",
	Long: `hclq is a tool for querying the values of HCL files. Queries can return either single or multiple values,
which means that hclq commands work over ALL results of the query. This allows for the retrieval and modification of
multiple values at once.`,
}

var RootFlags = RootCmd.PersistentFlags()

func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	var format string
	var outFile string
	var inFile string
	var inPlace string
	RootFlags.StringVar(&format, "fmt", "json", "output format, json or go, default json")
	RootFlags.StringVar(&outFile, "out", "", "write output to this file, otherwise use stdout")
	RootFlags.StringVar(&inFile, "in", "", "read input from this file, otherwise use stdin")
	RootFlags.StringVar(&inPlace, "in-place", "", "read from this file and write changes back into it")
	RootCmd.AddCommand(GetCmd)
	RootCmd.AddCommand(SetCmd)
	RootCmd.AddCommand(StrReplaceCmd)
	RootCmd.AddCommand(PrependCmd)
	RootCmd.AddCommand(AppendCmd)
}
