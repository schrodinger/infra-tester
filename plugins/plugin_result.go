package plugins

import (
	"errors"
	"fmt"
	"testing"

	"github.com/schrodinger/infra-tester/utils"
	"github.com/schrodinger/infra-tester/utils/cmd"
)

const (
	SUCCESS       = 0
	ERROR         = 1
	INVALID_INPUT = 2
)

type PluginResult interface {
	CheckErrors() error
	Logf(t *testing.T)
}

type pipPluginRunnerResult struct {
	Error           bool              `json:"error"`
	Message         *string           `json:"message"`
	cmdRunnerResult cmd.CommandResult `json:"-"`
}

func (p *pipPluginRunnerResult) CheckErrors() error {
	exitCode := p.cmdRunnerResult.ExitCode()
	executedCommand := p.cmdRunnerResult.ExecutedCommand()
	stdout := p.cmdRunnerResult.Stdout()
	stderr := p.cmdRunnerResult.Stderr()

	// If the command returned an error, return the error.
	// This can happen if there are any issues with the plugin
	// or the plugin framework.
	if exitCode == ERROR || exitCode == INVALID_INPUT {
		return fmt.Errorf("error while executing plugin (%s): %s"+
			"Please raise an issue with the logs",
			executedCommand,
			stderr+stdout)
	}

	// If the command returned an invalid input error, it most likely
	// means there's something wrong with the inputs infra-tester
	// is passing to the plugin framework.
	if exitCode == INVALID_INPUT {
		return fmt.Errorf("error while executing plugin (%s): %s "+
			"Please raise an issue with the logs",
			executedCommand,
			stderr+stdout)
	}

	if p.cmdRunnerResult.Error() == nil && exitCode == SUCCESS {
		// At this point we know the plugin ran successfully. The plugin
		// might still have an error but the plugin framework successfully
		// executed the plugin.

		// Populate the message and error fields from stdout.
		err := utils.FromJSON(p.cmdRunnerResult.Stdout(), &p)
		if err != nil {
			return fmt.Errorf("error while unmarshalling plugin result: %s", err)
		}

		// If the plugin returned an error, return the message.
		if p.Error {
			return errors.New(*p.Message)
		}

		return nil
	}

	return fmt.Errorf("received unknown exit code '%d' while executing plugin (%s): %s",
		exitCode,
		executedCommand,
		stderr+stdout)
}

func (p *pipPluginRunnerResult) Logf(t *testing.T) {
	stdout := p.cmdRunnerResult.Stdout()
	stderr := p.cmdRunnerResult.Stderr()

	t.Log("INFO: Plugin stdout")
	t.Log(stdout)
	t.Log("INFO: Plugin stderr")
	t.Log(stderr)
}
