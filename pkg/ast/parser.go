package ast

import (
	"fmt"

	"github.com/andydunstall/minc/pkg/token"
)

func Parse(scanner *token.Scanner, debug bool) (f *File, err error) {
	p := newParser(scanner, debug)
	f = p.parseFile()
	return
}

type parser struct {
	tok token.Token
	lit string

	scanner *token.Scanner

	line   int
	indent int
	debug  bool
}

func newParser(scanner *token.Scanner, debug bool) *parser {
	p := &parser{
		scanner: scanner,
		line:    1,
		debug:   debug,
	}
	p.tok, p.lit = p.scanner.Scan()
	return p
}

func (p *parser) parseFile() *File {
	if p.debug {
		defer un(trace(p, "File"))
	}

	var decls []Decl
	for p.tok != token.EOF {
		decls = append(decls, p.parseDecl())
	}

	return &File{
		Decls: decls,
	}
}

// Expressions.

func (p *parser) parseExpr(minPrec int) Expr {
	if p.debug {
		defer un(trace(p, "Expr"))
	}

	l := p.parseFactor()
	for {
		prec := p.precedence(p.tok)
		if prec <= minPrec {
			break
		}

		if p.tok == token.ASSIGN {
			l = p.parseAssignExpr(l, prec)
		} else {
			l = p.parseBinaryExpr(l, prec)
		}
	}

	return l
}

func (p *parser) parseAssignExpr(l Expr, prec int) *AssignExpr {
	if p.debug {
		defer un(trace(p, "AssignExpr"))
	}

	p.expect(token.ASSIGN)
	return &AssignExpr{
		L: l,
		R: p.parseExpr(prec),
	}
}

func (p *parser) parseBinaryExpr(l Expr, prec int) *BinaryExpr {
	if p.debug {
		defer un(trace(p, "BinaryExpr"))
	}

	op := p.tok
	p.next()

	return &BinaryExpr{
		Op: op,
		L:  l,
		R:  p.parseExpr(prec + 1),
	}
}

func (p *parser) parseCallExpr(name string) *CallExpr {
	if p.debug {
		defer un(trace(p, "CallExpr"))
	}

	var args []Expr

	p.expect(token.LPAREN)
	for p.tok != token.RPAREN {
		args = append(args, p.parseExpr(0))

		if p.tok != token.RPAREN {
			p.expect(token.COMMA)
		}
	}
	p.expect(token.RPAREN)

	return &CallExpr{
		Func: name,
		Args: args,
	}
}

func (p *parser) parseFactor() Expr {
	if p.debug {
		defer un(trace(p, "Factor"))
	}

	switch p.tok {
	case token.INT:
		f := &BasicLitExpr{
			Kind:  p.tok,
			Value: p.lit,
		}
		p.next()
		return f
	case token.SUB, token.TILDE, token.NOT:
		op := p.tok
		p.next()
		expr := p.parseExpr(0)
		return &UnaryExpr{
			Op:   op,
			Expr: expr,
		}
	case token.LPAREN:
		p.next()
		expr := p.parseExpr(0)
		p.expect(token.RPAREN)
		return expr
	case token.IDENT:
		name := p.lit
		p.next()

		if p.tok == token.LPAREN {
			return p.parseCallExpr(name)
		} else {
			return &VarExpr{
				Name: name,
			}
		}
	default:
		panic("unknown: " + p.tok.String())
	}
}

// Statements.

func (p *parser) parseStmt() (s Stmt) {
	if p.debug {
		defer un(trace(p, "Stmt"))
	}

	switch p.tok {
	case token.LBRACE:
		s = p.parseBlockStmt()
	case token.RETURN:
		s = p.parseReturnStmt()
	case token.LET:
		s = p.parseDeclStmt()
	case token.IF:
		s = p.parseIfStmt()
	case token.LOOP:
		s = p.parseLoopStmt()
	case token.BREAK:
		s = p.parseBreakStmt()
	case token.CONTINUE:
		s = p.parseContinueStmt()
	default:
		s = p.parseExprStmt()
	}
	return
}

func (p *parser) parseBlockStmt() *BlockStmt {
	if p.debug {
		defer un(trace(p, "BlockStmt"))
	}

	p.expect(token.LBRACE)
	var list []Stmt
	for p.tok != token.RBRACE && p.tok != token.EOF {
		list = append(list, p.parseStmt())
	}
	p.expect(token.RBRACE)
	return &BlockStmt{
		List: list,
	}
}

func (p *parser) parseReturnStmt() *ReturnStmt {
	if p.debug {
		defer un(trace(p, "ReturnStmt"))
	}

	p.expect(token.RETURN)

	expr := p.parseExpr(0)
	p.expect(token.SEMICOLON)
	return &ReturnStmt{
		Result: expr,
	}
}

func (p *parser) parseExprStmt() *ExprStmt {
	if p.debug {
		defer un(trace(p, "ExprStmt"))
	}

	expr := p.parseExpr(0)
	p.expect(token.SEMICOLON)
	return &ExprStmt{
		E: expr,
	}
}

