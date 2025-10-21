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
	getBytes      bool
	getLines      bool
	getWords      bool
	getCharecters bool
)

func counter(file *os.File, splitFunc bufio.SplitFunc, resetFilePointer bool) int {
	count := 0

	scanner := bufio.NewScanner(file)
	scanner.Split(splitFunc)

	for scanner.Scan() {
		count++
	}

	if resetFilePointer {
		file.Seek(0, 0)
	}

	return count
}

var rootCmd = &cobra.Command{
	Use:   "ccwc",
	Short: "own version of the Unix command line tool wc!",
	Long:  ``,

	PreRunE: func(cmd *cobra.Command, args []string) error {
		count := cmd.Flags().NFlag()

		if count > 1 {
			return fmt.Errorf("number of flags is greater than 1")
		}

		return nil
	},

	RunE: func(cmd *cobra.Command, args []string) error {
		var filePath string

		if len(args) > 0 {
			filePath = args[0]
		} else {
			return fmt.Errorf("file path is not provided")
		}

		file, _ := os.Open(filePath)

		defer file.Close()

		if getBytes {
			fmt.Println(counter(file, bufio.ScanBytes, false), filePath)
		} else if getLines {
			fmt.Println(counter(file, bufio.ScanLines, false), filePath)
		} else if getWords {
			fmt.Println(counter(file, bufio.ScanWords, false), filePath)
		} else if getCharecters {
			fmt.Println(counter(file, bufio.ScanRunes, false), filePath)
		} else {
			numWords := counter(file, bufio.ScanWords, true)
			numLines := counter(file, bufio.ScanLines, true)
			numBytes := counter(file, bufio.ScanBytes, false)

			fmt.Printf("%d %d %d %s\n", numLines, numWords, numBytes, filePath)
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
	rootCmd.Flags().BoolVarP(&getBytes, "bytes", "c", false, "Get number of bytes of given file")
	rootCmd.Flags().BoolVarP(&getLines, "line", "l", false, "Get line count of given file")
	rootCmd.Flags().BoolVarP(&getWords, "word", "w", false, "Get word count of given file")
	rootCmd.Flags().BoolVarP(&getCharecters, "char", "m", false, "Get number of charecters of given file")
}
