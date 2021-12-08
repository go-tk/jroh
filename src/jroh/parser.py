from dataclasses import dataclass
from typing import Any, Optional, Type, TypeVar

import re2
import yaml

from .spec import (
    BOOL,
    ENUM,
    ENUM_UNDERLYING_TYPE_PATTERN,
    FIELD_TYPE_PATTERN,
    FLOAT32,
    FLOAT64,
    ID_PATTERN,
    INT32,
    INT64,
    MODEL_TYPE_PATTERN,
    REF_PATTERN,
    STRING,
    STRUCT,
    XPRIMIT,
    Constant,
    Enum,
    Error,
    ErrorCase,
    Field,
    FieldType,
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


@dataclass
class ParseFilesResults:
    ignored_node_uris: list[str]
    specs: list[Spec]


def parse_files(file_path_2_file_data: dict[str, str]) -> ParseFilesResults:
    parser = _Parser()
    for file_path, file_data in file_path_2_file_data.items():
        parser.parse_file(file_data, file_path)
    return ParseFilesResults(
        ignored_node_uris=parser.ignored_node_uris(),
        specs=parser.specs(),
    )


class _Parser:
    def __init__(self) -> None:
        self._ignored_node_uris: list[str] = []
        self._specs: list[Spec] = []

    def parse_file(self, file_data: str, file_path: str) -> None:
        raw_spec = yaml.safe_load(file_data)
        node_uri = file_path + "#/"
        spec = Spec(node_uri)
        self._parse_raw_spec(raw_spec, spec)
        self._specs.append(spec)

    def _parse_raw_spec(self, raw_spec, spec: Spec) -> None:
        raw_spec = _ensure_node_kind(raw_spec, dict, spec.node_uri)
        raw_spec = raw_spec.copy()
        if (namespace := raw_spec.pop("namespace", None)) is not None:
            node_uri = spec.node_uri + "namespace"
            namespace = _ensure_node_kind(namespace, str, node_uri)
            _check_id(namespace, node_uri)
            spec.namespace = namespace
        if (raw_services := raw_spec.pop("services", None)) is not None:
            self._parse_raw_services(
                raw_services,
                spec.services,
                spec.node_uri + "services",
            )
        if (raw_methods := raw_spec.pop("methods", None)) is not None:
            self._parse_raw_methods(
                raw_methods, spec.methods, spec.node_uri + "methods"
            )
        if (raw_models := raw_spec.pop("models", None)) is not None:
            self._parse_raw_models(raw_models, spec.models, spec.node_uri + "models")
        if (raw_errors := raw_spec.pop("errors", None)) is not None:
            self._parse_raw_errors(raw_errors, spec.errors, spec.node_uri + "errors")
        for key in raw_spec.keys():
            self._ignored_node_uris.append(spec.node_uri + key)

    def _parse_raw_services(
        self, raw_services, services: list[Service], node_uri: str
    ) -> None:
        raw_services = _ensure_node_kind(raw_services, dict, node_uri)
        _ensure_non_empty_mapping(raw_services, node_uri)
        for service_id, raw_service in raw_services.items():
            node_uri2 = node_uri + "/" + str(service_id)
            service_id = _ensure_node_kind(service_id, str, node_uri2)
            _check_id(service_id, node_uri2)
            service = Service(node_uri2, service_id)
            self._parse_raw_service(raw_service, service)
            services.append(service)

    def _parse_raw_service(self, raw_service, service: Service) -> None:
        raw_service = _ensure_node_kind(raw_service, dict, service.node_uri)
        raw_service = raw_service.copy()
        node_uri = service.node_uri + "/version"
        version = _pop_node(raw_service, "version", node_uri)
        version = _ensure_node_kind(version, str, node_uri)
        service.version = version
        if (description := raw_service.pop("description", None)) is not None:
            description = _ensure_node_kind(
                description, str, service.node_uri + "/description"
            )
            service.description = description
        if (
            rpc_path_template := raw_service.pop("rpc_path_template", None)
        ) is not None:
            rpc_path_template = _ensure_node_kind(
                rpc_path_template, str, service.node_uri + "/rpc_path_template"
            )
            service.rpc_path_template = rpc_path_template
        for key in raw_service.keys():
            self._ignored_node_uris.append(service.node_uri + "/" + key)

    def _parse_raw_methods(
        self, raw_methods, methods: list[Method], node_uri: str
    ) -> None:
        raw_methods = _ensure_node_kind(raw_methods, dict, node_uri)
        _ensure_non_empty_mapping(raw_methods, node_uri)
        for method_id, raw_method in raw_methods.items():
            node_uri2 = node_uri + "/" + str(method_id)
            method_id = _ensure_node_kind(method_id, str, node_uri2)
            _check_id(method_id, node_uri2)
            method = Method(node_uri2, method_id)
            self._parse_raw_method(raw_method, method)
            methods.append(method)

    def _parse_raw_method(self, raw_method, method: Method) -> None:
        raw_method = _ensure_node_kind(raw_method, dict, method.node_uri)
        raw_method = raw_method.copy()
        if (service_ids := raw_method.pop("service_ids", None)) is None:
            node_uri = method.node_uri + "/service_id"
            service_id = _pop_node(raw_method, "service_id", node_uri)
            service_id = _ensure_node_kind(service_id, str, node_uri)
            _check_id(service_id, node_uri)
            method.service_ids = [service_id]
        else:
            node_uri = method.node_uri + "/service_ids"
            service_ids = _ensure_node_kind(service_ids, list, node_uri)
            for i, service_id in enumerate(service_ids):
                node_uri2 = node_uri + f"[{i}]"
                service_id = _ensure_node_kind(service_id, str, node_uri2)
                _check_id(service_id, node_uri2)
            method.service_ids = service_ids
        if (summary := raw_method.pop("summary", None)) is not None:
            summary = _ensure_node_kind(summary, str, method.node_uri + "/summary")
            method.summary = summary
        if (description := raw_method.pop("description", None)) is not None:
            description = _ensure_node_kind(
                description, str, method.node_uri + "/description"
            )
            method.description = description
        if (raw_params := raw_method.pop("params", None)) is not None:
            params = Params(method.node_uri + "/params")
            self._parse_raw_params(raw_params, params)
            method.params = params
        if (raw_results := raw_method.pop("results", None)) is not None:
            results = Results(method.node_uri + "/results")
            self._parse_raw_results(raw_results, results)
            method.results = results
        if (raw_error_cases := raw_method.pop("error_cases", None)) is not None:
            self._parse_raw_error_cases(
                raw_error_cases, method.error_cases, method.node_uri + "/error_cases"
            )
        for key in raw_method.keys():
            self._ignored_node_uris.append(method.node_uri + "/" + key)

    def _parse_raw_params(self, raw_params, params: Params) -> None:
        self._parse_raw_fields(raw_params, params.fields, params.node_uri)

    def _parse_raw_results(self, raw_results, results: Results) -> None:
        self._parse_raw_fields(raw_results, results.fields, results.node_uri)

    def _parse_raw_error_cases(
        self, raw_error_cases, error_cases: list[ErrorCase], node_uri: str
    ) -> None:
        raw_error_cases = _ensure_node_kind(raw_error_cases, dict, node_uri)
        _ensure_non_empty_mapping(raw_error_cases, node_uri)
        for raw_error_ref, raw_error_case in raw_error_cases.items():
            node_uri2 = node_uri + "/" + str(raw_error_ref)
            raw_error_ref = _ensure_node_kind(raw_error_ref, str, node_uri2)
            error_ref = _parse_raw_ref(raw_error_ref, node_uri2)
            error_case = ErrorCase(node_uri2, error_ref)
            self._parse_raw_error_case(raw_error_case, error_case)
            error_cases.append(error_case)

    def _parse_raw_error_case(self, raw_error_case, error_case: ErrorCase) -> None:
        raw_error_case = _ensure_node_kind(raw_error_case, dict, error_case.node_uri)
        raw_error_case = raw_error_case.copy()
        if (description := raw_error_case.pop("description", None)) is not None:
            description = _ensure_node_kind(
                description, str, error_case.node_uri + "/description"
            )
            error_case.description = description
        for key in raw_error_case.keys():
            self._ignored_node_uris.append(error_case.node_uri + "/" + key)

    def _parse_raw_models(self, raw_models, models: list[Model], node_uri: str) -> None:
        raw_models = _ensure_node_kind(raw_models, dict, node_uri)
        _ensure_non_empty_mapping(raw_models, node_uri)
        for model_id, raw_model in raw_models.items():
            node_uri2 = node_uri + "/" + str(model_id)
            model_id = _ensure_node_kind(model_id, str, node_uri2)
            _check_id(model_id, node_uri2)
            model = Model(node_uri2, model_id)
            self._parse_raw_model(raw_model, model)
            models.append(model)

    def _parse_raw_model(self, raw_model, model: Model) -> None:
        raw_model = _ensure_node_kind(raw_model, dict, model.node_uri)
        raw_model = raw_model.copy()
        node_uri = model.node_uri + "/type"
        raw_model_type = _pop_node(raw_model, "type", node_uri)
        raw_model_type = _ensure_node_kind(raw_model_type, str, node_uri)
        _check_raw_model_type(raw_model_type, node_uri)
        if raw_model_type == STRUCT:
            struct = Struct()
            self._parse_raw_struct(raw_model, struct, model.node_uri)
            model.type = STRUCT
            model.definition = struct
        elif raw_model_type == ENUM:
            enum = Enum()
            self._parse_raw_enum(raw_model, enum, model.node_uri)
            model.type = ENUM
            model.definition = enum
        else:
            primitive_type = raw_model_type
            xprimit = Xprimit(primitive_type)
            self._parse_raw_xprimit(raw_model, xprimit, model.node_uri)
            model.type = XPRIMIT
            model.definition = xprimit
        if (description := raw_model.pop("description", None)) is not None:
            description = _ensure_node_kind(
                description, str, model.node_uri + "/description"
            )
            model.description = description
        for key in raw_model.keys():
            self._ignored_node_uris.append(model.node_uri + "/" + key)

    def _parse_raw_struct(
        self, raw_struct: dict, struct: Struct, node_uri: str
    ) -> None:
        if (raw_fields := raw_struct.pop("fields", None)) is not None:
            self._parse_raw_fields(raw_fields, struct.fields, node_uri + "/fields")

    def _parse_raw_fields(self, raw_fields, fields: list[Field], node_uri: str) -> None:
        raw_fields = _ensure_node_kind(raw_fields, dict, node_uri)
        _ensure_non_empty_mapping(raw_fields, node_uri)
        for field_id, raw_field in raw_fields.items():
            node_uri2 = node_uri + "/" + str(field_id)
            field_id = _ensure_node_kind(field_id, str, node_uri2)
            _check_id(field_id, node_uri2)
            field = Field(node_uri2, field_id)
            self._parse_raw_field(raw_field, field)
            fields.append(field)

    def _parse_raw_field(self, raw_field, field: Field) -> None:
        raw_field = _ensure_node_kind(raw_field, dict, field.node_uri)
        raw_field = raw_field.copy()
        node_uri = field.node_uri + "/type"
        raw_field_type = _pop_node(raw_field, "type", node_uri)
        raw_field_type = _ensure_node_kind(raw_field_type, str, node_uri)
        field_type = field.type
        _parse_raw_field_type(raw_field_type, field_type, node_uri)
        if field_type.is_primitive():
            primitive_type = field_type.primitive_type()
            _load_primitive_constraints(
                primitive_type, raw_field, field, field.node_uri
            )
        if (is_optional := raw_field.pop("is_optional", None)) is not None:
            is_optional = _ensure_node_kind(
                is_optional, bool, field.node_uri + "/is_optional"
            )
            field.is_optional = is_optional
        elif (is_repeated := raw_field.pop("is_repeated", None)) is not None:
            is_repeated = _ensure_node_kind(
                is_repeated, bool, field.node_uri + "/is_repeated"
            )
            field.is_repeated = is_repeated
            min_count = raw_field.pop("min_count", None)
            if min_count is not None:
                node_uri2 = field.node_uri + "/min_count"
                min_count = _ensure_node_kind(min_count, int, node_uri2)
                _check_number(min_count, 0, None, node_uri2)
                field.min_count = min_count
            max_count = raw_field.pop("max_count", None)
            if max_count is not None:
                node_uri2 = field.node_uri + "/max_count"
                max_count = _ensure_node_kind(max_count, int, node_uri2)
                _check_number(max_count, 1, None, node_uri2)
                field.max_count = max_count
            if (
                min_count is not None
                and max_count is not None
                and min_count > max_count
            ):
                raise InvalidSpecError(
                    f"invalid field, min_count > max_count: node_uri={field.node_uri!r} min_count={field.min_count} max_count={field.max_count}"
                )
        if (description := raw_field.pop("description", None)) is not None:
            description = _ensure_node_kind(
                description, str, field.node_uri + "/description"
            )
            field.description = description
        if (
            field_type.is_primitive()
            and (example := raw_field.pop("example", None)) is not None
        ):
            primitive_type = field_type.primitive_type()
            node_uri2 = field.node_uri + "/example"
            if field.is_repeated:
                example = _ensure_node_kind(example, list, node_uri2)
                _check_sequence_length(
                    len(example), field.min_count, field.max_count, node_uri2
                )
                for i, v in enumerate(example):
                    _check_primitive_value(
                        primitive_type, v, field, node_uri2 + f"[{i}]"
                    )
            else:
                _check_primitive_value(primitive_type, example, field, node_uri2)
            field.example = example
        for key in raw_field.keys():
            self._ignored_node_uris.append(field.node_uri + "/" + key)

    def _parse_raw_enum(self, raw_enum: dict, enum: Enum, node_uri: str) -> None:
        node_uri2 = node_uri + "/underlying_type"
        enum_underlying_type = _pop_node(raw_enum, "underlying_type", node_uri2)
        enum_underlying_type = _ensure_node_kind(enum_underlying_type, str, node_uri2)
        _check_enum_underlying_type(enum_underlying_type, node_uri2)
        enum.underlying_type = enum_underlying_type
        if (raw_fields := raw_enum.pop("constants", None)) is not None:
            type = {
                INT32: int,
                INT64: int,
                STRING: str,
            }[enum_underlying_type]
            self._parse_raw_constants(
                raw_fields, enum.constants, node_uri + "/constants", type
            )

    def _parse_raw_constants(
        self, raw_constants, constants: list[Constant], node_uri: str, type: Type
    ) -> None:
        raw_constants = _ensure_node_kind(raw_constants, dict, node_uri)
        _ensure_non_empty_mapping(raw_constants, node_uri)
        for constant_id, raw_constant in raw_constants.items():
            node_uri2 = node_uri + "/" + str(constant_id)
            constant_id = _ensure_node_kind(constant_id, str, node_uri2)
            _check_id(constant_id, node_uri2)
            constant = Constant(node_uri2, constant_id)
            self._parse_raw_constant(raw_constant, constant, type)
            constants.append(constant)

    def _parse_raw_constant(self, raw_constant, constant: Constant, type: Type) -> None:
        raw_constant = _ensure_node_kind(raw_constant, dict, constant.node_uri)
        raw_constant = raw_constant.copy()
        node_uri = constant.node_uri + "/value"
        value = _pop_node(raw_constant, "value", node_uri)
        value = _ensure_node_kind(value, type, node_uri)
        constant.value = value
        if (description := raw_constant.pop("description", None)) is not None:
            description = _ensure_node_kind(
                description, str, constant.node_uri + "/description"
            )
            constant.description = description
        for key in raw_constant.keys():
            self._ignored_node_uris.append(constant.node_uri + "/" + key)

    def _parse_raw_xprimit(
        self, raw_xprimit: dict, xprimit: Xprimit, node_uri: str
    ) -> None:
        _load_primitive_constraints(
            xprimit.primitive_type, raw_xprimit, xprimit, node_uri
        )
        if (example := raw_xprimit.pop("example", None)) is not None:
            _check_primitive_value(
                xprimit.primitive_type, example, xprimit, node_uri + "/example"
            )
            xprimit.example = example

    def _parse_raw_errors(self, raw_errors, errors: list[Error], node_uri: str) -> None:
        raw_errors = _ensure_node_kind(raw_errors, dict, node_uri)
        _ensure_non_empty_mapping(raw_errors, node_uri)
        for error_id, raw_error in raw_errors.items():
            node_uri2 = node_uri + "/" + str(error_id)
            error_id = _ensure_node_kind(error_id, str, node_uri2)
            _check_id(error_id, node_uri2)
            error = Error(node_uri2, error_id)
            self._parse_raw_error(raw_error, error)
            errors.append(error)

    def _parse_raw_error(self, raw_error, error: Error) -> None:
        raw_error = _ensure_node_kind(raw_error, dict, error.node_uri)
        raw_error = raw_error.copy()
        node_uri = error.node_uri + "/code"
        error_code = _pop_node(raw_error, "code", node_uri)
        error_code = _ensure_node_kind(error_code, int, node_uri)
        _check_number(error_code, 1000, _MAX_INT32, node_uri)
        error.code = error_code
        node_uri = error.node_uri + "/status_code"
        status_code = _pop_node(raw_error, "status_code", node_uri)
        status_code = _ensure_node_kind(status_code, int, node_uri)
        _check_number(status_code, 100, 599, node_uri)
        error.status_code = status_code
        if (description := raw_error.pop("description", None)) is not None:
            description = _ensure_node_kind(
                description, str, error.node_uri + "/description"
            )
            error.description = description
        for key in raw_error.keys():
            self._ignored_node_uris.append(error.node_uri + "/" + key)

    def ignored_node_uris(self) -> list[str]:
        return self._ignored_node_uris

    def specs(self) -> list[Spec]:
        return self._specs


class InvalidSpecError(Exception):
    def __init__(self, message: str) -> None:
        super().__init__("invalid specification: " + message)


def _pop_node(mapping: dict, key: str, node_uri: str):
    node_value = mapping.pop(key, None)
    if node_value is None:
        raise InvalidSpecError(f"missing node: node_uri={node_uri!r}")
    return node_value


_node_kinds: dict[Type, str] = {
    bool: "boolean",
    int: "integer",
    float: "floating-point",
    str: "string",
    list: "sequence",
    dict: "mapping",
}


_T = TypeVar("_T")


def _ensure_node_kind(node_value, expected_node_type: Type[_T], node_uri: str) -> _T:
    if not isinstance(node_value, expected_node_type):
        raise InvalidSpecError(
            f"invalid node kind: node_uri={node_uri!r} node_kind={_node_kinds[type(node_value)]} expected_node_kind={_node_kinds[expected_node_type]}"
        )
    return node_value


def _check_id(id: str, node_uri: str) -> None:
    if ID_PATTERN.fullmatch(id) is None:
        raise InvalidSpecError(
            f"invalid id; node_uri={node_uri!r} id={id!r} expected_pattern={ID_PATTERN.pattern!r}"
        )


def _ensure_non_empty_mapping(mapping: dict, node_uri: str) -> None:
    if len(mapping) == 0:
        raise InvalidSpecError(f"non-empty mapping required: node_uri={node_uri!r}")


_T = TypeVar("_T", int, float)


def _check_number(
    number: _T,
    min_number: Optional[_T],
    max_number: Optional[_T],
    node_uri: str,
    *,
    min_number_is_exclusive: bool = False,
    max_number_is_exclusive: bool = False,
) -> None:
    if min_number is not None:
        if min_number_is_exclusive:
            if number <= min_number:
                raise InvalidSpecError(
                    f"number too small: node_uri={node_uri!r} number={number} exclusive_min_number={min_number}"
                )
        else:
            if number < min_number:
                raise InvalidSpecError(
                    f"number too small: node_uri={node_uri!r} number={number} min_number={min_number}"
                )
    if max_number is not None:
        if max_number_is_exclusive:
            if number >= max_number:
                raise InvalidSpecError(
                    f"number too large: node_uri={node_uri!r} number={number} exclusive_max_number={max_number}"
                )
        else:
            if number > max_number:
                raise InvalidSpecError(
                    f"number too large: node_uri={node_uri!r} number={number} max_number={max_number}"
                )


def _check_string(string: str, pattern, node_uri: str) -> None:
    if pattern.fullmatch(string) is None:
        raise InvalidSpecError(
            f"unexpected string: node_uri={node_uri!r} string={string!r} expected_pattern={pattern.pattern!r}"
        )


def _check_string_length(
    string_length: int,
    min_string_length: int,
    max_string_length: Optional[int],
    node_uri: str,
) -> None:
    _check_length(
        "string", string_length, min_string_length, max_string_length, node_uri
    )


def _check_sequence_length(
    sequence_length: int,
    min_sequence_length: int,
    max_sequence_length: Optional[int],
    node_uri: str,
) -> None:
    _check_length(
        "sequence", sequence_length, min_sequence_length, max_sequence_length, node_uri
    )


def _check_length(
    kind: str,
    length: int,
    min_length: int,
    max_length: Optional[int],
    node_uri: str,
) -> None:
    if length < min_length:
        raise InvalidSpecError(
            f"{kind} too short: node_uri={node_uri!r} {kind}_length={length} min_{kind}_length={min_length}"
        )
    if max_length is not None and length > max_length:
        raise InvalidSpecError(
            f"{kind} too long: node_uri={node_uri!r} {kind}_length={length} max_{kind}_length={max_length}"
        )


def _check_raw_model_type(raw_model_type: str, node_uri: str) -> None:
    if MODEL_TYPE_PATTERN.fullmatch(raw_model_type) is None:
        raise InvalidSpecError(
            f"invalid model type; node_uri={node_uri!r} model_type={raw_model_type!r} expected_pattern={MODEL_TYPE_PATTERN.pattern!r}"
        )


def _parse_raw_field_type(
    raw_field_type: str, field_type: FieldType, node_uri: str
) -> None:
    _check_raw_field_type(raw_field_type, node_uri)
    if (i := raw_field_type.find(".")) < 0:
        if raw_field_type[0].islower():
            field_type.value = raw_field_type
        else:
            field_type.value = Ref(namespace=None, id=raw_field_type)
    else:
        field_type.value = Ref(namespace=raw_field_type[:i], id=raw_field_type[i + 1 :])


def _check_raw_field_type(raw_field_type: str, node_uri: str) -> None:
    if FIELD_TYPE_PATTERN.fullmatch(raw_field_type) is None:
        raise InvalidSpecError(
            f"invalid field type; node_uri={node_uri!r} field_type={raw_field_type!r} expected_pattern={FIELD_TYPE_PATTERN.pattern!r}"
        )


def _check_enum_underlying_type(enum_underlying_type: str, node_uri: str) -> None:
    if ENUM_UNDERLYING_TYPE_PATTERN.fullmatch(enum_underlying_type) is None:
        raise InvalidSpecError(
            f"invalid enum underlying type; node_uri={node_uri!r} enum_underlying_type={enum_underlying_type!r} expected_pattern={ENUM_UNDERLYING_TYPE_PATTERN.pattern!r}"
        )


_MIN_INT32 = -(1 << 31)
_MAX_INT32 = (1 << 31) - 1
_MIN_INT64 = -(1 << 63)
_MAX_INT64 = (1 << 63) - 1


def _load_primitive_constraints(
    primitive_type: str,
    raw_object: dict,
    primitive_constraints: PrimitiveConstraints,
    node_uri: str,
) -> None:
    if primitive_type == BOOL:
        pass
    elif primitive_type in (INT32, INT64):
        min = raw_object.pop("min", None)
        if min is not None:
            node_uri2 = node_uri + "/min"
            min = _ensure_node_kind(min, int, node_uri2)
            if primitive_type == INT32:
                _check_number(min, _MIN_INT32, _MAX_INT32, node_uri2)
            else:
                _check_number(min, _MIN_INT64, _MAX_INT64, node_uri2)
            primitive_constraints.min = min
        max = raw_object.pop("max", None)
        if max is not None:
            node_uri2 = node_uri + "/max"
            max = _ensure_node_kind(max, int, node_uri2)
            if primitive_type == INT32:
                _check_number(max, _MIN_INT32, _MAX_INT32, node_uri2)
            else:
                _check_number(max, _MIN_INT64, _MAX_INT64, node_uri2)
            primitive_constraints.max = max
        if min is not None and max is not None and min > max:
            raise InvalidSpecError(
                f"invalid primitive constraints, min > max: node_uri={node_uri!r} min={min} max={max}"
            )
    elif primitive_type in (FLOAT32, FLOAT64):
        min = raw_object.pop("min", None)
        if min is not None:
            min = _ensure_node_kind(min, float, node_uri + "/min")
            primitive_constraints.min = min
            if (
                min_is_exclusive := raw_object.pop("min_is_exclusive", None)
            ) is not None:
                min_is_exclusive = _ensure_node_kind(
                    min_is_exclusive,
                    bool,
                    node_uri + "/min_is_exclusive",
                )
                primitive_constraints.min_is_exclusive = min_is_exclusive
        max = raw_object.pop("max", None)
        if max is not None:
            max = _ensure_node_kind(max, float, node_uri + "/max")
            primitive_constraints.max = max
            if (
                max_is_exclusive := raw_object.pop("max_is_exclusive", None)
            ) is not None:
                max_is_exclusive = _ensure_node_kind(
                    max_is_exclusive,
                    bool,
                    node_uri + "/max_is_exclusive",
                )
                primitive_constraints.max_is_exclusive = max_is_exclusive
        if min is not None and max is not None:
            if min == max:
                if (
                    primitive_constraints.min_is_exclusive
                    or primitive_constraints.max_is_exclusive
                ):
                    raise InvalidSpecError(
                        f"invalid primitive constraints, min > max: node_uri={node_uri!r}"
                        f" {'exclusive_' if primitive_constraints.min_is_exclusive else ''}min={min}"
                        f" {'exclusive_' if primitive_constraints.max_is_exclusive else ''}max={max}"
                    )
            elif min > max:
                raise InvalidSpecError(
                    f"invalid primitive constraints, min > max: node_uri={node_uri!r} min={min} max={max}"
                )
    elif primitive_type == STRING:
        min_length = raw_object.pop("min_length", None)
        if min_length is not None:
            node_uri2 = node_uri + "/min_length"
            min_length = _ensure_node_kind(min_length, int, node_uri2)
            _check_number(min_length, 0, None, node_uri2)
            primitive_constraints.min_length = min_length
        max_length = raw_object.pop("max_length", None)
        if max_length is not None:
            node_uri2 = node_uri + "/max_length"
            max_length = _ensure_node_kind(max_length, int, node_uri2)
            _check_number(max_length, 1, None, node_uri2)
            primitive_constraints.max_length = max_length
        if (
            min_length is not None
            and max_length is not None
            and min_length > max_length
        ):
            raise InvalidSpecError(
                f"invalid primitive constraints, min_length > max_length: node_uri={node_uri!r} min_length={min_length} max_length={max_length}"
            )
        if (pattern := raw_object.pop("pattern", None)) is not None:
            node_uri2 = node_uri + "/pattern"
            pattern = _ensure_node_kind(pattern, str, node_uri2)
            _check_string_length(len(pattern), 1, None, node_uri2)
            try:
                re2.compile(pattern)
            except re2.error:
                raise InvalidSpecError(
                    f"invalid regular expression (re2): node_uri={node_uri!r} reg_exp={pattern!r}"
                ) from None
            primitive_constraints.pattern = pattern
    else:
        assert False, primitive_type


def _check_primitive_value(
    primitive_type: str,
    primitive_value: Any,
    primitive_constraints: PrimitiveConstraints,
    node_uri: str,
) -> None:
    if primitive_type == BOOL:
        _ensure_node_kind(primitive_value, bool, node_uri)
    elif primitive_type in (INT32, INT64):
        min = primitive_constraints.min
        if min is None:
            if primitive_type == INT32:
                min = _MIN_INT32
            else:
                min = _MIN_INT64
        max = primitive_constraints.max
        if max is None:
            if primitive_type == INT32:
                max = _MAX_INT32
            else:
                max = _MAX_INT64
        pv = _ensure_node_kind(primitive_value, int, node_uri)
        _check_number(pv, min, max, node_uri)
    elif primitive_type in (FLOAT32, FLOAT64):
        pv = _ensure_node_kind(primitive_value, float, node_uri)
        _check_number(
            pv,
            primitive_constraints.min,
            primitive_constraints.max,
            node_uri,
            min_number_is_exclusive=primitive_constraints.min_is_exclusive,
            max_number_is_exclusive=primitive_constraints.max_is_exclusive,
        )
    elif primitive_type == STRING:
        pv = _ensure_node_kind(primitive_value, str, node_uri)
        _check_string_length(
            len(pv.encode()),
            primitive_constraints.min_length,
            primitive_constraints.max_length,
            node_uri,
        )
        if primitive_constraints.pattern != "":
            _check_string(pv, re2.compile(primitive_constraints.pattern), node_uri)
    else:
        assert False, primitive_type


def _parse_raw_ref(raw_ref: str, node_uri: str) -> Ref:
    _check_raw_ref(raw_ref, node_uri)
    if (i := raw_ref.find(".")) < 0:
        ref = Ref(namespace=None, id=raw_ref)
    else:
        ref = Ref(namespace=raw_ref[:i], id=raw_ref[i + 1 :])
    return ref


def _check_raw_ref(raw_ref: str, node_uri: str) -> None:
    if REF_PATTERN.fullmatch(raw_ref) is None:
        raise InvalidSpecError(
            f"invalid ref; node_uri={node_uri!r} ref={raw_ref!r} expected_pattern={REF_PATTERN.pattern!r}"
        )
