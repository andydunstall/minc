package assembly

import (
	"github.com/andydunstall/minc/pkg/token"
)

func Fix(root *File, debug bool) *File {
	r := newFixer(debug)
	r.fix(root)
	return root
}

type fixer struct{}

func newFixer(debug bool) *fixer {
	return &fixer{}
}

func (f *fixer) fix(root *File) {
	for _, decl := range root.Decls {
		fn := decl.(*FuncDecl)
		fn.Insts = f.fixInsts(fn.Insts)
	}
}

func (f *fixer) fixInsts(insts []Inst) []Inst {
	insts, offset := f.replacePseudos(insts)

	var updatedInsts []Inst

	updatedInsts = append(updatedInsts, &AllocateStackInst{
		N: roundUpToNextMultipleOf16(-offset),
	})
	for _, inst := range insts {
		switch v := inst.(type) {
		case *MovInst:
			if _, ok := v.L.(*StackOperand); !ok {
				break
			}
			if _, ok := v.R.(*StackOperand); !ok {
				break
			}

			// Mov can't move a value from one memory address to another.

			updatedInsts = append(updatedInsts, &MovInst{
				L: v.L,
				R: &RegisterOperand{
					Reg: "R10",
				},
			})
			updatedInsts = append(updatedInsts, &MovInst{
				L: &RegisterOperand{
					Reg: "R10",
				},
				R: v.R,
			})

			continue
		case *IdivInst:
			if _, ok := v.V.(*ImmOperand); !ok {
				break
			}

			// Idiv can't operate on constants.

			updatedInsts = append(updatedInsts, &MovInst{
				L: v.V,
				R: &RegisterOperand{
					Reg: "R10",
				},
			})
			updatedInsts = append(updatedInsts, &IdivInst{
				V: &RegisterOperand{
					Reg: "R10",
				},
			})

			continue
		case *BinaryInst:
			switch v.Op {
			case token.ADD, token.SUB:
				if _, ok := v.Src.(*StackOperand); !ok {
					break
				}
				if _, ok := v.Dest.(*StackOperand); !ok {
					break
				}

				// Add/Sub can't use memory addresses for both operands.

				updatedInsts = append(updatedInsts, &MovInst{
					L: v.Src,
					R: &RegisterOperand{
						Reg: "R10",
					},
				})
				updatedInsts = append(updatedInsts, &BinaryInst{
					Op: v.Op,
					Src: &RegisterOperand{
						Reg: "R10",
					},
					Dest: v.Dest,
				})

				continue
			case token.MUL:
				if _, ok := v.Dest.(*StackOperand); !ok {
					break
				}

				// Destination of Mult can't be in memory.

				updatedInsts = append(updatedInsts, &MovInst{
					L: v.Dest,
					R: &RegisterOperand{
						Reg: "R11",
					},
				})
				updatedInsts = append(updatedInsts, &BinaryInst{
					Op:  v.Op,
					Src: v.Src,
					Dest: &RegisterOperand{
						Reg: "R11",
					},
				})
				updatedInsts = append(updatedInsts, &MovInst{
					L: &RegisterOperand{
						Reg: "R11",
					},
					R: v.Dest,
				})

				continue
			}
		case *CmpInst:
			_, ok1 := v.C.(*StackOperand)
			_, ok2 := v.V.(*StackOperand)
			if ok1 && ok2 {
				updatedInsts = append(updatedInsts, &MovInst{
					L: v.C,
					R: &RegisterOperand{
						Reg: "R10",
					},
				})
				updatedInsts = append(updatedInsts, &CmpInst{
					C: &RegisterOperand{
						Reg: "R10",
					},
					V: v.V,
				})

				continue
			}

			if _, ok := v.V.(*ImmOperand); ok {
				updatedInsts = append(updatedInsts, &MovInst{
					L: v.V,
					R: &RegisterOperand{
						Reg: "R11",
					},
				})
				updatedInsts = append(updatedInsts, &CmpInst{
					C: v.C,
					V: &RegisterOperand{
						Reg: "R11",
					},
				})

				continue
			}
		}

		updatedInsts = append(updatedInsts, inst)
	}

	return updatedInsts
}

func (f *fixer) replacePseudos(insts []Inst) ([]Inst, int32) {
	// TODO(andydunstall): Add support for AST walking instead of checking
	// each instruction.

	var lastOffset int32
	offsets := make(map[string]int32)

	replace := func(op Operand) Operand {
		pseudo, ok := op.(*PseudoOperand)
		if !ok {
			return op
		}

		if off, ok := offsets[pseudo.V]; ok {
			return &StackOperand{
				Offset: off,
			}
		} else {
			lastOffset -= 4
			offsets[pseudo.V] = lastOffset
			return &StackOperand{
				Offset: lastOffset,
			}
		}
	}

	var updatedInsts []Inst
	for _, inst := range insts {
		switch v := inst.(type) {
		case *MovInst:
			inst = &MovInst{
				L: replace(v.L),
				R: replace(v.R),
			}
		case *UnaryInst:
			inst = &UnaryInst{
				Op: v.Op,
				V:  replace(v.V),
			}
		case *BinaryInst:
			inst = &BinaryInst{
				Op:   v.Op,
				Src:  replace(v.Src),
				Dest: replace(v.Dest),
			}
		case *IdivInst:
			inst = &IdivInst{
				V: replace(v.V),
			}
		case *CmpInst:
			inst = &CmpInst{
				C: replace(v.C),
				V: replace(v.V),
			}
		case *SetCCInst:
			inst = &SetCCInst{
				C: v.C,
				V: replace(v.V),
			}
		case *PushInst:
			inst = &PushInst{
				V: replace(v.V),
			}
		}
		updatedInsts = append(updatedInsts, inst)
	}

	return updatedInsts, lastOffset
}

func roundUpToNextMultipleOf16(n int32) int32 {
	remainder := n % 16
	if remainder == 0 {
		return n
	}
	return n + (16 - remainder)
}
