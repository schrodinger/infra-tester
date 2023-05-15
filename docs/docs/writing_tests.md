# Writing a Test From Scratch

In this section, we'll write a simple Terraform code to create some
local resources. We'll then write a configuration to run tests using
*infra-tester*.

Before proceeding any further, make sure Terraform and *infra-tester*
are installed in your environment. Refer to [Install *infra-tester*](./index.md#install-infra-tester) section on how to install it.

## Terraform

Since we'd like this example to be simple and easy to follow, we'll
use the [`time`](https://registry.terraform.io/providers/hashicorp/time/latest) provider to create a time "resource" and then we'll see how we can test it.

Let's start by creating a Terraform file and then adding the terraform
block with version constraints and required providers.

```terraform title="test.tf" linenums="1"
terraform {
  required_version = ">= 0.12"
  required_providers {
    time = {
      source  = "hashicorp/time"
      version = ">= 0.8"
    }
  }
}
```

Now let's define a basic `time_static` resource and then add its value to the
outputs.

```terraform title="test.tf" linenums="10"

# Create a time resource.
resource "time_static" "my_time" {}

# Show the current time in the RFC-3339 format.
output "current_time" {
  value = time_static.my_time.rfc3339
}
```

And that's it for the Terraform code. Let's try it out.
```shell
terraform init
terraform apply # Review your plan and approve the changes.
# You should see the `current_time` output. You can also run
terraform show # to see the outputs.
terraform destroy # To destroy the resources.
```

## Writing Tests Using *infra-tester*

Let's imagine that we use this Terraform code or module to generate a time stamp which is then consumed by other modules. Maybe the other modules expect
the output to be in a certain format. In this case, we'd like to ensure that
the output adheres to the RFC 3339 format no matter what underlying provider
we use to generate the `current_time` output.

A basic test would be a regular expression matching to make sure the output
adheres to the RFC 3339 format. Let's see how we can write such a test using
*infra-tester*.

Let's create a `config.yaml` file in the same directory where we created the
Terraform file and copy the below code into it. See the annotation next
to the code to understand what it does.

```yaml title="config.yaml" linenums="1"
test_plan:
  name: Time #(1)
  tests: # (2)
    - name: CurrentTimeOutputTests # (3)
      apply: # (4)
        assertions: # (5)
          - name: TimeStringMatchesRFC3339 # (6)
            type: OutputMatchesRegex # (7)
            output_name: current_time # (8)
            regex: ^((?:(\d{4}-\d{2}-\d{2})T(\d{2}:\d{2}:\d{2}(?:\.\d+)?))(Z|[\+-]\d{2}:\d{2})?)$ # (9)
```

1. This is the test plan name. All the tests are grouped under this test plan name.
It's best to name it the same as the module or component name.
2. The `tests` block can contain a list of tests that are to be run.
3. Each test must have a unique name. This name will show up in the final test output.
4. The `apply` block can contain assertions that will be run after `terraform apply`
was run. Similarly, there's also a `plan` block which can also contain a list of
assertions that are to be run after a `terraform plan`.
5. The `assertions` block can contain a list of assertions to be run under the `plan` or `apply` step depending on whether it's defined under `plan` or `apply`
block.
6. An assertion can optionally have a name. If a name is not provided the `type` of
assertion will be used to generate a name for the assertion in the test output.
7. `type` of assertion determines what assertion will be run. This must be a
valid assertion. *infra-tester* provides several inbuilt assertion types.
It also supports a plugin system to introduce custom assertions as well.
8. `output_name` is an input field specific to the `OutputMatchesRegex` assertion.
The assertion captures the string value of the output named `output_name`
to match the regular expression.
9. `OutputMatchesRegex` specific input field which will be used to match against the output value.

The above configuration is all that's required to test the use case we mentioned before. Now let's run the tests using *infra-tester*.

## Running the Tests

Change the working directory to the same directory where you created the Terraform
file and the `config.yaml` file and run

```shell
infra-tester -test.v
```

The `-test.v` flag can be used to run tests in verbose mode.

You should see the logs appear as the test runs, and finally, the test output is
printed.

```
--- PASS: TestMain (3.35s)
    --- PASS: TestMain/Time (2.97s)
        --- PASS: TestMain/Time/CurrentTimeOutputTests (1.03s)
            --- PASS: TestMain/Time/CurrentTimeOutputTests/Apply (1.03s)
                --- PASS: TestMain/Time/CurrentTimeOutputTests/Apply/TimeStringMatchesRFC3339 (0.09s)
PASS
```

## Try Breaking It!

Let's modify the regular expression so that it's invalid and see what happens.

To make it invalid let's remove the first two opening brackets, so the
line would then be:

```yaml title="config.yaml" linenums="10"
            regex: ^?:(\d{4}-\d{2}-\d{2})T(\d{2}:\d{2}:\d{2}(?:\.\d+)?))(Z|[\+-]\d{2}:\d{2})?)$

```

Let's run `infra-tester -test.v` to see what happens.

```
$ infra-tester -test.v
=== RUN   TestMain
    assertions.go:117: ERROR: Failure during test validation: test 'CurrentTimeOutputTests' failed validation: assertion 'OutputMatchesRegex' for apply step failed validation because - invalid regular expression
--- FAIL: TestMain (0.36s)
FAIL
```
As you can see, *infra-tester* runs test validation before running any of the
tests, and in this specific case, it caught the invalid regular expression.

Catching issues early on is very important, especially in the case of
Infrastructure as Code, as it reduces the time wasted on silly typos and
easy-to-catch issues. This leads to a better developer experience.

Let us now try to pass it a valid regular expression but one that doesn't
adhere to RFC 3339. Let's modify the line to the following:

```yaml title="config.yaml" linenums="10"
            regex: ^((?:T(\d{2}:\d{2}:\d{2}(?:\.\d+)?))(Z|[\+-]\d{2}:\d{2})?)$
```

And now if you run *infra-tester* again, the validation passes, but the test
fails as expected:

```
--- FAIL: TestMain (1.60s)
    --- FAIL: TestMain/Time (1.27s)
        --- FAIL: TestMain/Time/CurrentTimeOutputTests (0.47s)
            --- FAIL: TestMain/Time/CurrentTimeOutputTests/Apply (0.47s)
                --- FAIL: TestMain/Time/CurrentTimeOutputTests/Apply/TimeStringMatchesRFC3339 (0.08s)
FAIL
```

## What's More?

This section covered the basics of *infra-tester*. There are several more
features like passing variable inputs through the YAML configuration,
partially matching complex Terraform outputs, creating custom assertions
to extend *infra-tester*'s capabilities and so on, all of which are
extensively documented in this documentation site.
