# Apply Assertions

Apply assertions will be run after running **`terraform apply`**.

### ApplySucceeds

Asserts that **`terraform apply`** succeeds.

=== "Schema"
    ```yaml
    - name: <name>
      type: ApplySucceeds
    ```

=== "Example"
    ```yaml
    - name: ApplyShouldSucceed
      type: ApplySucceeds
    ```

| Inputs | Description            | Type   | Required |
| ------ | ---------------------- | ------ | -------- |
| `name` | Name for the assertion | String | No       |

### OutputEqual

Compares a specified terraform output with an expected value.

=== "Schema"
    ```yaml
    - name: <name>
      type: OutputEqual
      output_name: <output_name>
      value: <value>
    ```

=== "Example 1"
    ```yaml
    # `value` is of string type
    - name: ASimpleOutputEqualExample
      type: OutputEqual
      output_name: sample_output
      value: a sample value

    # `value` is of boolean type
    - name: OutputEqualExampleForBool
      type: OutputEqual
      output_name: a_boolean_output
      value: true

    # `value` is of float type
    - name: OutputEqualExampleForFloat
      type: OutputEqual
      output_name: a_float_output
      value: 123.11
    ```

=== "Example 2"
    ```yaml
    # `value` is of map type
    - name: OutputEqualExampleForMap
      type: OutputEqual
      output_name: a_map_output
      value:
        key: value
    ```

=== "Example 3"
    ```yaml
    # `value` is of sequence type
    - name: OutputEqualExampleForList
      type: OutputEqual
      output_name: a_list_output
      value:
        - a
        - b
        - c
    ```

=== "Example 4"
    ```yaml
    # `value` is of an object type.
    # Note this example does not have `complete_match` enabled and
    # so it only checks the terraform output for keys and their values
    # for the specified fields (In this case, it checks values of "seq",
    # "map", "nested_map", "nested_key", and "boolean" and ignores other
    # fields).
    - name: OutputEqualExampleForComplexOutput
      type: OutputEqual
      output_name: a_complex_output
      value:
        seq:
          - a
          - b
          - c
        map:
          key: value
          nested_map:
            nested_key: nested_value
        boolean: true
    ```

=== "Example 5"
    ```yaml
    # `value` is of an object type, and a complete match is done
    - name: OutputEqualExampleForComplexOutputWithCompleteMatch
      type: OutputEqual
      output_name: a_complex_output
      complete_match: true      # Complete match is set to true
      value:
        natural_number: 100
        float: 123.11
        seq:
          - a
          - b
          - c
        str: hello
        map:
          key: value
          nested_map:
            nested_key: nested_value
        boolean: true
    ```




| Inputs           | Description                                                                                                                                        | Type                                                   | Required |
| ---------------- | -------------------------------------------------------------------------------------------------------------------------------------------------- | ------------------------------------------------------ | -------- |
| `name`           | Name for the assertion                                                                                                                             | String                                                 | No       |
| `output_name`    | Terraform output name to compare                                                                                                                   | String                                                 | **Yes**  |
| `value`          | The expected value - value can be of complex data type including objects consisting of maps, sequences, booleans, floats, integers, etc            | String, Integer, Float, Boolean, Map, Sequence, Object | **Yes**  |
| `complete_match` | Whether to match the full output by making sure the terraform output value has exactly the same fields specified in `value` - **false by default** | Boolean                                                | No       |

### OutputsAreEqual

Compares multiple terraform outputs and asserts that all the specified outputs have the same value.

=== "Schema"
    ```yaml
    - name: <name>
      type: OutputsAreEqual
      output_names:
        - <output_name_1>
        - <output_name_2>
        ...
    ```

=== "Example 1"
    ```yaml
    - name: TwoOutputsMustMatch
      type: OutputsAreEqual
      output_names:
        - sample_output
        - another_output
    ```

=== "Example 2"
    ```yaml
    - name: MultipleOutputsMustMatch
      type: OutputsAreEqual
      output_names:
        - sample_output
        - another_output
        - yet_another_output
    ```

