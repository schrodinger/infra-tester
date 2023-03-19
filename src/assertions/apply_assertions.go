package assertions

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/mitchellh/mapstructure"
	"github.com/stretchr/testify/assert"
)

type ApplyAssertions struct {
	EnsureIdempotent bool `mapstructure:"ensure_idempotent"`
	Assertions       []Assertion
}

type ApplyMetadata struct {
	CmdOut string
	Err    error
}

var ValidApplyAssertions = map[string]AssertionImplementation{
	"ApplySucceeds": {
		ValidateFunction: func(a Assertion) error { return nil },
		RunFunction:      AssertApplySucceeds,
	},

	// Output asserts

	"OutputEqual": {
		ValidateFunction: validateOutputEqualAssertion,
		RunFunction:      AssertOutputEqual,
	},
	"OutputsAreEqual": {
		ValidateFunction: validateOutputsAreEqualAssertion,
		RunFunction:      AssertOutputsAreEqual,
	},
	"OutputContains": {
		ValidateFunction: validateOutputContainsAssertion,
		RunFunction:      AssertOutputContains,
	},

	"OutputMatchesRegex": {
		ValidateFunction: validateOutputMatchesRegexAssertion,
		RunFunction:      AssertOutputMatchesRegex,
	},

	// Resource count asserts

	"ResourcesAffected": {
		ValidateFunction: validateResourcesModified,
		RunFunction:      AssertResourcesAffected,
	},
	"NoResourcesAffected": {
		ValidateFunction: func(a Assertion) error { return nil },
		RunFunction:      AssertNoResourcesAffected,
	},
}

// ------------------------------------------------------------------------------------------------------------------------------

func AssertApplySucceeds(t *testing.T, terraformOptions *terraform.Options, assertion Assertion, stepMetadata interface{}) {
	var applyMetadata ApplyMetadata
	var ok bool
	if applyMetadata, ok = stepMetadata.(ApplyMetadata); !ok {
		ErrorAndSkip(t, "stepMetadata is not of type ApplyMetadata")
	}

	if applyMetadata.Err != nil {
		ErrorAndSkip(t, "Terraform apply is expected to succeed.")
	}
}

// ------------------------------------------------------------------------------------------------------------------------------

type outputEqualMetadata struct {
	OutputName    string                      `mapstructure:"output_name"`
	CompleteMatch bool                        `mapstructure:"complete_match"`
	Value         map[interface{}]interface{} `mapstructure:",remain"`
}

func validateOutputEqualAssertion(assertion Assertion) error {
	var outputEqualMetadata outputEqualMetadata

	err := mapstructure.Decode(assertion.Metadata, &outputEqualMetadata)
	if err != nil {
		return fmt.Errorf("error decoding assertion metadata: %s", err)
	}

	if outputEqualMetadata.OutputName == "" {
		return fmt.Errorf("output_name is either not defined or is empty")
	}

	val, ok := outputEqualMetadata.Value["value"]
	if !ok {
		return fmt.Errorf("value is not defined")
	}

	if val == nil {
		return fmt.Errorf("value can not be empty")
	}

	for key := range outputEqualMetadata.Value {
		if key != "value" {
			return fmt.Errorf("unexpected key: %s", key)
		}
	}

	return nil
}

func AssertOutputEqual(t *testing.T, terraformOptions *terraform.Options, assertion Assertion, stepMetadata interface{}) {
	var outputEqualMetadata outputEqualMetadata

	err := mapstructure.Decode(assertion.Metadata, &outputEqualMetadata)
	if err != nil {
		ErrorAndSkipf(t, "error decoding assertion metadata: %s", err)
	}

	// Get properties
	outputName := outputEqualMetadata.OutputName
	expectedValue := outputEqualMetadata.Value["value"]
	outputValue := terraform.OutputAll(t, terraformOptions)[outputName]

	partialComparisonResult := partialDeepCompare(expectedValue, outputValue)
	if partialComparisonResult != nil {
		ErrorAndSkipf(t, "The property %s has an unexpected value.\n\nExpected value:\n%+v\n\nActual Value:\n%+v\n\nReason: %s", outputName, expectedValue, outputValue, partialComparisonResult.Error())
	}

	if !outputEqualMetadata.CompleteMatch {
		return
	}

	fullComparisonResult := partialDeepCompare(outputValue, expectedValue)
	if fullComparisonResult != nil {
		ErrorAndSkipf(t, "The property %s has an unexpected value.\n\nExpected following value(s):\n%+v\n\nActual value:\n%+v\n\nReason: %s", outputName, expectedValue, outputValue, fullComparisonResult.Error())
	}
}

