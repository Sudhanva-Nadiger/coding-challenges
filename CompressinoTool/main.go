package main

import (
	"compressiontool/util"
	"fmt"
	"log"
	"os"
	"slices"
	"strings"
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

	fmt.Println(args)

	if argsLen != 3 {
		log.Fatalf("Too many arguments expected 3 but receieved %d\n", argsLen)
		os.Exit(1)
	}

	operation := ""

	for i := range args {
		if strings.HasPrefix(args[i], "--") {
			operation = args[i][2:]
			args = slices.Delete(args, i, i+1)
			break
		}
	}

	if operation != "compress" && operation != "decompress" {
		log.Fatalf("invalid operation %v. Allowed --compress & --decompress", operation)
		os.Exit(1)
	}

	inputFileName := args[0]
	outputFileName := args[1]

	if operation == "compress" {
		err := util.CreateCompressedFile(inputFileName, outputFileName)
		if err != nil {
			log.Panic(err)
		}

		return
	}

	err := util.DecompressFile(inputFileName, outputFileName)
	if err != nil {
		log.Panic(err)
	}
}
