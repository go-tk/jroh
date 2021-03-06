from dataclasses import dataclass
from typing import Optional

import yaml

from . import utils
from .spec import (
    BOOL,
    ENUM,
    FLOAT32,
    FLOAT64,
    INT32,
    INT64,
    STRING,
    STRUCT,
    XPRIMIT,
    Constant,
    Enum,
    Error,
    ErrorCase,
    Field,
    Method,
    Model,
    Params,
    PrimitiveConstraints,
    Ref,
    Results,
    Service,
    Spec,
    Struct,
    Xprimit,
)


def _str_representer(dumper: yaml.Dumper, data: str):
    if "\n" in data:
        return dumper.represent_scalar("tag:yaml.org,2002:str", data, style="|")
    return dumper.represent_scalar("tag:yaml.org,2002:str", data)


yaml.add_representer(str, _str_representer)


@dataclass
class TranslateSpecsResults:
    file_path_2_file_data: dict[str, str]


def translate_specs(
    specs: list[Spec], test_mode: bool = False
) -> TranslateSpecsResults:
    translateer = _Translator(test_mode)
    translateer.translate_specs(specs)
    return TranslateSpecsResults(
        file_path_2_file_data=translateer.file_path_2_file_data(),
    )


class _Translator:
    def __init__(self, test_mode: bool) -> None:
        self._test_mode = test_mode
        self._namespace: str = ""
        self._common_file_path: str = ""
        self._schemas: dict[str, dict] = {}
        self._file_path_2_file_data: dict[str, str] = {}

    def translate_specs(self, specs: list[Spec]) -> None:
        for spec in specs:
            self._namespace = spec.namespace
            self._common_file_path = "../" + _COMMON_YAML
            open_apis: dict[str, dict] = {}
            schemas = {}
            self._schemas = schemas
            for service in spec.services:
                open_api = {}
                file_name = _SERVICE_YAML_TEMPLATE.format(utils.snake_case(service.id))
                open_apis[file_name] = open_api
                self._translate_service(service, open_api)
            self._translate_models(spec.models, schemas)
            if len(schemas) >= 1:
                open_api = {}
                open_apis[_MODELS_YAML] = open_api
                _save_schemas(schemas, open_api)
            for file_name, open_api in open_apis.items():
                file_path = utils.snake_case(spec.namespace) + "/" + file_name
                _fix_dollar_refs(open_api)
                self._file_path_2_file_data[file_path] = yaml.dump(
                    open_api, sort_keys=False
                )
        if not self._test_mode:
            self._file_path_2_file_data[_COMMON_YAML] = yaml.dump(
                _COMMON_OPEN_API, sort_keys=False
            )

    def _translate_service(self, service: Service, open_api: dict) -> None:
        open_api["openapi"] = "3.0.0"
        info = {
            "title": utils.title_case(service.id) + " Service",
            "version": service.version,
        }
        open_api["info"] = info
        if service.description is not None:
            info["description"] = service.description
        paths = {}
        open_api["paths"] = paths
        for i, method in enumerate(service.methods):
            rpc_path = service.rpc_paths[i]
            operation = {}
            paths[rpc_path] = {
                "post": operation,
            }
            self._translate_method(method, operation)

    def _translate_method(self, method: Method, operation: dict) -> None:
        operation["operationId"] = utils.camel_case(method.id)
        if method.summary is not None:
            operation["summary"] = method.summary
        if method.description is not None:
            operation["description"] = method.description
        operation["requestBody"] = self._emit_request_body(method.params, method.id)
        error_cases = [_NOT_IMPLEMENTED_ERROR_CASE]
        if method.params is not None:
            error_cases.append(_INVALID_PARAMS_ERROR_CASE)
        error_cases.extend(method.error_cases)

        def get_error_code(error_case: ErrorCase) -> int:
            error = error_case.error
            assert error is not None
            return error.code

        error_cases.sort(key=get_error_code)
        operation["responses"] = self._emit_responses(
            error_cases, method.results, method.id
        )

    def _emit_request_body(self, params: Optional[Params], method_id: str) -> dict:
        schema = {}
        request_body = {
            "content": {
                "application/json": {
                    "schema": schema,
                },
            },
        }
        self._translate_params(params, method_id, schema)
        return request_body

    def _translate_params(
        self, params: Optional[Params], method_id: str, schema: dict
    ) -> None:
        if params is None:
            schema["type"] = "object"
            return
        schema_id = utils.camel_case(method_id) + "Params"
        schema["$ref"] = _MODELS_YAML + "#/components/schemas/" + schema_id
        if schema_id in self._schemas:
            return
        schema2 = {}
        self._schemas[schema_id] = schema2
        self._translate_fields(params.fields, schema2)

    def _emit_responses(
        self, error_cases: list[ErrorCase], results: Optional[Results], method_id: str
    ) -> dict:
        markdown = _translate_error_cases(error_cases)
        description_parts = ["## Error Cases\n\n" + markdown]
        schema = {}
        responses = {
            "200": {
                "description": "\n\n".join(description_parts),
                "headers": {
                    "Jroh-Trace-Id": {
                        "description": "The trace identifier.",
                        "schema": {
                            "type": "string",
                        },
                    },
                    "Jroh-Error-Code": {
                        "description": "The error code. This header is present only if error occurs.",
                        "schema": {
                            "type": "int32",
                        },
                    },
                },
                "content": {
                    "application/json": {
                        "schema": {
                            "oneOf": [
                                schema,
                                {
                                    "$ref": self._common_file_path
                                    + "#/components/schemas/error"
                                },
                            ]
                        },
                    },
                },
            },
        }
        self._translate_results(results, method_id, schema)
        return responses

    def _translate_results(
        self, results: Optional[Results], method_id: str, schema: dict
    ) -> None:
        if results is None:
            schema["type"] = "object"
            return
        schema_id = utils.camel_case(method_id) + "Results"
        schema["$ref"] = _MODELS_YAML + "#/components/schemas/" + schema_id
        if schema_id in self._schemas:
            return
        schema2 = {}
        self._schemas[schema_id] = schema2
        self._translate_fields(results.fields, schema2)

    def _translate_models(self, models: list[Model], schemas: dict[str, dict]) -> None:
        for model in models:
            schema_id = utils.camel_case(model.id)
            schema = {}
            schemas[schema_id] = schema
            if model.type == STRUCT:
                self._translate_struct(model.struct(), schema)
            elif model.type == ENUM:
                self._translate_enum(model.enum(), schema)
            elif model.type == XPRIMIT:
                self._translate_xprimit(model.xprimit(), schema)
            else:
                assert False, model.type

    def _translate_struct(self, struct: Struct, schema: dict) -> None:
        self._translate_fields(struct.fields, schema)

    def _translate_fields(self, fields: list[Field], schema: dict) -> None:
        schema["type"] = "object"
        if len(fields) == 0:
            return
        properties = {}
        schema["properties"] = properties
        required_property_ids = []
        for field in fields:
            property_id = utils.camel_case(field.id)
            schema2 = {}
            properties[property_id] = schema2
            self._translate_field(field, schema2)
            if not (field.is_optional or (field.is_repeated and field.min_count == 0)):
                required_property_ids.append(property_id)
        if len(required_property_ids) >= 1:
            schema["required"] = required_property_ids

    def _translate_field(self, field: Field, schema: dict) -> None:
        if field.is_repeated:
            schema["type"] = "array"
            if field.min_count >= 1:
                schema["minItems"] = field.min_count
            if field.max_count is not None:
                schema["maxItems"] = field.max_count
            schema2 = {}
            schema["items"] = schema2
        else:
            schema2 = schema
        field_type = field.type
        if field_type.is_primitive():
            _translate_primitive_type_and_constraints(
                field_type.primitive_type(), field, schema2
            )
        else:
            model_ref = field_type.model_ref()
            namespace = model_ref.namespace
            if namespace is None:
                namespace = self._namespace
            if namespace == self._namespace:
                schemas_file_path = ""
            else:
                schemas_file_path = (
                    "../" + utils.snake_case(namespace) + "/" + _MODELS_YAML
                )
            schema2["$ref"] = (
                schemas_file_path
                + "#/components/schemas/"
                + utils.camel_case(model_ref.id)
            )
        if field_type.is_primitive():
            model = None
        else:
            model = field_type.model
            assert model is not None
        description_parts = []
        if field.description is None:
            if model is not None and model.description is not None:
                description_parts.append(model.description)
        else:
            description_part = field.description
            if description_part.startswith("+"):
                description_part = description_part[1:]
                if model is not None and model.description is not None:
                    description_part = model.description + description_part
            description_parts.append(description_part)
        if (
            model is not None
            and model.type == ENUM
            and len(constants := model.enum().constants) >= 1
        ):
            markdown = _translate_constants(constants)
            description_parts.append("Constants:\n\n" + markdown)
        if len(description_parts) >= 1:
            schema2["description"] = "\n\n".join(description_parts)
        if field.example is not None:
            schema["example"] = field.example

    def _translate_enum(self, enum: Enum, schema: dict) -> None:
        schema.update(
            {
                INT32: {"type": "integer", "format": "int32"},
                INT64: {"type": "integer", "format": "int32"},
                STRING: {"type": "string"},
            }[enum.underlying_type]
        )
        values = []
        for constant in enum.constants:
            values.append(constant.value)
        if len(values) >= 1:
            schema["enum"] = values

    def _translate_xprimit(self, xprimit: Xprimit, schema: dict) -> None:
        _translate_primitive_type_and_constraints(
            xprimit.primitive_type, xprimit, schema
        )
        if xprimit.example is not None:
            schema["example"] = xprimit.example

    def file_path_2_file_data(self) -> dict[str, str]:
        return self._file_path_2_file_data