func (p *parser) parseDeclStmt() *DeclStmt {
	if p.debug {
		defer un(trace(p, "DeclStmt"))
	}

	return &DeclStmt{
		Decl: p.parseDecl(),
	}
}

func (p *parser) parseIfStmt() *IfStmt {
	if p.debug {
		defer un(trace(p, "IfStmt"))
	}

	p.expect(token.IF)
	p.expect(token.LPAREN)
	cond := p.parseExpr(0)
	p.expect(token.RPAREN)
	thenStmt := p.parseStmt()

	var elseStmt Stmt
	if p.tok == token.ELSE {
		p.next()
		elseStmt = p.parseStmt()
	}

	return &IfStmt{
		Cond: cond,
		Then: thenStmt,
		Else: elseStmt,
	}
}

func (p *parser) parseLoopStmt() *LoopStmt {
	if p.debug {
		defer un(trace(p, "LoopStmt"))
	}

	p.expect(token.LOOP)
	p.expect(token.LPAREN)
	cond := p.parseExpr(0)
	p.expect(token.RPAREN)
	body := p.parseBlockStmt()
	return &LoopStmt{
		Cond: cond,
		Body: body,
	}
}

func (p *parser) parseBreakStmt() *BreakStmt {
	if p.debug {
		defer un(trace(p, "BreakStmt"))
	}

	p.expect(token.BREAK)
	p.expect(token.SEMICOLON)

	return &BreakStmt{}
}

func (p *parser) parseContinueStmt() *ContinueStmt {
	if p.debug {
		defer un(trace(p, "ContinueStmt"))
	}

	p.expect(token.CONTINUE)
	p.expect(token.SEMICOLON)

	return &ContinueStmt{}
}

// Declaration.

func (p *parser) parseDecl() Decl {
	if p.debug {
		defer un(trace(p, "Decl"))
	}

	switch p.tok {
	case token.FN:
		return p.parseFuncDecl()
	case token.LET:
		return p.parseVarDecl()
	default:
		panic("unsupported decl")
	}
}

func (p *parser) parseFuncDecl() *FuncDecl {
	if p.debug {
		defer un(trace(p, "FuncDecl"))
	}

	p.expect(token.FN)
	funcName := p.parseIdent()

	var funcType FuncType

	p.expect(token.LPAREN)
	for p.tok != token.RPAREN {
		paramType := p.lit
		p.expect(token.IDENT)
		if paramType != "int" {
			panic("unsupported type: " + paramType)
		}

		name := p.lit
		p.expect(token.IDENT)

		funcType.Params = append(funcType.Params, name)

		if p.tok != token.RPAREN {
			p.expect(token.COMMA)
		}
	}
	p.expect(token.RPAREN)

	body := p.parseBlockStmt()
	return &FuncDecl{
		Name: funcName,
		Type: &funcType,
		Body: body,
	}
}

func (p *parser) parseVarDecl() *VarDecl {
	if p.debug {
		defer un(trace(p, "VarDecl"))
	}

	p.expect(token.LET)
	name := p.lit
	p.expect(token.IDENT)
	p.expect(token.ASSIGN)

	expr := p.parseExpr(0)
	p.expect(token.SEMICOLON)

	return &VarDecl{
		Name: name,
		Expr: expr,
	}
}

func (p *parser) parseIdent() string {
	ident := p.lit
	p.expect(token.IDENT)
	return ident
}

func (p *parser) expect(tok token.Token) {
	if p.tok != tok {
		panic("unexpected token: " + p.tok.String())
	}
	p.next()
}

func (p *parser) next() {
	if p.debug {
		s := p.tok.String()
		switch {
		case p.tok.IsLiteral():
			p.printTrace(s, p.lit)
		case p.tok.IsOperator(), p.tok.IsKeyword():
			p.printTrace("\"" + s + "\"")
		default:
			p.printTrace(s)
		}
	}

	p.tok, p.lit = p.scanner.Scan()
}

func (p *parser) precedence(tok token.Token) int {
	switch tok {
	case token.MUL, token.QUO, token.REM:
		return 50
	case token.ADD, token.SUB:
		return 45
	case token.LSS, token.LEQ, token.GTR, token.GEQ:
		return 35
	case token.EQL, token.NEQ:
		return 30
	case token.LAND:
		return 10
	case token.LOR:
		return 5
	case token.ASSIGN:
		return 1
	default:
		return -1
	}
}

func (p *parser) printTrace(a ...any) {
	fmt.Printf("%6d  ", p.line)

	const dots = ". . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . "
	const n = len(dots)
	i := 2 * p.indent
	for i > n {
		fmt.Print(dots)
		i -= n
	}
	// i <= n
	fmt.Print(dots[0:i])
	fmt.Println(a...)

	p.line++
}

func trace(p *parser, msg string) *parser {
	p.printTrace(msg, "(")
	p.indent++
	return p
}

// Usage pattern: defer un(trace(p, "..."))
func un(p *parser) {
	p.indent--
	p.printTrace(")")
}
