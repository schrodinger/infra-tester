package test

import (
	"fmt"
	"testing"

	"github.com/gruntwork-io/terratest/modules/terraform"
)

type AssertionImplementation struct {
	ValidateFunction func(Assertion) error
	RunFunction      func(t *testing.T, terraformOptions *terraform.Options, assertion Assertion, stepMetadata interface{})
}

func getAssertionImplementation(assertionType string, step string) (AssertionImplementation, error) {
	var assertionImplementation AssertionImplementation
	var ok bool
	if step == "plan" {
		if assertionImplementation, ok = ValidPlanAssertions[assertionType]; !ok {
			if _, ok = ValidApplyAssertions[assertionType]; ok {
				return AssertionImplementation{}, fmt.Errorf("'%s' is only valid for 'apply' tests", assertionType)
			}

			return AssertionImplementation{}, fmt.Errorf("assertion type '%s' is invalid", assertionType)
		}
	} else if step == "apply" {
		if assertionImplementation, ok = ValidApplyAssertions[assertionType]; !ok {
			if _, ok = ValidPlanAssertions[assertionType]; ok {
				return AssertionImplementation{}, fmt.Errorf("'%s' is only valid for 'plan' tests", assertionType)
			}

			return AssertionImplementation{}, fmt.Errorf("assertion type '%s' is invalid", assertionType)
		}
	} else {
		return AssertionImplementation{}, fmt.Errorf("step '%s' is invalid", step)
	}

	return assertionImplementation, nil
}

func validateAssertion(assertion Assertion, step string) error {
	AssertionImplementation, err := getAssertionImplementation(assertion.Type, step)
	if err != nil {
		return err
	}

	validateFunction := AssertionImplementation.ValidateFunction
	return validateFunction(assertion)
}

func runAssertion(t *testing.T, terraformOptions *terraform.Options, assertion Assertion, step string, stepMetadata interface{}) {
	assertionType := assertion.Type
	var assertionImplementation AssertionImplementation

	assertionImplementation, err := getAssertionImplementation(assertionType, step)
	if err != nil {
		// This shouldn't happen as we are validating the tests before running them
		t.Fatalf("ERROR: Failure while running assertion: %s.\n", err)
	}

	runFunction := assertionImplementation.RunFunction
	runFunction(t, terraformOptions, assertion, stepMetadata)
}