_SERVICE_YAML_TEMPLATE = "{}_service.yaml"
_MODELS_YAML = "models.yaml"
_COMMON_YAML = "common.yaml"
_COMMON_OPEN_API = {
    "openapi": "3.0.0",
    "info": {
        "title": "Common",
        "version": "",
    },
    "paths": {},
    "components": {
        "schemas": {
            "error": {
                "type": "object",
                "properties": {
                    "message": {
                        "type": "string",
                        "description": "A short description of the error.",
                        "example": "parse error",
                    },
                    "details": {
                        "type": "string",
                        "description": "Detailed information about the error.",
                        "example": "unexpected end of JSON input",
                    },
                    "data": {
                        "type": "object",
                        "additionalProperties": True,
                        "description": "Additional information about the error.",
                        "example": {
                            "key": "value",
                        },
                    },
                },
                "required": ["message"],
            },
        },
    },
}

_NOT_IMPLEMENTED_ERROR_CASE = ErrorCase("", Ref(namespace="", id=""))
_NOT_IMPLEMENTED_ERROR_CASE.error = Error("", "Not-Implemented")
_NOT_IMPLEMENTED_ERROR_CASE.error.code = 1
_NOT_IMPLEMENTED_ERROR_CASE.error.status_code = 501
_NOT_IMPLEMENTED_ERROR_CASE.error.description = "The method is not implemented."

