package plugins

import (
	"fmt"
	"testing"

	"github.com/schrodinger/infra-tester/utils"
	"github.com/schrodinger/infra-tester/utils/cmd"
)

const (
	ACTION_VALIDATE_INPUTS = "validate_inputs"
	ACTION_RUN_ASSERTION   = "run_assertion"
	ACTION_CLEANUP         = "cleanup"
)

type PluginRunner interface {
	// Retrieves the name of the plugin.
	PluginName() string

	// Validates the inputs for the plugin.
	ValidateInputs(inputs utils.GenericMappable) error

	// Runs the plugin.
	Run(t *testing.T,
		inputs utils.GenericMappable,
		state *string) error

	// Executes cleanup for the plugin. This function
	// should ideally be deferred in the same function
	// where the plugin is run.
	Cleanup(t *testing.T,
		inputs utils.GenericMappable,
		state *string) error
}

type pipPluginRunner struct {
	pluginName string
	cmdRunner  cmd.CommandRunner
}

// Creates a new PluginRunner for the given plugin. PluginRunner can be used
// to validate inputs, run or execute cleanup for a given plugin.
func NewPluginRunner(commandRunner cmd.CommandRunner, pluginName string) PluginRunner {
	return &pipPluginRunner{
		pluginName: pluginName,
		cmdRunner:  commandRunner,
	}
}

func (p *pipPluginRunner) PluginName() string {
	return p.pluginName
}

func (p *pipPluginRunner) ValidateInputs(
	inputs utils.GenericMappable) error {
	res, err := p.execute(PLUGIN_RUNNER_EXECUTABLE, ACTION_VALIDATE_INPUTS, inputs, nil)
	if err != nil {
		return err
	}

	return res.CheckErrors()
}

func (p *pipPluginRunner) Run(
	t *testing.T,
	inputs utils.GenericMappable,
	state *string) error {
	res, err := p.execute(PLUGIN_RUNNER_EXECUTABLE, ACTION_RUN_ASSERTION, inputs, state)
	if err != nil {
		return err
	}

	res.Logf(t)

	return res.CheckErrors()
}

func (p *pipPluginRunner) Cleanup(
	t *testing.T,
	inputs utils.GenericMappable,
	state *string) error {
	res, err := p.execute(PLUGIN_RUNNER_EXECUTABLE, ACTION_CLEANUP, inputs, state)
	if err != nil {
		return err
	}

	return res.CheckErrors()
}

// Executes the plugin with the given command and args. The error returned
// by this function does not indicate if the command ran successfully or not.
// The error returned by this function indicates if there's utils.GenericMappable error with the
// inputs to the plugin.
func (p *pipPluginRunner) execute(
	command string,
	action string,
	inputs utils.GenericMappable,
	state *string) (PluginResult, error) {

	args, err := p.buildArgs(action, inputs, state)
	if err != nil {
		return nil, fmt.Errorf("error while building args for %s: %s ."+
			"Please raise an issue with the logs", p.pluginName, err)
	}

	return &pipPluginRunnerResult{
		cmdRunnerResult: p.cmdRunner.RunCommand(command, args...),
	}, nil
}

func (p *pipPluginRunner) buildArgs(
	action string,
	inputs utils.GenericMappable,
	state *string) ([]string, error) {

	args := []string{
		"--name", p.pluginName,
		"--action", action,
	}

	if inputs != nil {
		jsonInputs, err := utils.ToJSON(inputs)

		if err != nil {
			return nil, fmt.Errorf("error while converting inputs %#v to json: %s", inputs, err)
		}

		args = append(args, "--inputs", jsonInputs)
	}

	// state is optional and already a JSON string.
	if state != nil {
		args = append(args, "--state", *state)
	}

	return args, nil
}
