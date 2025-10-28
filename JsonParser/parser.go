package main

import (
	"fmt"
	"log"
	"os"
	"unicode"
)

type Jsonparser struct {
	input  string
	cursor int
}

func (p *Jsonparser) Parse() {
	p.ConsumeWhiteSpace()

	p.parseValue()

	p.ConsumeWhiteSpace()

	if p.HasNext() {
		log.Fatalf("Unexpected token %v at position %v", p.CurrentToken(), p.cursor)
	}

	fmt.Println("Successfully parsed the value")

	os.Exit(0)
}

func (p *Jsonparser) parseValue() JSONValue {
	switch p.CurrentToken() {
	case rune(BEGIN_OBJECT):
		return p.parseObject()
	}

	return ""
}

func (p *Jsonparser) parseObject() JSONObject {
	var obj JSONObject

	p.Consume(BEGIN_OBJECT, true)

	hasMorePair := false

	for p.CurrentToken() != rune(END_OBJECT) || hasMorePair {
		pair := p.parsePair()

		obj[pair.key] = pair.value

		currToken := p.CurrentToken()

		if currToken == rune(COMMA) {
			p.Consume(COMMA, true)
			hasMorePair = true
		} else if currToken != rune(END_OBJECT) {
			hasMorePair = false
			log.Fatalf("Invalid object at %v", p.cursor)
			os.Exit(1)
		} else {
			hasMorePair = false
		}
	}

	p.Consume(END_OBJECT, true)

	return obj
}

func (p *Jsonparser) parsePair() *KeyValuePair {
	key := p.parseString()

	p.Consume(SEMI_COLON, true)

	value := p.parseValue()

	return &KeyValuePair{
		key:   key,
		value: value,
	}
}

func (p *Jsonparser) parseString() string {
	str := ""

	p.Consume(QUOTE, true)

	for currToken := p.CurrentToken(); currToken != rune(QUOTE); {
		if currToken == rune(ESCAPE) {
			str += p.parseEscape()
		} else {
			if unicode.IsControl(currToken) {
				log.Fatalf("Invalid character at %v. control chars should be escaped", p.cursor)
				os.Exit(1)
			}

			str += string(currToken)
			p.cursor++
		}
	}

	p.Consume(QUOTE, true)

	return str
}

func (p *Jsonparser) parseEscape() string {
	p.Consume(ESCAPE, false)
	return ""
}

func (p *Jsonparser) Consume(expected any, skip bool) {
	if expected != nil && rune(expected.(Token)) != p.CurrentToken() {
		log.Fatalf("Expected %v but found %v at position %v", expected, p.CurrentToken(), p.cursor)
		os.Exit(1)
	}

	p.cursor++

	if skip {
		for p.cursor < len(p.input) && unicode.IsSpace(p.CurrentToken()) {
			p.cursor++
		}
	}
}

func (p *Jsonparser) ConsumeWhiteSpace() {
	for p.cursor < len(p.input) && unicode.IsSpace(p.CurrentToken()) {
		p.Consume(nil, true)
	}
}

func (p *Jsonparser) CurrentToken() rune {
	if p.cursor >= len(p.input) {
		return 0 // Return null character if we're at the end
	}
	return rune(p.input[p.cursor])
}

func (p *Jsonparser) HasNext() bool {
	p.ConsumeWhiteSpace()
	return p.cursor < len(p.input)
}
