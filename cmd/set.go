package cmd

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/hashicorp/hcl/hcl/ast"
	"github.com/hashicorp/hcl/hcl/printer"
	"github.com/hashicorp/hcl/hcl/token"
	jsonParser "github.com/hashicorp/hcl/json/parser"
	"github.com/mattolenik/hclq/query"
	"github.com/spf13/cobra"
)

var modify bool
var replaceN int

type listAction = func(list *ast.ListType, newNodes []ast.Node)
type valueAction = func(token *token.Token, newValue string)

// SetCmd cobra command
var SetCmd = &cobra.Command{
	Use:   "set <query> <newValue>",
	Short: "set matching value(s), specify a string, number, or JSON object or array",
	Args:  cobra.ExactArgs(2),

	RunE: func(cmd *cobra.Command, args []string) error {
		newValue := args[1]
		return setImpl(cmd, args[0],
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
		return setImpl(cmd, args[0],
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
		return setImpl(cmd, args[0],
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
		return setImpl(cmd, args[0],
			func(list *ast.ListType) error {
				panic("Not implemented")
			}, func(tok *token.Token) error {
				tok.Text = `"` + strings.Replace(trimToken(tok.Text), args[1], args[2], replaceN) + `"`
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
	cmd *cobra.Command,
	queryString string,
	listAction func(*ast.ListType) error,
	valueAction func(*token.Token) error) error {

	queryNodes, err := query.Parse(queryString)
	if err != nil {
		return err
	}
	reader := os.Stdin
	if inFile != "" {
		reader, err = os.Open(inFile)
		if err != nil {
			return err
		}
	}
	resultPairs, isList, docRoot, err := query.HCL(reader, queryNodes)
	if isList {
		list, ok := resultPairs.Node.(*ast.ListType)
		if !ok {
			return fmt.Errorf("Expected ListType as query result")
		}
		listAction(list)
	} else {
		item, ok := resultPairs.Node.(*ast.LiteralType)
		if !ok {
			return fmt.Errorf("Expected LiteralType in query results")
		}
		valueAction(&item.Token)
	}

	writer := os.Stdout
	if outFile != "" {
		writer, err = os.Create(outFile)
	}
	return printer.Fprint(writer, docRoot)
}

func init() {
	SetCmd.PersistentFlags().BoolVarP(&modify, "modify", "m", false, "modify the input file rather than printing output, conflicts with --out")
	RootCmd.AddCommand(SetCmd)
	SetCmd.AddCommand(AppendCmd)
	SetCmd.AddCommand(PrependCmd)
	SetCmd.AddCommand(ReplaceCmd)
	ReplaceCmd.Flags().IntVarP(&replaceN, "replace-n", "n", -1, "Limit replacements to n occurrences")
}
