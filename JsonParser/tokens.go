package main

type NumberToken rune
type EscapeToken rune
type Token rune

const (
	ZERO             NumberToken = '0'
	ONE              NumberToken = '1'
	TWO              NumberToken = '2'
	THREE            NumberToken = '3'
	FOUR             NumberToken = '4'
	FIVE             NumberToken = '5'
	SIX              NumberToken = '6'
	SEVEN            NumberToken = '7'
	EIGHT            NumberToken = '8'
	NINE             NumberToken = '9'
	NUMBER_MINUS     NumberToken = '-'
	DOT              NumberToken = '.'
	SMALL_EXPONENT   NumberToken = 'e'
	CAPITAL_EXPONENT NumberToken = 'E'
	PLUS             NumberToken = '+'
)

const (
	ESCAPE_QUOTE          EscapeToken = '"'
	REVERSE_SOLIDUS       EscapeToken = '\\'
	SOLIDUS               EscapeToken = '/'
	ESCAPE_BACKSPACE      EscapeToken = 'b'
	ESCAPE_FORM_FEED      EscapeToken = 'f'
	ESCAPE_LINE_FEED      EscapeToken = 'n'
	ESCAPE_CAR_RETURN     EscapeToken = 'r'
	ESCAPE_HORIZONTAL_TAB EscapeToken = 't'
	ESCAPE_HEX            EscapeToken = 'u'
)

const (
	BEGIN_OBJECT   Token = '{'
	END_OBJECT     Token = '}'
	BEGIN_TRUE     Token = 't'
	BEGIN_FALSE    Token = 'f'
	BEGIN_NULL     Token = 'n'
	BEGIN_ARRAY    Token = '['
	END_ARRAY      Token = ']'
	COMMA          Token = ','
	QUOTE          Token = '"'
	MINUS          Token = '-'
	SEMI_COLON     Token = ':'
	ESCAPE         Token = '\\'
	BACKSPACE      Token = '\b'
	FORM_FEED      Token = '\f'
	LINE_FEED      Token = '\n'
	CAR_RETURN     Token = '\r'
	HORIZONTAL_TAB Token = '\t'
	HEX            Token = 'u'
)
