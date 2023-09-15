# Introduction

Testing Infrastructure as Code (**IaC**) can be very complex. Usually, you pick one of the several testing frameworks or libraries available in your programming language of choice such as Go or Python. The vast majority of the infrastructure (unit) tests mostly make sure:

 - A module deploys without any failure for valid inputs.
 - Guard rails work as expected in catching invalid inputs or states.
 - The outputs of a module are as expected.

With *infra-tester*, these could be achieved without writing tests in a programming language such as Go or Python. You can define the tests using [**YAML**](https://yaml.org/) configuration. This reduces the barrier in testing infrastructure by not having to worry about maintaining lots of code just for testing. *infra-tester* provides several assertions that you can use and we'll add even more as more people use it.

## Getting Started

Terraform must be already installed on your system and available in `$PATH` as *infra-tester* does not bundle Terraform.
See [**official Terraform documentation**](https://developer.hashicorp.com/terraform/tutorials/aws-get-started/install-cli#install-terraform) on how to install it.

#### Install *infra-tester*

!!! info "Install"

    === "Install using Go"
        ```shell
        go install github.com/schrodinger/infra-tester@latest
        ```

    === "Install from GitHub Release - Linux and MacOS"

        ```shell
        PLATFORM=linux
        PLATFORM=macos
        # Download the latest release binary
        $ curl -L https://github.com/schrodinger/infra-tester/releases/latest/download/infra-tester-${PLATFORM}-x86_64 -o infra-tester

        # Make it executable
        $ chmod +x infra-tester

        # Move it to a directory in the $PATH
        $ sudo mv infra-tester /usr/local/bin
        ```

    === "Install from GitHub Release - Windows"

        1. Download the latest Windows release binary from the below URL:
            ```
            https://github.com/schrodinger/infra-tester/releases/latest/download/infra-tester-windows-x86_64.exe
            ```

        2. Move it to a directory under `PATH`, or add the directory where you'd like to keep the executable to `PATH`.

    === "Build From Source"

        ```shell
        # Clone the repo
        $ git clone git@github.com:schrodinger/infra-tester.git

        # Build the executable and move the binary to a directory in the $PATH
        $ cd infra-tester && go build -o bin/infra-tester
        $ sudo mv infra-tester /usr/local/bin

        # OR you may run go install
        $ go install

        ```

#### Writing Your First *infra-tester* Test Configuration

The [Writing a Test From Scratch](./writing_tests.md) page provides a simple hands-on
tutorial where we write a simple Terraform module and then write a test config
to test it.

#### Use *infra-tester* to run tests

Once *infra-tester* is set up, run the [example tests](https://github.com/schrodinger/infra-tester/tree/main/example):

!!! example "Running Example Tests"

    ```shell
    # Clone the repo
    $ git clone git@github.com:schrodinger/infra-tester.git

    # Change directory to example tests
    $ cd example/

    # Install plugin example assertions to run test
    $ pip install \
        "git+https://github.com/schrodinger/infra-tester.git#subdirectory=python-plugins/" \
         "git+https://github.com/schrodinger/infra-tester.git#subdirectory=example/plugin-example/"

    # Run the tests
    $ infra-tester -test.v

    ```

#### Extending *infra-tester* Using Plugins

If you are interested in extending *infra-tester*'s assertion library by writing
plugins, we provide extensive documentation on how developers can [write plugins](./extending_infra_tester.md) in Python and then use them in the test configuration.
