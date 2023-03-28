# Test Output

A test summary is generated at the end of the test if *infra-tester* is run with the `-test.v` flag.

```title="Test Summary"
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

  - **`TestPlanName`** is obtained from the `name` property of `test_plan` in the YAML config.
  - **`TestName`** corresponds to the `name` of each test defined in the test plan.
  - **`PlanAssertion1`**, **`PlanAssertion2`**, and so on refer to the name (if defined, else assertion type) of the assertions in the plan step.
  - **`ApplyAssertion1`**, **`ApplyAssertion2`**, and so on refer to the name (if defined, else assertion type) of the assertions in the apply step.

As seen in the test summary, Plan and Apply tests are separated so you can run them separately using **`-test.run`** flag.

!!! warning

    If a test is dependent (e.g, by using a test as a "stage") on the resultant Terraform state of a previous test, then selectively running a test that has such a dependency will obviously fail. In this case, you might want to name the test and its dependency test in such a way that, when you selectively run the test with a test name pattern, both the tests will be selected.
