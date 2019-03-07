package cmd

import (
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/hashicorp/hcl/hcl/ast"
	"github.com/hashicorp/hcl/hcl/token"
	"github.com/mattolenik/hclq/config"
	"github.com/mattolenik/hclq/hclq"
)

func Set(reader io.Reader, writer io.Writer, query, newValue string) error {
	return performSet(reader, writer, query,
		func(list *ast.ListType) error {
			listNode, err := hclq.HclListFromJSON(newValue)
			if err != nil {
				return err
			}
			list.List = listNode.List
			return nil
		}, func(tok *token.Token) error {
			tok.Text = `"` + newValue + `"`
			tok.Type = getTokenType(newValue)
			return nil
		})
}

func Append(reader io.Reader, writer io.Writer, query, newValue string) error {
	return performSet(reader, writer, query,
		func(list *ast.ListType) error {
			node, err := hclq.HclFromJSON(newValue)
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
}

//Short: "prepend value(s) to a list or a string",
func Prepend(reader io.Reader, writer io.Writer, query, newValue string) error {
	return performSet(reader, writer, query,
		func(list *ast.ListType) error {
			node, err := hclq.HclFromJSON(newValue)
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
}

//Short: "find and replace a subsequence of items (or chars for strings)",
//ReplaceCmd.Flags().IntVarP(&config.ReplaceNTimes, "replace-n", "n", -1, "Limit replacements to n occurrences")
func Replace(reader io.Reader, writer io.Writer, query, oldValue, newValue string, n int) error {
	return performSet(reader, writer, query,
		func(list *ast.ListType) error {
			// TODO: implement replace on lists
			panic("replace on lists is not implemented yet")
		}, func(tok *token.Token) error {
			tok.Text = `"` + strings.Replace(trimToken(tok.Text), oldValue, newValue, n) + `"`
			tok.Type = token.STRING
			return nil
		})
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

func init() {
}

func performSet(reader io.Reader, writer io.Writer, queryString string, listAction func(*ast.ListType) error, valueAction func(*token.Token) error) error {
	err := validateOutputConfig()
	if err != nil {
		return err
	}
	doc, err := hclq.FromReader(reader)
	if err != nil {
		return err
	}
	err = doc.Set(queryString, listAction, valueAction)
	if err != nil {
		return err
	}
	doc.Print(writer)
	return nil
}
func validateOutputConfig() error {
	if config.ModifyInPlace && len(config.OutputFile) > 0 {
		return fmt.Errorf("cannot use both --modify and --out at the same time")
	}
	if config.ModifyInPlace && len(config.InputFile) == 0 {
		return fmt.Errorf("cannot use --modify without specifying an input file with --in")
	}
	return nil
}
