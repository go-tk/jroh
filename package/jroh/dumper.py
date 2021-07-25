from dataclasses import dataclass
from typing import Optional

import yaml

from . import utils
from .spec import (FIELD_BOOL, FIELD_FLOAT32, FIELD_FLOAT64, FIELD_INT32,
                   FIELD_INT64, FIELD_STRING, GLOBAL, MODEL_ENUM, MODEL_STRUCT,
                   Constant, Enum, ErrorCase, Field, Method, Model, Params,
                   Result, Service, Spec, Struct)


def _str_representer(dumper: yaml.Dumper, data: str):
    if "\n" in data:
        return dumper.represent_scalar("tag:yaml.org,2002:str", data, style="|")
    return dumper.represent_scalar("tag:yaml.org,2002:str", data)


yaml.add_representer(str, _str_representer)


@dataclass
class DumpSpecsResult:
    file_path_2_file_data: dict[str, str]


def dump_specs(specs: list[Spec]) -> DumpSpecsResult:
    dumper = _Dumper()
    dumper.dump_specs(specs)
    return DumpSpecsResult(
        file_path_2_file_data=dumper.file_path_2_file_data(),
    )


class _Dumper:
    def __init__(self) -> None:
        self._namespace: str = ""
        self._builtins_file_path: str = ""
        self._file_name: str = ""
        self._schemas: dict[str, dict] = {}
        self._file_path_2_file_data: dict[str, str] = {}

    def dump_specs(self, specs: list[Spec]) -> None:
        for spec in specs:
            self._namespace = spec.namespace
            if spec.namespace == GLOBAL:
                self._builtins_file_path = _BUILTINS_YAML
            else:
                self._builtins_file_path = "../" + _BUILTINS_YAML
            open_apis: dict[str, dict] = {}
            for service in spec.services:
                self._file_name = _SERVICE_YAML_TEMPLATE.format(
                    utils.snake_case(service.id)
                )
                open_api = {}
                open_apis[self._file_name] = open_api
                self._dump_service(service, open_api)
            if len(spec.models) >= 1:
                self._file_name = _MODELS_YAML
                open_api = {}
                open_apis[self._file_name] = open_api
                self._dump_models(spec.models, open_api)
            for file_name, open_api in open_apis.items():
                if spec.namespace == GLOBAL:
                    file_path = file_name
                else:
                    file_path = utils.snake_case(spec.namespace) + "/" + file_name
                self._fix_dollar_refs(open_api)
                self._file_path_2_file_data[file_path] = yaml.dump(
                    open_api, sort_keys=False
                )

    def _dump_service(self, service: Service, open_api: dict) -> None:
        open_api["openapi"] = "3.0.0"
        open_api["info"] = {
            "title": utils.title_case(service.id) + " API",
            "version": service.version,
        }
        paths = {}
        open_api["paths"] = paths
        schemas = {}
        self._schemas = schemas
        for method in service.methods:
            path = "/" + service.method_path_template.lstrip("/").format(
                namespace=utils.pascal_case(self._namespace),
                service_id=utils.pascal_case(service.id),
                method_id=utils.pascal_case(method.id),
            )
            operation = {}
            paths[path] = {
                "post": operation,
            }
            self._dump_method(method, operation)
        if len(schemas) >= 1:
            open_api["components"] = {"schemas": schemas}

    def _dump_method(self, method: Method, operation: dict) -> None:
        operation["operationId"] = utils.camel_case(method.id)
        if method.summary is not None:
            operation["summary"] = method.summary
        if method.description is not None:
            operation["description"] = method.description
        if method.params is not None:
            operation["requestBody"] = self._make_request_body(method.params, method.id)
        operation["responses"] = self._make_responses(
            method.error_cases, method.result, method.id
        )

    def _make_request_body(self, params: Params, method_id: str) -> dict:
        schema = {}
        request_body = {
            "content": {
                "application/json": {
                    "schema": schema,
                },
            },
        }
        self._dump_params(params, method_id, schema)
        return request_body

    def _dump_params(self, params: Params, method_id: str, schema: dict) -> None:
        schema_id = utils.camel_case(method_id) + "Params"
        schema["$ref"] = "#/components/schemas/" + schema_id
        schema2 = {}
        self._schemas[schema_id] = schema2
        self._dump_fields(params.fields, schema2)

    def _make_responses(
        self, error_cases: list[ErrorCase], result: Optional[Result], method_id: str
    ) -> dict:
        description_parts = []
        if len(error_cases) >= 1:
            markdown = _dump_error_cases(error_cases)
            description_parts.append("## Error Cases\n\n" + markdown)
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
        self._dump_resp(result, method_id, schema)
        return responses

    def _dump_resp(
        self, result: Optional[Result], method_id: str, schema: dict
    ) -> None:
        if result is None:
            schema["$ref"] = (
                self._builtins_file_path + "#/components/schemas/rpcRespWithoutResult"
            )
            return
        schema_id = utils.camel_case(method_id) + "Resp"
        schema["$ref"] = "#/components/schemas/" + schema_id
        properties = {
            "id": {
                "type": "integer",
                "format": "int64",
                "description": "The RPC identifier generated by the server.",
            },
            "error": {
                "$ref": self._builtins_file_path + "#/components/schemas/rpcError",
                "description": "The RPC error encountered. This field is mutually exclusive of the `result` field.",
            },
        }
        schema2 = {
            "type": "object",
            "properties": properties,
            "required": ["id"],
        }
        self._schemas[schema_id] = schema2
        schema3 = {}
        properties["result"] = schema3
        self._dump_result(result, method_id, schema3)

    def _dump_result(self, result: Result, method_id: str, schema: dict) -> None:
        schema_id = utils.camel_case(method_id) + "Result"
        schema["$ref"] = "#/components/schemas/" + schema_id
        schema[
            "description"
        ] = "The RPC result returned. This field is mutually exclusive of the `error` field."
        schema2 = {}
        self._schemas[schema_id] = schema2
        self._dump_fields(result.fields, schema2)

    def _dump_models(self, models: list[Model], open_api: dict) -> None:
        open_api["openapi"] = "3.0.0"
        open_api["info"] = {
            "title": "Models",
            "version": "",
        }
        schemas = {}
        open_api["components"] = {
            "schemas": schemas,
        }
        for model in models:
            schema_id = utils.camel_case(model.id)
            schema = {}
            schemas[schema_id] = schema
            if model.type == MODEL_STRUCT:
                self._dump_struct(model.struct(), schema)
            elif model.type == MODEL_ENUM:
                self._dump_enum(model.enum(), schema)
            else:
                assert False, model.type

    def _dump_struct(self, struct: Struct, schema: dict) -> None:
        self._dump_fields(struct.fields, schema)

    def _dump_fields(self, fields: list[Field], schema: dict) -> None:
        schema["type"] = "object"
        properties = {}
        schema["properties"] = properties
        required_property_ids = []
        for field in fields:
            property_id = utils.camel_case(field.id)
            property = {}
            properties[property_id] = property
            self._dump_field(field, property)
            field_type = field.type
            if not field_type.is_optional:
                required_property_ids.append(property_id)
        if len(required_property_ids) >= 1:
            schema["required"] = required_property_ids

    def _dump_field(self, field: Field, property: dict) -> None:
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
                    FIELD_INT64: {"type": "integer", "format": "int32"},
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
                if self._file_name == _MODELS_YAML:
                    models_file_path = ""
                else:
                    models_file_path = _MODELS_YAML
            else:
                if namespace == GLOBAL:
                    models_file_path = "../" + _MODELS_YAML
                else:
                    models_file_path = (
                        "../" + utils.snake_case(namespace) + "/" + _MODELS_YAML
                    )
            property2["$ref"] = (
                models_file_path
                + "#/components/schemas/"
                + utils.camel_case(model_ref.id)
            )
        description_parts = []
        if field.description is not None:
            description_parts.append(field.description)
        if (
            field_type.model_ref is not None
            and (model := field_type.model).type == MODEL_ENUM
            and len(constants := model.enum().constants) >= 1
        ):
            markdown = _dump_constants(constants)
            description_parts.append("Constants:\n\n" + markdown)
        if len(description_parts) >= 1:
            property["description"] = "\n\n".join(description_parts)
        if field.example is not None:
            property["example"] = field.example

    def _dump_enum(self, enum: Enum, schema: dict) -> None:
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


_SERVICE_YAML_TEMPLATE = "{}.service.yaml"
_MODELS_YAML = "models.yaml"
_BUILTINS_YAML = "builtins.yaml"


def _dump_error_cases(error_cases: list[ErrorCase]) -> str:
    lines = []
    lines.append("| Code | Message | Description |")
    lines.append("| - | - | - |")
    for error_case in error_cases:
        error = error_case.error
        assert error is not None
        message = error.id.lower().replace("-", " ")
        if error_case.description is None:
            description = ""
        else:
            description = error_case.description
        lines.append(f"| {error.code} | {message} | {description} |")
    return "\n".join(lines)


def _dump_constants(constants: list[Constant]) -> str:
    lines = []
    for constant in constants:
        line = f"- {utils.macro_case(constant.id)}({constant.value})"
        if constant.description is not None:
            line += f": {constant.description}"
        lines.append(line)
    return "\n".join(lines)
