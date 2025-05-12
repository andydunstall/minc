package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func newCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:          "minc [command] (flags)",
		SilenceUsage: true,
		Long: `Minc is a mini C compiler.

The compiler currently only supports a single C file, with no includes.

Compile a C file with:

  $ minc compile ./program.c

Which will output the x86 assembly to ./main.s (or specify the output file
with -o/--output).

`,
		CompletionOptions: cobra.CompletionOptions{
			DisableDefaultCmd: true,
		},
	}

	cmd.AddCommand(
		newCompileCommand(),
	)

	return cmd
}

func exitError(err error) {
	fmt.Fprintf(os.Stderr, "%s\n", err.Error())
	os.Exit(1)
}

func init() {
	cobra.EnableCommandSorting = false
}
