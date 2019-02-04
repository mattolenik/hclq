package cmd

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/hashicorp/hcl/hcl/ast"
	"github.com/hashicorp/hcl/hcl/printer"
	"github.com/hashicorp/hcl/hcl/token"
	jsonParser "github.com/hashicorp/hcl/json/parser"
	"github.com/mattolenik/hclq/config"
	"github.com/mattolenik/hclq/hclq"
	"github.com/mattolenik/hclq/query"
	"github.com/spf13/cobra"
)

type listAction = func(list *ast.ListType, newNodes []ast.Node)
type valueAction = func(token *token.Token, newValue string)

// SetCmd cobra command
var SetCmd = &cobra.Command{
	Use:   "set <query> <newValue>",
	Short: "set matching value(s), specify a string, number, or JSON object or array",
	Args:  cobra.ExactArgs(2),

	RunE: func(cmd *cobra.Command, args []string) error {
		newValue := args[1]
		return setImpl(args[0],
			func(list *ast.ListType) error {
				hcl, err := getValueFromJSON(newValue)
				if err != nil {
					return err
				}
				list.List = hcl.(*ast.ListType).List
				return nil
			}, func(tok *token.Token) error {
				tok.Text = `"` + newValue + `"`
				tok.Type = getTokenType(newValue)
				return nil
			})
	},
}

// AppendCmd cobra command
var AppendCmd = &cobra.Command{
	Use:   "append <query> <newValue>",
	Short: "append value(s) to a list or a string",
	Args:  cobra.ExactArgs(2),

	RunE: func(cmd *cobra.Command, args []string) error {
		newValue := args[1]
		return setImpl(args[0],
			func(list *ast.ListType) error {
				node, err := getValueFromJSON(newValue)
				if err != nil {
					return err
				}
				list.List = append(list.List, node)
				return nil
			}, func(tok *token.Token) error {
				tok.Text = `"` + trimToken(tok.Text) + newValue + `"`
				tok.Type = token.STRING
				return nil
			})
	},
}

// PrependCmd cobra command
var PrependCmd = &cobra.Command{
	Use:   "prepend <query> <newValue>",
	Short: "prepend value(s) to a list or a string",
	Args:  cobra.ExactArgs(2),

	RunE: func(cmd *cobra.Command, args []string) error {
		newValue := args[1]
		return setImpl(args[0],
			func(list *ast.ListType) error {
				node, err := getValueFromJSON(newValue)
				if err != nil {
					return err
				}
				list.List = append(node.(*ast.ListType).List, list.List...)
				return nil
			}, func(tok *token.Token) error {
				tok.Text = `"` + newValue + trimToken(tok.Text) + `"`
				tok.Type = token.STRING
				return nil
			})
	},
}

// ReplaceCmd cobra command
var ReplaceCmd = &cobra.Command{
	Use:   "replace <query> <oldSequence> <newSequence>",
	Short: "find and replace a subsequence of items (or chars for strings)",
	Args:  cobra.ExactArgs(3),
	RunE: func(cmd *cobra.Command, args []string) error {
		return setImpl(args[0],
			func(list *ast.ListType) error {
				panic("Not implemented")
			}, func(tok *token.Token) error {
				tok.Text = `"` + strings.Replace(trimToken(tok.Text), args[1], args[2], config.ReplaceNTimes) + `"`
				tok.Type = token.STRING
				return nil
			})
	},
}

func trimToken(tok string) string {
	return strings.Trim(tok, `"`)
}

func getTokenType(val string) token.Type {
	_, err := strconv.ParseInt(val, 0, 64)
	if err == nil {
		return token.NUMBER
	}
	_, err = strconv.ParseFloat(val, 64)
	if err == nil {
		return token.FLOAT
	}
	_, err = strconv.ParseBool(val)
	if err == nil {
		return token.BOOL
	}
	return token.STRING
	// TODO: support HEREDOC
}

func getValueFromJSON(json string) (ast.Node, error) {
	json = fmt.Sprintf(`{"root":%s}`, json) // parser requires `root` key
	hcl, err := jsonParser.Parse([]byte(json))
	if err != nil {
		return nil, err
	}
	node := hcl.Node.(*ast.ObjectList).Items[0].Val
	return node, nil
}

func setImpl(
	queryString string,
	listAction func(*ast.ListType) error,
	valueAction func(*token.Token) error) error {

	queryNodes, err := query.Parse(queryString)
	if err != nil {
		return err
	}
	reader, err := getInputReader()
	if err != nil {
		return err
	}

	doc := hclq.FromReader(reader)
	resultPairs, err := doc.Query(queryNodes)
	if err != nil {
		return err
	}

	for _, pair := range resultPairs {
		list, ok := pair.Node.(*ast.ListType)
		if ok {
			listAction(list)
			continue
		}
		literal, ok := pair.Node.(*ast.LiteralType)
		if ok {
			valueAction(&literal.Token)
			continue
		}
	}

	writer, err := getOutputWriter()
	if err != nil {
		return err
	}
	return printer.Fprint(writer, doc.FileNode)
}

func init() {
	SetCmd.PersistentFlags().BoolVarP(&config.ModifyInPlace, "modify", "m", false, "modify the input file rather than printing output, conflicts with --out")
	RootCmd.AddCommand(SetCmd)
	SetCmd.AddCommand(AppendCmd)
	SetCmd.AddCommand(PrependCmd)
	SetCmd.AddCommand(ReplaceCmd)
	ReplaceCmd.Flags().IntVarP(&config.ReplaceNTimes, "replace-n", "n", -1, "Limit replacements to n occurrences")
}
