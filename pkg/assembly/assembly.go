package assembly

import "github.com/andydunstall/minc/pkg/token"

type CondCode int

const (
	CondCodeE CondCode = iota + 1
	CondCodeNE
	CondCodeG
	CondCodeGE
	CondCodeL
	CondCodeLE
)

type Node interface {
	node()
}

// Operands.

type Operand interface {
	Node
	operandNode()
}

type ImmOperand struct {
	V string
}

func (n *ImmOperand) node()        {}
func (n *ImmOperand) operandNode() {}

type PseudoOperand struct {
	V string
}

func (n *PseudoOperand) node()        {}
func (n *PseudoOperand) operandNode() {}

type StackOperand struct {
	Offset int32
}

func (n *StackOperand) node()        {}
func (n *StackOperand) operandNode() {}

type RegisterOperand struct {
	Reg string
}

func (n *RegisterOperand) node()        {}
func (n *RegisterOperand) operandNode() {}

// Declarations.

type Decl interface {
	Node
	declNode()
}

type FuncDecl struct {
	Name  string
	Insts []Inst
}

func (n *FuncDecl) node()     {}
func (n *FuncDecl) declNode() {}

// Instructions.

type Inst interface {
	Node
	instNode()
}

type MovInst struct {
	L Operand
	R Operand
}

func (n *MovInst) node()     {}
func (n *MovInst) instNode() {}

type RetInst struct{}

func (n *RetInst) node()     {}
func (n *RetInst) instNode() {}

type UnaryInst struct {
	Op token.Token
	V  Operand
}

func (n *UnaryInst) node()     {}
func (n *UnaryInst) instNode() {}

type BinaryInst struct {
	Op   token.Token
	Src  Operand
	Dest Operand
}

func (n *BinaryInst) node()     {}
func (n *BinaryInst) instNode() {}

type IdivInst struct {
	V Operand
}

func (n *IdivInst) node()     {}
func (n *IdivInst) instNode() {}

type CDQInst struct{}

func (n *CDQInst) node()     {}
func (n *CDQInst) instNode() {}

type JmpInst struct {
	Label string
}

func (n *JmpInst) node()     {}
func (n *JmpInst) instNode() {}

type SetCCInst struct {
	C CondCode
	V Operand
}

func (n *SetCCInst) node()     {}
func (n *SetCCInst) instNode() {}

type JmpCCInst struct {
	C     CondCode
	Label string
}

func (n *JmpCCInst) node()     {}
func (n *JmpCCInst) instNode() {}

type CmpInst struct {
	C Operand
	V Operand
}

func (n *CmpInst) node()     {}
func (n *CmpInst) instNode() {}

type LabelInst struct {
	Name string
}

func (n *LabelInst) node()     {}
func (n *LabelInst) instNode() {}

type PushInst struct {
	V Operand
}

func (n *PushInst) node()     {}
func (n *PushInst) instNode() {}

type AllocateStackInst struct {
	N int32
}

func (n *AllocateStackInst) node()     {}
func (n *AllocateStackInst) instNode() {}

type DeallocateStackInst struct {
	N int32
}

func (n *DeallocateStackInst) node()     {}
func (n *DeallocateStackInst) instNode() {}

type CallInst struct {
	Func string
}

func (n *CallInst) node()     {}
func (n *CallInst) instNode() {}

type File struct {
	Decls []Decl
}

func (n *File) node() {}
