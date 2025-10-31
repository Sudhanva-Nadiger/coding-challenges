package jsonparser

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"unicode"
)

const MAX_DEPTH = 20

type Jsonparser struct {
	input  string
	cursor int
	depth  int
}

func NewJSONParser(input string) *Jsonparser {
	return &Jsonparser{
		input:  input,
		cursor: 0,
		depth:  0,
	}
}

func (p *Jsonparser) Parse() {
	p.consumeWhiteSpace()

	obj := p.parseValue()

	// Validate that top-level value is an object or array
	switch obj.(type) {
	case JSONObject, JSONArray:
		// Valid top-level types
	default:
		log.Fatalf("Top-level JSON value must be an object or array, not %T", obj)
		os.Exit(1)
	}

	p.consumeWhiteSpace()

	if p.HasNext() {
		log.Fatalf("Unexpected token %v at position %v", p.currentToken(), p.cursor)
	}

	fmt.Println("Successfully parsed the value")
	fmt.Println(obj)

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
	case rune(BEGIN_ARRAY):
		return p.parseArray()
	default:
		log.Fatalf("Unexpected token %v at position %v", p.currentToken(), p.cursor)
		os.Exit(1)
		return ""
	}
}

func (p *Jsonparser) parseObject() JSONObject {
	// Check depth limit before entering
	p.depth++
	if p.depth >= MAX_DEPTH {
		log.Fatalf("Maximum nesting depth of %d exceeded at position %v", MAX_DEPTH, p.cursor)
		os.Exit(1)
	}
	defer func() { p.depth-- }() // Decrement depth when exiting

	obj := make(map[string]JSONValue)

	p.consume(BEGIN_OBJECT, true)

	hasMorePair := false

	for p.currentToken() != rune(END_OBJECT) || hasMorePair {
		pair := p.parsePair()

		obj[pair.key] = pair.value

		currToken := p.currentToken()

		if currToken == rune(COMMA) {
			p.consume(COMMA, true)
			hasMorePair = true
		} else if currToken != rune(END_OBJECT) {
			hasMorePair = false
			log.Fatalf("Invalid object at %v", p.cursor)
			os.Exit(1)
		} else {
			hasMorePair = false
		}
	}

	p.consume(END_OBJECT, true)

	return obj
}

func (p *Jsonparser) parseArray() JSONArray {
	// Check depth limit before entering
	p.depth++
	if p.depth >= MAX_DEPTH {
		log.Fatalf("Maximum nesting depth of %d exceeded at position %v", MAX_DEPTH, p.cursor)
		os.Exit(1)
	}
	defer func() { p.depth-- }() // Decrement depth when exiting

	arr := []JSONValue{}

	p.consume(BEGIN_ARRAY, true)

	hasMoreElements := false

	for p.currentToken() != rune(END_ARRAY) || hasMoreElements {
		value := p.parseValue()

		arr = append(arr, value)

		if p.currentToken() == rune(COMMA) {
			p.consume(COMMA, true)
			hasMoreElements = true
		} else if p.currentToken() != rune(END_ARRAY) {
			hasMoreElements = false
			log.Fatalf("Invalid object at %v", p.cursor)
			os.Exit(1)
		} else {
			hasMoreElements = false
		}
	}

	p.consume(END_ARRAY, true)

	return arr
}

func (p *Jsonparser) parsePair() *KeyValuePair {

	key := p.parseString()

	p.consume(SEMI_COLON, true)

	value := p.parseValue()

	return &KeyValuePair{
		key:   key,
		value: value,
	}
}

func (p *Jsonparser) parseString() string {
	str := ""

	p.consume(QUOTE, true)

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

	p.consume(QUOTE, true)
	return str
}

