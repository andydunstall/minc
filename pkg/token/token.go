package token

import "fmt"

type Token int

const (
	ILLEGAL Token = iota
	EOF

	// Identifiers and basic type literals
	literal_beg
	IDENT // main
	INT   // 12345
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

	EQL    // ==
	NEQ    // !=
	LSS    // <
	GTR    // >
	ASSIGN // =
	NOT    // !
	LEQ    // <=
	GEQ    // >=

	LPAREN    // (
	LBRACE    // {
	RPAREN    // )
	RBRACE    // }
	SEMICOLON // ;
	operator_end

	// Keywords
	keyword_beg
	LET

	RETURN

	IF
	ELSE

	LOOP
	CONTINUE
	BREAK
	keyword_end

	// Additional tokens
	additional_beg
	TILDE // ~
	additional_end
)

var tokens = [...]string{
	ILLEGAL: "ILLEGAL",
	EOF:     "EOF",

	IDENT: "IDENT",
	INT:   "INT",

	ADD: "+",
	SUB: "-",
	MUL: "*",
	QUO: "/",
	REM: "%",

	LAND: "&&",
	LOR:  "||",

	EQL:    "==",
	NEQ:    "!=",
	LSS:    "<",
	GTR:    ">",
	ASSIGN: "=",
	NOT:    "!",
	LEQ:    "<=",
	GEQ:    ">=",

	LPAREN:    "(",
	LBRACE:    "{",
	RPAREN:    ")",
	RBRACE:    "}",
	SEMICOLON: ";",

	LET: "let",

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

// Lookup maps an identifier to its keyword token or [IDENT] (if not a
// keyword).
func Lookup(ident string) Token {
	if tok, is_keyword := keywords[ident]; is_keyword {
		return tok
	}
	return IDENT
}

func (tok Token) IsLiteral() bool {
	return literal_beg < tok && tok < literal_end
}

func (tok Token) IsOperator() bool {
	return (operator_beg < tok && tok < operator_end) || tok == TILDE
}

func (tok Token) IsKeyword() bool {
	return keyword_beg < tok && tok < keyword_end
}
