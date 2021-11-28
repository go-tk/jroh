import re
from dataclasses import dataclass
from typing import Any, Optional, Union

DEFAULT = "Default"

_WORD_PATTERN = re.compile(r"[A-Z]([A-Z0-9]*|[a-z0-9]*)s?")
ID_PATTERN = re.compile(
    r"{}(-{})*".format(_WORD_PATTERN.pattern, _WORD_PATTERN.pattern)
)

REF_PATTERN = re.compile(
    r"({}\.)?{}".format(f"({ID_PATTERN.pattern})", f"({ID_PATTERN.pattern})")
)

BOOL = "bool"
INT32 = "int32"
INT64 = "int64"
FLOAT32 = "float32"
FLOAT64 = "float64"
STRING = "string"

FIELD_TYPE_PATTERN = re.compile(
    r"|".join(
        (BOOL, INT32, INT64, FLOAT32, FLOAT64, STRING, f"({REF_PATTERN.pattern})")
    )
)

STRUCT = "struct"
ENUM = "enum"
XPRIMIT = "xprimit"

MODEL_TYPE_PATTERN = re.compile(
    r"|".join((STRUCT, ENUM, INT32, INT64, FLOAT32, FLOAT64, STRING))
)

ENUM_UNDERLYING_TYPE_PATTERN = re.compile(r"|".join((INT32, INT64, STRING)))


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
        self.rpc_path_template: str = "/rpc/{namespace}.{service_id}.{method_id}"

        # resolution
        self.methods: list[Method] = []
        self.rpc_paths: list[str] = []


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

    def xprimit(self) -> "Xprimit":
        assert isinstance(self.definition, Xprimit)
        return self.definition


class PrimitiveConstraints:
    def __init__(self) -> None:
        # parse
        self.min: Any = None  # only used for INT32, INT64, FLOAT32, FLOAT64
        self.min_is_exclusive: bool = False  # only used for FLOAT32, FLOAT64
        self.max: Any = None  # only used for INT32, INT64, FLOAT32, FLOAT64
        self.max_is_exclusive: bool = False  # only used for FLOAT32, FLOAT64
        self.min_length: int = 0  # only used for STRING
        self.max_length: Optional[int] = None  # only used for STRING
        self.pattern: str = ""  # only used for STRING

    def is_limited(self) -> bool:
        return (
            self.min is not None
            or self.max is not None
            or self.min_length >= 1
            or self.max_length is not None
            or self.pattern != ""
        )


class Struct:
    def __init__(self) -> None:
        # parse
        self.fields: list[Field] = []


class Field(PrimitiveConstraints):
    def __init__(self, node_uri: str, id: str) -> None:
        super().__init__()

        # parse
        self.node_uri: str = node_uri
        self.id: str = id

        self.type: FieldType = FieldType()
        self.is_optional: bool = False
        self.is_repeated: bool = False
        self.min_count: int = 0  # only used if is_repeated
        self.max_count: Optional[int] = None  # only used if is_repeated
        self.description: Optional[str] = None
        self.example: Any = None

    def count_is_limited(self) -> bool:
        return self.min_count >= 1 or self.max_count is not None


class FieldType:
    def __init__(self) -> None:
        # parse
        self.value: Union[str, Ref, None] = None

        # resolution
        self.model: Optional[Model] = None

    def is_primitive(self) -> bool:
        return isinstance(self.value, str)

    def primitive_type(self) -> str:
        assert isinstance(self.value, str)
        return self.value

    def model_ref(self) -> "Ref":
        assert isinstance(self.value, Ref)
        return self.value


class Enum:
    def __init__(self) -> None:
        # parse
        self.underlying_type: str = ""
        self.constants: list[Constant] = []


class Constant:
    def __init__(self, node_uri: str, id: str) -> None:
        # parse
        self.node_uri: str = node_uri
        self.id: str = id

        self.value: Any = None
        self.description: Optional[str] = None


class Xprimit(PrimitiveConstraints):
    def __init__(self, primitive_type: str) -> None:
        super().__init__()

        # parse
        self.primitive_type = primitive_type
        self.example: Any = None


class Error:
    def __init__(self, node_uri: str, id: str) -> None:
        # parse
        self.node_uri: str = node_uri
        self.id: str = id

        self.code: int = 0
        self.status_code: int = 0
        self.description: Optional[str] = None

        # resolution
        self.ref_count: int = 0


@dataclass
class Ref:
    namespace: Optional[str]
    id: str
