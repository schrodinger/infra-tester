package plugins

import (
	"fmt"
	"strings"

	"github.com/schrodinger/infra-tester/utils/cmd"
)

type PluginManager interface {
	// Retrieves a list of available plugins.
	ListPlugins() (map[string]bool, error)
	// Retrieves a PluginRunner for the given plugin. PluginRunner
	// can be used to validate inputs, run the plugin and cleanup the
	// resources created by the plugin.
	GetPluginRunner(pluginName string) (PluginRunner, error)
}

type pipPluginManager struct {
	pipPath          string
	availablePlugins map[string]bool
	commandRunner    cmd.CommandRunner
	pluginRunner     PluginRunner
}

// NewPipPluginManager creates a new PluginManager that uses pip3 to manage
// plugins. As of now infra-tester-plugin-manager is used to just list the
// available plugins.
func NewPipPluginManager(commandRunner cmd.CommandRunner) (PluginManager, error) {
	// Get the path to pip3 executable.
	pipPath, err := commandRunner.LookPath("pip3")
	if err != nil {
		// This shouldn't happen because we check if pip3 is installed
		// before running any plugins. This is just a sanity check.

		return nil, fmt.Errorf("error while creating plugin manager: %s", err)
	}

	pluginManager := pipPluginManager{
		pipPath:       pipPath,
		commandRunner: commandRunner,
	}

	if err != nil {
		return nil, fmt.Errorf("error while creating plugin manager: %s", err)
	}

	return &pluginManager, nil
}

func (p *pipPluginManager) ListPlugins() (map[string]bool, error) {
	if p.availablePlugins != nil {
		return p.availablePlugins, nil
	}

	res := p.commandRunner.RunCommand(PLUGIN_MANAGER_EXECUTABLE, "--list")
	err := res.Error()
	if err != nil {
		return nil, fmt.Errorf("error while listing available plugins: %s", err)
	}

	p.availablePlugins = parseAvailablePlugins(res.Stdout())

	return p.availablePlugins, nil
}

func (p *pipPluginManager) GetPluginRunner(pluginName string) (PluginRunner, error) {
	if p.pluginRunner != nil {
		return p.pluginRunner, nil
	}

	p.pluginRunner = NewPluginRunner(p.commandRunner, pluginName)

	return p.pluginRunner, nil
}

// Parses the output of infra-tester-plugin-manager --list and returns a list
// of available plugins.
func parseAvailablePlugins(out string) map[string]bool {
	availablePlugins := map[string]bool{}
	lines := strings.Split(strings.ReplaceAll(out, "\r\n", "\n"), "\n")
	for _, line := range lines {
		// infra-tester-plugin-manager only prints the names of
		// available plugins, so we just need to split the output
		// by lines, but just to be safe, we ignore empty lines
		// and trim the spaces.
		if len(line) > 0 {
			availablePlugins[strings.Trim(line, " ")] = true
		}
	}

	return availablePlugins
}
