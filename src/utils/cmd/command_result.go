package cmd

import "os/exec"

type CommandResult interface {
	// Returns the exit code of the command.
	ExitCode() int

	// Returns the error returned by the command if any.
	Error() error

	// Returns the stdout of the command.
	Stdout() string

	// Returns the stderr of the command.
	Stderr() string

	// Returns the command that was executed. This can be helpful for debugging.
	ExecutedCommand() string
}

type commandResult struct {
	cmd    *exec.Cmd
	err    error
	stdout string
	stderr string
}

func (c *commandResult) ExitCode() int {
	return c.cmd.ProcessState.ExitCode()
}

func (c *commandResult) Error() error {
	return c.err
}

func (c *commandResult) Stdout() string {
	return c.stdout
}

func (c *commandResult) Stderr() string {
	return c.stderr
}

func (c *commandResult) ExecutedCommand() string {
	return c.cmd.String()
}
