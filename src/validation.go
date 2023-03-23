package test

import (
	"fmt"
	"log"

	"schrodinger.com/infra-tester/assertions"
)

func validateTest(test Test) error {
	for _, assertion := range test.PlanAssertions.Assertions {
		if err := assertions.ValidateAssertion(assertion, "plan"); err != nil {
			return fmt.Errorf("assertion '%s' for plan step failed validation because - %s", assertion.Type, err)
		}
	}

	for _, assertion := range test.ApplyAssertions.Assertions {
		if err := assertions.ValidateAssertion(assertion, "apply"); err != nil {
			return fmt.Errorf("assertion '%s' for apply step failed validation because - %s", assertion.Type, err)
		}
	}

	return nil
}

func validateTests(tests []Test) error {
	validatedTests := make(map[string]Test)

	for _, test := range tests {
		testName := test.Name

		if test.Name == "" {
			return fmt.Errorf("test name is not defined")
		}

		// check for duplicate test names
		if _, ok := validatedTests[testName]; ok {
			return fmt.Errorf("test name '%s' is already defined previously - tests with same name are not allowed", testName)
		}

		if err := validateTest(test); err != nil {
			return fmt.Errorf("test '%s' failed validation: %s", testName, err)
		}

		validatedTests[testName] = test
	}

	log.Println("INFO: All tests and assertions are valid.")

	return nil
}
