package compiler

type Stage string

const (
	// StageTokenize parses the program into a token stream.
	StageTokenize Stage = "tokenize"

	// StageParse parses the program into an AST.
	StageParse Stage = "parse"

	// StageValidate parses the program into an AST and performs
	// semantic analysis.
	StageValidate Stage = "validate"

	// StageIR parses the program into an IR.
	StageIR Stage = "ir"

	// StageAssemble parses the program into assembly.
	StageAssemble Stage = "assemble"
)

var Stages = []Stage{
	StageTokenize,
	StageParse,
	StageValidate,
	StageIR,
	StageAssemble,
}
