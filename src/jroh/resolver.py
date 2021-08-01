from dataclasses import dataclass

from .spec import (
    MODEL_STRUCT,
    Error,
    ErrorCase,
    Field,
    Method,
    Model,
    Params,
    Results,
    Service,
    Spec,
)


@dataclass
class ResolveSpecsResults:
    unused_node_uris: list[str]
    merged_specs: list[Spec]


def resolve_specs(specs: list[Spec]) -> ResolveSpecsResults:
    resolver = _Resolver()
    resolver.resolve_specs(specs)
    return ResolveSpecsResults(
        unused_node_uris=resolver.unused_node_uris(),
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
        self._unused_node_uris: list[str] = []
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
                raise InvalidSpecError(
                    f"duplicate service id; node_uri1={service.node_uri!r} node_uri2={service2.node_uri!r}",
                )
            self._services[(spec.namespace, service.id)] = service
        for method in spec.methods:
            if (method2 := self._methods.get((spec.namespace, method.id))) is not None:
                raise InvalidSpecError(
                    f"duplicate method id; node_uri1={method.node_uri!r} node_uri2={method2.node_uri!r}",
                )
            self._methods[(spec.namespace, method.id)] = method
        for model in spec.models:
            if (model2 := self._models.get((spec.namespace, model.id))) is not None:
                raise InvalidSpecError(
                    f"duplicate model id; node_uri1={model.node_uri!r} node_uri2={model2.node_uri!r}",
                )
            model.namespace = spec.namespace
            self._models[(spec.namespace, model.id)] = model
        for error in spec.errors:
            if (
                error2 := self._errors_by_id.get((spec.namespace, error.id))
            ) is not None:
                raise InvalidSpecError(
                    f"duplicate error id; node_uri1={error.node_uri!r} node_uri2={error2.node_uri!r}",
                )
            self._errors_by_id[(spec.namespace, error.id)] = error
        for error in spec.errors:
            if (
                error2 := self._errors_by_code.get((spec.namespace, error.code))
            ) is not None:
                node_uri1 = error.node_uri + "/code"
                node_uri2 = error2.node_uri + "/code"
                raise InvalidSpecError(
                    f"duplicate error code; node_uri1={node_uri1!r} node_uri2={node_uri2!r} error_code={error.code}",
                )
            self._errors_by_code[(spec.namespace, error.code)] = error

    def _resolve_spec(self, spec: Spec) -> None:
        self._namespace = spec.namespace
        for method in spec.methods:
            self._resolve_method(method)

    def _resolve_method(self, method: Method) -> None:
        for i, service_id in enumerate(method.service_ids):
            service = self._services.get((self._namespace, service_id))
            if service is None:
                node_uri = method.node_uri + f"/service_ids[{i}]"
                raise InvalidSpecError(
                    f"service not found; node_uri={node_uri!r} service_id={service_id!r}",
                )
            service.methods.append(method)
        if method.params is not None:
            self._resolve_params(method.params)
        if method.results is not None:
            self._resolve_results(method.results)
        for error_case in method.error_cases:
            self._resolve_error_case(error_case)
        error_cases: dict[int, ErrorCase] = {}
        for error_case in method.error_cases:
            error = error_case.error
            assert error is not None
            if (error_case2 := error_cases.get(error.code)) is not None:
                raise InvalidSpecError(
                    f"error code conflict; node_uri1={error_case.node_uri!r} node_uri2={error_case2.node_uri!r} error_code={error.code!r}",
                )
            error_cases[error.code] = error_case

    def _resolve_params(self, params: Params) -> None:
        for field in params.fields:
            self._resolve_field(field)

    def _resolve_results(self, results: Results) -> None:
        for field in results.fields:
            self._resolve_field(field)

    def _resolve_field(self, field: Field) -> None:
        model_ref = field.type.model_ref
        if model_ref is None:
            return
        namespace = model_ref.namespace
        if namespace is None:
            namespace = self._namespace
        model_id = model_ref.id
        model = self._models.get((namespace, model_id))
        if model is None:
            node_uri = field.node_uri + "/type"
            raise InvalidSpecError(
                f"model not found; node_uri={node_uri!r} namespace={namespace!r} model_id={model_id!r}",
            )
        field.type.model = model
        model.ref_count += 1
        if model.ref_count == 1:
            self._resolve_model(model)

    def _resolve_model(self, model: Model) -> None:
        if model.type == MODEL_STRUCT:
            namespace = self._namespace
            self._namespace = model.namespace
            for field in model.struct().fields:
                self._resolve_field(field)
            self._namespace = namespace

    def _resolve_error_case(self, error_case: ErrorCase) -> None:
        error_ref = error_case.error_ref
        namespace = error_ref.namespace
        if namespace is None:
            namespace = self._namespace
        error_id = error_ref.id
        error = self._errors_by_id.get((namespace, error_id))
        if error is None:
            error_case.node_uri + "/type"
            raise InvalidSpecError(
                f"error not found; node_uri={error_case.node_uri!r} namespace={namespace!r} error_id={error_id!r}",
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
                self._unused_node_uris.append(service.node_uri)
            else:
                spec = get_spec(namespace)
                spec.services.append(service)
        for (namespace, _), method in self._methods.items():
            spec = get_spec(namespace)
            spec.methods.append(method)
        for (namespace, _), model in self._models.items():
            if model.ref_count == 0:
                self._unused_node_uris.append(model.node_uri)
            else:
                spec = get_spec(namespace)
                spec.models.append(model)
        for (namespace, _), error in self._errors_by_id.items():
            if error.ref_count == 0:
                self._unused_node_uris.append(error.node_uri)
            else:
                spec = get_spec(namespace)
                spec.errors.append(error)
        self._merged_specs = list(merged_specs.values())

    def unused_node_uris(self) -> list[str]:
        return self._unused_node_uris

    def merged_specs(self) -> list[Spec]:
        return self._merged_specs


class InvalidSpecError(Exception):
    def __init__(self, message: str) -> None:
        super().__init__("invalid spec: " + message)
