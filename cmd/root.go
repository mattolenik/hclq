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
	Short:        "Query and modify HashiCorp Configuration Language files. Like sed for HCL.",
	Long: `hclq is a tool for manipulating the config files used by HashiCorp tools.

hclq uses a "breadcrumb" or "path" style query. Given the HCL:
    data "foo" "bar" {
        id = "100"
        other = [1, 2, 3]
    }

A query for 'data.foo.bar.id' would return 100. Arrays/lists must be matched
with the [] suffix, e.g. 'data.foo.bar.other[]' or 'data.foo.bar.other[1]'.

Match types:
    literal     Match a literal value.
	list[]      Match a list and retrieve all items.
	list[1]     Match a list and retrieve a specific item.
	/regex/     Match anything according to the specified regex.
	/regex/[]   Match a list according to the regex and retrieve all items.
	/regex/[1]  Match a list according to the regex and retrieve a specific item.
	*           Match anything.

Queries can return either single or multiple values. If a query matches e.g.
multiple arrays across multiple objects, a list of arrays will be returned.
If this query is used with a set command, ALL of those matching arrays will be
set.
`,
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
