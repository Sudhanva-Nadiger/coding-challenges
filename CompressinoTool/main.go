package main

import (
	"compressiontool/huffmantree"
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

	if argsLen != 1 {
		log.Fatalf("Too many arguments expected 1 but receieved %d\n", argsLen)
		os.Exit(1)
	}

	fileName := args[0]

	content, err := os.ReadFile(fileName)

	if err != nil {
		log.Fatalf("Could not read file at path %v", err)
		os.Exit(1)
	}

	frequencyMap := make(map[rune]int64)

	for _, char := range string(content) {
		frequencyMap[char]++
	}

	huffmanTree := huffmantree.NewHuffManTree(frequencyMap)

	huffmanTree.BuildHuffManTree()

}