_INVALID_PARAMS_ERROR_CASE = ErrorCase("", Ref(namespace="", id=""))
_INVALID_PARAMS_ERROR_CASE.error = Error("", "Invalid-Params")
_INVALID_PARAMS_ERROR_CASE.error.code = 2
_INVALID_PARAMS_ERROR_CASE.error.status_code = 422
_INVALID_PARAMS_ERROR_CASE.error.description = "Invalid method parameter(s)."


def _translate_primitive_type_and_constraints(
    primitive_type: str, primitive_constraints: PrimitiveConstraints, schema: dict
) -> None:
    schema.update(
        {
            BOOL: {"type": "boolean"},
            INT32: {"type": "integer", "format": "int32"},
            INT64: {"type": "integer", "format": "int64"},
            FLOAT32: {"type": "number", "format": "float"},
            FLOAT64: {"type": "number", "format": "double"},
            STRING: {"type": "string"},
        }[primitive_type]
    )
    if primitive_type in (
        INT32,
        INT64,
        FLOAT32,
        FLOAT64,
    ):
        if primitive_constraints.min is not None:
            schema["minimum"] = primitive_constraints.min
            if (
                primitive_type in (FLOAT32, FLOAT64)
                and primitive_constraints.min_is_exclusive
            ):
                schema["exclusiveMinimum"] = True
        if primitive_constraints.max is not None:
            schema["maximum"] = primitive_constraints.max
            if (
                primitive_type in (FLOAT32, FLOAT64)
                and primitive_constraints.max_is_exclusive
            ):
                schema["exclusiveMaximum"] = True
    elif primitive_type == STRING:
        if primitive_constraints.min_length >= 1:
            schema["minLength"] = primitive_constraints.min_length
        if primitive_constraints.max_length is not None:
            schema["maxLength"] = primitive_constraints.max_length
        if primitive_constraints.pattern != "":
            schema["pattern"] = primitive_constraints.pattern


def _translate_error_cases(error_cases: list[ErrorCase]) -> str:
    lines = []
    lines.append("| Error Code | Status Code | Message | Description |")
    lines.append("| - | - | - | - |")
    lines.append("| -1 | ... | ... | Low-level error. |")
    for error_case in error_cases:
        error = error_case.error
        assert error is not None
        message = error.id.lower().replace("-", " ")
        if error_case.description is None:
            if error.description is None:
                description = ""
            else:
                description = error.description
        else:
            description = error_case.description
            if description.startswith("+"):
                description = description[1:]
                if error.description is not None:
                    description = error.description + description
        lines.append(
            f"| {error.code} | {error.status_code} | {message} | {description} |"
        )
    return "\n".join(lines)


def _translate_constants(constants: list[Constant]) -> str:
    lines = []
    for constant in constants:
        if isinstance(constant.value, int):
            value_repr = str(constant.value)
        elif isinstance(constant.value, str):
            value_repr = '"' + constant.value + '"'
        else:
            assert False, type(constant.value)
        line = f"- {utils.macro_case(constant.id)}({value_repr})"
        if constant.description is not None:
            line += f": {constant.description}"
        lines.append(line)
    return "\n".join(lines)


def _save_schemas(schemas: dict[str, dict], open_api: dict) -> None:
    open_api["openapi"] = "3.0.0"
    open_api["info"] = {
        "title": "Models",
        "version": "",
    }
    open_api["paths"] = {}
    open_api["components"] = {"schemas": schemas}


def _fix_dollar_refs(open_api: dict) -> None:
    def walk_mapping(m: dict) -> None:
        for k, v in m.items():
            if k == "$ref":
                continue
            if isinstance(v, dict):
                walk_mapping(v)
            elif isinstance(v, list):
                for v2 in v:
                    if isinstance(v2, dict):
                        walk_mapping(v2)
        if "$ref" in m and len(m) >= 2:
            dollar_ref = m["$ref"]
            d2 = {k: v for k, v in m.items() if k != "$ref"}
            m.clear()
            m["allOf"] = [{"$ref": dollar_ref}, d2]

    walk_mapping(open_api)
