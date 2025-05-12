package ast

import (
	"fmt"
)

type varEntry struct {
	name      string
	fromScope bool
}

// Validate performs semantic analysis on the AST:
// - Verify variables are defined
// - Map variables to a unique name
// - Add a label for each loop
func Validate(root Node, debug bool) (Node, error) {
	v := newValidator(debug)
	return v.validate(root)
}

type validator struct {
	vars map[string]varEntry

	varCounter int
	// A loop count of 0 means not in a loop.
	loopCounter int
}

func newValidator(debug bool) *validator {
	return &validator{
		vars:        make(map[string]varEntry),
		varCounter:  1,
		loopCounter: 0,
	}
}

func (v *validator) validate(n Node) (Node, error) {
	switch n := n.(type) {
	case *File:
		var decls []Decl
		for _, decl := range n.Decls {
			decls = append(decls, v.validateDecl(decl))
		}
		return &File{
			Decls: decls,
		}, nil
	default:
		panic("unsupported node type")
	}
}

// Expressions.

func (v *validator) validateExpr(expr Expr) Expr {
	switch expr := expr.(type) {
	case *VarExpr:
		e, ok := v.vars[expr.Name]
		if !ok {
			panic("undeclared variable: " + expr.Name)
		}
		// Map the variable to its updated name.
		expr.Name = e.name
	case *AssignExpr:
		if _, ok := expr.L.(*VarExpr); !ok {
			panic("expected variable")
		}
		expr.L = v.validateExpr(expr.L)
		expr.R = v.validateExpr(expr.R)
	case *UnaryExpr:
		expr.Expr = v.validateExpr(expr.Expr)
	case *BinaryExpr:
		expr.L = v.validateExpr(expr.L)
		expr.R = v.validateExpr(expr.R)
	}
	return expr
}

// Statements.

func (v *validator) validateStmt(stmt Stmt) Stmt {
	switch stmt := stmt.(type) {
	case *DeclStmt:
		stmt.Decl = v.validateDecl(stmt.Decl)
	case *ReturnStmt:
		stmt.Result = v.validateExpr(stmt.Result)
	case *ExprStmt:
		stmt.E = v.validateExpr(stmt.E)
	case *IfStmt:
		stmt.Cond = v.validateExpr(stmt.Cond)
		stmt.Then = v.validateStmt(stmt.Then)
		if stmt.Else != nil {
			stmt.Else = v.validateStmt(stmt.Else)
		}
	case *LoopStmt:
		// Add a unique label to each loop.

		existingLoop := v.loopCounter
		v.loopCounter++

		stmt.Cond = v.validateExpr(stmt.Cond)
		stmt.Body = v.validateBlockStmt(stmt.Body)

		stmt.Label = v.loopLabel()
		v.loopCounter = existingLoop
	case *ContinueStmt:
		if v.loopCounter == 0 {
			panic("not in loop")
		}
		// Point to closing loop.
		stmt.Label = v.loopLabel()
	case *BreakStmt:
		if v.loopCounter == 0 {
			panic("not in loop")
		}
		// Point to closing loop.
		stmt.Label = v.loopLabel()
	case *BlockStmt:
		return v.validateBlockStmt(stmt)
	}
	return stmt
}

func (v *validator) validateBlockStmt(block *BlockStmt) *BlockStmt {
	// Block statements create a new scope. Therefore store the current
	// variables, and create a new scope for the block. After the block,
	// reset to the existing scope.

	existingVars := make(map[string]varEntry)
	for k, e := range v.vars {
		existingVars[k] = e
		e.fromScope = false
		v.vars[k] = e
	}

	for i, stmt := range block.List {
		block.List[i] = v.validateStmt(stmt)
	}

	v.vars = existingVars

	return block
}

// Declarations.

func (v *validator) validateDecl(decl Decl) Decl {
	switch decl := decl.(type) {
	case *FuncDecl:
		decl.Body = v.validateBlockStmt(decl.Body)
	case *VarDecl:
		v.validateVarDecl(decl)
	}
	return decl
}

func (v *validator) validateVarDecl(decl *VarDecl) *VarDecl {
	e, ok := v.vars[decl.Name]
	if ok && e.fromScope {
		panic("duplicate declaration: " + decl.Name)
	}

	updatedName := v.nextVar(decl.Name)
	v.vars[decl.Name] = varEntry{
		name:      updatedName,
		fromScope: true,
	}

	decl.Name = updatedName
	decl.Expr = v.validateExpr(decl.Expr)
	return decl
}

func (v *validator) loopLabel() string {
	return fmt.Sprintf("loop.%d", v.loopCounter)
}

func (v *validator) nextVar(name string) string {
	n := fmt.Sprintf("%s.%d", name, v.varCounter)
	v.varCounter++
	return n
}
