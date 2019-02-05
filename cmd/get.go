package cmd

import (
	"fmt"
	"os"

	"github.com/mattolenik/hclq/config"
	"github.com/mattolenik/hclq/hclq"
	"github.com/spf13/cobra"
)

// GetCmd command
var GetCmd = &cobra.Command{
	Use:   "get <query>",
	Short: "retrieve values matching <query>",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		input, err := getInputReader()
		if err != nil {
			return err
		}
		doc, err := hclq.FromReader(input)
		if err != nil {
			return err
		}
		result, err := doc.Get(args[0])
		if err != nil {
			return err
		}
		output, err := getOutput(result, config.UseRawOutput)
		if err != nil {
			return err
		}
		if config.OutputFile != "" {
			file, err := os.Create(config.OutputFile)
			if err != nil {
				return err
			}
			defer file.Close()
			fmt.Fprintln(file, output)
		} else {
			fmt.Println(output)
		}
		return nil
	},
}

func init() {
	GetCmd.PersistentFlags().BoolVarP(&config.UseRawOutput, "raw", "r", false, "output raw format instead of JSON")
	RootCmd.AddCommand(GetCmd)
}