func (p *Jsonparser) parseEscape() string {
	p.consume(ESCAPE, false)

	if unicode.IsControl(p.currentToken()) {
		log.Fatalf("Invalid character at %v. control chars should be escaped", p.cursor)
		os.Exit(1)
	}

	switch currToken := p.currentToken(); currToken {
	case rune(ESCAPE_QUOTE), rune(REVERSE_SOLIDUS), rune(SOLIDUS):
		p.consume(nil, true)
		return string(currToken)
	case rune(ESCAPE_BACKSPACE),
		rune(ESCAPE_CAR_RETURN),
		rune(ESCAPE_FORM_FEED),
		rune(ESCAPE_LINE_FEED),
		rune(ESCAPE_HORIZONTAL_TAB):
		p.consume(nil, true)
		return string(escapeTokenToTokenMap[EscapeToken(currToken)])
	case rune(ESCAPE_HEX):
		p.consume(nil, false) // Move past the 'u'
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
	p.consume(BEGIN_TRUE, false)
	p.consume(TRUE_R, false)
	p.consume(TRUE_U, false)
	p.consume(TRUE_E, true)

	return true
}

func (p *Jsonparser) parseFalse() bool {
	p.consume(BEGIN_FALSE, false)
	p.consume(FALSE_A, false)
	p.consume(FALSE_L, false)
	p.consume(FALSE_S, false)
	p.consume(FALSE_E, true)

	return false
}

func (p *Jsonparser) parseNull() any {
	p.consume(BEGIN_NULL, false)
	p.consume(NULL_U, false)
	p.consume(NULL_L, false)
	p.consume(NULL_L, true)

	return nil
}

func (p *Jsonparser) parseNumber() float64 {
	str := ""

	if p.currentToken() == rune(NUMBER_MINUS) {
		str += string(NUMBER_MINUS)
		p.consume(NUMBER_MINUS, false)
	}

	str += p.parseDigits(false)

	if p.currentToken() == rune(DOT) {
		str += string(p.currentToken())
		p.consume(DOT, false)
		str += p.parseDigits(true)
	}

	if p.currentToken() == rune(SMALL_EXPONENT) || p.currentToken() == rune(CAPITAL_EXPONENT) {
		str += string(p.currentToken())
		p.consume(nil, false)

		if p.currentToken() == rune(PLUS) || p.currentToken() == rune(MINUS) {
			str += string(p.currentToken())
			p.consume(nil, false)
		}

		str += p.parseDigits(true)
	}

	p.consumeWhiteSpace()

	num, _ := strconv.ParseFloat(str, 64)
	return num
}

func (p *Jsonparser) parseDigits(allowLeadingZero bool) string {
	str := ""

	if allowLeadingZero {
		// For fractional or exponent parts, must have at least one digit
		if p.currentToken() < rune(ZERO) || p.currentToken() > rune(NINE) {
			log.Fatalf("Expected digit but found %v at position %v", string(p.currentToken()), p.cursor)
			os.Exit(1)
		}

		// Consume all digits (including leading zeros)
		for p.currentToken() >= rune(ZERO) && p.currentToken() <= rune(NINE) {
			str += string(p.currentToken())
			p.consume(nil, false)
		}
	} else {
		// For integer part: single zero OR non-zero digit followed by any digits
		if p.currentToken() == rune(ZERO) {
			str += string(p.currentToken())
			p.consume(ZERO, false)
			// Don't consume more digits after a leading zero (per JSON spec)
		} else if p.currentToken() >= rune(ONE) && p.currentToken() <= rune(NINE) {
			str += string(p.currentToken())
			p.consume(nil, false)

			for p.currentToken() >= rune(ZERO) && p.currentToken() <= rune(NINE) {
				str += string(p.currentToken())
				p.consume(nil, false)
			}
		} else {
			log.Fatalf("Invalid character %v at position %v", p.currentToken(), p.cursor)
			os.Exit(1)
		}
	}

	return str
}

func (p *Jsonparser) consume(expected any, skip bool) {
	if expected != nil {
		var expectedRune rune
		switch v := expected.(type) {
		case Token:
			expectedRune = rune(v)
		case NumberToken:
			expectedRune = rune(v)
		default:
			log.Fatalf("Invalid token type at position %v", p.cursor)
			os.Exit(1)
		}

		if expectedRune != p.currentToken() {
			log.Fatalf("Expected %v but found %v at position %v", string(expectedRune), string(p.currentToken()), p.cursor)
			os.Exit(1)
		}
	}

	p.cursor++

	if skip {
		for p.cursor < len(p.input) && unicode.IsSpace(p.currentToken()) {
			p.cursor++
		}
	}
}

func (p *Jsonparser) consumeWhiteSpace() {
	for p.cursor < len(p.input) && unicode.IsSpace(p.currentToken()) {
		p.consume(nil, true)
	}
}

func (p *Jsonparser) currentToken() rune {
	if p.cursor >= len(p.input) {
		return 0 // Return null character if we're at the end
	}
	return rune(p.input[p.cursor])
}

func (p *Jsonparser) HasNext() bool {
	p.consumeWhiteSpace()
	return p.cursor < len(p.input)
}
