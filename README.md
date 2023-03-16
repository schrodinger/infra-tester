# Terraform Module Tests

Run tests on Terraform code with just a configuration. `infra-tester` implements the test flow described in the diagram below. The concrete steps in the test flow are described by a **test plan** in YAML format (see the example `config.yaml` to find out more about the format of the test plan). It uses the [terratest](https://terratest.gruntwork.io/) library, but doesn't require any knowledge about terratest or Golang.

```mermaid
flowchart TB
  subgraph TestPlan[<b>Test Plan</b>]
    direction TB
    subgraph Test1[<b>Test 1</b>]
        direction LR
        Plan1(Plan\nResource\nCreation) --> AssertP1(Plan\nAssertions\n1..N<sub>1</sub>) --> Create1(Resource\nCreation) --> Assert1(Assertions\n1..M<sub>1</sub>)
    end
    subgraph Test2[<b>Test 2</b>]
        direction LR
        Plan2(Plan\nResource\nModification) --> AssertP2(Plan\nAssertions\n1..N<sub>2</sub>) --> Create2(Resource\nModification) --> Assert2(Assertions\n1..M<sub>2</sub>)
    end
    subgraph TestN[<b>Test N</b>]
        direction LR
        PlanN(Plan\nResource\nModification) --> AssertPN(Plan\nAssertions\n1..N<sub>k</sub>) --> CreateN(Resource\nModification) --> AssertN(Assertions\n1..M<sub>k</sub>)
    end
    Test1 --> Test2
    Test2 --> dots((...))
    dots((...)) --> TestN
  end
Start:::term --> Test1
TestN --> destroy(Destroy\nResources) --> End:::term

style TestPlan fill:#d3d3d3
classDef term fill:#00FF00,stroke:#333,stroke-width:2px
```

### Running the Tests

You can run the tests with `go test -v`

You can run specific tests with `go test -v -run <test name pattern>`. See more about `-run` [here](https://pkg.go.dev/testing#hdr-Subtests_and_Sub_benchmarks)

### Writing the Config

The configuration has the following structure:

```yaml
test_plan:
  name: <Name of the test plan, usually the resource or module name>
  tests:
    - name: <Name of the test>
      with_clean_state: true                            # Whether this test should be run in a clean state. If true, terraform destroy will be run before running the tests.
      vars:                                             # Any values for the var you'd like to pass to terraform. Support complex objects as well
        check_condition: true
        complex_object: &complex_object
          id: 0
          count: 0
          str: "string"
          seq:
            - one
            - two
          map:
            key: value
      plan:                                             # Any assertions that are to run during the plan step
        assertions:
          - type: PlanSucceeds
      apply:                                            # Any assertions that are to run during the plan step
        ensure_idempotent: true                         # If true, makes sure the plan is idempotent
        assertions:                                     # list of assertions
          - type: <AssertionType>                       # The type of assertion
            <Assertion Inputs>                          # Any inputs to the assertions

          - type: OutputEqual                           # An example assertion
            output_name: sample_output
            value: it's working
```

## Assertions

#### Plan Assertions

```yaml
      plan:
        assertions:
          - type: PlanSucceeds                                  # Makes sure plan succeeds

          - type: PlanFailsWithError                            # Makes sure plan fails and that the error message contains a specific string
            error_message_contains: Intended to fail            # The string that the error should contain
```

```yaml
    apply:
        assertions:
          - type: ApplySucceeds                                 # Makes sure apply succeeds

          - type: OutputEqual                                   # compare the value of an output
            output_name: sample_output                          # name of the output
            value: it's working                                 # the value it should be equal to

          - type: OutputsAreEqual                               # make sure values of multiple outputs are equal
            output_names:                                       # the list of output names, all of which should have same value
              - sample_output
              - another_output
              - yet_another_output

          - type: OutputContains                                # Makes sure the value of an output contains a string
            output_name: sample_output                          # name of the output
            value: working                                      # the substring value the output must contain

          - type: NoResourcesAffected                           # Makes sure no resources are affected

          - type: ResourcesAffected                             # assert number of resources affected
            added: 1                                            # Note that the test only checks for values explicitly specified here
            # changed: 0                                        # if "changed" key is omitted, it won't be checked
            # deleted: 5                                        # can also specify number of resources deleted

          - type: OutputMatchesRegex                            # asserts an output matches a regular expression
            output_name: an_output                              # output name
            regex: strings \w+ \d+ apple \d\s+\w+               # regular expression
```

## Test Output

A test summary is generated at the end of the test, which will look like:

```
-- PASS: TestMain (4.39s)
    --- PASS: TestMain/<TestPlanName> (4.39s)
        --- PASS: TestMain/<TestPlanName>/<TestName> (0.25s)
            --- PASS: TestMain/<TestPlanName>/<TestName>/Plan (0.25s)
                --- PASS: TestMain/<TestPlanName>/<TestName>/Plan/<PlanAssertion1> (0.00s)
                --- PASS: TestMain/<TestPlanName>/<TestName>/Plan/<PlanAssertion2> (0.00s)
                    ...
            --- PASS: TestMain/<TestPlanName>/<TestName>/Apply (0.25s)
                --- PASS: TestMain/<TestPlanName>/<TestName>/Apply/<ApplyAssertion1> (0.00s)
                --- PASS: TestMain/<TestPlanName>/<TestName>/Apply/<ApplyAssertion2> (0.00s)
                    ...
PASS
ok      schrodinger.com/infra-tester    4.890s
```

In the above test summary:
- `TestPlanName` is obtained from the `name` property of `test_plan` in the yaml config.
- TestName corresponds to the name of each test defined in the test plan.
- PlanAssertion1, PlanAssertion2, and so on refers to assertions in the plan step.
- ApplyAssertion1, ApplyAssertion2, and so on refers to assertions in the apply step.


As seen in the test summary, Plan and Apply tests are separated so you can run them separately using `-run`.
> **Note** If a test is dependent (e.g, by using a test as a "stage") on the resultant Terraform state of a previous test, then selectively running a test that has such a dependency will obviously fail. In this case, you might want to name the test and it's dependency test in such a way that, when you selectively run the test with a test name pattern, both the tests will be selected.
