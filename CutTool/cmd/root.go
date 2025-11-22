/*
Copyright © 2025 NAME HERE <sudhanvanadiger12@gmail.com>
*/
package cmd

import (
	"CutTool/utils"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	defaultDelimeter string = "\t"
	fieldsLen        int
)

var (
	fields    []int64
	delimiter string
	fileName  string
)

var rootCmd = &cobra.Command{
	Use:   "cut",
	Short: "Custom command line tool for unix cut",
	Long:  `The cut command in Unix and Linux is a powerful utility used to extract specific sections from each line of a file or input stream and write the result to standard output. It can select portions of text based on bytes, characters, or fields. `,

	PreRunE: func(cmd *cobra.Command, args []string) error {
		fieldsLen = len(fields)

		if fieldsLen == 0 && delimiter != defaultDelimeter {
			return fmt.Errorf("delimeter field is allowed with field flag")
		}

		if delimiter == "" {
			delimiter = defaultDelimeter
		}

		// values are provided with 1 based index
		for i := range fields {
			if fields[i] == 0 {
				return fmt.Errorf("field value may not be zero")
			}

			fields[i]--
		}

		if len(args) != 0 {
			fileName = args[0]
		}

		return nil
	},

	RunE: func(cmd *cobra.Command, args []string) error {
		scanner, err := utils.GetScanner(fileName)

		if err != nil {
			return err
		}

		for scanner.Scan() {
			line := scanner.Text()
			fmt.Println(utils.ProcessLine(&fields, &line, &delimiter))
		}

		return nil
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().FuncP("field", "f", "Selects specific fields, typically separated by a delimiter", func(s string) error {
		_fields, err := utils.ParseFields(s)

		if err != nil {
			return err
		}

		fields = _fields

		return nil
	})

	rootCmd.Flags().StringVarP(&delimiter, "delimiter", "d", "", "Specifies a custom delimiter for field extraction. The default delimiter is the tab character.")
}
