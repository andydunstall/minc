package assembly

import (
	"github.com/andydunstall/minc/pkg/ir"
	"github.com/andydunstall/minc/pkg/token"
)

var paramPassingRegs = []string{"DI", "SI", "DX", "CX", "R8", "R9"}

func Parse(root ir.Node, debug bool) (n Node, err error) {
	p := newParser(debug)
	n = p.parse(root)
	return
}

type parser struct {
}

func newParser(debug bool) *parser {
	return &parser{}
}

func (p *parser) parse(n ir.Node) Node {
	switch v := n.(type) {
	case ir.Decl:
		return p.parseDecl(v)
	case ir.Value:
		return p.parseValue(v)
	case *ir.File:
		var decls []Decl
		for _, decl := range v.Decls {
			decls = append(decls, p.parseDecl(decl))
		}
		return &File{
			Decls: decls,
		}
	default:
		panic("unsupported node type")
	}
}

// Operands.

func (p *parser) parseValue(v ir.Value) Operand {
	switch v := v.(type) {
	case *ir.ConstValue:
		return &ImmOperand{
			V: v.V,
		}
	case *ir.VarValue:
		return &PseudoOperand{
			V: v.V,
		}
	default:
		panic("unsupported value type")
	}
}

// Declarations.

func (p *parser) parseDecl(decl ir.Decl) Decl {
	switch decl := decl.(type) {
	case *ir.FuncDecl:
		return p.parseFuncDecl(decl)
	default:
		panic("unsupported decl type")
	}
}

func (p *parser) parseFuncDecl(decl *ir.FuncDecl) *FuncDecl {
	var insts []Inst

	for i := 0; i != 6; i++ {
		if i >= len(decl.Params) {
			break
		}

		insts = append(insts, &MovInst{
			L: &RegisterOperand{
				Reg: paramPassingRegs[i],
			},
			R: &PseudoOperand{
				V: decl.Params[i],
			},
		})
	}
	for i := 6; i < len(decl.Params); i++ {
		insts = append(insts, &MovInst{
			L: &StackOperand{
				Offset: int32(16 + 8*(i-6)),
			},
			R: &PseudoOperand{
				V: decl.Params[i],
			},
		})
	}

	for _, inst := range decl.Insts {
		insts = append(insts, p.parseInst(inst)...)
	}

	return &FuncDecl{
		Name:  decl.Name,
		Insts: insts,
	}
}

// Instructions.

func (p *parser) parseInst(inst ir.Inst) (insts []Inst) {
	switch v := inst.(type) {
	case *ir.RetInst:
		return p.parseRetInst(v)
	case *ir.UnaryInst:
		return p.parseUnaryInst(v)
	case *ir.BinaryInst:
		return p.parseBinaryInst(v)
	case *ir.CopyInst:
		return p.parseCopyInst(v)
	case *ir.JumpInst:
		return p.parseJumpInst(v)
	case *ir.JumpIfZeroInst:
		return p.parseJumpIfZeroInst(v)
	case *ir.JumpIfNotZeroInst:
		return p.parseJumpIfNotZeroInst(v)
	case *ir.CallInst:
		return p.parseCallInst(v)
	case *ir.LabelInst:
		return []Inst{
			&LabelInst{
				Name: v.Name,
			},
		}
	default:
		panic("unsupported inst type")
	}

	return
}

func (p *parser) parseRetInst(inst *ir.RetInst) []Inst {
	V := p.parseValue(inst.Value)
	return []Inst{
		&MovInst{
			L: V,
			R: &RegisterOperand{
				Reg: "AX",
			},
		},
		&RetInst{},
	}
}

func (p *parser) parseUnaryInst(inst *ir.UnaryInst) []Inst {
	src := p.parseValue(inst.Src)
	dest := p.parseValue(inst.Dest)

	if inst.Op == token.NOT {
		return []Inst{
			&CmpInst{
				C: &ImmOperand{
					V: "0",
				},
				V: src,
			},
			&MovInst{
				L: &ImmOperand{
					V: "0",
				},
				R: dest,
			},
			&SetCCInst{
				C: CondCodeE,
				V: dest,
			},
		}
	}

	return []Inst{
		&MovInst{
			L: src,
			R: dest,
		},
		&UnaryInst{
			Op: inst.Op,
			V:  dest,
		},
	}
}

