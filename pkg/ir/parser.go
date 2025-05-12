package ir

import (
	"fmt"

	"github.com/andydunstall/minc/pkg/ast"
	"github.com/andydunstall/minc/pkg/token"
)

func Parse(root ast.Node, debug bool) (n Node, err error) {
	p := newParser(debug)
	n = p.parse(root)
	return
}

type parser struct {
	counter int
}

func newParser(debug bool) *parser {
	return &parser{}
}

func (p *parser) parse(n ast.Node) Node {
	switch v := n.(type) {
	case *ast.File:
		var decls []Decl
		for _, decl := range v.Decls {
			decls = append(decls, p.parseFuncDecl(decl.(*ast.FuncDecl)))
		}
		return &File{
			Decls: decls,
		}
	default:
		panic("unsupported node type")
	}
}

// Expressions.

func (p *parser) parseExpr(expr ast.Expr) (Value, []Inst) {
	switch expr := expr.(type) {
	case *ast.UnaryExpr:
		return p.parseUnaryExpr(expr)
	case *ast.BinaryExpr:
		return p.parseBinaryExpr(expr)
	case *ast.VarExpr:
		return p.parseVarExpr(expr)
	case *ast.AssignExpr:
		return p.parseAssignExpr(expr)
	case *ast.BasicLitExpr:
		return p.parseBasicLitExpr(expr)
	default:
		panic("unsupported expr type")
	}
}

func (p *parser) parseUnaryExpr(e *ast.UnaryExpr) (Value, []Inst) {
	src, insts := p.parseExpr(e.Expr)
	dest := &VarValue{
		V: p.nextVar(),
	}
	insts = append(insts, &UnaryInst{
		Op:   e.Op,
		Src:  src,
		Dest: dest,
	})
	return dest, insts
}

func (p *parser) parseBinaryExpr(e *ast.BinaryExpr) (Value, []Inst) {
	if e.Op == token.LAND {
		falseLabel := p.nextLabel("and_false")
		endLabel := p.nextLabel("and_end")

		var insts []Inst

		v1, insts1 := p.parseExpr(e.L)
		insts = append(insts, insts1...)
		insts = append(insts, &JumpIfZeroInst{
			V:     v1,
			Label: falseLabel,
		})

		v2, insts2 := p.parseExpr(e.R)
		insts = append(insts, insts2...)
		insts = append(insts, &JumpIfZeroInst{
			V:     v2,
			Label: falseLabel,
		})

		dest := p.nextVar()
		insts = append(insts, &CopyInst{
			L: &ConstValue{
				V: "1",
			},
			R: &VarValue{
				V: dest,
			},
		})
		insts = append(insts, &JumpInst{
			Label: endLabel,
		})

		insts = append(insts, &LabelInst{
			Name: falseLabel,
		})
		insts = append(insts, &CopyInst{
			L: &ConstValue{
				V: "0",
			},
			R: &VarValue{
				V: dest,
			},
		})
		insts = append(insts, &LabelInst{
			Name: endLabel,
		})

		return &VarValue{
			V: dest,
		}, insts
	} else if e.Op == token.LOR {
		trueLabel := p.nextLabel("or_true")
		endLabel := p.nextLabel("or_end")

		var insts []Inst

		v1, insts1 := p.parseExpr(e.L)
		insts = append(insts, insts1...)
		insts = append(insts, &JumpIfNotZeroInst{
			V:     v1,
			Label: trueLabel,
		})

		v2, insts2 := p.parseExpr(e.R)
		insts = append(insts, insts2...)
		insts = append(insts, &JumpIfNotZeroInst{
			V:     v2,
			Label: trueLabel,
		})

		dest := p.nextVar()
		insts = append(insts, &CopyInst{
			L: &ConstValue{
				V: "0",
			},
			R: &VarValue{
				V: dest,
			},
		})
		insts = append(insts, &JumpInst{
			Label: endLabel,
		})

		insts = append(insts, &LabelInst{
			Name: trueLabel,
		})
		insts = append(insts, &CopyInst{
			L: &ConstValue{
				V: "0",
			},
			R: &VarValue{
				V: dest,
			},
		})
		insts = append(insts, &LabelInst{
			Name: endLabel,
		})

		return &VarValue{
			V: dest,
		}, insts
	}

	v1, insts1 := p.parseExpr(e.L)
	v2, insts2 := p.parseExpr(e.R)
	dest := &VarValue{
		V: p.nextVar(),
	}

	insts := append(insts1, insts2...)
	insts = append(insts, &BinaryInst{
		Op:   e.Op,
		V1:   v1,
		V2:   v2,
		Dest: dest,
	})
	return dest, insts
}

