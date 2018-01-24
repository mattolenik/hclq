package cmd

import (
	"fmt"
	"os"
	"reflect"

	"github.com/hashicorp/hcl/hcl/ast"
	"github.com/hashicorp/hcl/json/parser"
	"github.com/mattolenik/hclq/query"
	"github.com/spf13/cobra"
)

// SetCmd cobra command
var SetCmd = &cobra.Command{
	Use:   "set <query> <valueAsJSON>",
	Short: "set matching value(s), the new value should be valid JSON",
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
		//raw := cmd.Flag("raw").Value.String() == "true"

		resultPairs, isList, err := query.HCL(reader, queryNodes)
		if isList {
			newValue := fmt.Sprintf(`{"root":%s}`, args[1]) // parser requires `root` key
			_, err := parser.Parse([]byte(newValue))
			if err != nil {
				return err
			}
			if len(resultPairs) != 1 {
				return fmt.Errorf("Expected exactly 1 result when retrieving a list")
			}
			fmt.Println(reflect.TypeOf(resultPairs[0].Node))
			switch node := resultPairs[0].Node.(type) {
			case *ast.ListType:
				for _, item := range node.List
					switch n := item.(type) {
					case *ast.LiteralType:
						n.
					}
				}
			}
		}
		results := []string{} // Requires empty slice declaration, not nil declaration
		for _, pair := range resultPairs {
			results = append(results, pair.Serialized)
		}
		return nil
	},
}
