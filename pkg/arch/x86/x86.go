package x86

import (
	"fmt"

	"github.com/andydunstall/minc/pkg/assembly"
	"github.com/andydunstall/minc/pkg/token"
)

func Emit(n assembly.Node) string {
	switch n.(type) {
	case *assembly.File:
		return emitFile(n.(*assembly.File))
	case assembly.Decl:
		return emitDecl(n.(assembly.Decl))
	case assembly.Inst:
		return emitInst(n.(assembly.Inst))
	case assembly.Operand:
		return emitOperand(n.(assembly.Operand))
	default:
		panic("unsupported node type")
	}
}

func emitFile(f *assembly.File) string {
	var s string
	for _, decl := range f.Decls {
		s += emitDecl(decl)
	}
	s += "\t.section .note.GNU-stack,\"\",@progbits\n"
	return s
}

func emitDecl(decl assembly.Decl) string {
	switch decl.(type) {
	case *assembly.FuncDecl:
		return emitFuncDecl(decl.(*assembly.FuncDecl))
	default:
		panic("unsupported decl type")
	}
}

func emitFuncDecl(decl *assembly.FuncDecl) string {
	var s string
	s += fmt.Sprintf("\t.global %s\n", decl.Name)
	s += fmt.Sprintf("%s:\n", decl.Name)
	s += "\tpushq %rbp\n"
	s += "\tmovq %rsp, %rbp\n"
	for _, inst := range decl.Insts {
		s += emitInst(inst)
	}
	return s
}

func emitInst(inst assembly.Inst) string {
	switch v := inst.(type) {
	case *assembly.MovInst:
		return emitMovInst(v)
	case *assembly.UnaryInst:
		return emitUnaryInst(v)
	case *assembly.BinaryInst:
		return emitBinaryInst(v)
	case *assembly.RetInst:
		return "\tmovq %rbp, %rsp\n\tpopq %rbp\n\tret\n"
	case *assembly.IdivInst:
		return fmt.Sprintf("\tidiv %s\n", emitOperand(v.V))
	case *assembly.CDQInst:
		return "\tcdq\n"
	case *assembly.AllocateStackInst:
		return fmt.Sprintf("\tsubq $%d, %%rsp\n", v.N)
	case *assembly.LabelInst:
		return fmt.Sprintf(".L%s:\n", v.Name)
	case *assembly.CmpInst:
		return fmt.Sprintf("\tcmpl %s, %s\n", emitOperand(v.C), emitOperand(v.V))
	case *assembly.SetCCInst:
		return fmt.Sprintf("\tset%s %s\n", emitCondCode(v.C), emitOperand(v.V))
	case *assembly.JmpInst:
		return fmt.Sprintf("\tjmp .L%s\n", v.Label)
	case *assembly.JmpCCInst:
		return fmt.Sprintf("\tj%s .L%s\n", emitCondCode(v.C), v.Label)
	default:
		fmt.Printf("%#v\n", inst)
		panic("unsupported inst type")
	}
}

func emitMovInst(inst *assembly.MovInst) string {
	return fmt.Sprintf(
		"\tmovl %s, %s\n",
		emitOperand(inst.L),
		emitOperand(inst.R),
	)
}

func emitUnaryInst(inst *assembly.UnaryInst) string {
	return fmt.Sprintf(
		"\t%s %s\n",
		emitUnaryOperator(inst.Op),
		emitOperand(inst.V),
	)
}

func emitBinaryInst(inst *assembly.BinaryInst) string {
	return fmt.Sprintf(
		"\t%s %s, %s\n",
		emitBinaryOperator(inst.Op),
		emitOperand(inst.Src),
		emitOperand(inst.Dest),
	)
}

func emitOperand(op assembly.Operand) string {
	switch v := op.(type) {
	case *assembly.RegisterOperand:
		if v.Reg == "AX" {
			return "%eax"
		} else if v.Reg == "R10" {
			return "%r10d"
		} else if v.Reg == "R11" {
			return "%r11d"
		} else if v.Reg == "DX" {
			return "%edx"
		} else {
			panic("unsupported register: " + v.Reg)
		}
	case *assembly.ImmOperand:
		return "$" + v.V
	case *assembly.PseudoOperand:
		// The IR pass should remove all pseudo operands.
		panic("pseudo operand")
	case *assembly.StackOperand:
		return fmt.Sprintf("-%d(%%rbp)", v.Offset)
	default:
		panic("unsupported operand type")
	}
}

func emitCondCode(cc assembly.CondCode) string {
	switch cc {
	case assembly.CondCodeE:
		return "e"
	case assembly.CondCodeNE:
		return "ne"
	case assembly.CondCodeG:
		return "g"
	case assembly.CondCodeGE:
		return "ge"
	case assembly.CondCodeL:
		return "l"
	case assembly.CondCodeLE:
		return "le"
	default:
		panic("unknown cond code")
	}
}

func emitUnaryOperator(op token.Token) string {
	switch op {
	case token.TILDE:
		return "notl"
	case token.SUB:
		return "negl"
	default:
		panic("unsupported unary operator: " + op.String())
	}
}

func emitBinaryOperator(op token.Token) string {
	switch op {
	case token.ADD:
		return "addl"
	case token.SUB:
		return "subl"
	case token.MUL:
		return "imull"
	default:
		panic("unsupported binary operator: " + op.String())
	}
}
