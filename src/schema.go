package test

import "schrodinger.com/infra-tester/assertions"

type Config struct {
	TestPlan TestPlan                    `mapstructure:"test_plan"`
	Metadata map[interface{}]interface{} `mapstructure:",remain"`
}

type TestPlan struct {
	Name  string
	Tests []Test
}

type Test struct {
	Name            string
	WithCleanState  bool `mapstructure:"with_clean_state"`
	Vars            map[string]interface{}
	PlanAssertions  assertions.PlanAssertions  `mapstructure:"plan"`
	ApplyAssertions assertions.ApplyAssertions `mapstructure:"apply"`
}
