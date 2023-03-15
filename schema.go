package test

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
	PlanAssertions  PlanAssertions  `mapstructure:"plan"`
	ApplyAssertions ApplyAssertions `mapstructure:"apply"`
}

type PlanAssertions struct {
	Assertions []Assertion
}

type ApplyAssertions struct {
	EnsureIdempotent bool `mapstructure:"ensure_idempotent"`
	Assertions       []Assertion
}

type Assertion struct {
	Type     string
	Metadata map[interface{}]interface{} `mapstructure:",remain"`
}
