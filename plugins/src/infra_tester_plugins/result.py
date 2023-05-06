import json
from typing import Union


class PluginResult():
    def __init__(self, message: Union[str, None] = None) -> None:
        self.error = message is not None
        self.message = message

    def __str__(self) -> str:
        return json.dumps({
            "error": self.error,
            "message": self.message
        }, indent=4)

    def __repr__(self) -> str:
        return self.__str__()
