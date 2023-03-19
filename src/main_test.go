package test

import (
	"fmt"
	"os"
	"testing"

	"github.com/davecgh/go-spew/spew"
	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/mitchellh/mapstructure"
	"gopkg.in/yaml.v3"
	"schrodinger.com/infra-tester/assertions"
)

func TestMain(t *testing.T) {
	terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{})
	(*terraformOptions).NoColor = true

	testPlan, err := getTests()
	if err != nil {
		t.Fatalf("ERROR: Failed to process all tests: %s", err)
	}

	// // fmt.Printf("DEBUG: ")
	spew.Dump(testPlan)

	// Validate the tests
	if err = validateTests(testPlan.Tests); err != nil {
		assertions.ErrorAndSkipf(t, "ERROR: Failure during test validation: %s", err)
	}

	t.Run(testPlan.Name, func(t *testing.T) {
		_, err = terraform.InitE(t, terraformOptions)
		if err != nil {
			assertions.ErrorAndSkipf(t, "ERROR: Failure during terraform init: %s", err)
		}

		runTests(t, terraformOptions, testPlan)
	})
}

func runTests(t *testing.T, terraformOptions *terraform.Options, testPlan TestPlan) {
	// Run destroy regardless of test results to clean up any left overs
	defer terraform.Destroy(t, terraformOptions)

	for _, test := range testPlan.Tests {
		t.Run(test.Name, func(t *testing.T) {
			skipApplyTests := false

			// Run all plan assertions first
			if test.PlanAssertions.Assertions == nil {
				t.Logf("No plan assertions for %s", test.Name)
			} else {
				runPlanAssertions(t, test, terraformOptions)

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
				runApplyAssertions(t, test, terraformOptions, skipApplyTests)
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

func runPlanAssertions(t *testing.T, test Test, terraformOptions *terraform.Options) {
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
				assertions.RunAssertion(t, terraformOptions, assertion, "plan", planMetadata)
			})
		}
	})
}

func runApplyAssertions(t *testing.T, test Test, terraformOptions *terraform.Options, skipTests bool) {
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
				assertions.RunAssertion(t, terraformOptions, assertion, "apply", applyMetadata)
			})
		}
	})
}

func getTests() (TestPlan, error) {
	yamlConfig, err := os.ReadFile("config.yaml")
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
