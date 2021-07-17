from dataclasses import dataclass
from typing import Type, TypeVar, Union

import yaml

from .spec import (ENUM_UNDERLYING_TYPE_PATTERN, FIELD_BOOL, FIELD_FLOAT32,
                   FIELD_FLOAT64, FIELD_INT32, FIELD_INT64, FIELD_STRING,
                   ID_PATTERN, MODEL_ENUM, MODEL_STRUCT, MODEL_TYPE_PATTERN,
                   NAMESPACE_PATTERN, RAW_FIELD_TYPE_PATTERN, Constant, Enum,
                   Error, ErrorCase, Field, FieldType, Method, Model, Params,
                   Result, Service, Spec, Struct)


@dataclass
class ParseFilesResult:
    ignored_prop_uris: list[str]
    specs: list[Spec]


def parse_files(file_path_2_file_data: dict[str, str]) -> ParseFilesResult:
    parser = _Parser()
    for file_path, file_data in file_path_2_file_data.items():
        parser.parse_file(file_data, file_path)
    return ParseFilesResult(
        ignored_prop_uris=parser.ignored_prop_uris(),
        specs=parser.specs(),
    )


class _Parser:
    def __init__(self) -> None:
        self._ignored_prop_uris: list[str] = []
        self._specs: list[Spec] = []

    def parse_file(self, file_data: str, file_path: str) -> None:
        raw_spec = yaml.safe_load(file_data)
        prop_uri = file_path + "#/"
        spec = Spec(prop_uri)
        self._parse_raw_spec(raw_spec, spec)
        self._specs.append(spec)

    def _parse_raw_spec(self, raw_spec, spec: Spec) -> None:
        raw_spec = _ensure_prop_type(raw_spec, dict, spec.prop_uri)
        raw_spec = raw_spec.copy()
        if (namespace := raw_spec.pop("namespace", None)) is not None:
            prop_uri = spec.prop_uri + "namespace"
            namespace = _ensure_prop_type(namespace, str, prop_uri)
            _check_namespace(namespace, prop_uri)
            spec.namespace = namespace
        if (raw_services := raw_spec.pop("services", None)) is not None:
            self._parse_raw_services(
                raw_services,
                spec.services,
                spec.prop_uri + "services",
            )
        if (raw_methods := raw_spec.pop("methods", None)) is not None:
            self._parse_raw_methods(
                raw_methods, spec.methods, spec.prop_uri + "methods"
            )
        if (raw_models := raw_spec.pop("models", None)) is not None:
            self._parse_raw_models(raw_models, spec.models, spec.prop_uri + "models")
        if (raw_errors := raw_spec.pop("errors", None)) is not None:
            self._parse_raw_errors(raw_errors, spec.errors, spec.prop_uri + "errors")
        for prop_name in raw_spec.keys():
            self._ignored_prop_uris.append(spec.prop_uri + prop_name)

    def _parse_raw_services(
        self, raw_services, services: list[Service], prop_uri: str
    ) -> None:
        raw_services = _ensure_prop_type(raw_services, dict, prop_uri)
        _ensure_not_empty(raw_services, prop_uri)
        for service_id, raw_service in raw_services.items():
            prop_uri2 = prop_uri + "/" + str(service_id)
            service_id = _ensure_prop_type(service_id, str, prop_uri2)
            _check_id(service_id, prop_uri2)
            service = Service(prop_uri2, service_id)
            self._parse_raw_service(raw_service, service)
            services.append(service)

    def _parse_raw_service(self, raw_service, service: Service) -> None:
        raw_service = _ensure_prop_type(raw_service, dict, service.prop_uri)
        raw_service = raw_service.copy()
        prop_uri = service.prop_uri + "/version"
        version = _pop_prop(raw_service, "version", prop_uri)
        version = _ensure_prop_type(version, str, prop_uri)
        service.version = version
        if (description := raw_service.pop("description", None)) is not None:
            description = _ensure_prop_type(
                description, str, service.prop_uri + "/description"
            )
            service.description = description
        if (
            method_path_template := raw_service.pop("method_path_template", None)
        ) is not None:
            method_path_template = _ensure_prop_type(
                method_path_template, str, service.prop_uri + "/method_path_template"
            )
            service.method_path_template = method_path_template
        for prop_name in raw_service.keys():
            self._ignored_prop_uris.append(service.prop_uri + "/" + prop_name)

    def _parse_raw_methods(
        self, raw_methods, methods: list[Method], prop_uri: str
    ) -> None:
        raw_methods = _ensure_prop_type(raw_methods, dict, prop_uri)
        _ensure_not_empty(raw_methods, prop_uri)
        for method_id, raw_method in raw_methods.items():
            prop_uri2 = prop_uri + "/" + str(method_id)
            method_id = _ensure_prop_type(method_id, str, prop_uri2)
            _check_id(method_id, prop_uri2)
            method = Method(prop_uri2, method_id)
            self._parse_raw_method(raw_method, method)
            methods.append(method)

    def _parse_raw_method(self, raw_method, method: Method) -> None:
        raw_method = _ensure_prop_type(raw_method, dict, method.prop_uri)
        raw_method = raw_method.copy()
        prop_uri = method.prop_uri + "/service_id"
        service_id = _pop_prop(raw_method, "service_id", prop_uri)
        service_id = _ensure_prop_type(service_id, str, prop_uri)
        _check_id(service_id, prop_uri)
        method.service_id = service_id
        if (summary := raw_method.pop("summary", None)) is not None:
            summary = _ensure_prop_type(summary, str, method.prop_uri + "/summary")
            method.summary = summary
        if (description := raw_method.pop("description", None)) is not None:
            description = _ensure_prop_type(
                description, str, method.prop_uri + "/description"
            )
            method.description = description
        if (raw_params := raw_method.pop("params", None)) is not None:
            params = Params(method.prop_uri + "/params")
            self._parse_raw_params(raw_params, params)
            method.params = params
        if (raw_result := raw_method.pop("result", None)) is not None:
            result = Result(method.prop_uri + "/result")
            self._parse_raw_result(raw_result, result)
            method.result = result
        if (raw_error_cases := raw_method.pop("error_cases", None)) is not None:
            self._parse_raw_error_cases(
                raw_error_cases, method.error_cases, method.prop_uri + "/error_cases"
            )
        for prop_name in raw_method.keys():
            self._ignored_prop_uris.append(method.prop_uri + "/" + prop_name)

    def _parse_raw_params(self, raw_params, params: Params) -> None:
        self._parse_raw_fields(raw_params, params.fields, params.prop_uri)

    def _parse_raw_result(self, raw_result, result: Result) -> None:
        self._parse_raw_fields(raw_result, result.fields, result.prop_uri)

    def _parse_raw_error_cases(
        self, raw_error_cases, error_cases: list[ErrorCase], prop_uri: str
    ) -> None:
        raw_error_cases = _ensure_prop_type(raw_error_cases, dict, prop_uri)
        _ensure_not_empty(raw_error_cases, prop_uri)
        for error_id, raw_error_case in raw_error_cases.items():
            prop_uri2 = prop_uri + "/" + str(error_id)
            error_id = _ensure_prop_type(error_id, str, prop_uri2)
            _check_id(error_id, prop_uri2)
            error_case = ErrorCase(prop_uri2, error_id)
            self._parse_raw_error_case(raw_error_case, error_case)
            error_cases.append(error_case)

    def _parse_raw_error_case(self, raw_error_case, error_case: ErrorCase) -> None:
        raw_error_case = _ensure_prop_type(raw_error_case, dict, error_case.prop_uri)
        raw_error_case = raw_error_case.copy()
        if (description := raw_error_case.pop("description", None)) is not None:
            description = _ensure_prop_type(
                description, str, error_case.prop_uri + "/description"
            )
            error_case.description = description
        for prop_name in raw_error_case.keys():
            self._ignored_prop_uris.append(error_case.prop_uri + "/" + prop_name)

    def _parse_raw_models(self, raw_models, models: list[Model], prop_uri: str) -> None:
        raw_models = _ensure_prop_type(raw_models, dict, prop_uri)
        _ensure_not_empty(raw_models, prop_uri)
        for model_id, raw_model in raw_models.items():
            prop_uri2 = prop_uri + "/" + str(model_id)
            model_id = _ensure_prop_type(model_id, str, prop_uri2)
            _check_id(model_id, prop_uri2)
            model = Model(prop_uri2, model_id)
            self._parse_raw_model(raw_model, model)
            models.append(model)

    def _parse_raw_model(self, raw_model, model: Model) -> None:
        raw_model = _ensure_prop_type(raw_model, dict, model.prop_uri)
        raw_model = raw_model.copy()
        prop_uri = model.prop_uri + "/type"
        model_type = _pop_prop(raw_model, "type", prop_uri)
        model_type = _ensure_prop_type(model_type, str, prop_uri)
        _check_model_type(model_type, prop_uri)
        model.type = model_type
        if model_type == MODEL_STRUCT:
            struct = Struct()
            self._parse_raw_struct(raw_model, struct, model.prop_uri)
            model.definition = struct
        elif model_type == MODEL_ENUM:
            enum = Enum()
            self._parse_raw_enum(raw_model, enum, model.prop_uri)
            model.definition = enum
        else:
            assert False, model_type
        for prop_name in raw_model.keys():
            self._ignored_prop_uris.append(model.prop_uri + "/" + prop_name)

    def _parse_raw_struct(
        self, raw_struct: dict, struct: Struct, prop_uri: str
    ) -> None:
        if (raw_fields := raw_struct.pop("fields", None)) is not None:
            self._parse_raw_fields(raw_fields, struct.fields, prop_uri + "/fields")

    def _parse_raw_fields(self, raw_fields, fields: list[Field], prop_uri: str) -> None:
        raw_fields = _ensure_prop_type(raw_fields, dict, prop_uri)
        _ensure_not_empty(raw_fields, prop_uri)
        for field_id, raw_field in raw_fields.items():
            prop_uri2 = prop_uri + "/" + str(field_id)
            field_id = _ensure_prop_type(field_id, str, prop_uri2)
            _check_id(field_id, prop_uri2)
            field = Field(prop_uri2, field_id)
            self._parse_raw_field(raw_field, field)
            fields.append(field)

    def _parse_raw_field(self, raw_field, field: Field) -> None:
        raw_field = _ensure_prop_type(raw_field, dict, field.prop_uri)
        raw_field = raw_field.copy()
        prop_uri = field.prop_uri + "/type"
        raw_field_type = _pop_prop(raw_field, "type", prop_uri)
        raw_field_type = _ensure_prop_type(raw_field_type, str, prop_uri)
        field_type = FieldType(prop_uri)
        _parse_raw_field_type(raw_field_type, field_type)
        field.type = field_type
        if (description := raw_field.pop("description", None)) is not None:
            description = _ensure_prop_type(
                description, str, field.prop_uri + "/description"
            )
            field.description = description
        if (
            not field_type.is_model_ref
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
            prop_uri2 = field.prop_uri + "/example"
            if field_type.is_repeated:
                example = _ensure_prop_type(example, list, prop_uri2)
                if len(example) == 0:
                    _ensure_not_empty(example, prop_uri2)
                for i, v in enumerate(example):
                    _ensure_prop_type(v, type, prop_uri2 + f"[{i}]")
            else:
                _ensure_prop_type(example, type, prop_uri2)
            field.example = example
        for prop_name in raw_field.keys():
            self._ignored_prop_uris.append(field.prop_uri + "/" + prop_name)

    def _parse_raw_enum(self, raw_enum: dict, enum: Enum, prop_uri: str) -> None:
        prop_uri2 = prop_uri + "/underlying_type"
        enum_underlying_type = _pop_prop(raw_enum, "underlying_type", prop_uri2)
        _ensure_prop_type(enum_underlying_type, str, prop_uri2)
        _check_enum_underlying_type(enum_underlying_type, prop_uri2)
        enum.underlying_type = enum_underlying_type
        if (raw_fields := raw_enum.pop("constants", None)) is not None:
            type = {
                FIELD_INT32: int,
                FIELD_INT64: int,
                FIELD_STRING: str,
            }[enum_underlying_type]
            self._parse_raw_constants(
                raw_fields, enum.constants, prop_uri + "/constants", type
            )

    def _parse_raw_constants(
        self, raw_constants, constants: list[Constant], prop_uri: str, type: Type
    ) -> None:
        raw_constants = _ensure_prop_type(raw_constants, dict, prop_uri)
        _ensure_not_empty(raw_constants, prop_uri)
        for constant_id, raw_constant in raw_constants.items():
            prop_uri2 = prop_uri + "/" + str(constant_id)
            constant_id = _ensure_prop_type(constant_id, str, prop_uri2)
            _check_id(constant_id, prop_uri2)
            constant = Constant(prop_uri2, constant_id)
            self._parse_raw_constant(raw_constant, constant, type)
            constants.append(constant)

    def _parse_raw_constant(self, raw_constant, constant: Constant, type: Type) -> None:
        raw_constant = _ensure_prop_type(raw_constant, dict, constant.prop_uri)
        raw_constant = raw_constant.copy()
        prop_uri = constant.prop_uri + "/value"
        value = _pop_prop(raw_constant, "value", prop_uri)
        value = _ensure_prop_type(value, type, prop_uri)
        constant.value = value
        if (description := raw_constant.pop("description", None)) is not None:
            description = _ensure_prop_type(
                description, str, constant.prop_uri + "/description"
            )
            constant.description = description
        for prop_name in raw_constant.keys():
            self._ignored_prop_uris.append(constant.prop_uri + "/" + prop_name)

    def _parse_raw_errors(self, raw_errors, errors: list[Error], prop_uri: str) -> None:
        raw_errors = _ensure_prop_type(raw_errors, dict, prop_uri)
        _ensure_not_empty(raw_errors, prop_uri)
        for error_id, raw_error in raw_errors.items():
            prop_uri2 = prop_uri + "/" + str(error_id)
            error_id = _ensure_prop_type(error_id, str, prop_uri2)
            _check_id(error_id, prop_uri2)
            error = Error(prop_uri2, error_id)
            self._parse_raw_error(raw_error, error)
            errors.append(error)

    def _parse_raw_error(self, raw_error, error: Error) -> None:
        raw_error = _ensure_prop_type(raw_error, dict, error.prop_uri)
        raw_error = raw_error.copy()
        prop_uri = error.prop_uri + "/code"
        code = _pop_prop(raw_error, "code", prop_uri)
        code = _ensure_prop_type(code, int, prop_uri)
        error.code = code
        for prop_name in raw_error.keys():
            self._ignored_prop_uris.append(error.prop_uri + "/" + prop_name)

    def ignored_prop_uris(self) -> list[str]:
        return self._ignored_prop_uris

    def specs(self) -> list[Spec]:
        return self._specs


class InvalidSpecError(Exception):
    def __init__(self, message: str) -> None:
        super().__init__("invalid specification: " + message)


def _pop_prop(object: dict, prop_name: str, prop_uri: str):
    prop_value = object.pop(prop_name, None)
    if prop_value is None:
        raise InvalidSpecError(f"missing property: prop_uri={prop_uri!r}")
    return prop_value


_T = TypeVar("_T")


def _ensure_prop_type(prop_value, expected_prop_type: Type[_T], prop_uri: str) -> _T:
    if not isinstance(prop_value, expected_prop_type):
        raise InvalidSpecError(
            f"invalid property type: prop_uri={prop_uri!r} prop_type={type(prop_value).__name__} expected_prop_type={expected_prop_type.__name__}"
        )
    return prop_value


def _ensure_not_empty(prop_value: Union[list, dict], prop_uri: str):
    if len(prop_value) == 0:
        raise InvalidSpecError(f"property should not be empty: prop_uri={prop_uri!r}")


def _check_namespace(namespace: str, prop_uri: str) -> None:
    if NAMESPACE_PATTERN.fullmatch(namespace) is None:
        raise InvalidSpecError(
            f"invalid namespace; prop_uri={prop_uri!r} namespace={namespace!r} expected_pattern={NAMESPACE_PATTERN.pattern!r}"
        )


def _check_id(id: str, prop_uri: str) -> None:
    if ID_PATTERN.fullmatch(id) is None:
        raise InvalidSpecError(
            f"invalid id; prop_uri={prop_uri!r} id={id!r} expected_pattern={ID_PATTERN.pattern!r}"
        )


def _check_model_type(model_type: str, prop_uri: str) -> None:
    if MODEL_TYPE_PATTERN.fullmatch(model_type) is None:
        raise InvalidSpecError(
            f"invalid model type; prop_uri={prop_uri!r} model_type={model_type!r} expected_pattern={MODEL_TYPE_PATTERN.pattern!r}"
        )


def _parse_raw_field_type(raw_field_type: str, field_type: FieldType) -> None:
    _check_raw_field_type(raw_field_type, field_type.prop_uri)
    s = raw_field_type
    if s[0].isupper():
        field_type.is_model_ref = True
        if (i := s.find(".")) >= 0:
            field_type.namespace = s[:i]
            s = s[i + 1 :]
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
    field_type.value = s


def _check_raw_field_type(raw_field_type: str, prop_uri: str) -> None:
    if RAW_FIELD_TYPE_PATTERN.fullmatch(raw_field_type) is None:
        raise InvalidSpecError(
            f"invalid field type; prop_uri={prop_uri!r} field_type={raw_field_type!r} expected_pattern={RAW_FIELD_TYPE_PATTERN.pattern!r}"
        )


def _check_enum_underlying_type(enum_underlying_type: str, prop_uri: str) -> None:
    if ENUM_UNDERLYING_TYPE_PATTERN.fullmatch(enum_underlying_type) is None:
        raise InvalidSpecError(
            f"invalid enum underlying type; prop_uri={prop_uri!r} enum_underlying_type={enum_underlying_type!r} expected_pattern={ENUM_UNDERLYING_TYPE_PATTERN.pattern!r}"
        )
