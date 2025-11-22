/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"bufio"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	prependLineNumber                 bool
	prependLineNumberExcludeBlankLine bool
)

func shouldAcceptStreams(fileNames []string) bool {
	n := len(fileNames)

	return n == 0 || (n == 1 && fileNames[0] == "-")
}

func printLineWithNumber(line string, num *int) {
	fmt.Printf("%d %s\n", *num, line)
	(*num)++
}

func shouldPrintLineWithNumber(isEmptyLine bool) bool {
	return prependLineNumber || (prependLineNumberExcludeBlankLine && !isEmptyLine)
}

func flushContent(scanner *bufio.Scanner) {
	lineNumber := 1
	for scanner.Scan() {
		if shouldPrintLineWithNumber(scanner.Text() == "") {
			printLineWithNumber(scanner.Text(), &lineNumber)
		} else {
			fmt.Println(scanner.Text())
		}
	}
}

var rootCmd = &cobra.Command{
	Use:   "cccat",
	Short: "custom tool for cat command",
	Long:  ``,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if cmd.Flags().NFlag() > 1 {
			return fmt.Errorf("too many flags recieved. expected only one")
		}

		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		var scanner *bufio.Scanner = nil

		if shouldAcceptStreams(args) {
			scanner = bufio.NewScanner(os.Stdin)

			flushContent(scanner)
		} else {
			for _, fileName := range args {
				file, err := os.Open(fileName)

				if err != nil {
					return err
				}

				scanner = bufio.NewScanner(file)
				flushContent(scanner)

				file.Close()
			}
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

	rootCmd.Flags().BoolVarP(&prependLineNumber, "number", "n", false, "number each line")
	rootCmd.Flags().BoolVarP(&prependLineNumberExcludeBlankLine, "", "b", false, "dont number blank lines")
}
