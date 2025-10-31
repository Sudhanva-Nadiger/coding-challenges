package main

import (
	"jsonparser/jsonparser"
	"log"
	"os"
	"strings"
)

func main() {
	args := os.Args

	args = args[1:]

	if len(args) < 1 {
		log.Fatalf("File path not provided")
		os.Exit(1)
	}

	filePath := args[0]

	if !strings.HasSuffix(filePath, ".json") {
		log.Fatalf("Invalid json file")
		os.Exit(1)
	}

	file, _ := os.ReadFile(filePath)

	p := jsonparser.NewJSONParser(string(file))
	p.Parse()
}
