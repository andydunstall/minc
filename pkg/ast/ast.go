package ast

import "github.com/andydunstall/minc/pkg/token"

type Node interface {
	node()
}

// Expressions.

type Expr interface {
	Node
	exprNode()
}

type UnaryExpr struct {
	Op   token.Token
	Expr Expr
}

func (n *UnaryExpr) node()     {}
func (n *UnaryExpr) exprNode() {}

type BinaryExpr struct {
	Op token.Token
	L  Expr
	R  Expr
}

func (n *BinaryExpr) node()     {}
func (n *BinaryExpr) exprNode() {}

type VarExpr struct {
	Name string
}

func (n *VarExpr) node()     {}
func (n *VarExpr) exprNode() {}

type AssignExpr struct {
	L Expr
	R Expr
}

func (n *AssignExpr) node()     {}
func (n *AssignExpr) exprNode() {}

type CallExpr struct {
	Func string
	Args []Expr
}

func (n *CallExpr) node()     {}
func (n *CallExpr) exprNode() {}

type BasicLitExpr struct {
	Kind  token.Token
	Value string
}

func (n *BasicLitExpr) node()     {}
func (n *BasicLitExpr) exprNode() {}

// Statements.

type Stmt interface {
	Node
	stmtNode()
}

type BlockStmt struct {
	List []Stmt
}

func (n *BlockStmt) node()     {}
func (n *BlockStmt) stmtNode() {}

type ReturnStmt struct {
	Result Expr
}

func (n *ReturnStmt) node()     {}
func (n *ReturnStmt) stmtNode() {}

type ExprStmt struct {
	E Expr
}

func (n *ExprStmt) node()     {}
func (n *ExprStmt) stmtNode() {}

type DeclStmt struct {
	Decl Decl
}

func (n *DeclStmt) node()     {}
func (n *DeclStmt) stmtNode() {}

type IfStmt struct {
	Cond Expr
	Then Stmt
	Else Stmt
}

func (n *IfStmt) node()     {}
func (n *IfStmt) stmtNode() {}

type LoopStmt struct {
	Cond Expr
	Body *BlockStmt

	Label string
}

func (n *LoopStmt) node()     {}
func (n *LoopStmt) stmtNode() {}

type BreakStmt struct {
	Label string
}

func (n *BreakStmt) node()     {}
func (n *BreakStmt) stmtNode() {}

type ContinueStmt struct {
	Label string
}

func (n *ContinueStmt) node()     {}
func (n *ContinueStmt) stmtNode() {}

// Declarations.

type Decl interface {
	Node
	declNode()
}

type VarDecl struct {
	Name string
	Expr Expr
}

func (n *VarDecl) node()     {}
func (n *VarDecl) declNode() {}

type FuncType struct {
	// TODO(andydunstall): No types yet, only identifier names.
	Params []string
}

func (n *FuncType) node() {}

type FuncDecl struct {
	Name string
	Type *FuncType
	Body *BlockStmt
}

func (n *FuncDecl) node()     {}
func (n *FuncDecl) declNode() {}

type File struct {
	Decls []Decl
}

func (n *File) node() {}
