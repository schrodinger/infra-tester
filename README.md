# infra-tester

Run tests on Terraform code with just a configuration. It hides the boilerplate code for common infrastructure test patterns and lets you focus on defining the tests using simple YAML configurations.

## Introduction

Testing Infrastructure as Code (**IaC**) can be very complex. Usually, you pick one of the several testing frameworks or libraries available in your programming language of choice such as Go or Python. The vast majority of the infrastructure (unit) tests mostly make sure:

 - A module deploys without any failure for valid inputs.
 - Guard rails work as expected in catching invalid inputs or states.
 - The outputs of a module are as expected.

With *infra-tester*, these could be achieved without writing tests in a programming language such as Go or Python. You can define the tests using [**YAML**](https://yaml.org/) configuration. This reduces the barrier in testing infrastructure by not having to worry about maintaining lots of code just for testing. *infra-tester* provides several assertions that you can use and we'll add even more as more people use it.

## Documentation

You can find extensive documentation on *infra-tester* [here](https://jubilant-disco-ey46y7z.pages.github.io/).
