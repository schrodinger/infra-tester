package main

import (
	"fmt"
	"os"
	"testing"

	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/mitchellh/mapstructure"
	"gopkg.in/yaml.v3"
	"schrodinger.com/infra-tester/assertions"
	"schrodinger.com/infra-tester/plugins"
	"schrodinger.com/infra-tester/utils/cmd"
)

func main() {
	testing.Main(
		nil,
		[]testing.InternalTest{
			{
				Name: "TestMain",
				F:    TestMain,
			},
		},
		nil, nil,
	)
}

func TestMain(t *testing.T) {
	terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{})
	(*terraformOptions).NoColor = true

	testPlan, err := getTests()
	if err != nil {
		t.Fatalf("ERROR: Failed to process all tests: %s", err)
	}

	// Build assertion context.
	assertionContext := buildAssertionContext(t)

	// Validate the tests.
	if err = validateTests(testPlan, assertionContext); err != nil {
		assertions.ErrorAndSkipf(t, "ERROR: Failure during test validation: %s", err)
	}

	// Run the tests.
	t.Run(testPlan.Name, func(t *testing.T) {
		_, err = terraform.InitE(t, terraformOptions)
		if err != nil {
			assertions.ErrorAndSkipf(t, "ERROR: Failure during terraform init: %s", err)
		}

		runTests(t, terraformOptions, testPlan, assertionContext)
	})
}

func buildAssertionContext(t *testing.T) *assertions.AssertionContext {
	// Build assertion context.
	assertionContext := assertions.AssertionContext{}

	// Setup plugins
	setupPlugins(t, &assertionContext)

	return &assertionContext
}

func setupPlugins(t *testing.T, assertionContext *assertions.AssertionContext) {
	// Check if plugins are supported in the current environment.
	cmdRunner := cmd.NewCmdRunner()
	if err := plugins.CanRunPlugins(cmdRunner); err == nil {
		t.Log("INFO: Plugin framework is installed.")
		pluginManager, err := plugins.NewPipPluginManager(cmdRunner)

		if err != nil {
			t.Fatalf("ERROR: Failed to create plugin manager: %s."+
				"Please file an issue with the logs.", err)
		}

		assertionContext.PluginManager = &pluginManager
		assertionContext.AvailablePlugins, err = pluginManager.ListPlugins()
		if err != nil {
			t.Fatalf("ERROR: Failed to list plugins: %s. "+
				"Please raise an issue with the logs", err)
		}
	} else {
		t.Logf("INFO: Can not use plugins in this environment: %s", err)
		assertionContext.AvailablePlugins = map[string]bool{}
		assertionContext.PluginManager = nil
	}
}

func runTests(
	t *testing.T,
	terraformOptions *terraform.Options,
	testPlan TestPlan,
	assertionContext *assertions.AssertionContext) {
	// Run destroy regardless of test results to clean up any left overs
	defer terraform.Destroy(t, terraformOptions)

	for _, test := range testPlan.Tests {
		t.Run(test.Name, func(t *testing.T) {
			skipApplyTests := false

			// Run all plan assertions first
			if test.PlanAssertions.Assertions == nil {
				t.Logf("No plan assertions for %s", test.Name)
			} else {
				runPlanAssertions(t, test, terraformOptions, assertionContext)

				// Skip all the apply assertions if any plan assertions failed. We may want to provide this as an
				// option in the future.
				if t.Failed() {
					t.Logf("Plan assertions failed for %s, skipping Apply assertions", test.Name)

					skipApplyTests = true
				}
			}

			// Run all apply assertions
			if test.ApplyAssertions.Assertions == nil {
				t.Logf("No apply assertions for %s", test.Name)
			} else {
				runApplyAssertions(t, test, terraformOptions, skipApplyTests, assertionContext)
			}
		})
	}

	t.Log("A final destroy will be called to cleanup any left over resources")
	// Set terraform vars to destroy vars if they are provided
	if testPlan.DestroyVars != nil {
		t.Log("Using destroy_vars for final destroy")
		terraformOptions.Vars = testPlan.DestroyVars
	}
}

func runPlanAssertions(
	t *testing.T,
	test Test,
	terraformOptions *terraform.Options,
	assertionContext *assertions.AssertionContext) {
	t.Run("Plan", func(t *testing.T) {
		if test.WithCleanState {
			t.Logf("INFO: with_clean_state enabled - running destroy before plan for %s", test.Name)
			_, err := terraform.DestroyE(t, terraformOptions)
			if err != nil {
				assertions.ErrorAndSkipf(t, "ERROR: Failure during terraform destroy: %s", err)
			}
		}

		if test.Vars != nil {
			terraformOptions.Vars = test.Vars
		}

		stdOutErr, err := terraform.PlanE(t, terraformOptions)
		planMetadata := assertions.PlanMetadata{CmdOut: stdOutErr, Err: err}

		for _, assertion := range test.PlanAssertions.Assertions {
			subTestName := assertion.Type
			if assertion.Name != "" {
				subTestName = assertion.Name
			}

			t.Run(subTestName, func(t *testing.T) {
				assertions.RunAssertion(
					t,
					terraformOptions,
					assertion,
					"plan",
					planMetadata,
					assertionContext)
			})
		}
	})
}

func runApplyAssertions(
	t *testing.T,
	test Test,
	terraformOptions *terraform.Options,
	skipTests bool,
	assertionContext *assertions.AssertionContext) {
	t.Run("Apply", func(t *testing.T) {
		if skipTests {
			t.SkipNow()
		}

		if test.WithCleanState {
			t.Logf("INFO: with_clean_state enabled - running destroy before apply for %s", test.Name)
			_, err := terraform.DestroyE(t, terraformOptions)
			if err != nil {
				assertions.ErrorAndSkipf(t, "ERROR: Failure during terraform destroy: %s", err)
			}
		}

		if test.Vars != nil {
			terraformOptions.Vars = test.Vars
		}

		var stdOutErr string
		var err error
		if test.ApplyAssertions.EnsureIdempotent {
			stdOutErr, err = terraform.ApplyAndIdempotentE(t, terraformOptions)
		} else {
			stdOutErr, err = terraform.ApplyE(t, terraformOptions)
		}
		applyMetadata := assertions.ApplyMetadata{CmdOut: stdOutErr, Err: err}

		for _, assertion := range test.ApplyAssertions.Assertions {
			subTestName := assertion.Type
			if assertion.Name != "" {
				subTestName = assertion.Name
			}

			t.Run(subTestName, func(t *testing.T) {
				assertions.RunAssertion(
					t,
					terraformOptions,
					assertion,
					"apply",
					applyMetadata,
					assertionContext)
			})
		}
	})
}

func getTests() (TestPlan, error) {
	yamlConfig, err := os.ReadFile(".infra-tester-config.yaml")
	if err != nil {
		return TestPlan{}, fmt.Errorf("failed to read yaml config: %v", err)
	}

	var mapStruct map[string]interface{}
	err = yaml.Unmarshal(yamlConfig, &mapStruct)
	if err != nil {
		return TestPlan{}, fmt.Errorf("failed to unmarshal yaml config: %v", err)
	}

	var config Config
	err = mapstructure.Decode(mapStruct, &config)
	if err != nil {
		return TestPlan{}, fmt.Errorf("failed to decode map structure: %v", err)
	}

	return config.TestPlan, nil
}
