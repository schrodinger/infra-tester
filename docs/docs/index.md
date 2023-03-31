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

    === "Build From Source"

        ```shell
        # Clone the repo
        $ git clone git@github.com:schrodinger/infra-tester.git

        # Go to src directory and build the executable
        $ cd src/
        $ go test -c -o infra-tester

        # Move the binary to a directory in the $PATH
        $ sudo mv infra-tester /usr/local/bin
        ```

    === "Install Latest Linux Release"

        ```shell
        # Download the latest release binary
        $ curl --location --silent --fail --show-error -o infra_tester infra_tester \
          https://github.com/schrodinger/infra-tester/releases/latest/download/infra-tester-linux-x86_64

        # Make it executable
        $ chmod +x infra_tester

        # Move it to a directory in the $PATH
        $ sudo mv infra_tester /usr/local/bin
        ```

    === "Install Latest MacOS Release"

        ```shell
        # Download the latest release binary
        $ curl --location --silent --fail --show-error -o infra_tester infra_tester \
          https://github.com/schrodinger/infra-tester/releases/latest/download/infra-tester-macos-x86_64

        # Make it executable
        $ chmod +x infra_tester

        # Move it to a directory in the $PATH
        $ sudo mv infra_tester /usr/local/bin
        ```

    === "Install Latest Windows Release"

        1. Download the latest Windows release binary from the below URL:
            ```
            https://github.com/schrodinger/infra-tester/releases/latest/download/infra-tester-windows-x86_64
            ```

        2. Move it to a directory under `PATH`, or add the directory where you'd like to keep the executable to `PATH`.

#### Use *infra-tester* to run tests

Once *infra-tester* is set up, run the [example tests](https://github.com/schrodinger/infra-tester/tree/main/example):

!!! example "Running Example Tests"

    ```shell
    # Clone the repo
    $ git clone git@github.com:schrodinger/infra-tester.git

    # Change directory to example tests
    $ cd example/

    # Run the tests
    $ infra-tester -test.v

    ```