// ------------------------------------------------------------------------------------------------------------------------------

type resourcesModifiedMetadata struct {
	Added     int
	Changed   int
	Destroyed int
}

func validateResourcesModified(assertion Assertion) error {
	var resourcesModifiedMetadata resourcesModifiedMetadata

	decoderMetadata, err := decodeWithMetadata(assertion, &resourcesModifiedMetadata)
	if err != nil {
		return fmt.Errorf("error decoding assertion metadata: %s", err)
	}

	if len(decoderMetadata.Keys) == 0 {
		return fmt.Errorf("at least one of the following keys must be specified: Added, Changed, Destroyed")
	}

	return nil
}

func AssertResourcesAffected(t *testing.T, terraformOptions *terraform.Options, assertion Assertion, stepMetadata interface{}) {
	var resourcesModifiedMetadata resourcesModifiedMetadata
	decoderMetadata, err := decodeWithMetadata(assertion, &resourcesModifiedMetadata)
	if err != nil {
		ErrorAndSkipf(t, "error while decoding assertion metadata: %s", err)
	}

	// cast stepMetadata to ApplyMetadata
	var applyMetadata ApplyMetadata
	var ok bool
	if applyMetadata, ok = stepMetadata.(ApplyMetadata); !ok {
		ErrorAndSkip(t, "stepMetadata is not of type ApplyMetadata")
	}

	resourcesCount := terraform.GetResourceCount(t, applyMetadata.CmdOut)

	// Only check for keys explicitly specified in yaml config
	for _, key := range decoderMetadata.Keys {
		if key == "Added" {
			assert.Equal(t, resourcesModifiedMetadata.Added, resourcesCount.Add, "Unexpected number of resources were added.")
		} else if key == "Changed" {
			assert.Equal(t, resourcesModifiedMetadata.Changed, resourcesCount.Change, "Unexpected number of resources were changed.")
		} else if key == "Destroyed" {
			assert.Equal(t, resourcesModifiedMetadata.Destroyed, resourcesCount.Destroy, "Unexpected number of resources were destroyed.")
		}
	}
}

func AssertNoResourcesAffected(t *testing.T, terraformOptions *terraform.Options, assertion Assertion, stepMetadata interface{}) {
	assertion.Metadata = map[interface{}]interface{}{
		"added":     0,
		"changed":   0,
		"destroyed": 0,
	}

	AssertResourcesAffected(t, terraformOptions, assertion, stepMetadata)
}

// ------------------------------------------------------------------------------------------------------------------------------

type outputsAreEqualMetadata struct {
	OutputNames []string `mapstructure:"output_names"`
}

func validateOutputsAreEqualAssertion(assertion Assertion) error {
	var outputsAreEqualMetadata outputsAreEqualMetadata

	err := mapstructure.Decode(assertion.Metadata, &outputsAreEqualMetadata)
	if err != nil {
		return fmt.Errorf("error decoding assertion metadata: %s", err)
	}

	if len(outputsAreEqualMetadata.OutputNames) < 2 {
		return fmt.Errorf("there must be at least two output names for comparison")
	}

	return nil
}

