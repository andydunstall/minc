package scanner

import (
	"github.com/andydunstall/minc/pkg/token"
)

const (
	eof = 0xff // end of file
)

type Scanner struct {
	// immutable state
	src []byte

	// scanning state
	ch     byte // current character
	offset int  // character offset
}

func New(src []byte) *Scanner {
	ch := byte(eof)
	if len(src) > 0 {
		ch = src[0]
	}
	return &Scanner{
		src:    src,
		ch:     ch,
		offset: 0,
	}
}

func (s *Scanner) Scan() (tok token.Token, lit string) {
	s.skipWhitespace()

	switch ch := s.ch; {
	case isLetter(ch):
		lit = s.scanIdentifier()
		tok = token.Lookup(lit)
	case isDecimal(ch):
		lit = s.scanNumber()
		tok = token.INT
	default:
		s.next()
		switch ch {
		case '+':
			if s.ch == '+' {
				tok = token.INC
				s.next()
			} else {
				tok = token.ADD
			}
		case '-':
			if s.ch == '-' {
				tok = token.DEC
				s.next()
			} else if s.ch == '>' {
				tok = token.ARROW
				s.next()
			} else {
				tok = token.SUB
			}
		case '*':
			tok = token.MUL
		case '/':
			tok = token.QUO
		case '%':
			tok = token.REM
		case '=':
			if s.ch == '=' {
				tok = token.EQL
				s.next()
			} else {
				tok = token.ASSIGN
			}
		case '!':
			if s.ch == '=' {
				tok = token.NEQ
				s.next()
			} else {
				tok = token.NOT
			}
		case '<':
			if s.ch == '=' {
				tok = token.LEQ
				s.next()
			} else {
				tok = token.LSS
			}
		case '>':
			if s.ch == '=' {
				tok = token.GEQ
				s.next()
			} else {
				tok = token.GTR
			}
		case '(':
			tok = token.LPAREN
		case '[':
			tok = token.LBRACK
		case '{':
			tok = token.LBRACE
		case ')':
			tok = token.RPAREN
		case ']':
			tok = token.RBRACK
		case '}':
			tok = token.RBRACE
		case ';':
			tok = token.SEMICOLON
		case ':':
			tok = token.COLON
		case '?':
			tok = token.QMARK
		case ',':
			tok = token.COMMA
		case '.':
			tok = token.PERIOD
		default:
			tok = token.ILLEGAL
		}
	}

	return
}

func (s *Scanner) next() {
	if s.offset < len(s.src)-1 {
		s.offset++
		s.ch = s.src[s.offset]
	} else {
		s.offset = len(s.src)
		s.ch = eof
	}
}

func (s *Scanner) scanIdentifier() string {
	for i, b := range s.src[s.offset:] {
		if 'a' <= b && b <= 'z' || 'A' <= b && b <= 'Z' || b == '_' || '0' <= b && b <= '9' {
			continue
		}

		ident := s.src[s.offset : s.offset+i]
		s.offset += i
		s.ch = s.src[s.offset]
		return string(ident)
	}

	panic("eof")
}

func (s *Scanner) scanNumber() string {
	for i, b := range s.src[s.offset:] {
		if '0' <= b && b <= '9' {
			continue
		}

		ident := s.src[s.offset : s.offset+i]
		s.offset += i
		s.ch = s.src[s.offset]
		return string(ident)
	}

	panic("eof")
}

func (s *Scanner) skipWhitespace() {
	for s.ch == ' ' || s.ch == '\t' || s.ch == '\n' || s.ch == '\r' {
		s.next()
	}
}

func isLetter(ch byte) bool {
	return 'a' <= lower(ch) && lower(ch) <= 'z' || ch == '_'
}

func isDecimal(ch byte) bool { return '0' <= ch && ch <= '9' }

func lower(ch byte) byte { return ('a' - 'A') | ch }
