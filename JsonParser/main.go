package main

import (
	"jsonparser/jsonparser"
	"os"
)

func main() {
	file, _ := os.ReadFile("tests/step3/valid.json")

	p := jsonparser.NewJSONParser(string(file))
	p.Parse()
}
