import sys
from typing import Any, Dict, Union


class BaseAssertionPlugin(object):
    def validate_inputs(self, inputs: Dict[Any, Any]) -> Union[str, None]:
        """
        Validate the inputs provided to the plugin. The inputs
        will be in the form of a dictionary with the key being the
        name of the input and the value being the value of the input.

        This method should return None if the inputs are valid.
        Otherwise, it should return a string describing the error.

        Any exception thrown by this method will be treated as an
        implementation error and will be logged as such.

        Args:
            inputs (Dict[str, object]): The inputs provided to the
            plugin.

        Raises:
            NotImplementedError: If the plugin does not implement
            this method.

        Returns:
            Union[str, None]: If the inputs are valid, return None.
            Otherwise, return a string describing the error.
        """

        # It's better to encourage input validation by the plugin
        # developers by forcing them to implement this method.

        raise NotImplementedError(
            f"Plugin must implement the method \
              '{sys._getframe().f_code.co_name}'"
        )

    def run_assertion(
        self, inputs: Dict[Any, Any], state: Dict[Any, Any]
    ) -> Union[str, None]:
        """
        This method should contain the logic to run the assertion
        and return the result.

        If the assertion fails, this method should return a
        string describing the error. Otherwise, it should return
        None.

        This method should not raise any exception. Any exception
        thrown by this method will be treated as an implementation
        error and will be logged as such.

        Args:
            inputs (Dict[str, object]): The inputs provided to the
            plugin.

            state (Dict[str, object]): The current Terraform state
            as a dictionary.

        Raises:
            NotImplementedError: If the plugin does not implement
            this method.

        Returns:
            Union[str, None]: None if the assertion passes. Otherwise
            return a string describing why the assertion failed.
        """

        raise NotImplementedError(
            f"Plugin must implement the method \
             '{sys._getframe().f_code.co_name}'"
        )

    def cleanup(self, inputs: Dict[Any, Any], state: Dict[Any, Any]):
        """
        Cleanup any resources used by the plugin. This method
        will be called after the plugin has been run regardless
        of whether the run was successful or not. It is optional
        and does not need to be implemented if there's nothing to
        clean up.

        If this method is implemented, it should be
        idempotent, i.e, cleanup should always have the same end
        result regardless of how many times it is called. Cleanup
        is considered successful if it does not raise any exception.

        Args:
            inputs (Dict[str, object]): The inputs provided to the
            plugin.

            state (Dict[str, object]): The current Terraform state
            as a dictionary.
        """

        pass
