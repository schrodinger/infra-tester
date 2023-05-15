from typing import Dict
from urllib import error, request
from urllib.parse import urlparse

from infra_tester_plugins import BaseAssertionPlugin


class URLReachableAssertionPlugin(BaseAssertionPlugin):
    FIELD_URL = "url"
    FIELD_STATUS_CODE = "status_code"
    FIELD_FROM_OUTPUTS = "from_outputs"

    @classmethod
    def validate_url(cls, plugin_inputs: Dict[str, object]):
        if cls.FIELD_URL not in plugin_inputs:
            return "url is a required input."

        if cls.FIELD_FROM_OUTPUTS in plugin_inputs:
            from_outputs = plugin_inputs[cls.FIELD_FROM_OUTPUTS]

            if not isinstance(from_outputs, bool):
                return "from_outputs must be a boolean."

            if from_outputs:
                return None

        url = plugin_inputs[cls.FIELD_URL]

        if not isinstance(url, str):
            return "url must be a string."

        try:
            parse_res = urlparse(url)

            if not parse_res.scheme:
                return "url must have a scheme."

            if not parse_res.netloc:
                return "url must have a valid network location."
        except ValueError:
            return "url must be a valid."

        return None

    @classmethod
    def validate_status_code(cls, plugin_inputs: Dict[str, object]):
        # status_code is optional
        if cls.FIELD_STATUS_CODE not in plugin_inputs:
            return None

        status_code = plugin_inputs[cls.FIELD_STATUS_CODE]

        if not isinstance(status_code, int):
            return "status_code must be an integer."

        # Not being strict about the status code range here.
        if status_code < 0:
            return (
                f"status_code must be valid. "
                f"Received {status_code} of type {type(status_code)}."
            )

        return None

    def validate_inputs(self, inputs: Dict[str, object]):
        print("Running validate_inputs for URLReachable")

        url_validation = self.validate_url(inputs)
        if url_validation:
            return url_validation

        status_code_validation = self.validate_status_code(inputs)
        if status_code_validation:
            return status_code_validation

        return None

    def assert_url_reachable(self, url: str, expected_status_code: int):
        try:
            res = request.urlopen(url)

            if res.getcode() != expected_status_code:
                return (
                    f"Unexpected status code for {url}: {res.getcode()}. "
                    f"Expected {expected_status_code}."
                )
        except error.HTTPError as e:
            if e.code == expected_status_code:
                return None

            return (
                f"Unexpected status code for {url}: {e.code}. "
                f"Expected {expected_status_code}."
            )
        except error.URLError as e:
            return f"Failed to reach {url}: {e.reason}."
        except Exception as e:
            return f"Failed to reach {url}: {e}"

        return None

    def get_url_from_outputs(self, output_name: str, state: Dict[str, object]):
        if "values" not in state or "outputs" not in state["values"]:
            raise ValueError(
                (
                    "Could not find 'values' in state. "
                    "Make sure terraform apply has been "
                    "run successfully."
                )
            )

        outputs = state["values"]["outputs"]
        if output_name not in outputs:
            raise ValueError(
                (
                    f"Could not find {output_name} in outputs. "
                    "Make sure your terraform code has an "
                    f"output named {output_name}."
                )
            )

        return outputs[output_name]["value"]

    def run_assertion(
        self, inputs: Dict[str, object], state: Dict[str, object]
    ):
        print("Running run_assertion for URLReachable")

        plugin_inputs = inputs
        url = plugin_inputs[self.FIELD_URL]
        from_outputs = plugin_inputs.get(self.FIELD_FROM_OUTPUTS, False)

        if from_outputs:
            url = self.get_url_from_outputs(url, state)

        expected_status_code = inputs.get(
            self.FIELD_STATUS_CODE, 200
        )

        return self.assert_url_reachable(url, expected_status_code)


def load_plugin():
    return URLReachableAssertionPlugin()
