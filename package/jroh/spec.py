import re
from typing import Any, Optional

_WORD_PATTERN = re.compile(r"[A-Z]([A-Z0-9]*|[a-z0-9]*)s?")
ID_PATTERN = re.compile(
    r"{}(-{})*".format(_WORD_PATTERN.pattern, _WORD_PATTERN.pattern)
)

GLOBAL = "global"
NAMESPACE_PATTERN = re.compile(GLOBAL + r"|" + ID_PATTERN.pattern)

FIELD_BOOL = "bool"
FIELD_INT32 = "int32"
FIELD_INT64 = "int64"
FIELD_FLOAT32 = "float32"
FIELD_FLOAT64 = "float64"
FIELD_STRING = "string"
RAW_FIELD_TYPE_PATTERN = re.compile(
    r"({}|({}\.)?{})(\?|\+|\*)?".format(
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
        NAMESPACE_PATTERN.pattern,
        ID_PATTERN.pattern,
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

        self.namespace: str = GLOBAL
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
        self.method_path_template: str = "/rpc/{service_id}.{method_id}"

        # resolution
        self.methods: list[Method] = []


class Method:
    def __init__(self, node_uri: str, id: str) -> None:
        # parse
        self.node_uri: str = node_uri
        self.id: str = id

        self.service_id: str = ""
        self.summary: Optional[str] = None
        self.description: Optional[str] = None
        self.params: Optional[Params] = None
        self.result: Optional[Result] = None
        self.error_cases: list[ErrorCase] = []


class Params:
    def __init__(self, node_uri: str) -> None:
        # parse
        self.node_uri: str = node_uri

        self.fields: list[Field] = []


class Result:
    def __init__(self, node_uri: str) -> None:
        # parse
        self.node_uri: str = node_uri

        self.fields: list[Field] = []


class ErrorCase:
    def __init__(self, node_uri: str, error_id: str) -> None:
        # parse
        self.node_uri: str = node_uri
        self.error_id: str = error_id

        self.description: Optional[str] = None

        # resolution
        self.error: Optional[Error] = None


class Model:
    def __init__(self, node_uri: str, id: str) -> None:
        # parse
        self.node_uri: str = node_uri
        self.id: str = id

        self.type: str = ""
        self.definition: Any = None

        # resolution
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

        self.value: str = ""
        self.is_model_ref: bool = False
        self.namespace: Optional[str] = None
        self.is_optional: bool = False
        self.is_repeated: bool = False

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

        # resolution
        self.ref_count: int = 0
