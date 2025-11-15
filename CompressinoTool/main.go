package main

import (
	"compressiontool/util"
	"log"
	"os"
)

func main() {
	args := os.Args

	// removing first arg which is default program name
	args = args[1:]

	argsLen := len(args)

	if argsLen == 0 {
		log.Fatalln("File path is not provided")
		os.Exit(1)
	}

	if argsLen != 2 {
		log.Fatalf("Too many arguments expected 2 but receieved %d\n", argsLen)
		os.Exit(1)
	}

	inputFileName := args[0]
	outputFileName := args[1]

	util.CreateCompressedFile(inputFileName, outputFileName)

}
