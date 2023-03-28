# Plan Assertions

Plan assertions will be run after running **`terraform plan`**.

### PlanSucceeds

Asserts that **`terraform plan`** succeeds.

=== "Schema"
    ```yaml
    - name: <name>
      type: PlanSucceeds
    ```

=== "Example"
    ```yaml
    - name: PlanShouldSucceed
      type: PlanSucceeds
    ```

| Inputs | Description            | Type   | Required |
| ------ | ---------------------- | ------ | -------- |
| `name` | Name for the assertion | String | No       |

### PlanFails

Asserts that **`terraform plan`** fails.

=== "Schema"
    ```yaml
    - name: <name>
      type: PlanFails
    ```

=== "Example"
    ```yaml
    - name: PlanMustFail
      type: PlanFails
    ```

| Inputs | Description            | Type   | Required |
| ------ | ---------------------- | ------ | -------- |
| `name` | Name for the assertion | String | No       |

### PlanFailsWithError

Asserts that **`terraform plan`** with an error containing a certain substring.

=== "Schema"
    ```yaml
    - name: <name>
      type: PlanFailsWithError
      error_message_contains: <error_message_contains>
    ```

=== "Example"
    ```yaml
    - name: MustFailWithSampleError
      type: PlanFailsWithError
      error_message_contains: Failed with sample error
    ```

| Inputs                   | Description                                        | Type   | Required |
| ------------------------ | -------------------------------------------------- | ------ | -------- |
| `name`                   | Name for the assertion                             | String | No       |
| `error_message_contains` | String that should be present in the error message | String | **Yes**  |
