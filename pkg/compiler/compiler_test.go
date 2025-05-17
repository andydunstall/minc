package compiler_test

import (
	"os"
	"testing"

	"github.com/andydunstall/minc/pkg/arch/x86"
	"github.com/andydunstall/minc/pkg/assembly"
	"github.com/andydunstall/minc/pkg/ast"
	"github.com/andydunstall/minc/pkg/ir"
	"github.com/andydunstall/minc/pkg/token"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCompileX86(t *testing.T) {
	tests := []struct {
		Name string
		Path string
		Want string
	}{
		{
			Name: "return",
			Path: "return.c",
			Want: `	.global main
main:
	pushq %rbp
	movq %rsp, %rbp
	subq $0, %rsp
	movl $10, %eax
	movq %rbp, %rsp
	popq %rbp
	ret
	.section .note.GNU-stack,"",@progbits
`,
		},
		{
			Name: "unary",
			Path: "unary.c",
			Want: `	.global main
main:
	pushq %rbp
	movq %rsp, %rbp
	subq $16, %rsp
	movl $2, -4(%rbp)
	negl -4(%rbp)
	movl -4(%rbp), %r10d
	movl %r10d, -8(%rbp)
	notl -8(%rbp)
	movl -8(%rbp), %eax
	movq %rbp, %rsp
	popq %rbp
	ret
	.section .note.GNU-stack,"",@progbits
`,
		},
		{
			Name: "logical",
			Path: "logical.c",
			Want: `	.global main
main:
	pushq %rbp
	movq %rsp, %rbp
	subq $16, %rsp
	movl $1, %r11d
	cmpl $0, %r11d
	je .Land_false.2
	movl $0, %r11d
	cmpl $0, %r11d
	je .Land_false.2
	movl $1, -4(%rbp)
	jmp .Land_end.3
.Land_false.2:
	movl $0, -4(%rbp)
.Land_end.3:
	cmpl $0, -4(%rbp)
	jne .Lor_true.0
	movl $3, %r11d
	cmpl $0, %r11d
	jne .Lor_true.0
	movl $0, -8(%rbp)
	jmp .Lor_end.1
.Lor_true.0:
	movl $0, -8(%rbp)
.Lor_end.1:
	cmpl $3, -8(%rbp)
	movl $0, -12(%rbp)
	setne -12(%rbp)
	cmpl $0, -12(%rbp)
	movl $0, -16(%rbp)
	sete -16(%rbp)
	movl -16(%rbp), %eax
	movq %rbp, %rsp
	popq %rbp
	ret
	.section .note.GNU-stack,"",@progbits
`,
		},
		{
			Name: "variables",
			Path: "variables.c",
			Want: `	.global main
main:
	pushq %rbp
	movq %rsp, %rbp
	subq $32, %rsp
	movl $1, -4(%rbp)
	movl $2, -8(%rbp)
	movl $3, -12(%rbp)
	movl $4, -8(%rbp)
	movl -4(%rbp), %r10d
	movl %r10d, -16(%rbp)
	movl -8(%rbp), %r10d
	addl %r10d, -16(%rbp)
	movl -16(%rbp), %r10d
	movl %r10d, -20(%rbp)
	movl -12(%rbp), %r10d
	addl %r10d, -20(%rbp)
	movl -20(%rbp), %eax
	movq %rbp, %rsp
	popq %rbp
	ret
	.section .note.GNU-stack,"",@progbits
`,
		},
		{
			Name: "conditional",
			Path: "conditional.c",
			Want: `	.global main
main:
	pushq %rbp
	movq %rsp, %rbp
	subq $16, %rsp
	movl $3, -4(%rbp)
	cmpl $2, -4(%rbp)
	movl $0, -8(%rbp)
	setl -8(%rbp)
	cmpl $0, -8(%rbp)
	je .Lelse.0
	movl $5, -12(%rbp)
	movl -12(%rbp), %r10d
	movl %r10d, -4(%rbp)
	jmp .Lif_end.1
.Lelse.0:
	cmpl $2, -4(%rbp)
	movl $0, -16(%rbp)
	sete -16(%rbp)
	cmpl $0, -16(%rbp)
	je .Lelse.3
	movl $6, -4(%rbp)
	jmp .Lif_end.4
.Lelse.3:
	movl $7, -4(%rbp)
.Lif_end.4:
.Lif_end.1:
	movl -4(%rbp), %eax
	movq %rbp, %rsp
	popq %rbp
	ret
	.section .note.GNU-stack,"",@progbits
`,
		},
		{
			Name: "loops",
			Path: "loops.c",
			Want: `	.global main
main:
	pushq %rbp
	movq %rsp, %rbp
	subq $32, %rsp
	movl $1, -4(%rbp)
.Lcontinue.loop.1:
	cmpl $5, -4(%rbp)
	movl $0, -8(%rbp)
	setl -8(%rbp)
	cmpl $0, -8(%rbp)
	je .Lbreak.loop.1
	movl -4(%rbp), %r10d
	movl %r10d, -12(%rbp)
	addl $1, -12(%rbp)
	movl -12(%rbp), %r10d
	movl %r10d, -4(%rbp)
	cmpl $12, -4(%rbp)
	movl $0, -16(%rbp)
	setg -16(%rbp)
	cmpl $0, -16(%rbp)
	je .Lelse.2
	jmp .Lbreak.loop.1
	jmp .Lif_end.3
.Lelse.2:
	cmpl $7, -4(%rbp)
	movl $0, -20(%rbp)
	setg -20(%rbp)
	cmpl $0, -20(%rbp)
	je .Lelse.5
	jmp .Lcontinue.loop.1
	jmp .Lif_end.6
.Lelse.5:
.Lif_end.6:
.Lif_end.3:
	jmp .Lcontinue.loop.1
.Lbreak.loop.1:
	movl -4(%rbp), %eax
	movq %rbp, %rsp
	popq %rbp
	ret
	.section .note.GNU-stack,"",@progbits
`,
		},
		{
			Name: "functions",
			Path: "functions.c",
			Want: `	.global two
two:
	pushq %rbp
	movq %rsp, %rbp
	subq $0, %rsp
	movl $2, %eax
	movq %rbp, %rsp
	popq %rbp
	ret
	.global addFive
addFive:
	pushq %rbp
	movq %rsp, %rbp
	subq $16, %rsp
	movl %edi, -4(%rbp)
	movl -4(%rbp), %r10d
	movl %r10d, -8(%rbp)
	addl $5, -8(%rbp)
	movl -8(%rbp), %eax
	movq %rbp, %rsp
	popq %rbp
	ret
	.global addTen
addTen:
	pushq %rbp
	movq %rsp, %rbp
	subq $16, %rsp
	movl %edi, -4(%rbp)
	subq $8, %rsp
	movl -4(%rbp), %edi
	call addFive
	addq $16, %rsp
	movl %eax, -8(%rbp)
	subq $8, %rsp
	movl -8(%rbp), %edi
	call addFive
	addq $16, %rsp
	movl %eax, -12(%rbp)
	movl -12(%rbp), %eax
	movq %rbp, %rsp
	popq %rbp
	ret
	.global main
main:
	pushq %rbp
	movq %rsp, %rbp
	subq $32, %rsp
	call two
	addq $0, %rsp
	movl %eax, -4(%rbp)
	call two
	addq $0, %rsp
	movl %eax, -8(%rbp)
	movl -8(%rbp), %r10d
	movl %r10d, -12(%rbp)
	addl $1, -12(%rbp)
	subq $8, %rsp
	movl -12(%rbp), %edi
	call addTen
	addq $16, %rsp
	movl %eax, -16(%rbp)
	movl -4(%rbp), %r10d
	movl %r10d, -20(%rbp)
	movl -16(%rbp), %r10d
	addl %r10d, -20(%rbp)
	subq $8, %rsp
	movl $5, %edi
	call addTen
	addq $16, %rsp
	movl %eax, -24(%rbp)
	movl -20(%rbp), %r10d
	movl %r10d, -28(%rbp)
	movl -24(%rbp), %r10d
	addl %r10d, -28(%rbp)
	movl -28(%rbp), %eax
	movq %rbp, %rsp
	popq %rbp
	ret
	.section .note.GNU-stack,"",@progbits
`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			assert.Equal(t, tt.Want, compileX86("../../testdata/"+tt.Path, t))
		})
	}
}

func compileX86(path string, t *testing.T) string {
	src, err := os.ReadFile(path)
	require.NoError(t, err)

	scanner := token.NewScanner(src)

	fileAST, err := ast.Parse(scanner, false)
	require.NoError(t, err)

	validatedAST, err := ast.Validate(fileAST, false)
	require.NoError(t, err)

	irFile, err := ir.Parse(validatedAST, false)
	require.NoError(t, err)

	assem, err := assembly.Parse(irFile, false)
	require.NoError(t, err)
	assem = assembly.Fix(assem.(*assembly.File), false)

	return x86.Emit(assem)
}
