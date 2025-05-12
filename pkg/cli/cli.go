package cli

func Start() error {
	cmd := newCommand()
	return cmd.Execute()
}
