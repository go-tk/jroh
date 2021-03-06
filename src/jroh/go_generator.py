import os
import subprocess
import sys
from dataclasses import dataclass
from typing import Callable, Optional

from mako.template import Template

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
    Error,
    Field,
    Method,
    Model,
    Service,
    Spec,
)
from .version import VERSION


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
        self._patterns: list[str] = []
        self._buffer: list[str] = []
        self._file_path_2_file_data: dict[str, str] = {}

    def generate_code(self, specs: list[Spec]) -> None:
        for spec in specs:
            self._generate_package_code(spec)

    def _generate_package_code(self, spec: Spec) -> None:
        self._namespace = spec.namespace
        for service in spec.services:
            self._generate_service_code(service)
        for service in spec.services:
            self._generate_client_code(service)
        if len(spec.methods) + len(spec.models) >= 1:
            self._generate_models_code(spec.methods, spec.models)
        if len(spec.errors) >= 1:
            self._generate_errors_code(spec.errors)
        if len(spec.services) >= 1:
            self._generate_misc_code(spec.services)

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
    apicommon = g._import_package("apicommon", "github.com/go-tk/jroh/go/apicommon")
    http = g._import_package("http", "net/http")

    namespace = utils.pascal_case(g._namespace)
    service_name = utils.pascal_case(service.id)
