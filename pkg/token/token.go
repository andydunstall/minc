// Package token defines the lexical tokens of C.
package token

import "fmt"

type Token int

const (
	ILLEGAL Token = iota

	// Identifiers and basic type literals
	literal_beg
	IDENT  // main
	INT    // 12345
	CHAR   // 'a'
	STRING // "abc"
	literal_end

	// Operators and delimiters
	operator_beg
	ADD // +
	SUB // -
	MUL // *
	QUO // /
	REM // %

	LAND // &&
	LOR  // ||

	INC // ++
	DEC // --

	EQL    // ==
	NEQ    // !=
	LSS    // <
	GTR    // >
	ASSIGN // =
	NOT    // !
	LEQ    // <=
	GEQ    // >=

	LPAREN    // (
	LBRACK    // [
	LBRACE    // {
	RPAREN    // )
	RBRACK    // ]
	RBRACE    // }
	SEMICOLON // ;
	COLON     // :
	QMARK     // ?
	ARROW     // ->

	COMMA  // ,
	PERIOD // .
	operator_end

	// Keywords
	keyword_beg
	LET
	MUT

	FN
	RETURN

	IF
	ELSE

	LOOP
	CONTINUE
	BREAK
	keyword_end

	// Additional tokens
	additional_beg
	TILDE
	additional_end
)

var tokens = [...]string{
	ILLEGAL: "ILLEGAL",

	IDENT:  "IDENT",
	INT:    "INT",
	CHAR:   "CHAR",
	STRING: "STRING",

	ADD: "+",
	SUB: "-",
	MUL: "*",
	QUO: "/",
	REM: "%",

	LAND: "&&",
	LOR:  "||",

	INC: "++",
	DEC: "--",

	EQL:    "==",
	NEQ:    "!=",
	LSS:    "<",
	GTR:    ">",
	ASSIGN: "=",
	NOT:    "!",
	LEQ:    "<=",
	GEQ:    ">=",

	LPAREN:    "(",
	LBRACK:    "[",
	LBRACE:    "{",
	RPAREN:    ")",
	RBRACK:    "]",
	RBRACE:    "}",
	SEMICOLON: ";",
	COLON:     ":",
	QMARK:     "?",
	ARROW:     "->",

	COMMA:  ",",
	PERIOD: ".",

	LET: "let",
	MUT: "mut",

	FN:     "fn",
	RETURN: "return",

	IF:   "if",
	ELSE: "else",

	LOOP:     "loop",
	CONTINUE: "continue",
	BREAK:    "break",

	TILDE: "~",
}

func (tok Token) String() string {
	if 0 <= tok && tok < Token(len(tokens)) {
		return tokens[tok]
	}
	panic(fmt.Sprintf("unknown token: %d", tok))
}

var keywords map[string]Token

func init() {
	keywords = make(map[string]Token, keyword_end-(keyword_beg+1))
	for i := keyword_beg + 1; i < keyword_end; i++ {
		keywords[tokens[i]] = i
	}
}

// Lookup maps an identifier to its keyword token or [IDENT] (if not a keyword).
func Lookup(ident string) Token {
	if tok, is_keyword := keywords[ident]; is_keyword {
		return tok
	}
	return IDENT
}