func AssertOutputsAreEqual(t *testing.T, terraformOptions *terraform.Options, assertion Assertion, stepMetadata interface{}) {
	var outputsAreEqualMetadata outputsAreEqualMetadata

	err := mapstructure.Decode(assertion.Metadata, &outputsAreEqualMetadata)
	if err != nil {
		ErrorAndSkipf(t, "error decoding assertion metadata: %s", err)
	}

	// get output names and values
	for i := range outputsAreEqualMetadata.OutputNames {
		if i == 0 {
			continue
		}

		outputName1 := outputsAreEqualMetadata.OutputNames[i-1]
		outputName2 := outputsAreEqualMetadata.OutputNames[i]

		outputValue1 := terraform.Output(t, terraformOptions, outputName1)
		outputValue2 := terraform.Output(t, terraformOptions, outputName2)

		assert.Equal(t, outputValue1, outputValue2, "The values for output \""+outputName1+"\" ("+outputValue1+") and \""+outputName2+"\" ("+outputValue2+") do not match")
	}
}

// ------------------------------------------------------------------------------------------------------------------------------

type outputContainsMetadata struct {
	OutputName string `mapstructure:"output_name"`
	Value      string
}

func validateOutputContainsAssertion(assertion Assertion) error {
	var outputContainsMetadata outputContainsMetadata

	err := mapstructure.Decode(assertion.Metadata, &outputContainsMetadata)
	if err != nil {
		return fmt.Errorf("error decoding assertion metadata: %s", err)
	}

	if outputContainsMetadata.OutputName == "" {
		return fmt.Errorf("output_name is either not defined or is empty")
	}

	if outputContainsMetadata.Value == "" {
		return fmt.Errorf("value is either not defined or is empty")
	}

	return nil
}

func AssertOutputContains(t *testing.T, terraformOptions *terraform.Options, assertion Assertion, stepMetadata interface{}) {
	var outputContainsMetadata outputContainsMetadata

	err := mapstructure.Decode(assertion.Metadata, &outputContainsMetadata)
	if err != nil {
		ErrorAndSkipf(t, "error decoding assertion metadata: %s", err)
	}

	// Get properties
	outputName := outputContainsMetadata.OutputName
	shouldContain := outputContainsMetadata.Value
	outputValue := terraform.Output(t, terraformOptions, outputName)

	assert.Contains(t, outputValue, shouldContain, "The property \""+outputName+"\" has an unexpected value. It does not contain the string \""+shouldContain+"\". Value is: \""+outputValue+"\".")
}

// ------------------------------------------------------------------------------------------------------------------------------

type outputMatchesRegexMetadata struct {
	OutputName string `mapstructure:"output_name"`
	Regex      string
}

func validateOutputMatchesRegexAssertion(assertion Assertion) error {
	var outputMatchesMetadata outputMatchesRegexMetadata

	err := mapstructure.Decode(assertion.Metadata, &outputMatchesMetadata)
	if err != nil {
		return fmt.Errorf("error decoding assertion metadata: %s", err)
	}

	if outputMatchesMetadata.OutputName == "" {
		return fmt.Errorf("output_name is either not defined or is empty")
	}

	if outputMatchesMetadata.Regex == "" {
		return fmt.Errorf("regex is either not defined or is empty")
	}

	if _, err := regexp.Compile(outputMatchesMetadata.Regex); err != nil {
		return fmt.Errorf("invalid regular expression")
	}

	return nil
}

func AssertOutputMatchesRegex(t *testing.T, terraformOptions *terraform.Options, assertion Assertion, stepMetadata interface{}) {
	var outputMatchesMetadata outputMatchesRegexMetadata

	err := mapstructure.Decode(assertion.Metadata, &outputMatchesMetadata)
	if err != nil {
		ErrorAndSkipf(t, "error decoding assertion metadata: %s", err)
	}

	// Get properties
	outputName := outputMatchesMetadata.OutputName
	regex := outputMatchesMetadata.Regex
	outputValue := terraform.Output(t, terraformOptions, outputName)

	// regexp is already validated in validateOutputMatchesAssertion so there shouldn't be any panic
	assert.Regexp(t, regexp.MustCompile(regex), outputValue, "The property \""+outputName+"\" has an unexpected value. It does not match the regular expression \""+regex+"\". Value is: \""+outputValue+"\".")
}