%>\
type ${service_name}Actor interface {
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

func Register${service_name}Actor(a ${service_name}Actor, router *${apicommon()}.Router, options ${apicommon()}.ActorOptions) {
    options.Sanitize()
    var rpcFiltersTable [NumberOf${service_name}Methods][]${apicommon()}.IncomingRPCHandler
    ${apicommon()}.FillIncomingRPCFiltersTable(rpcFiltersTable[:], options.RPCFilters)
% for i, method in enumerate(service.methods):
<%
    method_name = utils.pascal_case(method.id)
    full_method_name = f"{namespace}.{service_name}.{method_name}"
    rpc_path = service.rpc_paths[i]
%>\
    {
        rpcFilters := rpcFiltersTable[${service_name}_${method_name}]
        handler := ${http()}.HandlerFunc(func (w ${http()}.ResponseWriter, r *${http()}.Request) {
            var s struct {
                rpc ${apicommon()}.IncomingRPC
                params \
    % if method.params is None:
${apicommon()}.DummyModel
    % else:
${method_name}Params
    % endif
                results \
    % if method.results is None:
${apicommon()}.DummyModel
    % else:
${method_name}Results
    % endif
            }
            s.rpc.Namespace = "${namespace}"
            s.rpc.ServiceName = "${service_name}"
            s.rpc.MethodName = "${method_name}"
            s.rpc.FullMethodName = "${full_method_name}"
            s.rpc.MethodIndex = ${service_name}_${method_name}
            s.rpc.Params = &s.params
            s.rpc.Results = &s.results
            s.rpc.SetHandler(func(ctx ${context1()}.Context, rpc *${apicommon()}.IncomingRPC) error {
                return a.${method_name}(ctx\
    % if method.params is not None:
, rpc.Params.(*${method_name}Params)\
    % endif
    % if method.results is not None:
, rpc.Results.(*${method_name}Results)\
    % endif
)
            })
            s.rpc.SetFilters(rpcFilters)
            ${apicommon()}.HandleRequest(r, &s.rpc, options.TraceIDGenerator, w)
        })
        router.AddRoute(${utils.quote(rpc_path)}, handler, "${full_method_name}", rpcFilters)
    }
% endfor
}

type ${service_name}ActorFuncs struct {
% for method in service.methods:
<%
    method_name = utils.pascal_case(method.id)
%>\
    ${method_name}Func func(${context1()}.Context\
    % if method.params is not None:
, *${method_name}Params\
    % endif
    % if method.results is not None:
, *${method_name}Results\
    % endif
) error
% endfor
}

var _ ${service_name}Actor = (*${service_name}ActorFuncs)(nil)
% for method in service.methods:
<%
    method_name = utils.pascal_case(method.id)
%>\

func (sf *${service_name}ActorFuncs) ${method_name}(ctx ${context1()}.Context\
    % if method.params is not None:
, params *${method_name}Params\
    % endif
    % if method.results is not None:
, results *${method_name}Results\
    % endif
) error {
    if f := sf.${method_name}Func; f != nil {
        return f(ctx\
    % if method.params is not None:
, params\
    % endif
    % if method.results is not None:
, results\
    % endif
)
    }
    return ${apicommon()}.NewNotImplementedError()
}
% endfor
"""
            ).render(
                utils=utils,
                g=self,
                service=service,
            )
        )
        self._buffer[1] = self._generate_imports_code()
        file_name = f"{utils.flat_case(service.id)}actor_generated.go"
        self._flush(file_name)

    def _generate_client_code(self, service: Service) -> None:
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
    apicommon = g._import_package("apicommon", "github.com/go-tk/jroh/go/apicommon")
    http = g._import_package("http", "net/http")
    fmt = g._import_package("fmt", "fmt")

    namespace = utils.pascal_case(g._namespace)
    service_name = utils.pascal_case(service.id)
    service_name2 = utils.camel_case(service.id)
%>\
type ${service_name}Client interface {
% for i, method in enumerate(service.methods):
<%
    method_name = utils.pascal_case(method.id)
    rpc_path = service.rpc_paths[i]
%>\
    ${method_name}(ctx ${context1()}.Context\
    % if method.params is not None:
, params *${method_name}Params\
    % endif
) (\
    % if method.results is not None:
results *${method_name}Results,\
    % endif
err error)
% endfor
}

type ${service_name2}Client struct {
    rpcBaseURL string
    options ${apicommon()}.ClientOptions
    rpcFiltersTable [NumberOf${service_name}Methods][]${apicommon()}.OutgoingRPCHandler
}

func New${service_name}Client(rpcBaseURL string, options ${apicommon()}.ClientOptions) ${service_name}Client {
    var c ${service_name2}Client
    c.rpcBaseURL = rpcBaseURL
    c.options = options
    c.options.Sanitize()
    ${apicommon()}.FillOutgoingRPCFiltersTable(c.rpcFiltersTable[:], options.RPCFilters)
    return &c
}
% for i, method in enumerate(service.methods):
<%
    method_name = utils.pascal_case(method.id)
    full_method_name = f"{namespace}.{service_name}.{method_name}"
    rpc_path = service.rpc_paths[i]
%>\

func (c *${service_name2}Client) ${method_name}(ctx ${context1()}.Context\
    % if method.params is not None:
, params *${method_name}Params\
    % endif
) (\
    % if method.results is not None:
*${method_name}Results,\
    % endif
error) {
    var s struct {
        rpc ${apicommon()}.OutgoingRPC
        params \
    % if method.params is None:
${apicommon()}.DummyModel
    % else:
${method_name}Params
    % endif
        results \
    % if method.results is None:
${apicommon()}.DummyModel
    % else:
${method_name}Results
    % endif
    }
    s.rpc.Namespace = "${namespace}"
    s.rpc.ServiceName = "${service_name}"
    s.rpc.MethodName = "${method_name}"
    s.rpc.FullMethodName = "${full_method_name}"
    s.rpc.MethodIndex = ${service_name}_${method_name}
    % if method.params is not None:
    s.params = *params
    % endif
    s.rpc.Params = &s.params
    s.rpc.Results = &s.results
    if err := c.doRPC(ctx, &s.rpc, ${utils.quote(rpc_path)}); err != nil {
        return \
    % if method.results is not None:
nil, \
    % endif
${fmt()}.Errorf("do rpc; fullMethodName=${utils.quote(utils.quote(full_method_name))[1:-1]} traceID=%q: %w", s.rpc.TraceID, err)
    }
    return \
    % if method.results is not None:
&s.results, \
    % endif
nil
}
% endfor

func (c *${service_name2}Client) doRPC(ctx ${context1()}.Context, rpc *${apicommon()}.OutgoingRPC, rpcPath string) error {
    if timeout := c.options.Timeout; timeout >= 1 {
        var cancel ${context1()}.CancelFunc
        ctx, cancel = ${context1()}.WithTimeout(ctx, timeout)
        defer cancel()
    }
    rpc.Transport = c.options.Transport
    rpc.URL = c.rpcBaseURL + rpcPath
    rpc.SetHandler(${apicommon()}.HandleOutgoingRPC)
    rpc.SetFilters(c.rpcFiltersTable[rpc.MethodIndex])
    return rpc.Do(ctx)
}

type ${service_name}ClientFuncs struct {
% for method in service.methods:
<%
    method_name = utils.pascal_case(method.id)
    full_method_name = f"{namespace}.{service_name}.{method_name}"
%>\
    ${method_name}Func func(${context1()}.Context\
    % if method.params is not None:
, *${method_name}Params\
    % endif
) \
    % if method.results is None:
error
    % else:
(*${method_name}Results, error)
    % endif
% endfor
}

var _ ${service_name}Client = (*${service_name}ClientFuncs)(nil)
% for method in service.methods:
<%
    method_name = utils.pascal_case(method.id)
%>\

func (cf *${service_name}ClientFuncs) ${method_name}(ctx ${context1()}.Context\
    % if method.params is not None:
, params *${method_name}Params\
    % endif
) \
    % if method.results is None:
error {
    % else:
(*${method_name}Results, error) {
    % endif
    if f := cf.${method_name}Func; f != nil {
        return f(ctx\
        % if method.params is not None:
, params\
        % endif
)
    }
    err := ${apicommon()}.NewNotImplementedError()
    return \
    % if method.results is not None:
nil, \
    % endif
${fmt()}.Errorf("do rpc; fullMethodName=${utils.quote(utils.quote(full_method_name))[1:-1]}: %w", err)
}
% endfor
"""
            ).render(
                utils=utils,
                g=self,
                service=service,
            )
        )
        self._buffer[1] = self._generate_imports_code()
        file_name = f"{utils.flat_case(service.id)}client_generated.go"
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
    apicommon = g._import_package("apicommon", "github.com/go-tk/jroh/go/apicommon")
    strconv = g._import_package("strconv", "strconv")
    fmt = g._import_package("fmt", "fmt")