func (p *parser) parseVarExpr(e *ast.VarExpr) (Value, []Inst) {
	return &VarValue{
		V: e.Name,
	}, nil
}

func (p *parser) parseAssignExpr(e *ast.AssignExpr) (Value, []Inst) {
	name := e.L.(*ast.VarExpr).Name
	dest, insts := p.parseExpr(e.R)
	v := &VarValue{
		V: name,
	}
	insts = append(insts, &CopyInst{
		L: dest,
		R: v,
	})
	return v, insts
}

func (p *parser) parseBasicLitExpr(e *ast.BasicLitExpr) (Value, []Inst) {
	return &ConstValue{
		V: e.Value,
	}, nil
}

// Statements.

func (p *parser) parseStmt(stmt ast.Stmt) []Inst {
	switch stmt := stmt.(type) {
	case *ast.BlockStmt:
		return p.parseBlockStmt(stmt)
	case *ast.ReturnStmt:
		return p.parseReturnStmt(stmt)
	case *ast.ExprStmt:
		_, insts := p.parseExpr(stmt.E)
		return insts
	case *ast.DeclStmt:
		return p.parseDecl(stmt.Decl)
	case *ast.IfStmt:
		return p.parseIfStmt(stmt)
	case *ast.LoopStmt:
		return p.parseLoopStmt(stmt)
	case *ast.BreakStmt:
		return []Inst{&JumpInst{
			Label: "break." + stmt.Label,
		}}
	case *ast.ContinueStmt:
		return []Inst{&JumpInst{
			Label: "continue." + stmt.Label,
		}}
	default:
		fmt.Printf("%#v\n", stmt)
		panic("unsupported stmt type")
	}
}

func (p *parser) parseBlockStmt(stmt *ast.BlockStmt) []Inst {
	var insts []Inst
	for _, stmt := range stmt.List {
		insts = append(insts, p.parseStmt(stmt)...)
	}
	return insts
}

func (p *parser) parseReturnStmt(stmt *ast.ReturnStmt) []Inst {
	value, insts := p.parseExpr(stmt.Result)
	return append(insts, &RetInst{
		Value: value,
	})
}

func (p *parser) parseIfStmt(stmt *ast.IfStmt) []Inst {
	elseLabel := p.nextLabel("else")
	endLabel := p.nextLabel("if_end")

	c, insts := p.parseExpr(stmt.Cond)
	insts = append(insts, &JumpIfZeroInst{
		V:     c,
		Label: elseLabel,
	})

	insts = append(insts, p.parseStmt(stmt.Then)...)
	insts = append(insts, &JumpInst{
		Label: endLabel,
	})

	insts = append(insts, &LabelInst{
		Name: elseLabel,
	})
	if stmt.Else != nil {
		insts = append(insts, p.parseStmt(stmt.Else)...)
	}

	insts = append(insts, &LabelInst{
		Name: endLabel,
	})
	return insts
}

func (p *parser) parseLoopStmt(stmt *ast.LoopStmt) []Inst {
	var insts []Inst

	continueLabel := "continue." + stmt.Label
	breakLabel := "break." + stmt.Label

	insts = append(insts, &LabelInst{
		Name: continueLabel,
	})

	c, condInsts := p.parseExpr(stmt.Cond)
	insts = append(insts, condInsts...)
	insts = append(insts, &JumpIfZeroInst{
		V:     c,
		Label: breakLabel,
	})

	insts = append(insts, p.parseStmt(stmt.Body)...)

	insts = append(insts, &JumpInst{
		Label: continueLabel,
	})

	insts = append(insts, &LabelInst{
		Name: breakLabel,
	})

	return insts
}

// Declarations.

func (p *parser) parseDecl(decl ast.Decl) []Inst {
	switch decl := decl.(type) {
	case *ast.VarDecl:
		_, insts := p.parseExpr(&ast.AssignExpr{
			L: &ast.VarExpr{
				Name: decl.Name,
			},
			R: decl.Expr,
		})
		return insts
	default:
		panic("unsupported decl type")
	}
}

func (p *parser) parseFuncDecl(decl *ast.FuncDecl) Decl {
	return &FuncDecl{
		Name:  decl.Name,
		Insts: p.parseBlockStmt(decl.Body),
	}
}

func (p *parser) nextVar() string {
	s := fmt.Sprintf("tmp.%d", p.counter)
	p.counter++
	return s
}

func (p *parser) nextLabel(name string) string {
	s := fmt.Sprintf("%s.%d", name, p.counter)
	p.counter++
	return s
}
