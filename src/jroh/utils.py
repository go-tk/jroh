def title_case(id: str) -> str:
    return id.replace("-", " ")


def pascal_case(id: str) -> str:
    return id.replace("-", "")


def snake_case(id: str) -> str:
    return id.lower().replace("-", "_")


def macro_case(id: str) -> str:
    return id.upper().replace("-", "_")


def camel_case(id: str) -> str:
    words = id.split("-")
    return words[0].lower() + "".join(words[1:])
