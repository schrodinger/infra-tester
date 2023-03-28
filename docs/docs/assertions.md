# Assertions

Assertions are the core of *infra-tester*. *infra-tester* provides several assertions to define your tests. "Plan assertions" run after the **`terraform plan`** step and "Apply Assertions" run after the **`terraform apply`** step.

Assertions generally have the following schema:

```yaml
- name: <An optional name for the assertion>
  type: <Type of the assertion>
  <Inputs specific to the assertion>
```

### **`name`**

This an optional field which allows you to set a custom meaningful name for the assertion.
If a value is set for `name`, it will be used in generating the test summary.

### **`type`**

The value of `type` must be one of the valid assertion types available.
You can refer to [**plan**](/plan_assertions) and [**apply**](/apply_assertions) assertions for the list of valid assertions.

### Assertion Inputs

Some assertions may require inputs, and different assertions will have different inputs.
