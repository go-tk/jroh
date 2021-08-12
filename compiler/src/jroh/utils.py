import json
import os
import pkgutil


def title_case(id: str) -> str:
    return id.replace("-", " ")


def pascal_case(id: str) -> str:
    return id.replace("-", "")


def flat_case(id: str) -> str:
    return id.lower().replace("-", "")


def snake_case(id: str) -> str:
    return id.lower().replace("-", "_")


def macro_case(id: str) -> str:
    return id.upper().replace("-", "_")


def camel_case(id: str) -> str:
    words = id.split("-")
    return words[0].lower() + "".join(words[1:])


def quote(s: str) -> str:
    return json.dumps(s)


def get_data(file_path: str) -> bytes:
    data = pkgutil.get_data("jroh", file_path)
    if data is not None:
        return data
    file_path = os.path.join(os.path.dirname(__file__), file_path)
    with open(file_path, "rb") as f:
        data = f.read()
    return data
