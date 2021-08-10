import os
from dataclasses import dataclass
from typing import Callable, Optional

from mako.template import Template

from . import utils
from .spec import MODEL_ENUM, MODEL_STRUCT, Error, Field, Method, Model, Service, Spec


@dataclass
class GenerateCodeResults:
    file_path_2_file_data: dict[str, str]


def generate_code(output_package_path: str, specs: list[Spec]) -> GenerateCodeResults:
    generator = _Generator(output_package_path)
    generator.generate_code(specs)
    return GenerateCodeResults(
        file_path_2_file_data=generator.file_path_2_file_data(),
    )


class _Generator:
    def __init__(self, output_package_path: str) -> None:
        self._output_package_path: str = output_package_path
        self._namespace: str = ""
        self._imports: dict[str, _Import] = {}
        self._buffer: list[str] = []
        self._file_path_2_file_data: dict[str, str] = {}

    def generate_code(self, specs: list[Spec]) -> None:
        for spec in specs:
            self._generate_package_code(spec)
        for file_path in (
            "apicommon/errors.go",
            "apicommon/rpcinfo.go",
            "apicommon/utils.go",
        ):
            self._file_path_2_file_data[file_path] = utils.get_data(
                "data/go/" + file_path
            ).decode()

    def _generate_package_code(self, spec: Spec) -> None:
        self._namespace = spec.namespace
        for service in spec.services:
            self._generate_service_code(service)
        if len(spec.methods) + len(spec.models) >= 1:
            self._generate_models_code(spec.methods, spec.models)
        if len(spec.errors) >= 1:
            self._generate_errors_code(spec.errors)

    def _generate_service_code(self, service: Service) -> None:
        self._buffer.append(
            f"""\
package {self._make_package_name()}
"""
        )
        self._buffer.append("")
        self._buffer.append(
            Template(
                r"""\

<%
    context1 = g._import_package("context", "context")
    http = g._import_package("http", "net/http")
    apicommon = g._import_package("apicommon", "../apicommon")

    service_name = utils.pascal_case(service.id)
%>\
type ${service_name}Service interface {
% for method in service.methods:
<%
    method_name = utils.pascal_case(method.id)
%>\
    ${method_name}(ctx ${context1()}.Context\
% if method.params is not None:
, params *${method_name}Params\
% endif
% if method.results is not None:
, results *${method_name}Results\
% endif
) (err error)
% endfor
}

type Dummy${service_name}Service struct {}

var _ ${service_name}Service = Dummy${service_name}Service{}

% for method in service.methods:
<%
    method_name = utils.pascal_case(method.id)
%>\
func (Dummy${service_name}Service) ${method_name}(${context1()}.Context\
% if method.params is not None:
, *${method_name}Params\
% endif
% if method.results is not None:
, *${method_name}Results\
% endif
) error { return nil }
% endfor

func ForEachRPCHandlerOf${service_name}Service(service ${service_name}Service, callback func(
    rpcPath string,
    rpcHandler ${http()}.HandlerFunc,
    rpcInfoFactory ${apicommon()}.RPCInfoFactory,
) (ok bool)) {
<%
    namespace = utils.pascal_case(g._namespace)
    service_name = utils.pascal_case(service.id)
%>\
% for i, method in enumerate(service.methods):
<%
    method_name = utils.pascal_case(method.id)
    rpc_path = service.rpc_paths[i]
%>\
    {
        rpcHandler := func(w ${http()}.ResponseWriter, r *${http()}.Request) {
            ctx := r.Context()
            rpcInfo := ${apicommon()}.RPCInfoFromContext(ctx)
            var data struct {
% if method.params is not None:
                Params ${method_name}Params
% endif
                Error ${apicommon()}.Error
% if method.results is not None:
                Results ${method_name}Results
% endif
                Resp ${method_name}Resp
            }
% if method.params is not None:
            rpcInfo.SetParams(&data.Params)
% endif
            rpcInfo.SetError(&data.Error)
% if method.results is not None:
            rpcInfo.SetResults(&data.Results)
% endif
            defer func() {
                if panicValue := recover(); panicValue != nil {
                    ${apicommon()}.SavePanicValue(panicValue, rpcInfo)
                }
                data.Resp.ID = rpcInfo.ID()
                if data.Error.Code == 0 {
                    rpcInfo.SetError(nil)
% if method.results is not None:
                    data.Resp.Results = &data.Results
% endif
                } else {
% if method.results is not None:
                    rpcInfo.SetResults(nil)
% endif
                    data.Resp.Error = &data.Error
                }
                ${apicommon()}.WriteResp(&data.Resp, w, rpcInfo)
            }()
% if method.params is not None:
            if !${apicommon()}.ReadParams(r.Body, rpcInfo) {
                return
            }
% endif
            err := service.${method_name}(ctx\
% if method.params is not None:
, &data.Params\
% endif
% if method.results is not None:
, &data.Results\
% endif
)
            ${apicommon()}.SaveErr(err, rpcInfo)
        }
        rpcInfoFactory := func(id string) *${apicommon()}.RPCInfo { return ${apicommon()}.NewRPCInfo("${namespace}", "${service_name}", "${method_name}", id) }
        if !callback(${utils.quote(rpc_path)}, rpcHandler, rpcInfoFactory) {
            return
        }
    }
% endfor
}
"""
            ).render(
                utils=utils,
                g=self,
                service=service,
            )
        )
        self._buffer[1] = self._generate_imports_code()
        file_name = f"{utils.flat_case(service.id)}service.go"
        self._flush(file_name)

    def _generate_models_code(self, methods: list[Method], models: list[Model]) -> None:
        self._buffer.append(
            f"""\
package {self._make_package_name()}
"""
        )
        self._buffer.append("")
        self._buffer.append(
            Template(
                r"""\
<%
    apicommon = g._import_package("apicommon", "../apicommon")
%>\
% for method in methods:
<%
    method_name = utils.pascal_case(method.id)
%>\
% if method.params is not None:

type ${method_name}Params struct {
% for field in method.params.fields:
    ${g._generate_field_code(field)}
% endfor
}
% endif

type ${method_name}Resp struct {
    ID string `json:"id"`
    Error *${apicommon()}.Error `json:"error,omitempty"`
% if method.results is not None:
    Results *${method_name}Results `json:"results,omitempty"`
% endif
}
% if method.results is not None:

type ${method_name}Results struct {
% for field in method.results.fields:
    ${g._generate_field_code(field)}
% endfor
}
% endif
% endfor
% for model in models:
<%
    model_name = utils.pascal_case(model.id)
%>\

% if model.type == MODEL_STRUCT:
<%
    struct = model.struct()
%>\
type ${model_name} struct {
% for field in struct.fields:
    ${g._generate_field_code(field)}
% endfor
}
% elif model.type == MODEL_ENUM:
<%
    enum = model.enum()
%>\
type ${model_name} ${enum.underlying_type}

const (
% for constant in enum.constants:
        ${utils.pascal_case(constant.id)} ${model_name} = \
% if isinstance(constant.value, int):
${constant.value}
% elif isinstance(constant.value, str):
${utils.quote(constant.value)}
% else:
<%
    assert False, type(constant.value)
%>\
% endif
% endfor
)
% else:
<%
    assert False, model.type
%>\
% endif
% endfor
"""
            ).render(
                utils=utils,
                MODEL_STRUCT=MODEL_STRUCT,
                MODEL_ENUM=MODEL_ENUM,
                g=self,
                methods=methods,
                models=models,
            )
        )
        self._buffer[1] = self._generate_imports_code()
        self._flush("models.go")

    def _generate_errors_code(self, errors: list[Error]) -> None:
        self._buffer.append(
            f"""\
package {self._make_package_name()}
"""
        )
        self._buffer.append("")
        self._buffer.append(
            Template(
                r"""\
<%
    apicommon = g._import_package("apicommon", "../apicommon")
%>\
% for error in errors:
<%
    error_name = utils.pascal_case(error.id.removesuffix("-Error"))
%>\

const Error${error_name} ${apicommon()}.ErrorCode = ${error.code}

var Err${error_name} = &${apicommon()}.Error{
    Code:    Error${error_name},
    Message: "${error.id.lower().replace("-", " ")}",
}
% endfor
"""
            ).render(
                utils=utils,
                g=self,
                errors=errors,
            )
        )
        self._buffer[1] = self._generate_imports_code()
        self._flush("errors.go")

    def _generate_imports_code(self) -> str:
        imports = [
            (import_name, import1)
            for import_name, import1 in self._imports.items()
            if import1.ref_count >= 1
        ]
        if len(imports) == 0:
            return ""
        imports.sort(key=lambda x: x[0])
        imports_code = Template(
            r"""\

import (
% for import_name, import1 in imports:
    ${import_name} ${utils.quote(import1.package_path)}
% endfor
)
"""
        ).render(
            utils=utils,
            imports=imports,
        )
        return imports_code

    def _generate_field_code(self, field: Field) -> str:
        field_type = field.type
        field_type_code = ""
        if field_type.is_repeated:
            field_type_code += "[]"
        else:
            if field_type.is_optional:
                field_type_code += "*"
        if (model_ref := field_type.model_ref) is None:
            field_type_code += field_type.value
        else:
            namespace = model_ref.namespace
            if namespace is None:
                namespace = self._namespace
            if namespace != self._namespace:
                package_name = self._make_package_name(namespace)
                import_name = self._import_package(package_name, "../" + package_name)()
                field_type_code += import_name + "."
            field_type_code += utils.pascal_case(model_ref.id)
        json_tag = utils.camel_case(field.id)
        if field_type.is_optional:
            json_tag += ",omitempty"
        filed_code = (
            f'{utils.pascal_case(field.id)} {field_type_code} `json:"{json_tag}"`'
        )
        return filed_code

    def _make_package_name(self, namespace: Optional[str] = None) -> str:
        if namespace is None:
            namespace = self._namespace
        return utils.flat_case(namespace) + "api"

    def _import_package(
        self, package_name: str, package_path: str
    ) -> Callable[[], str]:
        if package_path.startswith("../"):
            package_path = self._output_package_path + package_path[2:]
        n = 1
        import_name = package_name
        while True:
            import1 = self._imports.get(import_name)
            if import1 is None:
                import1 = _Import(package_path, 0)
                self._imports[package_name] = import1
                break
            if import1.package_path == package_path:
                break
            n += 1
            import_name = package_name + str(n)

        def package() -> str:
            import1.ref_count += 1
            return import_name

        return package

    def _flush(self, file_name: str) -> None:
        file_path = os.path.join(self._make_package_name(), file_name)
        file_data = "".join(self._buffer)
        self._file_path_2_file_data[file_path] = file_data
        self._imports.clear()
        self._buffer.clear()

    def file_path_2_file_data(self) -> dict[str, str]:
        return self._file_path_2_file_data


@dataclass
class _Import:
    package_path: str
    ref_count: int
