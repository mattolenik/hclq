package main

import (
	"fmt"
	"io"
	"os"

	command "github.com/mattolenik/hclq/cmd"
	"github.com/mattolenik/hclq/common"
	"github.com/mattolenik/hclq/hclq"

	"github.com/jawher/mow.cli"
	"github.com/logrusorgru/aurora"
	"github.com/mattn/go-isatty"
)

var app *cli.Cli
var au aurora.Aurora

var version string

func init() {
	// Enable color only if running inside a terminal
	au = aurora.NewAurora(isatty.IsTerminal(os.Stdout.Fd()))

	app = cli.App("hclq", "querying and editing tool for HCL files (such as .tf)")
	app.LongDesc = LONG_DESC
}

func actionErr(action func() error) func() {
	return func() {
		err := action()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %s\n", err.Error())
			os.Exit(1)
		}
	}
}

func getIO(inFile, outFile *string, inPlace *bool) (reader io.Reader, writer io.WriteCloser, err error) {
	reader, err = common.GetInputReader(inFile)
	if err != nil {
		return
	}
	writer, err = common.GetOutputWriter(inFile, outFile, inPlace)
	return reader, writer, err
}

func main() {
	inFile := app.StringOpt("in", "", "read input from this file, otherwise use stdin")
	outFile := app.StringOpt("out", "", "read output from this file, otherwise use stdout")

	app.Command("get", "retrieve values matching query", func(cmd *cli.Cmd) {
		query := cmd.StringArg("QUERY", "", "query string matching value(s) to get")
		raw := cmd.BoolOpt("r raw", false, "output raw (Go-like) formatting, rather than JSON")
		inPlace := false  // always false for get commands

		cmd.Action = actionErr(func() error {
			reader, err := common.GetInputReader(inFile)
			if err != nil {
				return err
			}
			doc, err := hclq.FromReader(reader)
			if err != nil {
				return err
			}
			result, err := doc.Get(*query)
			if err != nil {
				return err
			}
			output, err := common.GetOutput(result, *raw)
			if err != nil {
				return err
			}
			writer, err := common.GetOutputWriter(inFile, outFile, &inPlace)
			if err != nil {
				return err
			}
			defer writer.Close()
			fmt.Fprintln(writer, output)
			return nil
		})
	})

	app.Command("get-keys", "retrieve keys matching query", func(cmd *cli.Cmd) {
		query := cmd.StringArg("QUERY", "", "query string matching key(s) to get")
		raw := cmd.BoolOpt("r raw", false, "output raw (Go-like) formatting, rather than JSON")
		inFile := cmd.StringOpt("in", "", "read from this file instead of stdin")
		outFile := cmd.StringOpt("out", "", "write to this file instead of stdout")
		inPlace := false  // always false for get commands

		cmd.Action = actionErr(func() error {
			reader, err := common.GetInputReader(inFile)
			if err != nil {
				return err
			}
			doc, err := hclq.FromReader(reader)
			if err != nil {
				return err
			}
			result, err := doc.GetKeys(*query)
			if err != nil {
				return err
			}
			output, err := common.GetOutput(result, *raw)
			if err != nil {
				return err
			}
			writer, err := common.GetOutputWriter(inFile, outFile, &inPlace)
			if err != nil {
				return err
			}
			defer writer.Close()
			fmt.Fprintln(writer, output)
			return nil
		})
	})

	app.Command("set", "set matching value(s)", func(cmd *cli.Cmd) {
		cmd.Spec = "QUERY VALUE [--in=<file>] [--out=<file>|--in-place --in=<file>]"
		query := cmd.StringArg("QUERY", "", "query string matching key(s) to set")
		newValue := cmd.StringArg("VALUE", "", "replacement value")
		inFile := cmd.StringOpt("in", "", "read from this file instead of stdin")
		outFile := cmd.StringOpt("out", "", "write to this file instead of stdout, conflicts with --in-place")
		inPlace := cmd.BoolOpt("in-place", false, "edit the input file in-place, conflicts with --out")

		cmd.Action = actionErr(func() error {
			reader, writer, err := getIO(inFile, outFile, inPlace)
			if err != nil {
				return err
			}
			return command.Set(reader, writer, *query, *newValue)
		})
	})

	app.Command("append", "append matching value(s)", func(cmd *cli.Cmd) {
		cmd.Spec = "QUERY VALUE [--in=<file>] [--out=<file>|--in-place --in=<file>]"
		query := cmd.StringArg("QUERY", "", "query string matching key(s) to append")
		newValue := cmd.StringArg("VALUE", "", "value to append")
		inFile := cmd.StringOpt("in", "", "read from this file instead of stdin")
		outFile := cmd.StringOpt("out", "", "write to this file instead of stdout, conflicts with --in-place")
		inPlace := cmd.BoolOpt("in-place", false, "edit the input file in-place, conflicts with --out")

		cmd.Action = actionErr(func() error {
			reader, writer, err := getIO(inFile, outFile, inPlace)
			if err != nil {
				return err
			}
			return command.Append(reader, writer, *query, *newValue)
		})
	})

	app.Command("prepend", "prepend matching value(s)", func(cmd *cli.Cmd) {
		cmd.Spec = "QUERY VALUE [--in=<file>] [--out=<file>|--in-place --in=<file>]"
		query := cmd.StringArg("QUERY", "", "query string matching key(s) to prepend")
		newValue := cmd.StringArg("VALUE", "", "value to prepend")
		inFile := cmd.StringOpt("in", "", "read from this file instead of stdin")
		outFile := cmd.StringOpt("out", "", "write to this file instead of stdout, conflicts with --in-place")
		inPlace := cmd.BoolOpt("in-place", false, "edit the input file in-place, conflicts with --out")

		cmd.Action = actionErr(func() error {
			reader, writer, err := getIO(inFile, outFile, inPlace)
			if err != nil {
				return err
			}
			return command.Prepend(reader, writer, *query, *newValue)
		})
	})


	app.Command("replace", "find/replace matching value(s)", func(cmd *cli.Cmd) {
		cmd.Spec = "QUERY OLD_VALUE NEW_VALUE [--in=<file>] [--out=<file>|--in-place --in=<file>]"
		query := cmd.StringArg("QUERY", "", "query string matching key(s) to prepend")
		oldValue := cmd.StringArg("OLD_VALUE", "", "value to find")
		newValue := cmd.StringArg("NEW_VALUE", "", "replacement value")
		inFile := cmd.StringOpt("in", "", "read from this file instead of stdin")
		outFile := cmd.StringOpt("out", "", "write to this file instead of stdout, conflicts with --in-place")
		inPlace := cmd.BoolOpt("in-place", false, "edit the input file in-place, conflicts with --out")
		n := cmd.IntOpt("n", 0, "replace at most n times")

		cmd.Action = actionErr(func() error {
			reader, writer, err := getIO(inFile, outFile, inPlace)
			if err != nil {
				return err
			}
			return command.Replace(reader, writer, *query, *oldValue, *newValue, *n)
		})
	})

	app.Version("version", version)

	app.Run(os.Args)
}

func fail(formatString string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, formatString, args...)
	fmt.Fprintf(os.Stderr, "\n")
	cli.Exit(1)
}

const LONG_DESC = `hclq uses a "breadcrumb" or "path" style query. Given the HCL:
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
`
