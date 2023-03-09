package test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/mitchellh/mapstructure"
	"github.com/stretchr/testify/assert"
)

type PlanMetadata struct {
	CmdOut string
	Err    error
}

var ValidPlanAssertions = map[string]AssertionImplementation{
	"PlanSucceeds": {
		ValidateFunction: func(a Assertion) error { return nil },
		RunFunction:      AssertPlanSucceeds,
	},
	"PlanFailsWithError": {
		ValidateFunction: validatePlanFailsWithErrorAssertion,
		RunFunction:      AssertPlanFailsWithError,
	},
}

// ------------------------------------------------------------------------------------------------------------------------------

func AssertPlanSucceeds(t *testing.T, terraformOptions *terraform.Options, assertion Assertion, stepMetadata interface{}) {
	// cast stepMetadata to PlanMetadata
	var planMetadata PlanMetadata
	var ok bool
	if planMetadata, ok = stepMetadata.(PlanMetadata); !ok {
		t.Fatal("stepMetadata is not of type PlanMetadata")
	}

	if planMetadata.Err != nil {
		t.Fatal("Terraform plan is expected to succeed.")
	}
}

// ------------------------------------------------------------------------------------------------------------------------------

type planFailsWithErrorMetadata struct {
	ErrorMessageContains string `mapstructure:"error_message_contains"`
}

func validatePlanFailsWithErrorAssertion(assertion Assertion) error {
	var planFailsWithErrorMetadata planFailsWithErrorMetadata

	err := mapstructure.Decode(assertion.Metadata, &planFailsWithErrorMetadata)
	if err != nil {
		return fmt.Errorf("error decoding assertion metadata: %s", err)
	}

	if planFailsWithErrorMetadata.ErrorMessageContains == "" {
		return fmt.Errorf("error_message_contains is either not defined or is empty")
	}

	return nil
}

func AssertPlanFailsWithError(t *testing.T, terraformOptions *terraform.Options, assertion Assertion, stepMetadata interface{}) {
	var planFailsWithErrorMetadata planFailsWithErrorMetadata
	err := mapstructure.Decode(assertion.Metadata, &planFailsWithErrorMetadata)
	if err != nil {
		t.FailNow()
	}

	// cast stepMetadata to PlanMetadata
	var planMetadata PlanMetadata
	var ok bool
	if planMetadata, ok = stepMetadata.(PlanMetadata); !ok {
		t.Fatal("stepMetadata is not of type PlanMetadata")
	}

	if planMetadata.Err == nil {
		t.Fatal("Terraform plan is expected to fail.")
	}

	receivedError := strings.Replace(planMetadata.Err.Error(), "\n", " ", -1)

	assert.Equal(t, strings.Contains(receivedError, planFailsWithErrorMetadata.ErrorMessageContains), true, "The expected error message ("+planFailsWithErrorMetadata.ErrorMessageContains+") is not contained in the error message.")
}

// ------------------------------------------------------------------------------------------------------------------------------
