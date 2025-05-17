package ir

import "github.com/andydunstall/minc/pkg/token"

type Node interface {
	node()
}

// Values.

type Value interface {
	Node
	valueNode()
}

type ConstValue struct {
	V string
}

func (n *ConstValue) node()      {}
func (n *ConstValue) valueNode() {}

type VarValue struct {
	V string
}

func (n *VarValue) node()      {}
func (n *VarValue) valueNode() {}

// Declarations.

type Decl interface {
	Node
	declNode()
}

type FuncDecl struct {
	Name   string
	Params []string
	Insts  []Inst
}

func (n *FuncDecl) node()     {}
func (n *FuncDecl) declNode() {}

// Instructions.

type Inst interface {
	Node
	instNode()
}

type RetInst struct {
	Value Value
}

func (n *RetInst) node()     {}
func (n *RetInst) instNode() {}

type UnaryInst struct {
	Op   token.Token
	Src  Value
	Dest Value
}

func (n *UnaryInst) node()     {}
func (n *UnaryInst) instNode() {}

type BinaryInst struct {
	Op   token.Token
	V1   Value
	V2   Value
	Dest Value
}

func (n *BinaryInst) node()     {}
func (n *BinaryInst) instNode() {}

type CopyInst struct {
	L Value
	R Value
}

func (n *CopyInst) node()     {}
func (n *CopyInst) instNode() {}

type JumpInst struct {
	Label string
}

func (n *JumpInst) node()     {}
func (n *JumpInst) instNode() {}

type JumpIfZeroInst struct {
	V     Value
	Label string
}

func (n *JumpIfZeroInst) node()     {}
func (n *JumpIfZeroInst) instNode() {}

type JumpIfNotZeroInst struct {
	V     Value
	Label string
}

func (n *JumpIfNotZeroInst) node()     {}
func (n *JumpIfNotZeroInst) instNode() {}

type CallInst struct {
	Name string
	Args []Value
	Dest Value
}

func (n *CallInst) node()     {}
func (n *CallInst) instNode() {}

type LabelInst struct {
	Name string
}

func (n *LabelInst) node()     {}
func (n *LabelInst) instNode() {}

type File struct {
	Decls []Decl
}

func (n *File) node() {}