%>\
<%def name="struct_validate_func(model_name, fields)">\
func (m *${model_name}) Validate(validationContext *${apicommon()}.ValidationContext) bool {
    % for field in fields:
<%
    field_name = utils.pascal_case(field.id)
    field_type = field.type

    v = "m." + field_name
%>\
        % if field.is_optional:
            % if not field_type.is_primitive() or field.is_limited():
    if ${v} != nil {
                % if field_type.is_primitive():
        v := *${v}
<%
    v = "v"
%>\
                % endif
        validationContext.Enter("${utils.camel_case(field.id)}")
                % if field_type.is_primitive():
${validate_primitive_value("        ", v, field_type.primitive_type(), field)}\
                % else:
        if !${v}.Validate(validationContext) {
            return false
        }
                % endif
        validationContext.Leave()
    }
            % endif
        % elif field.is_repeated:
            % if field.count_is_limited() or not field_type.is_primitive() or field.is_limited():
    {
        validationContext.Enter("${utils.camel_case(field.id)}")
                % if field.min_count >= 1:
        if len(${v}) < ${field.min_count} {
            validationContext.SetErrorDetails("length < ${field.min_count}")
            return false
        }
                % endif
                % if field.max_count is not None:
        if len(${v}) > ${field.max_count} {
            validationContext.SetErrorDetails("length > ${field.max_count}")
            return false
        }
                % endif
                % if not field_type.is_primitive() or field.is_limited():
                    % if field_type.is_primitive():
        for i, v := range ${v} {
                    % else:
        for i := range ${v} {
            v := &${v}[i]
                    % endif
<%
    v = "v"
%>\
            validationContext.Enter(${strconv()}.Itoa(i))
                    % if field_type.is_primitive():
${validate_primitive_value("            ", v, field_type.primitive_type(), field)}\
                    % else:
            if !${v}.Validate(validationContext) {
                return false
            }
                    % endif
            validationContext.Leave()
        }
                % endif
        validationContext.Leave()
    }
            % endif
        % else:
            % if not field_type.is_primitive() or field.is_limited():
    {
        validationContext.Enter("${utils.camel_case(field.id)}")
                % if field_type.is_primitive():
${validate_primitive_value("        ", v, field_type.primitive_type(), field)}\
                % else:
        if !${v}.Validate(validationContext) {
            return false
        }
                % endif
        validationContext.Leave()
    }
            % endif
        % endif
    % endfor
    mm := struct {
        ${apicommon()}.DummyFurtherValidator
        *${model_name}
    }{${model_name}: m}
    return mm.FurtherValidate(validationContext)
}
</%def>\
<%def name="validate_primitive_value(indent, v, primitive_type, primitive_constraints)">\
    % if primitive_type in (INT32, INT64):
        % if primitive_constraints.min is not None:
${indent}if ${v} < ${primitive_constraints.min} {
${indent}    validationContext.SetErrorDetails("value < ${primitive_constraints.min}")
${indent}    return false
${indent}}
        % endif
        % if primitive_constraints.max is not None:
${indent}if ${v} > ${primitive_constraints.max} {
${indent}    validationContext.SetErrorDetails("value > ${primitive_constraints.max}")
${indent}    return false
${indent}}
        % endif
    % elif primitive_type in (FLOAT32, FLOAT64):
        % if primitive_constraints.min is not None:
<%
    min = str(primitive_constraints.min).removesuffix(".0")
%>\
            % if primitive_constraints.min_is_exclusive:
${indent}if ${v} <= ${min} {
${indent}    validationContext.SetErrorDetails("value <= ${min}")
${indent}    return false
${indent}}
            % else:
${indent}if ${v} < ${min} {
${indent}    validationContext.SetErrorDetails("value < ${min}")
${indent}    return false
${indent}}
            % endif
        % endif
        % if primitive_constraints.max is not None:
<%
    max = str(primitive_constraints.max).removesuffix(".0")
%>\
            % if primitive_constraints.max_is_exclusive:
${indent}if ${v} >= ${max} {
${indent}    validationContext.SetErrorDetails("value >= ${max}")
${indent}    return false
${indent}}
            % else:
${indent}if ${v} > ${max} {
${indent}    validationContext.SetErrorDetails("value > ${max}")
${indent}    return false
${indent}}
            % endif
        % endif
    % elif primitive_type == STRING:
        % if primitive_constraints.min_length >= 1:
${indent}if len(${v}) < ${primitive_constraints.min_length} {
${indent}    validationContext.SetErrorDetails("length < ${primitive_constraints.min_length}")
${indent}    return false
${indent}}
        % endif
        % if primitive_constraints.max_length is not None:
${indent}if len(${v}) > ${primitive_constraints.max_length} {
${indent}    validationContext.SetErrorDetails("length > ${primitive_constraints.max_length}")
${indent}    return false
${indent}}
        % endif
        % if primitive_constraints.pattern != "":
<%
    pattern_index = g._add_pattern(primitive_constraints.pattern)
%>\
${indent}if !patterns[${pattern_index}].MatchString(\
% if v == "m":
string(${v})\
% else:
${v}\
% endif
) {
${indent}    validationContext.SetErrorDetails("value not matched to ${utils.quote(utils.quote(primitive_constraints.pattern))[1:-1]}")
${indent}    return false
${indent}}
        % endif
    % endif
</%def>\
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

var _ ${apicommon()}.Model = (*${method_name}Params)(nil)

${struct_validate_func(method_name + "Params", method.params.fields)}\
    % endif
    % if method.results is not None:

type ${method_name}Results struct {
        % for field in method.results.fields:
    ${g._generate_field_code(field)}
        % endfor
}

var _ ${apicommon()}.Model = (*${method_name}Results)(nil)

${struct_validate_func(method_name + "Results", method.results.fields)}\
    % endif
% endfor
% for model in models:
<%
    model_name = utils.pascal_case(model.id)
%>\

    % if model.type == STRUCT:
<%
    struct = model.struct()
%>\
type ${model_name} struct {
    % for field in struct.fields:
    ${g._generate_field_code(field)}
    % endfor
}

var _ ${apicommon()}.Model = (*${model_name})(nil)

${struct_validate_func(model_name, struct.fields)}\
    % elif model.type == ENUM:
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

var _ ${fmt()}.Stringer = ${model_name}(${primitive_zero_literals[enum.underlying_type]})

func (m ${model_name}) String() string {
    switch m {
        % for constant in enum.constants:
        case ${utils.pascal_case(constant.id)}:
            return "${utils.pascal_case(constant.id)}"
        % endfor
        default:
        % if enum.underlying_type in (INT32, INT64):
            return "${model_name}(" + ${strconv()}.FormatInt(int64(m), 10) + ")"
        % elif enum.underlying_type == STRING:
            return "${model_name}(" + ${strconv()}.Quote(string(m)) + ")"
        % else:
<%
    assert False, enum.underlying_type
%>\
        % endif
    }
}

var _ ${apicommon()}.Model = ${model_name}(${primitive_zero_literals[enum.underlying_type]})

func (m ${model_name}) Validate(validationContext *${apicommon()}.ValidationContext) bool {
    switch m {
        % for constant in enum.constants:
        case ${utils.pascal_case(constant.id)}:
            return true
        % endfor
        default:
        % if enum.underlying_type in (INT32, INT64):
            validationContext.SetErrorDetails("value not in (${", ".join(str(constant.value) for constant in enum.constants)})")
        % elif enum.underlying_type == STRING:
            validationContext.SetErrorDetails("value not in (${", ".join(utils.quote(utils.quote(constant.value))[1:-1] for constant in enum.constants)})")
        % else:
<%
    assert False, enum.underlying_type
%>\
        % endif
            return false
    }
}
    % else:
<%
    xprimit = model.xprimit()
%>\
type ${model_name} ${xprimit.primitive_type}

var _ ${apicommon()}.Model = ${model_name}(${primitive_zero_literals[xprimit.primitive_type]})

func (m ${model_name}) Validate(validationContext *${apicommon()}.ValidationContext) bool {
${validate_primitive_value("    ", "m", xprimit.primitive_type, xprimit)}\
    mm := struct {
        ${apicommon()}.DummyFurtherValidator
        ${model_name}
    }{${model_name}: m}
    return mm.FurtherValidate(validationContext)
}
    % endif
% endfor
"""
            ).render(
                utils=utils,
                STRUCT=STRUCT,
                ENUM=ENUM,
                INT32=INT32,
                INT64=INT64,
                FLOAT32=FLOAT32,
                FLOAT64=FLOAT64,
                STRING=STRING,
                primitive_zero_literals=_primitive_zero_literals,
                g=self,
                methods=methods,
                models=models,
            )
        )
        self._buffer.append(self._generate_patterns_code())
        self._buffer[1] = self._generate_imports_code()
        self._flush("models_generated.go")

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
    apicommon = g._import_package("apicommon", "github.com/go-tk/jroh/go/apicommon")
%>\
% for error in errors:
<%
    error_name = utils.pascal_case(error.id.removesuffix("-Error"))
%>\

const Error${error_name} ${apicommon()}.ErrorCode = ${error.code}

func New${error_name}Error() *${apicommon()}.Error {
    return &${apicommon()}.Error{
        Code: Error${error_name},
        StatusCode: ${error.status_code},
        Message: "${error.id.lower().replace("-", " ")}",
    }
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
        self._flush("errors_generated.go")

    def _generate_misc_code(self, services: list[Service]) -> None:
        self._buffer.append(
            f"""\
package {self._make_package_name()}
"""
        )
        self._buffer.append("")
        self._buffer.append(
            Template(
                r"""\
% for service in services:
<%
    service_name = utils.pascal_case(service.id)
%>\

const (
    % for i, method in enumerate(service.methods):
<%
    method_name = utils.pascal_case(method.id)
%>\
    ${service_name}_${method_name} = ${i}
    % endfor
)

const NumberOf${service_name}Methods = ${len(service.methods)}
% endfor
"""
            ).render(
                utils=utils,
                g=self,
                services=services,
            )
        )
        self._buffer[1] = self._generate_imports_code()
        self._flush("misc_generated.go")

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

    def _generate_patterns_code(self) -> str:
        if len(self._patterns) == 0:
            return ""
        patterns_code = Template(
            r"""\
<%
    regexp = g._import_package("regexp", "regexp")
%>\

var patterns = [...]*regexp.Regexp {
% for i, pattern in enumerate(patterns):
    ${i}: ${regexp()}.MustCompile(${utils.quote(pattern)}),
% endfor
}
"""
        ).render(
            utils=utils,
            g=self,
            patterns=self._patterns,
        )
        return patterns_code

    def _generate_field_code(self, field: Field) -> str:
        field_type = field.type
        field_type_code = ""
        if field.is_repeated:
            field_type_code += "[]"
        else:
            if field.is_optional:
                field_type_code += "*"
        if field_type.is_primitive():
            primitive_type = field_type.primitive_type()
            field_type_code += primitive_type
        else:
            model_ref = field_type.model_ref()
            namespace = model_ref.namespace
            if namespace is None:
                namespace = self._namespace
            if namespace != self._namespace:
                package_name = self._make_package_name(namespace)
                import_name = self._import_package(package_name, "../" + package_name)()
                field_type_code += import_name + "."
            field_type_code += utils.pascal_case(model_ref.id)
        json_tag = utils.camel_case(field.id)
        if field.is_optional or (field.is_repeated and field.min_count == 0):
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

    def _add_pattern(self, pattern: str) -> int:
        if not pattern.startswith("^"):
            pattern = "^" + pattern
        if not pattern.endswith("$"):
            pattern += "$"
        for i, pattern2 in enumerate(self._patterns):
            if pattern2 == pattern:
                return i
        i = len(self._patterns)
        self._patterns.append(pattern)
        return i

    def _flush(self, file_name: str) -> None:
        file_path = os.path.join(self._make_package_name(), file_name)
        file_data = "".join(self._buffer)
        self._file_path_2_file_data[file_path] = file_data
        self._imports.clear()
        self._patterns.clear()
        self._buffer.clear()

    def file_path_2_file_data(self) -> dict[str, str]:
        return self._file_path_2_file_data


@dataclass
class _Import:
    package_path: str
    ref_count: int


_primitive_zero_literals: dict[str, str] = {
    BOOL: "false",
    INT32: "0",
    INT64: "0",
    FLOAT32: "0",
    FLOAT64: "0",
    STRING: '""',
}


def format_go_code(output_file_paths: list[str]) -> None:
    try:
        subprocess.run(["gofmt", "-w", *output_file_paths])
    except Exception as e:
        print(f"WARNING: format go code: {e}", file=sys.stderr)


def update_go_mod_file(output_dir_path: str, output_package_path: str) -> None:
    if output_package_path.startswith("github.com/go-tk/jroh/"):
        return
    try:
        go_mod_file_exists = (
            subprocess.run(
                ["go", "mod", "edit", "-print"],
                cwd=output_dir_path,
                stderr=subprocess.DEVNULL,
                stdout=subprocess.DEVNULL,
            ).returncode
            == 0
        )
        if not go_mod_file_exists:
            subprocess.run(
                ["go", "mod", "init", output_package_path],
                cwd=output_dir_path,
            )
        subprocess.run(
            ["go", "mod", "edit", "-require", "github.com/go-tk/jroh/go@v" + VERSION],
            cwd=output_dir_path,
        )
    except Exception as e:
        print(f"WARNING: update go.mod: {e}", file=sys.stderr)
