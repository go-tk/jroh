from dataclasses import dataclass

from .spec import (MODEL_STRUCT, Error, ErrorCase, Field, Method, Model,
                   Params, Result, Service, Spec)


@dataclass
class ResolveSpecsResult:
    unused_prop_uris: list[str]
    merged_specs: list[Spec]


def resolve_specs(specs: list[Spec]) -> ResolveSpecsResult:
    resolver = _Resolver()
    resolver.resolve_specs(specs)
    return ResolveSpecsResult(
        unused_prop_uris=resolver.unused_prop_uris(),
        merged_specs=resolver.merged_specs(),
    )


class _Resolver:
    def __init__(self):
        self._services: dict[tuple[str, str], Service] = {}
        self._methods: dict[tuple[str, str], Method] = {}
        self._models: dict[tuple[str, str], Model] = {}
        self._errors_by_id: dict[tuple[str, str], Error] = {}
        self._errors_by_code: dict[tuple[str, int], Error] = {}
        self._namespace: str = ""
        self._unused_prop_uris: list[str] = []
        self._merged_specs: list[Spec] = []

    def resolve_specs(self, specs: list[Spec]) -> None:
        for spec in specs:
            self._load_spec(spec)
        for spec in specs:
            self._resolve_spec(spec)
        self._merge_specs()

    def _load_spec(self, spec: Spec) -> None:
        for service in spec.services:
            if (
                service2 := self._services.get((spec.namespace, service.id))
            ) is not None:
                raise InvalidSpec(
                    f"duplicate service id; prop_uri1={service.prop_uri!r} prop_uri2={service2.prop_uri!r}",
                )
            self._services[(spec.namespace, service.id)] = service
        for method in spec.methods:
            if (method2 := self._methods.get((spec.namespace, method.id))) is not None:
                raise InvalidSpec(
                    f"duplicate method id; prop_uri1={method.prop_uri!r} prop_uri2={method2.prop_uri!r}",
                )
            self._methods[(spec.namespace, method.id)] = method
        for model in spec.models:
            if (model2 := self._models.get((spec.namespace, model.id))) is not None:
                raise InvalidSpec(
                    f"duplicate model id; prop_uri1={model.prop_uri!r} prop_uri2={model2.prop_uri!r}",
                )
            self._models[(spec.namespace, model.id)] = model
        for error in spec.errors:
            if (
                error2 := self._errors_by_id.get((spec.namespace, error.id))
            ) is not None:
                raise InvalidSpec(
                    f"duplicate error id; prop_uri1={error.prop_uri!r} prop_uri2={error2.prop_uri!r}",
                )
            self._errors_by_id[(spec.namespace, error.id)] = error
        for error in spec.errors:
            if (
                error2 := self._errors_by_code.get((spec.namespace, error.code))
            ) is not None:
                prop_uri1 = error.prop_uri + "/code"
                prop_uri2 = error2.prop_uri + "/code"
                raise InvalidSpec(
                    f"duplicate error code; prop_uri1={prop_uri1!r} prop_uri2={prop_uri2!r} error_code={error.code}",
                )
            self._errors_by_code[(spec.namespace, error.code)] = error

    def _resolve_spec(self, spec: Spec) -> None:
        self._namespace = spec.namespace
        for method in spec.methods:
            self._resolve_method(method)

    def _resolve_method(self, method: Method) -> None:
        service = self._services.get((self._namespace, method.service_id))
        if service is None:
            prop_uri = method.prop_uri + "/service_id"
            raise InvalidSpec(
                f"service not found; prop_uri={prop_uri!r} service_id={method.service_id!r}",
            )
        service.methods.append(method)
        if method.params is not None:
            self._resolve_params(method.params)
        if method.result is not None:
            self._resolve_result(method.result)
        for error_case in method.error_cases:
            self._resolve_error_case(error_case)

    def _resolve_params(self, params: Params) -> None:
        for field in params.fields:
            self._resolve_field(field)

    def _resolve_result(self, result: Result) -> None:
        for field in result.fields:
            self._resolve_field(field)

    def _resolve_field(self, field: Field) -> None:
        if not field.type.is_model_ref:
            return
        if field.type.namespace is None:
            namespace = self._namespace
        else:
            namespace = field.type.namespace
        model_id = field.type.value
        model = self._models.get((namespace, model_id))
        if model is None:
            prop_uri = field.prop_uri + "/type"
            raise InvalidSpec(
                f"model not found; prop_uri={prop_uri!r} namespace={namespace!r} model_id={model_id!r}",
            )
        field.type.model = model
        model.ref_count += 1
        if model.ref_count == 1:
            self._resolve_model(model)

    def _resolve_model(self, model: Model) -> None:
        if model.type == MODEL_STRUCT:
            for field in model.struct().fields:
                self._resolve_field(field)

    def _resolve_error_case(self, error_case: ErrorCase) -> None:
        error = self._errors_by_id.get((self._namespace, error_case.error_id))
        if error is None:
            error_case.prop_uri + "/type"
            raise InvalidSpec(
                f"error not found; prop_uri={error_case.prop_uri!r}",
            )
        error_case.error = error
        error.ref_count += 1

    def _merge_specs(self):
        merged_specs: dict[str, Spec] = {}

        def get_spec(namespace: str) -> Spec:
            spec = merged_specs.get(namespace)
            if spec is None:
                spec = Spec("")
                spec.namespace = namespace
                merged_specs[namespace] = spec
            return spec

        for (namespace, _), service in self._services.items():
            if len(service.methods) == 0:
                self._unused_prop_uris.append(service.prop_uri)
            else:
                spec = get_spec(namespace)
                spec.services.append(service)
        for (namespace, _), method in self._methods.items():
            spec = get_spec(namespace)
            spec.methods.append(method)
        for (namespace, _), model in self._models.items():
            if model.ref_count == 0:
                self._unused_prop_uris.append(model.prop_uri)
            else:
                spec = get_spec(namespace)
                spec.models.append(model)
        for (namespace, _), error in self._errors_by_id.items():
            if error.ref_count == 0:
                self._unused_prop_uris.append(error.prop_uri)
            else:
                spec = get_spec(namespace)
                spec.errors.append(error)
        self._merged_specs = list(merged_specs.values())

    def unused_prop_uris(self) -> list[str]:
        return self._unused_prop_uris

    def merged_specs(self) -> list[Spec]:
        return self._merged_specs


class InvalidSpec(Exception):
    def __init__(self, message: str) -> None:
        super().__init__("invalid spec: " + message)
