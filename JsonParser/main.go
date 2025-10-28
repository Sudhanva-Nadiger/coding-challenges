package main

func NewJSONParser(input string) *Jsonparser {
	return &Jsonparser{
		input:  input,
		cursor: 0,
	}
}

func main() {
	p := NewJSONParser("{}")
	p.Parse()
}
