package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"unicode"
)

type Jsonparser struct {
	input  string
	cursor int
}

func NewJSONParser(input string) *Jsonparser {
	return &Jsonparser{
		input:  input,
		cursor: 0,
	}
}

func (p *Jsonparser) Parse() {
	fmt.Println(p.input)

	p.ConsumeWhiteSpace()

	p.parseValue()

	p.ConsumeWhiteSpace()

	if p.HasNext() {
		log.Fatalf("Unexpected token %v at position %v", p.currentToken(), p.cursor)
	}

	fmt.Println("Successfully parsed the value")

	os.Exit(0)
}

func (p *Jsonparser) parseValue() JSONValue {
	switch p.currentToken() {
	case rune(BEGIN_OBJECT):
		return p.parseObject()
	case rune(QUOTE):
		return p.parseString()
	case rune(BEGIN_TRUE):
		return p.parseTrue()
	case rune(BEGIN_FALSE):
		return p.parseFalse()
	case rune(BEGIN_NULL):
		return p.parseNull()
	case rune(ZERO),
		rune(ONE),
		rune(TWO),
		rune(THREE),
		rune(FOUR),
		rune(FIVE),
		rune(SIX),
		rune(SEVEN),
		rune(EIGHT),
		rune(NINE),
		rune(NUMBER_MINUS):
		return p.parseNumber()
	default:
		return "default"
	}
}

func (p *Jsonparser) parseObject() JSONObject {
	obj := make(map[string]JSONValue)

	p.Consume(BEGIN_OBJECT, true)

	hasMorePair := false

	for p.currentToken() != rune(END_OBJECT) || hasMorePair {
		pair := p.parsePair()

		obj[pair.key] = pair.value

		currToken := p.currentToken()

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

	for p.currentToken() != rune(QUOTE) {

		if p.currentToken() == rune(ESCAPE) {
			str += p.parseEscape()
		} else {
			if unicode.IsControl(p.currentToken()) {
				log.Fatalf("Invalid character at %v. control chars should be escaped", p.cursor)
				os.Exit(1)
			}

			str += string(p.currentToken())
			p.cursor++
		}
	}

	p.Consume(QUOTE, true)
	return str
}

func (p *Jsonparser) parseEscape() string {
	p.Consume(ESCAPE, false)

	if unicode.IsControl(p.currentToken()) {
		log.Fatalf("Invalid character at %v. control chars should be escaped", p.cursor)
		os.Exit(1)
	}

	switch currToken := p.currentToken(); currToken {
	case rune(ESCAPE_QUOTE), rune(REVERSE_SOLIDUS), rune(SOLIDUS):
		p.Consume(nil, true)
		return string(currToken)
	case rune(ESCAPE_BACKSPACE),
		rune(ESCAPE_CAR_RETURN),
		rune(ESCAPE_FORM_FEED),
		rune(ESCAPE_LINE_FEED),
		rune(ESCAPE_HORIZONTAL_TAB):
		p.Consume(nil, true)
		return string(EscapeTokenToTokenMap[EscapeToken(currToken)])
	case rune(ESCAPE_HEX):
		hexStr := p.input[p.cursor : p.cursor+4]
		code, err := strconv.ParseInt(hexStr, 16, 32)

		if err != nil {
			fmt.Printf("invalid hex code '%s' at position %d", hexStr, p.cursor)
			os.Exit(1)
		}

		p.cursor += 4
		return fmt.Sprint(code)
	default:
		fmt.Printf("Invalid escape token at position %d", p.cursor)
		os.Exit(1)
	}

	return ""
}

func (p *Jsonparser) parseTrue() bool {
	p.Consume(BEGIN_TRUE, false)
	p.Consume(TRUE_R, false)
	p.Consume(TRUE_U, false)
	p.Consume(TRUE_E, false)

	return true
}

func (p *Jsonparser) parseFalse() bool {
	p.Consume(BEGIN_FALSE, false)
	p.Consume(FALSE_A, false)
	p.Consume(FALSE_L, false)
	p.Consume(FALSE_S, false)
	p.Consume(FALSE_E, false)

	return false
}

func (p *Jsonparser) parseNull() any {
	p.Consume(BEGIN_NULL, false)
	p.Consume(NULL_U, false)
	p.Consume(NULL_L, false)
	p.Consume(NULL_L, false)

	return nil
}

func (p *Jsonparser) parseNumber() float64 {
	str := ""

	if p.currentToken() == rune(NUMBER_MINUS) {
		str += string(NUMBER_MINUS)
		p.Consume(NUMBER_MINUS, false)
	}

	str += p.parseDigits(false)

	num, _ := strconv.ParseFloat(str, 64)
	return num
}

func (p *Jsonparser) parseDigits(allowLeadingZero bool) string {
	return ""
}

func (p *Jsonparser) Consume(expected any, skip bool) {
	if expected != nil && rune(expected.(Token)) != p.currentToken() {
		log.Fatalf("Expected %v but found %v at position %v", string(expected.(Token)), string(p.currentToken()), p.cursor)
		os.Exit(1)
	}

	p.cursor++

	if skip {
		for p.cursor < len(p.input) && unicode.IsSpace(p.currentToken()) {
			p.cursor++
		}
	}
}

func (p *Jsonparser) ConsumeWhiteSpace() {
	for p.cursor < len(p.input) && unicode.IsSpace(p.currentToken()) {
		p.Consume(nil, true)
	}
}

func (p *Jsonparser) currentToken() rune {
	if p.cursor >= len(p.input) {
		return 0 // Return null character if we're at the end
	}
	return rune(p.input[p.cursor])
}

func (p *Jsonparser) HasNext() bool {
	p.ConsumeWhiteSpace()
	return p.cursor < len(p.input)
}
