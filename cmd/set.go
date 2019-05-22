package cmd

import (
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/hashicorp/hcl/hcl/ast"
	"github.com/hashicorp/hcl/hcl/token"
	"github.com/mattolenik/hclq/config"
	"github.com/mattolenik/hclq/hclq"
	"github.com/spf13/cobra"
)

// SetCmd cobra command
var SetCmd = &cobra.Command{
	Use:   "set <query> <newValue>",
	Short: "set matching value(s), specify a string, number, or JSON object or array",
	Args:  cobra.ExactArgs(2),

	RunE: func(cmd *cobra.Command, args []string) error {
		queryString := args[0]
		newValue := args[1]
		return performSet2(queryString, newValue)
	},
}

// AppendCmd cobra command
var AppendCmd = &cobra.Command{
	Use:   "append <query> <newValue>",
	Short: "append value(s) to a list or a string",
	Args:  cobra.ExactArgs(2),

	RunE: func(cmd *cobra.Command, args []string) error {
		newValue := args[1]
		return performSet(args[0],
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
	},
}

// PrependCmd cobra command
var PrependCmd = &cobra.Command{
	Use:   "prepend <query> <newValue>",
	Short: "prepend value(s) to a list or a string",
	Args:  cobra.ExactArgs(2),

	RunE: func(cmd *cobra.Command, args []string) error {
		newValue := args[1]
		return performSet(args[0],
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
	},
}

// ReplaceCmd cobra command
var ReplaceCmd = &cobra.Command{
	Use:   "replace <query> <oldSequence> <newSequence>",
	Short: "find and replace a subsequence of items (or chars for strings)",
	Args:  cobra.ExactArgs(3),
	RunE: func(cmd *cobra.Command, args []string) error {
		return performSet(args[0],
			func(list *ast.ListType) error {
				// TODO: implement replace on lists
				panic("replace on lists is not implemented yet")
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

func init() {
	SetCmd.PersistentFlags().BoolVarP(&config.ModifyInPlace, "in-place", "", false, "edit the input file in-place rather than printing to stdout, conflicts with --out")
	SetCmd.AddCommand(AppendCmd)
	SetCmd.AddCommand(PrependCmd)
	SetCmd.AddCommand(ReplaceCmd)
	RootCmd.AddCommand(SetCmd)
	ReplaceCmd.Flags().IntVarP(&config.ReplaceNTimes, "replace-n", "n", -1, "Limit replacements to n occurrences")
}

func performSet2(queryString string, newValue string) error {
	err := validateOutputConfig()
	if err != nil {
		return err
	}
	reader, err := getInputReader()
	if err != nil {
		return err
	}
	doc, err := hclq.FromReader(reader)
	if err != nil {
		return err
	}
	err = doc.Set2(queryString, newValue)
	if err != nil {
		return err
	}
	writer, err := getOutputWriter()
	if err != nil {
		return err
	}
	defer writer.Close()
	doc.Print(writer)
	return nil
}

func performSet(queryString string, listAction func(*ast.ListType) error, valueAction func(*token.Token) error) error {
	err := validateOutputConfig()
	if err != nil {
		return err
	}
	reader, err := getInputReader()
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
	writer, err := getOutputWriter()
	if err != nil {
		return err
	}
	defer writer.Close()
	doc.Print(writer)
	return nil
}

func getOutputWriter() (io.WriteCloser, error) {
	if config.OutputFile != "" {
		return os.Create(config.OutputFile)
	}
	if config.ModifyInPlace {
		return os.Create(config.InputFile)
	}
	return os.Stdout, nil
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
