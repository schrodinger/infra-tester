package assertions

import (
	"fmt"
	"testing"

	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/schrodinger/infra-tester/plugins"
	"github.com/schrodinger/infra-tester/utils"
)

type Assertion struct {
	Type     string
	Name     string
	Metadata map[interface{}]interface{} `mapstructure:",remain"`
}

// Converts the assertion to a generic map so that it can be serialized
// by json.Marshal().
func (assertion Assertion) ToGenericMap() map[string]interface{} {
	genericMap := map[string]interface{}{}
	genericMap["type"] = assertion.Type
	genericMap["name"] = assertion.Name
	genericMap["metadata"] = utils.ConvertToGenericInterface(assertion.Metadata)

	return genericMap
}

type AssertionContext struct {
	AvailablePlugins map[string]bool
	PluginManager    *plugins.PluginManager
}

type AssertionImplementation struct {
	ValidateFunction func(Assertion) error
	RunFunction      func(t *testing.T, terraformOptions *terraform.Options, assertion Assertion, stepMetadata interface{})
}

func GetAssertionImplementation(assertionType string, step string, assertionContext *AssertionContext) (AssertionImplementation, error) {
	var assertionImplementation AssertionImplementation
	var ok bool
	if step == "plan" {
		if assertionImplementation, ok = ValidPlanAssertions[assertionType]; !ok {
			// It could be a plugin assertion type.
			if assertionContext.PluginManager != nil {
				if _, ok := assertionContext.AvailablePlugins[assertionType]; ok {
					return GetCustomAssertionImplementation(assertionType, assertionContext.PluginManager)
				}
			}

			// Maybe users are trying to use an apply assertion in plan step.
			if _, ok = ValidApplyAssertions[assertionType]; ok {
				return AssertionImplementation{}, fmt.Errorf("'%s' is only valid for 'apply' tests", assertionType)
			}

			return AssertionImplementation{}, fmt.Errorf("assertion type '%s' is invalid", assertionType)
		}
	} else if step == "apply" {
		if assertionImplementation, ok = ValidApplyAssertions[assertionType]; !ok {
			// It could be a plugin assertion type.
			if assertionContext.PluginManager != nil {
				if _, ok := assertionContext.AvailablePlugins[assertionType]; ok {
					return GetCustomAssertionImplementation(assertionType, assertionContext.PluginManager)
				}
			}

			// Maybe users are trying to use a plan assertion in apply step.
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

func ValidateAssertion(assertion Assertion, step string, assertionContext *AssertionContext) error {
	AssertionImplementation, err := GetAssertionImplementation(assertion.Type, step, assertionContext)
	if err != nil {
		return err
	}

	validateFunction := AssertionImplementation.ValidateFunction
	return validateFunction(assertion)
}

func RunAssertion(
	t *testing.T,
	terraformOptions *terraform.Options,
	assertion Assertion,
	step string,
	stepMetadata interface{},
	assertionContext *AssertionContext) {
	assertionType := assertion.Type
	var assertionImplementation AssertionImplementation

	assertionImplementation, err := GetAssertionImplementation(assertionType, step, assertionContext)
	if err != nil {
		// This shouldn't happen as we are validating the tests before running them
		ErrorAndSkipf(t, "ERROR: Failure while running assertion: %s.\n", err)
	}

	runFunction := assertionImplementation.RunFunction
	runFunction(t, terraformOptions, assertion, stepMetadata)
}

func ErrorAndSkip(t *testing.T, args ...any) {
	t.Error(args...)
	t.SkipNow()
}

func ErrorAndSkipf(t *testing.T, format string, args ...any) {
	t.Errorf(format, args...)
	t.SkipNow()
}
