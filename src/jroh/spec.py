import re
from dataclasses import dataclass
from typing import Any, Optional

_WORD_PATTERN = re.compile(r"[A-Z]([A-Z0-9]*|[a-z0-9]*)s?")
ID_PATTERN = re.compile(
    r"{}(-{})*".format(_WORD_PATTERN.pattern, _WORD_PATTERN.pattern)
)

DEFAULT = "Default"
NAMESPACE_PATTERN = re.compile(DEFAULT + r"|" + ID_PATTERN.pattern)

REF_PATTERN = re.compile(
    r"({}\.)?{}".format(f"({NAMESPACE_PATTERN.pattern})", f"({ID_PATTERN.pattern})")
)

FIELD_BOOL = "bool"
FIELD_INT32 = "int32"
FIELD_INT64 = "int64"
FIELD_FLOAT32 = "float32"
FIELD_FLOAT64 = "float64"
FIELD_STRING = "string"
FIELD_TYPE_PATTERN = re.compile(
    r"({}|{})(\?|\+|\*)?".format(
        r"|".join(
            (
                FIELD_BOOL,
                FIELD_INT32,
                FIELD_INT64,
                FIELD_FLOAT32,
                FIELD_FLOAT64,
                FIELD_STRING,
            )
        ),
        f"({REF_PATTERN.pattern})",
    )
)

MODEL_STRUCT = "struct"
MODEL_ENUM = "enum"
MODEL_TYPE_PATTERN = re.compile(r"{}".format(r"|".join((MODEL_STRUCT, MODEL_ENUM))))

ENUM_UNDERLYING_TYPE_PATTERN = re.compile(
    r"|".join((FIELD_INT32, FIELD_INT64, FIELD_STRING))
)


class Spec:
    def __init__(self, node_uri: str) -> None:
        # parse
        self.node_uri: str = node_uri

        self.namespace: str = DEFAULT
        self.services: list[Service] = []
        self.methods: list[Method] = []
        self.models: list[Model] = []
        self.errors: list[Error] = []


class Service:
    def __init__(self, node_uri: str, id: str) -> None:
        # parse
        self.node_uri: str = node_uri
        self.id: str = id

        self.version: str = ""
        self.description: Optional[str] = None
        self.rpc_path_template: str = "/rpc/{service_id}.{method_id}"

        # resolution
        self.methods: list[Method] = []


class Method:
    def __init__(self, node_uri: str, id: str) -> None:
        # parse
        self.node_uri: str = node_uri
        self.id: str = id

        self.service_ids: list[str] = []
        self.summary: Optional[str] = None
        self.description: Optional[str] = None
        self.params: Optional[Params] = None
        self.results: Optional[Results] = None
        self.error_cases: list[ErrorCase] = []


class Params:
    def __init__(self, node_uri: str) -> None:
        # parse
        self.node_uri: str = node_uri

        self.fields: list[Field] = []


class Results:
    def __init__(self, node_uri: str) -> None:
        # parse
        self.node_uri: str = node_uri

        self.fields: list[Field] = []


class ErrorCase:
    def __init__(self, node_uri: str, error_ref: "Ref") -> None:
        # parse
        self.node_uri: str = node_uri
        self.error_ref: Ref = error_ref

        self.description: Optional[str] = None

        # resolution
        self.error: Optional[Error] = None


class Model:
    def __init__(self, node_uri: str, id: str) -> None:
        # parse
        self.node_uri: str = node_uri
        self.id: str = id

        self.type: str = ""
        self.description: Optional[str] = None
        self.definition: Any = None

        # resolution
        self.namespace: str = ""
        self.ref_count: int = 0

    def struct(self) -> "Struct":
        assert isinstance(self.definition, Struct)
        return self.definition

    def enum(self) -> "Enum":
        assert isinstance(self.definition, Enum)
        return self.definition


class Struct:
    def __init__(self) -> None:
        # parse
        self.fields: list[Field] = []


class Field:
    def __init__(self, node_uri: str, id: str) -> None:
        # parse
        self.node_uri: str = node_uri
        self.id: str = id

        self.type: FieldType = FieldType("")
        self.description: Optional[str] = None
        self.example: Any = None


class FieldType:
    def __init__(self, node_uri: str) -> None:
        # parse
        self.node_uri = node_uri

        self.is_optional: bool = False
        self.is_repeated: bool = False
        self.model_ref: Optional[Ref] = None
        self.value: str = ""

        # resolution
        self.model: Optional[Model] = None


class Enum:
    def __init__(self) -> None:
        # parse
        self.underlying_type: str
        self.constants: list[Constant] = []


class Constant:
    def __init__(self, node_uri: str, id: str) -> None:
        # parse
        self.node_uri: str = node_uri
        self.id: str = id

        self.value: Any = None
        self.description: Optional[str] = None


class Error:
    def __init__(self, node_uri: str, id: str) -> None:
        # parse
        self.node_uri: str = node_uri
        self.id: str = id

        self.code: int = 0
        self.description: Optional[str] = None

        # resolution
        self.ref_count: int = 0


@dataclass
class Ref:
    namespace: Optional[str]
    id: str
