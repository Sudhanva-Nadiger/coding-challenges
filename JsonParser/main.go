package main

import (
	"os"
)

func main() {
	file, _ := os.ReadFile("tests/step3/valid.json")

	p := NewJSONParser(string(file))
	p.Parse()
}
