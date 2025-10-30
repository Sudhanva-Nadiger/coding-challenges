package main

type NumberToken rune
type EscapeToken rune
type Token rune

var (
	EscapeTokenToTokenMap map[EscapeToken]Token
)

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
	BEGIN_OBJECT Token = '{'
	END_OBJECT   Token = '}'

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

	BEGIN_TRUE Token = 't'
	TRUE_R     Token = 'r'
	TRUE_U     Token = 'u'
	TRUE_E     Token = 'e'

	BEGIN_FALSE Token = 'f'
	FALSE_A     Token = 'a'
	FALSE_L     Token = 'l'
	FALSE_S     Token = 's'
	FALSE_E     Token = 'e'

	BEGIN_NULL Token = 'n'
	NULL_U     Token = 'u'
	NULL_L     Token = 'l'
)

func init() {
	// 1. You must 'make' the map before you can add to it
	EscapeTokenToTokenMap = make(map[EscapeToken]Token)

	// 2. Populate the map
	EscapeTokenToTokenMap[ESCAPE_QUOTE] = QUOTE
	EscapeTokenToTokenMap[ESCAPE_BACKSPACE] = BACKSPACE
	EscapeTokenToTokenMap[ESCAPE_FORM_FEED] = FORM_FEED
	EscapeTokenToTokenMap[ESCAPE_LINE_FEED] = LINE_FEED
	EscapeTokenToTokenMap[ESCAPE_CAR_RETURN] = CAR_RETURN
	EscapeTokenToTokenMap[ESCAPE_HORIZONTAL_TAB] = HORIZONTAL_TAB
	EscapeTokenToTokenMap[ESCAPE_HEX] = HEX
}
