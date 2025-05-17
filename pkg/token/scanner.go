package token

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

func NewScanner(src []byte) *Scanner {
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

func (s *Scanner) Scan() (tok Token, lit string) {
	s.skipWhitespace()

	switch ch := s.ch; {
	case isLetter(ch):
		lit = s.scanIdentifier()
		tok = Lookup(lit)
	case isDecimal(ch):
		lit = s.scanNumber()
		tok = INT
	default:
		s.next()
		switch ch {
		case '+':
			tok = ADD
		case '-':
			tok = SUB
		case '*':
			tok = MUL
		case '/':
			tok = QUO
		case '%':
			tok = REM
		case '&':
			if s.ch == '&' {
				tok = LAND
				s.next()
			} else {
				tok = ILLEGAL
			}
		case '|':
			if s.ch == '|' {
				tok = LOR
				s.next()
			} else {
				tok = ILLEGAL
			}
		case '=':
			if s.ch == '=' {
				tok = EQL
				s.next()
			} else {
				tok = ASSIGN
			}
		case '!':
			if s.ch == '=' {
				tok = NEQ
				s.next()
			} else {
				tok = NOT
			}
		case '<':
			if s.ch == '=' {
				tok = LEQ
				s.next()
			} else {
				tok = LSS
			}
		case '>':
			if s.ch == '=' {
				tok = GEQ
				s.next()
			} else {
				tok = GTR
			}
		case '(':
			tok = LPAREN
		case '{':
			tok = LBRACE
		case ')':
			tok = RPAREN
		case '}':
			tok = RBRACE
		case ';':
			tok = SEMICOLON
		case ',':
			tok = COMMA
		case '~':
			tok = TILDE
		case eof:
			tok = EOF
		default:
			tok = ILLEGAL
		}
	}

	return
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

func (s *Scanner) next() {
	if s.offset < len(s.src)-1 {
		s.offset++
		s.ch = s.src[s.offset]
	} else {
		s.offset = len(s.src)
		s.ch = eof
	}
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
