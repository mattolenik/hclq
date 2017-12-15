package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var RootCmd = &cobra.Command{
	Use:   "hclq",
	Short: "Query and modify HashiCorp HCL files",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {

	},
	Args: cobra.ExactArgs(1),
}

var getCmd = &cobra.Command{
	Use:   "get [file]",
	Short: "retrieve matching values, from file or stdin",
	Run: func(cmd *cobra.Command, args []string) {
	},
	Args: cobra.ExactArgs(1),
}

var setCmd = &cobra.Command{
	Use: "set <value> [--out-file | -o file] [--in-place | -i]",
	Short: "set matching values",
}

var appendCmd = &cobra.Command{
	Use: "append <value> [--out-file | -o out-file] [--in-place | -i]",
	Short: "append something to matching values",
}

var prependCmd = &cobra.Command{
	Use: "prepend <value> [--out-file | -o out-file] [--in-place | -i]",
	Short: "prepend something to matching values",
}

var replaceCmd = &cobra.Command{
	Use: "replace <str> <newStr> [--out-file | -o file] [--in-place | -i]",
	Short: "perform a string replace on matching values",
}

func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	var format string
	var outFile string
	var inPlace bool
	RootCmd.PersistentFlags().StringVarP(&format, "format", "f", "json", "output format, `json` or `go`, default json")
	RootCmd.AddCommand(getCmd)
	setCmd.Flags().StringVarP(&outFile, "out-file", "o", "", "write output to a file")
	setCmd.Flags().BoolVarP(&inPlace, "in-place", "-i", false, "write output back to the input file, modifying it in-place")
}
