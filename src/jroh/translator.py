from dataclasses import dataclass
from typing import Optional

import yaml

from . import utils
from .spec import (
    DEFAULT,
    FIELD_BOOL,
    FIELD_FLOAT32,
    FIELD_FLOAT64,
    FIELD_INT32,
    FIELD_INT64,
    FIELD_STRING,
    MODEL_ENUM,
    MODEL_STRUCT,
    Constant,
    Enum,
    Error,
    ErrorCase,
    Field,
    Method,
    Model,
    Params,
    Ref,
    Results,
    Service,
    Spec,
    Struct,
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
            if spec.namespace == DEFAULT:
                self._common_file_path = _COMMON_YAML
            else:
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
                self._save_schemas(schemas, open_api)
            for file_name, open_api in open_apis.items():
                file_path = utils.snake_case(spec.namespace) + "/" + file_name
                self._fix_dollar_refs(open_api)
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
        for method in service.methods:
            path = "/" + service.rpc_path_template.lstrip("/").format(
                namespace=utils.pascal_case(self._namespace),
                service_id=utils.pascal_case(service.id),
                method_id=utils.pascal_case(method.id),
            )
            operation = {}
            paths[path] = {
                "post": operation,
            }
            self._translate_method(method, operation)

    def _translate_method(self, method: Method, operation: dict) -> None:
        operation["operationId"] = utils.camel_case(method.id)
        if method.summary is not None:
            operation["summary"] = method.summary
        if method.description is not None:
            operation["description"] = method.description
        if method.params is not None:
            operation["requestBody"] = self._emit_request_body(method.params, method.id)
        error_cases = [_INTERNAL_ERROR_CASE]
        if method.params is not None:
            error_cases.append(_PARSE_ERROR_CASE)
            error_cases.append(_INVALID_PARAMS_ERROR_CASE)
        error_cases.extend(method.error_cases)
        error_cases.sort(key=lambda x: x.error.code)
        operation["responses"] = self._emit_responses(
            error_cases, method.results, method.id
        )

    def _emit_request_body(self, params: Params, method_id: str) -> dict:
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

    def _translate_params(self, params: Params, method_id: str, schema: dict) -> None:
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
                "content": {
                    "application/json": {
                        "schema": schema,
                    },
                },
            },
        }
        self._emit_resp(results, method_id, schema)
        return responses

    def _emit_resp(
        self, results: Optional[Results], method_id: str, schema: dict
    ) -> None:
        schema_id = utils.camel_case(method_id) + "Resp"
        schema["$ref"] = _MODELS_YAML + "#/components/schemas/" + schema_id
        if schema_id in self._schemas:
            return
        properties = {
            "id": {
                "type": "string",
                "description": "A String generated by the server to identify the RPC.",
                "example": "9m4e2mr0ui3e8a215n4g",
            },
            "error": {
                "$ref": self._common_file_path + "#/components/schemas/error",
                "description": "The error encountered."
                + (
                    ""
                    if results is None
                    else " This field is mutually exclusive of the `results` field."
                ),
            },
        }
        schema2 = {
            "type": "object",
            "properties": properties,
            "required": ["id"],
        }
        self._schemas[schema_id] = schema2
        if results is not None:
            schema3 = {}
            properties["results"] = schema3
            self._translate_results(results, method_id, schema3)

    def _translate_results(
        self, results: Results, method_id: str, schema: dict
    ) -> None:
        schema_id = utils.camel_case(method_id) + "Results"
        schema["$ref"] = "#/components/schemas/" + schema_id
        schema[
            "description"
        ] = "The results returned. This field is mutually exclusive of the `error` field."
        schema2 = {}
        self._schemas[schema_id] = schema2
        self._translate_fields(results.fields, schema2)

    def _translate_models(self, models: list[Model], schemas: dict[str, dict]) -> None:
        for model in models:
            schema_id = utils.camel_case(model.id)
            schema = {}
            schemas[schema_id] = schema
            if model.type == MODEL_STRUCT:
                self._translate_struct(model.struct(), schema)
            elif model.type == MODEL_ENUM:
                self._translate_enum(model.enum(), schema)
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
            property = {}
            properties[property_id] = property
            self._translate_field(field, property)
            field_type = field.type
            if not field_type.is_optional:
                required_property_ids.append(property_id)
        if len(required_property_ids) >= 1:
            schema["required"] = required_property_ids

    def _translate_field(self, field: Field, property: dict) -> None:
        field_type = field.type
        if field_type.is_repeated:
            property["type"] = "array"
            property2 = {}
            property["items"] = property2
        else:
            property2 = property
        if (model_ref := field_type.model_ref) is None:
            property2.update(
                {
                    FIELD_BOOL: {"type": "boolean"},
                    FIELD_INT32: {"type": "integer", "format": "int32"},
                    FIELD_INT64: {"type": "integer", "format": "int64"},
                    FIELD_FLOAT32: {"type": "number", "format": "float"},
                    FIELD_FLOAT64: {"type": "number", "format": "double"},
                    FIELD_STRING: {"type": "string"},
                }[field_type.value]
            )
        else:
            namespace = model_ref.namespace
            if namespace is None:
                namespace = self._namespace
            if namespace == self._namespace:
                schemas_file_path = ""
            else:
                schemas_file_path = (
                    "../" + utils.snake_case(namespace) + "/" + _MODELS_YAML
                )
            property2["$ref"] = (
                schemas_file_path
                + "#/components/schemas/"
                + utils.camel_case(model_ref.id)
            )
        if field_type.model_ref is None:
            model = None
        else:
            model = field_type.model
            assert model is not None
        description_parts = []
        if field.description is None:
            if model is not None and model.description is not None:
                description_parts.append(model.description)
        else:
            description_parts.append(field.description)
        if (
            model is not None
            and model.type == MODEL_ENUM
            and len(constants := model.enum().constants) >= 1
        ):
            markdown = _translate_constants(constants)
            description_parts.append("Constants:\n\n" + markdown)
        if len(description_parts) >= 1:
            if field.description is None:
                property3 = property2
            else:
                property3 = property
            property3["description"] = "\n\n".join(description_parts)
        if field.example is not None:
            property["example"] = field.example

    def _translate_enum(self, enum: Enum, schema: dict) -> None:
        schema.update(
            {
                FIELD_INT32: {"type": "integer", "format": "int32"},
                FIELD_INT64: {"type": "integer", "format": "int32"},
                FIELD_STRING: {"type": "string"},
            }[enum.underlying_type]
        )
        values = []
        for constant in enum.constants:
            values.append(constant.value)
        if len(values) >= 1:
            schema["enum"] = values

    def _save_schemas(self, schemas: dict[str, dict], open_api: dict) -> None:
        open_api["openapi"] = "3.0.0"
        open_api["info"] = {
            "title": "Models",
            "version": "",
        }
        open_api["paths"] = {}
        open_api["components"] = {"schemas": schemas}

    def _fix_dollar_refs(self, open_api: dict) -> None:
        def walk_mapping(m: dict):
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
                    "code": {
                        "type": "integer",
                        "format": "int32",
                        "description": "A Number that indicates the error type that occurred.",
                        "example": 10001,
                    },
                    "message": {
                        "type": "string",
                        "description": "A String providing a short description of the error.",
                        "example": "something wrong",
                    },
                    "data": {
                        "type": "object",
                        "additionalProperties": True,
                        "description": "A Structured value that contains additional information about the error.",
                    },
                },
                "required": ["code", "message"],
            },
        },
    },
}

_PARSE_ERROR_CASE = ErrorCase("", Ref(namespace="", id=""))
_PARSE_ERROR_CASE.error = Error("", "Parse-Error")
_PARSE_ERROR_CASE.error.code = -32700
_PARSE_ERROR_CASE.error.description = "Invalid JSON was received by the server."

_INVALID_PARAMS_ERROR_CASE = ErrorCase("", Ref(namespace="", id=""))
_INVALID_PARAMS_ERROR_CASE.error = Error("", "Invalid-Params")
_INVALID_PARAMS_ERROR_CASE.error.code = -32602
_INVALID_PARAMS_ERROR_CASE.error.description = "Invalid method parameter(s)."

_INTERNAL_ERROR_CASE = ErrorCase("", Ref(namespace="", id=""))
_INTERNAL_ERROR_CASE.error = Error("", "Internal-Error")
_INTERNAL_ERROR_CASE.error.code = -32603
_INTERNAL_ERROR_CASE.error.description = "Internal JSON-RPC error."


def _translate_error_cases(error_cases: list[ErrorCase]) -> str:
    lines = []
    lines.append("| Code | Message | Description |")
    lines.append("| - | - | - |")
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
        lines.append(f"| {error.code} | {message} | {description} |")
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
