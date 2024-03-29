import argparse
import contextlib
import json
import os
import sys
from enum import IntEnum
from importlib.metadata import entry_points
from typing import Callable, Union

from . import PLUGIN_GROUP, BaseAssertionPlugin
from .result import PluginResult


class ExitCodes(IntEnum):
    SUCCESS = 0
    ERROR = 1
    INVALID_INPUT = 2


def get_plugin(group: str, name: str) -> Callable[[], BaseAssertionPlugin]:
    """
    Get the callable that loads the plugin with the given name.

    Args:
        name (str): The name of the plugin to get.

    Raises:
        KeyError: If the plugin with the given name does not exist.

    Returns:
        Callable[[], BaseAssertionPlugin]: The callable that loads
        the plugin.
    """
    entry_point = [ep for ep in entry_points()[group] if ep.name == name]
    if len(entry_point) == 0:
        raise ValueError(f"No plugin with name '{name}' found.")
    elif len(entry_point) > 1:
        raise RuntimeError(
            (
                f"Multiple plugins with name '{name}' found. "
                "This should never happen."
            )
        )

    return entry_point[0].load()


def ensure_python_version() -> None:
    """
    Ensure that the Python version is supported.
    """
    if sys.version_info < (3, 8):
        print(
            "ERROR: (infra-tester-plugins) Python 3.7",
            "or higher is required to run this program.",
            f"Current version is {sys.version}.",
        )

        sys.exit(1)


def assertion_cli() -> Union[int, None]:
    """Entry point for the assertion CLI.

    Raises:
        ValueError: Raised when there's an error in the CLI arguments.

    Returns:
        int: Exit code of the CLI.
                0: Success
                1: Error
                2: Invalid input
    """

    ensure_python_version()

    parser = argparse.ArgumentParser(
        description=(
            "This program is part of infra-tester plugin framework. "
            "It is used to run plugins."
        )
    )

    parser.add_argument(
        "-n", "--name", type=str, help="Name of the plugin assertion to run."
    )

    parser.add_argument(
        "-a",
        "--action",
        type=str,
        choices=["validate_inputs", "run_assertion", "cleanup"],
        help="Action to run.",
    )

    parser.add_argument(
        "-i", "--inputs", type=str, default=None, help="Inputs in JSON."
    )

    parser.add_argument(
        "-s",
        "--state",
        type=str,
        default=None,
        help="Terraform state in JSON.",
    )

    args = parser.parse_args()

    if len(sys.argv) == 1:
        parser.print_help()
        sys.exit(1)

    assertion = None

    try:
        plugin_load_callable = get_plugin(PLUGIN_GROUP, args.name)
        assertion = plugin_load_callable()
    except ValueError:
        print(
            "ERROR: (infra-tester-plugins) ",
            f"Could not find plugin '{args.name}'.",
        )
        print(
            "INFO: Please make sure the PIP package ",
            "for the plugin is installed.",
        )
        print(
            "INFO: Run `infra-tester-plugin-manager --list` ",
            "to list all available plugins.",
        )

        # Exit with a non-zero exit code to indicate failure.
        return int(ExitCodes.ERROR)
    except Exception as e:
        print(
            "ERROR: (infra-tester-plugins) ",
            f"Failure while loading plugin '{args.name}': {e}.",
        )

        # Exit with a non-zero exit code to indicate failure.
        return int(ExitCodes.ERROR)

    return_value = None

    try:
        inputs = json.loads(args.inputs) if args.inputs is not None else None
    except json.JSONDecodeError as e:
        print(
            "ERROR: (infra-tester-plugins) ",
            f"Failure while parsing inputs: {e}.",
        )

        # Exit with a non-zero exit code to indicate failure.
        return int(ExitCodes.INVALID_INPUT)

    # The inputs from infra-tester are expected to be in the 'metadata`
    # key.
    if "metadata" not in inputs:
        print(
            "ERROR: (infra-tester-plugins) ",
            "Inputs must contain a 'metadata' key.",
        )

        # Exit with a non-zero exit code to indicate failure.
        return int(ExitCodes.INVALID_INPUT)

    # We don't care about the other keys in the inputs for now.
    inputs = inputs["metadata"]

    try:
        state = json.loads(args.state) if args.state is not None else None
    except json.JSONDecodeError as e:
        print(
            "ERROR: (infra-tester-plugins) ",
            f"Failure while parsing state: {e}.",
        )

        # Exit with a non-zero exit code to indicate failure.
        return int(ExitCodes.INVALID_INPUT)

    # We want to redirect stdout to stderr so that the output of the
    # plugin is not captured by the CLI. This is because infra-tester
    # communicates with the plugin via JSON. If the plugin prints
    # anything to stdout, it will be captured by the CLI and the JSON
    # string returned by the plugin will be invalid.
    with contextlib.redirect_stdout(sys.stderr):
        try:
            print(
                "INFO: (infra-tester-plugins) ",
                f"Running action {args.action}."
            )

            print(
                f"INFO: ------------------- {args.name}"
                f"-------------------{os.linesep}"
            )

            if args.action == "validate_inputs":
                return_value = assertion.validate_inputs(inputs)
            elif args.action == "run_assertion":
                return_value = assertion.run_assertion(inputs, state)
            elif args.action == "cleanup":
                assertion.cleanup(inputs, state)
            else:
                raise ValueError(f"Unknown action '{args.action}'")

            print(
                f"{os.linesep}INFO: ------------------- {args.name}",
                "-------------------",
            )

            print(
                "INFO: (infra-tester-plugins) ",
                "Action {args.action} completed.",
            )
        except Exception as e:
            print(
                "ERROR: (infra-tester-plugins) ",
                f"Failure while running action {args.action}: {e}.",
            )

            # Exit with a non-zero exit code to indicate failure.
            return int(ExitCodes.ERROR)

    plugin_result = PluginResult(return_value)
    print(plugin_result)


def manager_cli():
    """
    Entry point for the manager CLI.
    """

    ensure_python_version()

    parser = argparse.ArgumentParser(
        description="This program is part of infra-tester plugin framework."
    )

    parser.add_argument(
        "-l", "--list", action="store_true", help="List all available plugins."
    )

    parser.argument_default

    args = parser.parse_args()

    if len(sys.argv) == 1:
        parser.print_help()
        sys.exit(1)

    if args.list:
        for entry_point in entry_points().get(PLUGIN_GROUP, []):
            print(entry_point.name)
