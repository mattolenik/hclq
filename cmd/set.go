package cmd

import (
	"github.com/spf13/cobra"
	"io"
	"github.com/mattolenik/hclq/query"
	"io/ioutil"
	"github.com/mattolenik/hclq/utils"
	"github.com/hashicorp/hcl/hcl/parser"
	"github.com/hashicorp/hcl/hcl/ast"
	"os"
	"github.com/hashicorp/hcl/hcl/printer"
)

var SetCmd = &cobra.Command{
	Use: "set <query> <value>",
	Short: "set matching value(s)",
	Args: cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		queryNodes, _ := query.Parse(args[0])
		reader := os.Stdin
		if val := cmd.Flag("in").Value.String(); val != "" {
			var err error
			reader, err = os.Open(val)
			if err != nil {
				panic(err)
			}
		}
		writer := os.Stdout
		if val := cmd.Flag("out").Value.String(); val != "" {
			var err error
			writer, err = os.Open(val)
			if err != nil {
				panic(err)
			}
		}
		set(reader, writer, queryNodes, args[1])
	},
}

func set(reader io.Reader, writer io.Writer, qry []query.Node, value interface{}) error {
	bytes, err := ioutil.ReadAll(reader)
	if err != nil {
		return err
	}
	node, err := parser.Parse(bytes)
	if err != nil {
		return err
	}
	err = utils.Walk(node.Node, qry, 0, func(n ast.Node, queryNode query.Node) (err error) {
		if lit, ok := n.(*ast.LiteralType); ok {
			lit.Token.Text = value.(string)
		}
		//if lst, ok := n.(*ast.ListType); ok {
		//}
		return nil
	})
	if err != nil {
		return err
	}
	printer.Fprint(writer, node)
	return nil
}
