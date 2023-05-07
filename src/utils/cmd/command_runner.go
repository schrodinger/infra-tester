package cmd

import (
	"bytes"
	"os/exec"
)

type CommandRunner interface {
	// Returns the path to the executable for the given command.
	LookPath(commandName string) (string, error)

	// Runs the given command with the given args and returns the result.
	RunCommand(command string, args ...string) CommandResult
}

type cmdRunner struct{}

// Creates a new CommandRunner that can be used to run commands.
func NewCmdRunner() CommandRunner {
	return cmdRunner{}
}

func (c cmdRunner) LookPath(commandName string) (string, error) {
	return exec.LookPath(commandName)
}

func (c cmdRunner) RunCommand(command string, args ...string) CommandResult {
	cmd := exec.Command(command, args...)

	var stdoutBytes, stderrBytes bytes.Buffer
	cmd.Stdout = &stdoutBytes
	cmd.Stderr = &stderrBytes

	err := cmd.Run()

	stdout := stdoutBytes.String()
	stderr := stderrBytes.String()

	return &commandResult{
		cmd:    cmd,
		err:    err,
		stdout: stdout,
		stderr: stderr,
	}
}