func (p *parser) parseBinaryInst(inst *ir.BinaryInst) []Inst {
	v1 := p.parseValue(inst.V1)
	v2 := p.parseValue(inst.V2)
	dest := p.parseValue(inst.Dest)

	switch inst.Op {
	case token.QUO, token.REM:
		reg := "AX"
		if inst.Op == token.REM {
			reg = "DX"
		}
		return []Inst{
			&MovInst{
				L: v1,
				R: &RegisterOperand{
					Reg: "AX",
				},
			},
			&CDQInst{},
			&IdivInst{
				V: v2,
			},
			&MovInst{
				L: &RegisterOperand{
					Reg: reg,
				},
				R: dest,
			},
		}
	case token.EQL, token.NEQ, token.LSS, token.LEQ, token.GTR, token.GEQ:
		v1 := p.parseValue(inst.V1)
		v2 := p.parseValue(inst.V2)
		dest := p.parseValue(inst.Dest)

		var code CondCode
		switch inst.Op {
		case token.EQL:
			code = CondCodeE
		case token.NEQ:
			code = CondCodeNE
		case token.LSS:
			code = CondCodeL
		case token.LEQ:
			code = CondCodeLE
		case token.GTR:
			code = CondCodeG
		case token.GEQ:
			code = CondCodeGE
		}

		return []Inst{
			&CmpInst{
				C: v2,
				V: v1,
			},
			&MovInst{
				&ImmOperand{
					V: "0",
				},
				dest,
			},
			&SetCCInst{
				C: code,
				V: dest,
			},
		}
	default:
		return []Inst{
			&MovInst{
				L: v1,
				R: dest,
			},
			&BinaryInst{
				Op:   inst.Op,
				Src:  v2,
				Dest: dest,
			},
		}
	}
}

func (p *parser) parseCopyInst(inst *ir.CopyInst) []Inst {
	l := p.parseValue(inst.L)
	r := p.parseValue(inst.R)
	return []Inst{
		&MovInst{
			L: l,
			R: r,
		},
	}
}

func (p *parser) parseJumpInst(inst *ir.JumpInst) []Inst {
	return []Inst{
		&JmpInst{
			Label: inst.Label,
		},
	}
}

func (p *parser) parseJumpIfZeroInst(inst *ir.JumpIfZeroInst) []Inst {
	return []Inst{
		&CmpInst{
			C: &ImmOperand{
				V: "0",
			},
			V: p.parseValue(inst.V),
		},
		&JmpCCInst{
			C:     CondCodeE,
			Label: inst.Label,
		},
	}
}

func (p *parser) parseJumpIfNotZeroInst(inst *ir.JumpIfNotZeroInst) []Inst {
	return []Inst{
		&CmpInst{
			C: &ImmOperand{
				V: "0",
			},
			V: p.parseValue(inst.V),
		},
		&JmpCCInst{
			C:     CondCodeNE,
			Label: inst.Label,
		},
	}
}

func (p *parser) parseCallInst(inst *ir.CallInst) []Inst {
	var insts []Inst

	var padding int32
	if len(inst.Args)%2 != 0 {
		padding = 8
		// Padding.
		insts = append(insts, &AllocateStackInst{
			N: padding,
		})
	}

	// Push args to registers.
	n := len(inst.Args)
	if n > 6 {
		n = 6
	}
	for i := 0; i < n; i++ {
		insts = append(insts, &MovInst{
			L: p.parseValue(inst.Args[i]),
			R: &RegisterOperand{
				Reg: paramPassingRegs[i],
			},
		})
	}

	// Push args to stack.
	for i := 6; i < len(inst.Args); i++ {
		v := p.parseValue(inst.Args[i])
		switch v := v.(type) {
		case *ImmOperand, *RegisterOperand:
			insts = append(insts, &PushInst{
				V: v,
			})
		default:
			insts = append(insts, &MovInst{
				L: v,
				R: &RegisterOperand{
					Reg: "AX",
				},
			})
			insts = append(insts, &PushInst{
				V: &RegisterOperand{
					Reg: "AX",
				},
			})

		}
	}

	insts = append(insts, &CallInst{
		inst.Name,
	})

	insts = append(insts, &DeallocateStackInst{
		N: int32(len(inst.Args)*8) + padding,
	})

	dest := p.parseValue(inst.Dest)
	insts = append(insts, &MovInst{
		L: &RegisterOperand{
			Reg: "AX",
		},
		R: dest,
	})

	return insts
}
