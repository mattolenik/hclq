package cmd

import (
	"fmt"
	"os"

	"github.com/hashicorp/hcl/hcl/ast"
	"github.com/hashicorp/hcl/hcl/printer"
	jsonParser "github.com/hashicorp/hcl/json/parser"
	"github.com/mattolenik/hclq/query"
	"github.com/spf13/cobra"
)

// SetCmd cobra command
var SetCmd = &cobra.Command{
	Use:   "set <query> <newValue>",
	Short: "set matching value(s), specify a string, number, or JSON object or array",
	Args:  cobra.ExactArgs(2),

	RunE: func(cmd *cobra.Command, args []string) error {
		queryNodes, err := query.Parse(args[0])
		if err != nil {
			return err
		}
		reader := os.Stdin
		if val := cmd.Flag("in").Value.String(); val != "" {
			var err error
			reader, err = os.Open(val)
			if err != nil {
				return err
			}
		}
		newValueArg := args[1]
		newValueJSON := fmt.Sprintf(`{"root":%s}`, newValueArg) // parser requires `root` key

		newValue, err := jsonParser.Parse([]byte(newValueJSON))
		if err != nil {
			return err
		}
		resultPairs, isList, docRoot, err := query.HCL(reader, queryNodes)
		if isList {
			for _, pair := range resultPairs {
				list, ok := pair.Node.(*ast.ListType)
				if !ok {
					return fmt.Errorf("Expected ListType as query result")
				}
				list.List = newValue.Node.(*ast.ObjectList).Items[0].Val.(*ast.ListType).List
			}
		} else {
			for _, pair := range resultPairs {
				item, ok := pair.Node.(*ast.LiteralType)
				if !ok {
					return fmt.Errorf("Expected LiteralType in query results")
				}
				item.Token.Text = newValueArg
			}
		}
		return printer.Fprint(os.Stdout, docRoot)
	},
}

var appendCmd = &cobra.Command{
	Use:  "append",
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("append cmd")
		return nil
	},
}

func init() {
	var append bool
	var prepend bool
	var modify bool
	SetCmd.PersistentFlags().BoolVar(&append, "append", false, "append new value(s) to existing value")
	SetCmd.PersistentFlags().BoolVar(&prepend, "prepend", false, "prepend new value(s) to existing value")
	SetCmd.PersistentFlags().BoolVarP(&modify, "modify", "m", false, "modify the input file rather than printing output, conflicts with --out")
	RootCmd.AddCommand(SetCmd)
	SetCmd.AddCommand(AppendCmd)
}