| Inputs         | Description                    | Type               | Required |
| -------------- | ------------------------------ | ------------------ | -------- |
| `name`         | Name for the assertion         | String             | No       |
| `output_names` | List of terraform output names | Sequence of String | **Yes**  |

### OutputContains

Asserts that the specified terraform output contains a specified string.

=== "Schema"
    ```yaml
    - name: <name>
      type: OutputContains
      output_name: <output_name>
      value: <value>
    ```

=== "Example"
    ```yaml
    - name: OutputContainsACertainSubString
      type: OutputContains
      output_name: sample_output
      value: a certain substring
    ```

| Inputs        | Description                                        | Type   | Required |
| ------------- | -------------------------------------------------- | ------ | -------- |
| `name`        | Name for the assertion                             | String | No       |
| `output_name` | Name of the terraform output                       | String | **Yes**  |
| `value`       | A substring that the terraform output must contain | String | **Yes**  |

### OutputMatchesRegex

Asserts that the specified terraform output matches a specified regular expression.

=== "Schema"
    ```yaml
    - name: <name>
      type: OutputMatchesRegex
      output_name: <output_name>
      regex: <value>
    ```

=== "Example"
    ```yaml
    - name: OutputShouldMatchARegularExpression
      type: OutputMatchesRegex
      output_name: a_fourth_output
      regex: strings \w+ \d+ apple \d\s+\w+
    ```

| Inputs        | Description                                                       | Type   | Required |
| ------------- | ----------------------------------------------------------------- | ------ | -------- |
| `name`        | Name for the assertion                                            | String | No       |
| `output_name` | Name of the terraform output                                      | String | **Yes**  |
| `regex`       | Regular expression that the terraform output should match against | String | **Yes**  |

### ResourcesAffected

Asserts that **`terraform apply`** added, and/or changed, and/or destroyed a specified number of resources.

=== "Schema"
    ```yaml
    - name: <name>
      type: ResourcesAffected
      added: <added>
      changed: <changed>
      destroyed: <destroyed>
    ```

=== "Example 1"
    ```yaml
    # This only asserts that 1 resource has been added.
    - name: MustAddExactlyOneResource
      type: ResourcesAffected
      added: 1

    # This only asserts that 5 resources have been changed.
    - name: MustChangeFiveResource
      type: ResourcesAffected
      changed: 5

    # This only asserts that 1 resource has been destroyed.
    - name: MustDestroyOneResource
      type: ResourcesAffected
      destroyed: 1
    ```

=== "Example 2"
    ```yaml
    # This only asserts that 1 resource has been added, and 5 were changed.
    - name: MustAffectSpecificNumberOfResources
      type: ResourcesAffected
      added: 1
      changed: 5

    # This only asserts that 1 resource has been added, 5 were changed, and 0 were destroyed.
    - name: MustAffectSpecificNumberOfResources
      type: ResourcesAffected
      added: 1
      changed: 5
      destroyed: 0
    ```

| Inputs      | Description                                       | Type    | Required |
| ----------- | ------------------------------------------------- | ------- | -------- |
| `name`      | Name for the assertion                            | String  | No       |
| `added`     | Number of resources that must have been added     | Integer | No       |
| `changed`   | Number of resources that must have been changed   | Integer | No       |
| `destroyed` | Number of resources that must have been destroyed | Integer | No       |

!!! warning

    At least one of `added`, `changed`, or `destroyed` must be specified.
    If a field is not specified, *infra-tester* will not check against that
    specific field. Note that the default for unspecified fields is not zero.

### NoResourcesAffected

Similar to [ResourcesAffected](./apply_assertions.md#ResourcesAffected), but asserts that no resources have been added, changed, or destroyed.

=== "Schema"
    ```yaml
    - name: <name>
      type: NoResourcesAffected
    ```

=== "Example"
    ```yaml
    - name: MustAffectNoResource
      type: NoResourcesAffected
    ```
