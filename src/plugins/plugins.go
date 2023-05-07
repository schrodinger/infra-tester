package plugins

import (
	"fmt"

	"schrodinger.com/infra-tester/utils/cmd"
)

const (
	PLUGIN_FRAMEWORK_PACKAGE  = "infra-tester-plugins"
	PLUGIN_RUNNER_EXECUTABLE  = "infra-tester-run-plugin"
	PLUGIN_MANAGER_EXECUTABLE = "infra-tester-plugin-manager"
)

func CanRunPlugins(commandRunner cmd.CommandRunner) error {
	pipPath, err := commandRunner.LookPath("pip3")
	// pip3 is not installed.
	if err != nil {
		return fmt.Errorf(
			"could not find %s in PATH, please install '%s' for plugin support: %s",
			PLUGIN_RUNNER_EXECUTABLE,
			PLUGIN_FRAMEWORK_PACKAGE,
			err)
	}

	// check if plugin framework package is installed.
	res := commandRunner.RunCommand(pipPath, "show", PLUGIN_FRAMEWORK_PACKAGE)

	if err := res.Error(); err != nil {
		return fmt.Errorf("could not find plugin framework: (%s) %s", res.ExecutedCommand(), err)
	}

	return nil
}
