# Configuration

*infra-tester* uses [**YAML**](https://yaml.org/) as its configuration language and supports most of YAML 1.1 and 1.2, including support for anchors, tags, map merging, etc. When *infra-tester* is run, it looks for a `config.yaml` file in the current directory to run the tests.

The configuration for *infra-tester* has the following structure:

```yaml
test_plan:
  name: <Name of the test plan, usually the resource or module name>

  # Optional field, if defined passes vars defined
  # here to the final cleanup destroy.
  destroy_vars: 
    # By default, infra-tester passes vars in the last test to terraform 
    # destroy. If the last test has invalid vars, this field will be useful
    # to pass valid vars to successfully run destroy.
    ...

  # A list of tests to be run.
  tests:
    # Each test must have a unique name.
    - name: <Name of the test>

      # Whether this test should be run in a clean state. If true, terraform
      # destroy will be run before running the test. Default is false.
      with_clean_state: false

      # Any values to be passed as vars to terraform.
      # Support complex objects as well.
      vars:
        complex_object: &complex_object
          id: 0
          count: 0
          str: "string"
          seq:
            - one
            - two
          map:
            key: value

      # Any checks that are to be run during the plan step.
      plan:
        # You can check for as many assertions as you want.
        assertions:
          - type: <AssertionType>           # The type of assertion
            <Assertion Inputs>              # Any inputs to the assertions

          # Example
          - type: PlanSucceeds

      # Any assertions that are to be run during the apply step.
      apply:
        # If true, makes sure the apply is idempotent.
        ensure_idempotent: true
        assertions:                         # list of assertions
          - type: <AssertionType>           # The type of assertion
            <Assertion Inputs>              # Any inputs to the assertions

          # Example
          - type: OutputEqual               # An example assertion
            output_name: sample_output
            value: it's working
```

### **`test_plan`**

*infra-tester* looks for the `test_plan` key to figure out what tests to run.
All other top-level keys are not considered a part of the test configuration.
This means you can have custom YAML blocks at the top level that can be referred to within the `test_plan`.
This will be particularly useful if you'd like to keep the config DRY by reusing commonly used blocks or values.

### **`test_plan.name`**
The test plan must define a `name`. This will be used in the test summary.
It's recommended to name the tests as the resource or the module you are testing.

### **`test_plan.destroy_vars`**
You can optionally set the value of `destroy_vars` to the value that must be used for the final cleanup at the end
of all tests. This will be useful when the last test in the config may be set to fail intentionally due to invalid
values for the input variables. Since *infra-tester* uses the last used vars to perform the final cleanup, invalid
inputs may cause the final cleanup to fail, and so setting `destroy_vars` allows you to pass values specifically
for the final cleanup.

### **`test_plan.tests`**

The `test_plan.tests` key should contain a list of tests that will be run for the given test plan.

### **`test_plan.tests.name`**

Each test must have a unique name across a given test plan. These names will be used in tests summary generation.

### **`test_plan.tests.with_clean_state`**

*infra-tester* run terraform apply between each test to move from one test to another. Running destroy between each and
every test is possible, but more often than not, the tests within a test plan have lots of similarity and running destroy
and then apply to recreate essentially the same state with a slight difference is generally not efficient. There are use
cases which require tests to be run in a clean state. For this reason *infra-tester* runs apply between tests but still
provides an option to run a test with a clean state if that is absolutely required. This can be done by setting the value
of `with_clean_state` to `true`.

### **`test_plan.tests.vars`**

`test_plan.tests.vars` can be used to pass values for the terraform input variables for running `terraform plan` and `terraform apply`.
All data types are supported, and *infra-tester* will convert the values in YAML to appropriate terraform data types.


### **`test_plan.tests.plan`**

This contains the list of assertions that will be run after running **`terraform plan`** with the `test.vars` as input. Assertions can be
defined under the key `assertions`. Read more about assertions [here](/assertions).

### **`test_plan.tests.apply`**

This contains the list of assertions that will be run after running **`terraform apply`** with the `test.vars` as input. Assertions can be
defined under the key `assertions`. Read more about assertions [here](/assertions).

If the tests require terraform apply to be idempotent, you can set `ensure_idempotent` to `true` to make sure the apply does not
result in any more changes when run a second time after the first apply.

### **`test_plan.tests.(plan|apply).assertions`**

This key contains a list of assertions to be in the `plan` or `apply` step respectively. Each assertion should specify the `type` key.
It can optionally define a name, which if provided will be used in the test summary generation, else it uses the `type` value.

## Validation

*infra-tester* will validate the configuration before running any tests. Each assertion will have its own validation that checks for
required fields, the type of the value, whether regular expression is valid and so on. This provides a better experience when writing
a test configuration and minimizes the time lost chasing trivial bugs in the configuration.
