package assertions

import (
	"fmt"
	"testing"

	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/schrodinger/infra-tester/plugins"
	"github.com/stretchr/testify/assert"
)

func GetCustomAssertionImplementation(
	assertionType string,
	pluginManager *plugins.PluginManager) (AssertionImplementation, error) {

	// Get the plugin runner for the assertion type
	pluginRunner, err := (*pluginManager).GetPluginRunner(assertionType)
	if err != nil {
		return AssertionImplementation{}, fmt.Errorf("failed to get plugin runner for %s: %s", assertionType, err)
	}

	return AssertionImplementation{
		ValidateFunction: func(assertion Assertion) error {
			return pluginRunner.ValidateInputs(assertion)
		},
		RunFunction: func(t *testing.T, terraformOptions *terraform.Options, assertion Assertion, stepMetadata interface{}) {
			t.Log("INFO: Running custom assertion")
			terraformState, err := terraform.ShowE(t, terraformOptions)
			if err != nil {
				ErrorAndSkipf(t, "ERROR: Failed to get terraform state: %s", err)
			}

			err = pluginRunner.Run(t, assertion, &terraformState)
			assert.Nilf(t, err, "assertion '%s' failed: %s", assertion.Name, err)
		},
	}, nil
}
