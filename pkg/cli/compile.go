package cli

import (
	"fmt"
	"os"
	"slices"
	"strings"

	"github.com/andydunstall/minc/pkg/arch/x86"
	"github.com/andydunstall/minc/pkg/assembly"
	"github.com/andydunstall/minc/pkg/ast"
	"github.com/andydunstall/minc/pkg/compiler"
	"github.com/andydunstall/minc/pkg/ir"
	"github.com/andydunstall/minc/pkg/print"
	"github.com/andydunstall/minc/pkg/token"
	"github.com/spf13/cobra"
)

func newCompileCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "compile path [flags]",
		Short: "compile a C program",
		Long:  `...`,
	}

	var stages []string
	for _, s := range compiler.Stages {
		stages = append(stages, string(s))
	}

	var outputPath string
	cmd.Flags().StringVarP(
		&outputPath,
		"output",
		"o",
		"./dump.s",
		"output path",
	)

	var stage string
	cmd.Flags().StringVarP(
		&stage,
		"stage",
		"s",
		"",
		fmt.Sprintf("compiler stage (%s)", strings.Join(stages, ", ")),
	)

	var debug bool
	cmd.Flags().BoolVarP(
		&debug,
		"debug",
		"d",
		false,
		"whether to output debug logs and traces",
	)

	cmd.Run = func(_ *cobra.Command, args []string) {
		if len(args) == 0 {
			exitError(fmt.Errorf("compile: missing path"))
		}
		if len(args) > 1 {
			exitError(fmt.Errorf("compile: only one path is supported"))
		}

		if err := runCompile(args[0], outputPath, compiler.Stage(stage), debug); err != nil {
			exitError(fmt.Errorf("compile: %w", err))
		}
	}

	return cmd
}

func runCompile(path string, outputPath string, stage compiler.Stage, debug bool) error {
	if stage != "" && !slices.Contains(compiler.Stages, stage) {
		return fmt.Errorf("unsupported stage: %s", stage)
	}

	src, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("read: %s: %w", path, err)
	}

	if stage == compiler.StageTokenize || debug {
		scanner := token.NewScanner(src)

		if debug {
			fmt.Println("tokens:")
		}

		line := 1
		tok, lit := scanner.Scan()
		for tok != token.EOF {
			if lit == "" || lit == tok.String() {
				fmt.Printf("%6d %s\n", line, tok)
			} else {
				fmt.Printf("%6d %s (%s)\n", line, tok, lit)
			}
			tok, lit = scanner.Scan()
			line++
		}
		fmt.Println("")
		if stage == compiler.StageTokenize {
			return nil
		}
	}

	scanner := token.NewScanner(src)

	if debug {
		fmt.Println("parse:")
	}

	fileAST, err := ast.Parse(scanner, debug)
	if err != nil {
		return fmt.Errorf("parse ast: %w", err)
	}

	if debug {
		fmt.Println("")
	}

	if stage == compiler.StageParse || debug {
		if debug {
			fmt.Println("ast (unvalidated):")
		}

		print.Print(fileAST)
		if stage == compiler.StageParse {
			return nil
		}

		if debug {
			fmt.Println("")
		}
	}

	validatedAST, err := ast.Validate(fileAST, debug)
	if err != nil {
		return fmt.Errorf("validate ast: %w", err)
	}

	if stage == compiler.StageValidate || debug {
		if debug {
			fmt.Println("ast (validated):")
		}

		print.Print(validatedAST)
		if stage == compiler.StageValidate {
			return nil
		}

		if debug {
			fmt.Println("")
		}
	}

	irFile, err := ir.Parse(validatedAST, debug)
	if err != nil {
		return fmt.Errorf("parse ir: %w", err)
	}

	if stage == compiler.StageIR || debug {
		if debug {
			fmt.Println("ir:")
		}

		print.Print(irFile)
		if stage == compiler.StageIR {
			return nil
		}

		if debug {
			fmt.Println("")
		}
	}

	assem, err := assembly.Parse(irFile, debug)
	if err != nil {
		return fmt.Errorf("parse assembly: %w", err)
	}

	if stage == compiler.StageAssemble || debug {
		if debug {
			fmt.Println("assembly (unfixed):")
		}

		print.Print(assem)
		if stage == compiler.StageAssemble {
			return nil
		}

		if debug {
			fmt.Println("")
		}
	}

	assemFixed := assembly.Fix(assem.(*assembly.File), debug)
	if stage == compiler.StageAssemble || debug {
		if debug {
			fmt.Println("assembly (fixed):")
		}

		print.Print(assemFixed)
		if stage == compiler.StageAssemble {
			return nil
		}

		if debug {
			fmt.Println("")
		}
	}

	x86Assem := x86.Emit(assemFixed)

	if debug {
		fmt.Println("x86:")
		p := print.NewPrinter(os.Stdout)
		p.Write([]byte(x86Assem))
		fmt.Println("")
	}

	if err := os.WriteFile(outputPath, []byte(x86Assem), 0o666); err != nil {
		return fmt.Errorf("write: %w", err)
	}

	return nil
}
