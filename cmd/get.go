package cmd

import (
	"fmt"
	"os"

	"github.com/mattolenik/hclq/query"
	"github.com/spf13/cobra"
)

// GetCmd command
var GetCmd = &cobra.Command{
	Use:   "get <query>",
	Short: "retrieve values matching <query>",
	RunE: func(cmd *cobra.Command, args []string) error {
		qry, _ := query.Parse(args[0])
		reader := os.Stdin
		if val := cmd.Flag("in").Value.String(); val != "" {
			var err error
			reader, err = os.Open(val)
			if err != nil {
				return err
			}
		}
		raw := cmd.Flag("raw").Value.String() == "true"
		resultPairs, isList, _, err := query.HCL(reader, qry)
		results := []interface{}{}
		for _, pair := range resultPairs {
			results = append(results, pair.Value)
		}
		// The return type can be a list if: the queried object IS a list, or if the query matched multiple single items
		// So, return now if it's not a list and there is only one query result
		if !isList && len(results) == 1 {
			output, err := getOutput(results[0], raw)
			if err != nil {
				return err
			}
			fmt.Print(output)
			return err
		}
		output, err := getOutput(results, raw)
		if err != nil {
			return err
		}
		fmt.Print(output)
		return err
	},
}

func init() {
	var raw bool
	GetCmd.PersistentFlags().BoolVarP(&raw, "raw", "r", false, "output raw format instead of JSON")
	RootCmd.AddCommand(GetCmd)
}
