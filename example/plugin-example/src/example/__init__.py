from typing import Dict

from infra_tester_plugins import BaseAssertionPlugin


class ExampleAssertionPlugin(BaseAssertionPlugin):
    def validate_inputs(self, inputs: Dict[str, object]):
        print("Running validate_inputs from ExampleAssertionPlugin")
        print("Inputs:", inputs)

        if (
            "should_fail_validation" in inputs
            and inputs["should_fail_validation"]
        ):
            return "This is a validation error message."

        return None

    def run_assertion(
        self, inputs: Dict[str, object], state: Dict[str, object]
    ):
        print("Running run_assertion from ExampleAssertionPlugin")
        print("Inputs:", inputs)
        print("State:", state)

        if (
            "should_error" in inputs
            and inputs["should_error"]
        ):
            message = "This is an assertion error message."
            if "custom_message" in state and state["custom_message"]:
                message += f" Plus here's a custom message: \
                            {state['custom_message']}"

            return message

        return None

    def cleanup(self, inputs: Dict[str, object], state: Dict[str, object]):
        print("Running cleanup from ExampleAssertionPlugin")
        print("Inputs:", inputs)
        print("State:", state)


def load_plugin():
    return ExampleAssertionPlugin()
