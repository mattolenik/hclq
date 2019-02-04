package cmd

import (
	"os"

	"github.com/mattolenik/hclq/config"
	"github.com/spf13/cobra"
)

// Set by ldflags -X during build
var version string

// RootCmd command
var RootCmd = &cobra.Command{
	Use:     "hclq [flags] [command]",
	Version: version,
	// Don't print usage on error but still do so with --help and no args.
	SilenceUsage: true,
	Short:        "Query and modify HashiCorp HCL files",
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
		os.Exit(1)
	}
}

func init() {
	RootCmd.SetVersionTemplate("{{.Version}}\n")
	RootFlags.StringVarP(&config.OutputFile, "out", "o", "", "write output to this file, otherwise use stdout")
	RootFlags.StringVarP(&config.InputFile, "in", "i", "", "read input from this file, otherwise use stdin")
}
