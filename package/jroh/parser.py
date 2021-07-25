from dataclasses import dataclass
from typing import Type, TypeVar, Union

import yaml

from .spec import (ENUM_UNDERLYING_TYPE_PATTERN, FIELD_BOOL, FIELD_FLOAT32,
                   FIELD_FLOAT64, FIELD_INT32, FIELD_INT64, FIELD_STRING,
                   FIELD_TYPE_PATTERN, ID_PATTERN, MODEL_ENUM, MODEL_STRUCT,
                   MODEL_TYPE_PATTERN, NAMESPACE_PATTERN, REF_PATTERN,
                   Constant, Enum, Error, ErrorCase, Field, FieldType, Method,
                   Model, Params, Ref, Result, Service, Spec, Struct)


@dataclass
class ParseFilesResult:
    ignored_node_uris: list[str]
    specs: list[Spec]


def parse_files(file_path_2_file_data: dict[str, str]) -> ParseFilesResult:
    parser = _Parser()
    for file_path, file_data in file_path_2_file_data.items():
        parser.parse_file(file_data, file_path)
    return ParseFilesResult(
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
        raw_spec = _ensure_node_type(raw_spec, dict, spec.node_uri)
        raw_spec = raw_spec.copy()
        if (namespace := raw_spec.pop("namespace", None)) is not None:
            node_uri = spec.node_uri + "namespace"
            namespace = _ensure_node_type(namespace, str, node_uri)
            _check_namespace(namespace, node_uri)
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
        raw_services = _ensure_node_type(raw_services, dict, node_uri)
        _ensure_not_empty(raw_services, node_uri)
        for service_id, raw_service in raw_services.items():
            node_uri2 = node_uri + "/" + str(service_id)
            service_id = _ensure_node_type(service_id, str, node_uri2)
            _check_id(service_id, node_uri2)
            service = Service(node_uri2, service_id)
            self._parse_raw_service(raw_service, service)
            services.append(service)

    def _parse_raw_service(self, raw_service, service: Service) -> None:
        raw_service = _ensure_node_type(raw_service, dict, service.node_uri)
        raw_service = raw_service.copy()
        node_uri = service.node_uri + "/version"
        version = _pop_node(raw_service, "version", node_uri)
        version = _ensure_node_type(version, str, node_uri)
        service.version = version
        if (description := raw_service.pop("description", None)) is not None:
            description = _ensure_node_type(
                description, str, service.node_uri + "/description"
            )
            service.description = description
        if (
            method_path_template := raw_service.pop("method_path_template", None)
        ) is not None:
            method_path_template = _ensure_node_type(
                method_path_template, str, service.node_uri + "/method_path_template"
            )
            service.method_path_template = method_path_template
        for key in raw_service.keys():
            self._ignored_node_uris.append(service.node_uri + "/" + key)

    def _parse_raw_methods(
        self, raw_methods, methods: list[Method], node_uri: str
    ) -> None:
        raw_methods = _ensure_node_type(raw_methods, dict, node_uri)
        _ensure_not_empty(raw_methods, node_uri)
        for method_id, raw_method in raw_methods.items():
            node_uri2 = node_uri + "/" + str(method_id)
            method_id = _ensure_node_type(method_id, str, node_uri2)
            _check_id(method_id, node_uri2)
            method = Method(node_uri2, method_id)
            self._parse_raw_method(raw_method, method)
            methods.append(method)

    def _parse_raw_method(self, raw_method, method: Method) -> None:
        raw_method = _ensure_node_type(raw_method, dict, method.node_uri)
        raw_method = raw_method.copy()
        node_uri = method.node_uri + "/service_id"
        service_id = _pop_node(raw_method, "service_id", node_uri)
        service_id = _ensure_node_type(service_id, str, node_uri)
        _check_id(service_id, node_uri)
        method.service_id = service_id
        if (summary := raw_method.pop("summary", None)) is not None:
            summary = _ensure_node_type(summary, str, method.node_uri + "/summary")
            method.summary = summary
        if (description := raw_method.pop("description", None)) is not None:
            description = _ensure_node_type(
                description, str, method.node_uri + "/description"
            )
            method.description = description
        if (raw_params := raw_method.pop("params", None)) is not None:
            params = Params(method.node_uri + "/params")
            self._parse_raw_params(raw_params, params)
            method.params = params
        if (raw_result := raw_method.pop("result", None)) is not None:
            result = Result(method.node_uri + "/result")
            self._parse_raw_result(raw_result, result)
            method.result = result
        if (raw_error_cases := raw_method.pop("error_cases", None)) is not None:
            self._parse_raw_error_cases(
                raw_error_cases, method.error_cases, method.node_uri + "/error_cases"
            )
        for key in raw_method.keys():
            self._ignored_node_uris.append(method.node_uri + "/" + key)

    def _parse_raw_params(self, raw_params, params: Params) -> None:
        self._parse_raw_fields(raw_params, params.fields, params.node_uri)

    def _parse_raw_result(self, raw_result, result: Result) -> None:
        self._parse_raw_fields(raw_result, result.fields, result.node_uri)

    def _parse_raw_error_cases(
        self, raw_error_cases, error_cases: list[ErrorCase], node_uri: str
    ) -> None:
        raw_error_cases = _ensure_node_type(raw_error_cases, dict, node_uri)
        _ensure_not_empty(raw_error_cases, node_uri)
        for raw_error_ref, raw_error_case in raw_error_cases.items():
            node_uri2 = node_uri + "/" + str(raw_error_ref)
            raw_error_ref = _ensure_node_type(raw_error_ref, str, node_uri2)
            error_ref = _parse_raw_ref(raw_error_ref, node_uri2)
            error_case = ErrorCase(node_uri2, error_ref)
            self._parse_raw_error_case(raw_error_case, error_case)
            error_cases.append(error_case)

    def _parse_raw_error_case(self, raw_error_case, error_case: ErrorCase) -> None:
        raw_error_case = _ensure_node_type(raw_error_case, dict, error_case.node_uri)
        raw_error_case = raw_error_case.copy()
        if (description := raw_error_case.pop("description", None)) is not None:
            description = _ensure_node_type(
                description, str, error_case.node_uri + "/description"
            )
            error_case.description = description
        for key in raw_error_case.keys():
            self._ignored_node_uris.append(error_case.node_uri + "/" + key)

    def _parse_raw_models(self, raw_models, models: list[Model], node_uri: str) -> None:
        raw_models = _ensure_node_type(raw_models, dict, node_uri)
        _ensure_not_empty(raw_models, node_uri)
        for model_id, raw_model in raw_models.items():
            node_uri2 = node_uri + "/" + str(model_id)
            model_id = _ensure_node_type(model_id, str, node_uri2)
            _check_id(model_id, node_uri2)
            model = Model(node_uri2, model_id)
            self._parse_raw_model(raw_model, model)
            models.append(model)

    def _parse_raw_model(self, raw_model, model: Model) -> None:
        raw_model = _ensure_node_type(raw_model, dict, model.node_uri)
        raw_model = raw_model.copy()
        node_uri = model.node_uri + "/type"
        model_type = _pop_node(raw_model, "type", node_uri)
        model_type = _ensure_node_type(model_type, str, node_uri)
        _check_model_type(model_type, node_uri)
        model.type = model_type
        if model_type == MODEL_STRUCT:
            struct = Struct()
            self._parse_raw_struct(raw_model, struct, model.node_uri)
            model.definition = struct
        elif model_type == MODEL_ENUM:
            enum = Enum()
            self._parse_raw_enum(raw_model, enum, model.node_uri)
            model.definition = enum
        else:
            assert False, model_type
        for key in raw_model.keys():
            self._ignored_node_uris.append(model.node_uri + "/" + key)

    def _parse_raw_struct(
        self, raw_struct: dict, struct: Struct, node_uri: str
    ) -> None:
        if (raw_fields := raw_struct.pop("fields", None)) is not None:
            self._parse_raw_fields(raw_fields, struct.fields, node_uri + "/fields")

    def _parse_raw_fields(self, raw_fields, fields: list[Field], node_uri: str) -> None:
        raw_fields = _ensure_node_type(raw_fields, dict, node_uri)
        _ensure_not_empty(raw_fields, node_uri)
        for field_id, raw_field in raw_fields.items():
            node_uri2 = node_uri + "/" + str(field_id)
            field_id = _ensure_node_type(field_id, str, node_uri2)
            _check_id(field_id, node_uri2)
            field = Field(node_uri2, field_id)
            self._parse_raw_field(raw_field, field)
            fields.append(field)

    def _parse_raw_field(self, raw_field, field: Field) -> None:
        raw_field = _ensure_node_type(raw_field, dict, field.node_uri)
        raw_field = raw_field.copy()
        node_uri = field.node_uri + "/type"
        raw_field_type = _pop_node(raw_field, "type", node_uri)
        raw_field_type = _ensure_node_type(raw_field_type, str, node_uri)
        field_type = FieldType(node_uri)
        _parse_raw_field_type(raw_field_type, field_type)
        field.type = field_type
        if (description := raw_field.pop("description", None)) is not None:
            description = _ensure_node_type(
                description, str, field.node_uri + "/description"
            )
            field.description = description
        if (
            field_type.model_ref is None
            and (example := raw_field.pop("example", None)) is not None
        ):
            type = {
                FIELD_BOOL: bool,
                FIELD_INT32: int,
                FIELD_INT64: int,
                FIELD_FLOAT32: float,
                FIELD_FLOAT64: float,
                FIELD_STRING: str,
            }[field_type.value]
            node_uri2 = field.node_uri + "/example"
            if field_type.is_repeated:
                example = _ensure_node_type(example, list, node_uri2)
                if len(example) == 0:
                    _ensure_not_empty(example, node_uri2)
                for i, v in enumerate(example):
                    _ensure_node_type(v, type, node_uri2 + f"[{i}]")
            else:
                _ensure_node_type(example, type, node_uri2)
            field.example = example
        for key in raw_field.keys():
            self._ignored_node_uris.append(field.node_uri + "/" + key)

    def _parse_raw_enum(self, raw_enum: dict, enum: Enum, node_uri: str) -> None:
        node_uri2 = node_uri + "/underlying_type"
        enum_underlying_type = _pop_node(raw_enum, "underlying_type", node_uri2)
        _ensure_node_type(enum_underlying_type, str, node_uri2)
        _check_enum_underlying_type(enum_underlying_type, node_uri2)
        enum.underlying_type = enum_underlying_type
        if (raw_fields := raw_enum.pop("constants", None)) is not None:
            type = {
                FIELD_INT32: int,
                FIELD_INT64: int,
                FIELD_STRING: str,
            }[enum_underlying_type]
            self._parse_raw_constants(
                raw_fields, enum.constants, node_uri + "/constants", type
            )

    def _parse_raw_constants(
        self, raw_constants, constants: list[Constant], node_uri: str, type: Type
    ) -> None:
        raw_constants = _ensure_node_type(raw_constants, dict, node_uri)
        _ensure_not_empty(raw_constants, node_uri)
        for constant_id, raw_constant in raw_constants.items():
            node_uri2 = node_uri + "/" + str(constant_id)
            constant_id = _ensure_node_type(constant_id, str, node_uri2)
            _check_id(constant_id, node_uri2)
            constant = Constant(node_uri2, constant_id)
            self._parse_raw_constant(raw_constant, constant, type)
            constants.append(constant)

    def _parse_raw_constant(self, raw_constant, constant: Constant, type: Type) -> None:
        raw_constant = _ensure_node_type(raw_constant, dict, constant.node_uri)
        raw_constant = raw_constant.copy()
        node_uri = constant.node_uri + "/value"
        value = _pop_node(raw_constant, "value", node_uri)
        value = _ensure_node_type(value, type, node_uri)
        constant.value = value
        if (description := raw_constant.pop("description", None)) is not None:
            description = _ensure_node_type(
                description, str, constant.node_uri + "/description"
            )
            constant.description = description
        for key in raw_constant.keys():
            self._ignored_node_uris.append(constant.node_uri + "/" + key)

    def _parse_raw_errors(self, raw_errors, errors: list[Error], node_uri: str) -> None:
        raw_errors = _ensure_node_type(raw_errors, dict, node_uri)
        _ensure_not_empty(raw_errors, node_uri)
        for error_id, raw_error in raw_errors.items():
            node_uri2 = node_uri + "/" + str(error_id)
            error_id = _ensure_node_type(error_id, str, node_uri2)
            _check_id(error_id, node_uri2)
            error = Error(node_uri2, error_id)
            self._parse_raw_error(raw_error, error)
            errors.append(error)

    def _parse_raw_error(self, raw_error, error: Error) -> None:
        raw_error = _ensure_node_type(raw_error, dict, error.node_uri)
        raw_error = raw_error.copy()
        node_uri = error.node_uri + "/code"
        error_code = _pop_node(raw_error, "code", node_uri)
        error_code = _ensure_node_type(error_code, int, node_uri)
        error.code = error_code
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


_T = TypeVar("_T")


def _ensure_node_type(node_value, expected_node_type: Type[_T], node_uri: str) -> _T:
    if not isinstance(node_value, expected_node_type):
        raise InvalidSpecError(
            f"invalid node type: node_uri={node_uri!r} node_type={type(node_value).__name__} expected_node_type={expected_node_type.__name__}"
        )
    return node_value


def _ensure_not_empty(node_value: Union[list, dict], node_uri: str):
    if len(node_value) == 0:
        raise InvalidSpecError(f"node should not be empty: node_uri={node_uri!r}")


def _check_namespace(namespace: str, node_uri: str) -> None:
    if NAMESPACE_PATTERN.fullmatch(namespace) is None:
        raise InvalidSpecError(
            f"invalid namespace; node_uri={node_uri!r} namespace={namespace!r} expected_pattern={NAMESPACE_PATTERN.pattern!r}"
        )


def _check_id(id: str, node_uri: str) -> None:
    if ID_PATTERN.fullmatch(id) is None:
        raise InvalidSpecError(
            f"invalid id; node_uri={node_uri!r} id={id!r} expected_pattern={ID_PATTERN.pattern!r}"
        )


def _check_model_type(model_type: str, node_uri: str) -> None:
    if MODEL_TYPE_PATTERN.fullmatch(model_type) is None:
        raise InvalidSpecError(
            f"invalid model type; node_uri={node_uri!r} model_type={model_type!r} expected_pattern={MODEL_TYPE_PATTERN.pattern!r}"
        )


def _parse_raw_field_type(raw_field_type: str, field_type: FieldType) -> None:
    _check_raw_field_type(raw_field_type, field_type.node_uri)
    s = raw_field_type
    if (c := s[-1]) in ("?", "+", "*"):
        if c == "?":
            field_type.is_optional = True
        elif c == "+":
            field_type.is_repeated = True
        elif c == "*":
            field_type.is_optional = True
            field_type.is_repeated = True
        else:
            assert False, c
        s = s[:-1]
    if (i := s.find(".")) < 0:
        if s[0].islower():
            field_type.value = s
        else:
            field_type.model_ref = Ref(namespace=None, id=s)
    else:
        field_type.model_ref = Ref(namespace=s[:i], id=s[i + 1 :])


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
