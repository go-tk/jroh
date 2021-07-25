import unittest

from ..jroh.parser import InvalidSpecError
from . import common


class TestParser(unittest.TestCase):
    def test_spec(self):
        test_data_list = [
            common.TestData(
                in_file_path_2_file_data={
                    "foo.yaml": """
- 1
"""
                },
                out_exception_type=InvalidSpecError,
                out_exception_str=r"invalid specification: invalid property type: prop_uri='foo\.yaml#/' prop_type=list expected_prop_type=dict",
            ),
            common.TestData(
                in_file_path_2_file_data={
                    "foo.yaml": """
namespace: 1
"""
                },
                out_exception_type=InvalidSpecError,
                out_exception_str=r"invalid specification: invalid property type: prop_uri='foo\.yaml#/namespace' prop_type=int expected_prop_type=str",
            ),
            common.TestData(
                in_file_path_2_file_data={
                    "foo.yaml": """
namespace: aaa
"""
                },
                out_exception_type=InvalidSpecError,
                out_exception_str=r"invalid specification: invalid namespace; prop_uri='foo\.yaml#/namespace' namespace='aaa'",
            ),
            common.TestData(
                in_file_path_2_file_data={
                    "foo.yaml": """
services: 1
"""
                },
                out_exception_type=InvalidSpecError,
                out_exception_str=r"invalid specification: invalid property type: prop_uri='foo\.yaml#/services' prop_type=int expected_prop_type=dict",
            ),
            common.TestData(
                in_file_path_2_file_data={
                    "foo.yaml": """
services: {}
"""
                },
                out_exception_type=InvalidSpecError,
                out_exception_str=r"invalid specification: property should not be empty: prop_uri='foo\.yaml#/services'",
            ),
            common.TestData(
                in_file_path_2_file_data={
                    "foo.yaml": """
methods: 1
"""
                },
                out_exception_type=InvalidSpecError,
                out_exception_str=r"invalid specification: invalid property type: prop_uri='foo\.yaml#/methods' prop_type=int expected_prop_type=dict",
            ),
            common.TestData(
                in_file_path_2_file_data={
                    "foo.yaml": """
methods: {}
"""
                },
                out_exception_type=InvalidSpecError,
                out_exception_str=r"invalid specification: property should not be empty: prop_uri='foo\.yaml#/methods'",
            ),
            common.TestData(
                in_file_path_2_file_data={
                    "foo.yaml": """
models: 1
"""
                },
                out_exception_type=InvalidSpecError,
                out_exception_str=r"invalid specification: invalid property type: prop_uri='foo\.yaml#/models' prop_type=int expected_prop_type=dict",
            ),
            common.TestData(
                in_file_path_2_file_data={
                    "foo.yaml": """
models: {}
"""
                },
                out_exception_type=InvalidSpecError,
                out_exception_str=r"invalid specification: property should not be empty: prop_uri='foo\.yaml#/models'",
            ),
            common.TestData(
                in_file_path_2_file_data={
                    "foo.yaml": """
errors: 1
"""
                },
                out_exception_type=InvalidSpecError,
                out_exception_str=r"invalid specification: invalid property type: prop_uri='foo\.yaml#/errors' prop_type=int expected_prop_type=dict",
            ),
            common.TestData(
                in_file_path_2_file_data={
                    "foo.yaml": """
errors: {}
"""
                },
                out_exception_type=InvalidSpecError,
                out_exception_str=r"invalid specification: property should not be empty: prop_uri='foo\.yaml#/errors'",
            ),
            common.TestData(
                in_file_path_2_file_data={
                    "foo.yaml": """
xyz: 1
"""
                },
                out_ignored_prop_uris={"foo.yaml#/xyz"},
            ),
        ]
        common.test(self, test_data_list)

    def test_service(self):
        test_data_list = [
            common.TestData(
                in_file_path_2_file_data={
                    "foo.yaml": """
services:
  aaa: {}
"""
                },
                out_exception_type=InvalidSpecError,
                out_exception_str=r"invalid specification: invalid id; prop_uri='foo\.yaml#/services/aaa' id='aaa'",
            ),
            common.TestData(
                in_file_path_2_file_data={
                    "foo.yaml": """
services:
  AAA: 1
"""
                },
                out_exception_type=InvalidSpecError,
                out_exception_str=r"invalid specification: invalid property type: prop_uri='foo\.yaml#/services/AAA' prop_type=int expected_prop_type=dict",
            ),
            common.TestData(
                in_file_path_2_file_data={
                    "foo.yaml": """
services:
  AAA: {}
"""
                },
                out_exception_type=InvalidSpecError,
                out_exception_str=r"invalid specification: missing property: prop_uri='foo\.yaml#/services/AAA/version'",
            ),
            common.TestData(
                in_file_path_2_file_data={
                    "foo.yaml": """
services:
  AAA:
    version: 1
"""
                },
                out_exception_type=InvalidSpecError,
                out_exception_str=r"invalid specification: invalid property type: prop_uri='foo\.yaml#/services/AAA/version' prop_type=int expected_prop_type=str",
            ),
            common.TestData(
                in_file_path_2_file_data={
                    "foo.yaml": """
services:
  AAA:
    version: 1.0.0
    description: 1
"""
                },
                out_exception_type=InvalidSpecError,
                out_exception_str=r"invalid property type: prop_uri='foo\.yaml#/services/AAA/description' prop_type=int expected_prop_type=str",
            ),
            common.TestData(
                in_file_path_2_file_data={
                    "foo.yaml": """
services:
  AAA:
    version: 1.0.0
    method_path_template: 1
"""
                },
                out_exception_type=InvalidSpecError,
                out_exception_str=r"invalid property type: prop_uri='foo\.yaml#/services/AAA/method_path_template' prop_type=int expected_prop_type=str",
            ),
            common.TestData(
                in_file_path_2_file_data={
                    "foo.yaml": """
services:
  AAA:
    version: 1.0.0
    xyz: 1
"""
                },
                out_ignored_prop_uris={"foo.yaml#/services/AAA/xyz"},
            ),
        ]
        common.test(self, test_data_list)

    def test_method(self):
        test_data_list = [
            common.TestData(
                in_file_path_2_file_data={
                    "foo.yaml": """
methods:
  aaa: {}
"""
                },
                out_exception_type=InvalidSpecError,
                out_exception_str=r"invalid specification: invalid id; prop_uri='foo\.yaml#/methods/aaa' id='aaa'",
            ),
            common.TestData(
                in_file_path_2_file_data={
                    "foo.yaml": """
methods:
  AAA: 1
"""
                },
                out_exception_type=InvalidSpecError,
                out_exception_str=r"invalid specification: invalid property type: prop_uri='foo\.yaml#/methods/AAA' prop_type=int expected_prop_type=dict",
            ),
            common.TestData(
                in_file_path_2_file_data={
                    "foo.yaml": """
methods:
  AAA: {}
"""
                },
                out_exception_type=InvalidSpecError,
                out_exception_str=r"invalid specification: missing property: prop_uri='foo\.yaml#/methods/AAA/service_id'",
            ),
            common.TestData(
                in_file_path_2_file_data={
                    "foo.yaml": """
methods:
  AAA:
    service_id: 1
"""
                },
                out_exception_type=InvalidSpecError,
                out_exception_str=r"invalid specification: invalid property type: prop_uri='foo\.yaml#/methods/AAA/service_id' prop_type=int expected_prop_type=str",
            ),
            common.TestData(
                in_file_path_2_file_data={
                    "foo.yaml": """
methods:
  AAA:
    service_id: foo
"""
                },
                out_exception_type=InvalidSpecError,
                out_exception_str=r"invalid specification: invalid id; prop_uri='foo\.yaml#/methods/AAA/service_id' id='foo'",
            ),
            common.TestData(
                in_file_path_2_file_data={
                    "foo.yaml": """
methods:
  AAA:
    service_id: Foo
    summary: 1
"""
                },
                out_exception_type=InvalidSpecError,
                out_exception_str=r"invalid property type: prop_uri='foo\.yaml#/methods/AAA/summary' prop_type=int expected_prop_type=str",
            ),
            common.TestData(
                in_file_path_2_file_data={
                    "foo.yaml": """
methods:
  AAA:
    service_id: Foo
    description: 1
"""
                },
                out_exception_type=InvalidSpecError,
                out_exception_str=r"invalid property type: prop_uri='foo\.yaml#/methods/AAA/description' prop_type=int expected_prop_type=str",
            ),
            common.TestData(
                in_file_path_2_file_data={
                    "foo.yaml": """
methods:
  AAA:
    service_id: Foo
    params: 1
"""
                },
                out_exception_type=InvalidSpecError,
                out_exception_str=r"invalid property type: prop_uri='foo\.yaml#/methods/AAA/params' prop_type=int expected_prop_type=dict",
            ),
            common.TestData(
                in_file_path_2_file_data={
                    "foo.yaml": """
methods:
  AAA:
    service_id: Foo
    params: {}
"""
                },
                out_exception_type=InvalidSpecError,
                out_exception_str=r"invalid specification: property should not be empty: prop_uri='foo\.yaml#/methods/AAA/params'",
            ),
            common.TestData(
                in_file_path_2_file_data={
                    "foo.yaml": """
methods:
  AAA:
    service_id: Foo
    result: 1
"""
                },
                out_exception_type=InvalidSpecError,
                out_exception_str=r"invalid property type: prop_uri='foo\.yaml#/methods/AAA/result' prop_type=int expected_prop_type=dict",
            ),
            common.TestData(
                in_file_path_2_file_data={
                    "foo.yaml": """
methods:
  AAA:
    service_id: Foo
    result: {}
"""
                },
                out_exception_type=InvalidSpecError,
                out_exception_str=r"invalid specification: property should not be empty: prop_uri='foo\.yaml#/methods/AAA/result'",
            ),
            common.TestData(
                in_file_path_2_file_data={
                    "foo.yaml": """
methods:
  AAA:
    service_id: Foo
    error_cases: 1
"""
                },
                out_exception_type=InvalidSpecError,
                out_exception_str=r"invalid property type: prop_uri='foo\.yaml#/methods/AAA/error_cases' prop_type=int expected_prop_type=dict",
            ),
            common.TestData(
                in_file_path_2_file_data={
                    "foo.yaml": """
methods:
  AAA:
    service_id: Foo
    error_cases: {}
"""
                },
                out_exception_type=InvalidSpecError,
                out_exception_str=r"invalid specification: property should not be empty: prop_uri='foo\.yaml#/methods/AAA/error_cases'",
            ),
            common.TestData(
                in_file_path_2_file_data={
                    "foo.yaml": """
services:
  Foo:
    version: 1.0.1
methods:
  AAA:
    service_id: Foo
    xyz: 1
"""
                },
                out_ignored_prop_uris={"foo.yaml#/methods/AAA/xyz"},
            ),
        ]
        common.test(self, test_data_list)

    def test_params(self):
        test_data_list = [
            common.TestData(
                in_file_path_2_file_data={
                    "foo.yaml": """
methods:
  AAA:
    service_id: Foo
    params:
      aaa: {}
"""
                },
                out_exception_type=InvalidSpecError,
                out_exception_str=r"invalid specification: invalid id; prop_uri='foo\.yaml#/methods/AAA/params/aaa' id='aaa'",
            ),
        ]
        common.test(self, test_data_list)

    def test_result(self):
        test_data_list = [
            common.TestData(
                in_file_path_2_file_data={
                    "foo.yaml": """
methods:
  AAA:
    service_id: Foo
    result:
      aaa: {}
"""
                },
                out_exception_type=InvalidSpecError,
                out_exception_str=r"invalid specification: invalid id; prop_uri='foo\.yaml#/methods/AAA/result/aaa' id='aaa'",
            ),
        ]
        common.test(self, test_data_list)

    def test_error_case(self):
        test_data_list = [
            common.TestData(
                in_file_path_2_file_data={
                    "foo.yaml": """
methods:
  AAA:
    service_id: Foo
    error_cases:
      eee: {}
"""
                },
                out_exception_type=InvalidSpecError,
                out_exception_str=r"invalid specification: invalid id; prop_uri='foo\.yaml#/methods/AAA/error_cases/eee' id='eee'",
            ),
            common.TestData(
                in_file_path_2_file_data={
                    "foo.yaml": """
methods:
  AAA:
    service_id: Foo
    error_cases:
      EEE:
        description: 1
"""
                },
                out_exception_type=InvalidSpecError,
                out_exception_str=r"invalid specification: invalid property type: prop_uri='foo\.yaml#/methods/AAA/error_cases/EEE/description' prop_type=int expected_prop_type=str",
            ),
            common.TestData(
                in_file_path_2_file_data={
                    "foo.yaml": """
services:
  Foo:
    version: 1.1.1
methods:
  AAA:
    service_id: Foo
    error_cases:
      EEE:
        xyz: 1
errors:
  EEE:
    code: 123
"""
                },
                out_ignored_prop_uris={"foo.yaml#/methods/AAA/error_cases/EEE/xyz"},
            ),
        ]
        common.test(self, test_data_list)

    def test_model(self):
        test_data_list = [
            common.TestData(
                in_file_path_2_file_data={
                    "foo.yaml": """
models:
  mmm: {}
"""
                },
                out_exception_type=InvalidSpecError,
                out_exception_str=r"invalid specification: invalid id; prop_uri='foo\.yaml#/models/mmm' id='mmm'",
            ),
            common.TestData(
                in_file_path_2_file_data={
                    "foo.yaml": """
models:
  MMM: {}
"""
                },
                out_exception_type=InvalidSpecError,
                out_exception_str=r"invalid specification: missing property: prop_uri='foo\.yaml#/models/MMM/type'",
            ),
            common.TestData(
                in_file_path_2_file_data={
                    "foo.yaml": """
models:
  MMM:
    type: 1
"""
                },
                out_exception_type=InvalidSpecError,
                out_exception_str=r"invalid specification: invalid property type: prop_uri='foo\.yaml#/models/MMM/type' prop_type=int expected_prop_type=str",
            ),
            common.TestData(
                in_file_path_2_file_data={
                    "foo.yaml": """
models:
  MMM:
    type: ttt
"""
                },
                out_exception_type=InvalidSpecError,
                out_exception_str=r"invalid specification: invalid model type; prop_uri='foo\.yaml#/models/MMM/type' model_type='ttt'",
            ),
            common.TestData(
                in_file_path_2_file_data={
                    "foo.yaml": """
models:
  MMM:
    type: struct
    xyz: 1
"""
                },
                out_ignored_prop_uris={"foo.yaml#/models/MMM/xyz"},
            ),
        ]
        common.test(self, test_data_list)

    def test_struct(self):
        test_data_list = [
            common.TestData(
                in_file_path_2_file_data={
                    "foo.yaml": """
models:
  MMM:
    type: struct
    fields: 1
"""
                },
                out_exception_type=InvalidSpecError,
                out_exception_str=r"invalid specification: invalid property type: prop_uri='foo\.yaml#/models/MMM/fields' prop_type=int expected_prop_type=dict",
            ),
            common.TestData(
                in_file_path_2_file_data={
                    "foo.yaml": """
models:
  MMM:
    type: struct
    fields: {}
"""
                },
                out_exception_type=InvalidSpecError,
                out_exception_str=r"invalid specification: property should not be empty: prop_uri='foo\.yaml#/models/MMM/fields'",
            ),
        ]
        common.test(self, test_data_list)

    def test_field(self):
        test_data_list = [
            common.TestData(
                in_file_path_2_file_data={
                    "foo.yaml": """
models:
  MMM:
    type: struct
    fields:
      fff: {}
"""
                },
                out_exception_type=InvalidSpecError,
                out_exception_str=r"invalid specification: invalid id; prop_uri='foo\.yaml#/models/MMM/fields/fff' id='fff'",
            ),
            common.TestData(
                in_file_path_2_file_data={
                    "foo.yaml": """
models:
  MMM:
    type: struct
    fields:
      FFF: {}
"""
                },
                out_exception_type=InvalidSpecError,
                out_exception_str=r"invalid specification: missing property: prop_uri='foo\.yaml#/models/MMM/fields/FFF/type'",
            ),
            common.TestData(
                in_file_path_2_file_data={
                    "foo.yaml": """
models:
  MMM:
    type: struct
    fields:
      FFF:
        type: ttt
"""
                },
                out_exception_type=InvalidSpecError,
                out_exception_str=r"invalid specification: invalid field type; prop_uri='foo\.yaml#/models/MMM/fields/FFF/type' field_type='ttt'",
            ),
            common.TestData(
                in_file_path_2_file_data={
                    "foo.yaml": """
models:
  MMM:
    type: struct
    fields:
      FFF:
        type: bool
        example: 1.1
"""
                },
                out_exception_type=InvalidSpecError,
                out_exception_str=r"invalid specification: invalid property type: prop_uri='foo\.yaml#/models/MMM/fields/FFF/example' prop_type=float expected_prop_type=bool",
            ),
            common.TestData(
                in_file_path_2_file_data={
                    "foo.yaml": """
models:
  MMM:
    type: struct
    fields:
      FFF:
        type: int32
        example: 1.1
"""
                },
                out_exception_type=InvalidSpecError,
                out_exception_str=r"invalid specification: invalid property type: prop_uri='foo\.yaml#/models/MMM/fields/FFF/example' prop_type=float expected_prop_type=int",
            ),
            common.TestData(
                in_file_path_2_file_data={
                    "foo.yaml": """
models:
  MMM:
    type: struct
    fields:
      FFF:
        type: int64
        example: 1.1
"""
                },
                out_exception_type=InvalidSpecError,
                out_exception_str=r"invalid specification: invalid property type: prop_uri='foo\.yaml#/models/MMM/fields/FFF/example' prop_type=float expected_prop_type=int",
            ),
            common.TestData(
                in_file_path_2_file_data={
                    "foo.yaml": """
models:
  MMM:
    type: struct
    fields:
      FFF:
        type: float32
        example: sss
"""
                },
                out_exception_type=InvalidSpecError,
                out_exception_str=r"invalid specification: invalid property type: prop_uri='foo\.yaml#/models/MMM/fields/FFF/example' prop_type=str expected_prop_type=float",
            ),
            common.TestData(
                in_file_path_2_file_data={
                    "foo.yaml": """
models:
  MMM:
    type: struct
    fields:
      FFF:
        type: float64
        example: sss
"""
                },
                out_exception_type=InvalidSpecError,
                out_exception_str=r"invalid specification: invalid property type: prop_uri='foo\.yaml#/models/MMM/fields/FFF/example' prop_type=str expected_prop_type=float",
            ),
            common.TestData(
                in_file_path_2_file_data={
                    "foo.yaml": """
models:
  MMM:
    type: struct
    fields:
      FFF:
        type: string
        example: 1
"""
                },
                out_exception_type=InvalidSpecError,
                out_exception_str=r"invalid specification: invalid property type: prop_uri='foo\.yaml#/models/MMM/fields/FFF/example' prop_type=int expected_prop_type=str",
            ),
            common.TestData(
                in_file_path_2_file_data={
                    "foo.yaml": """
models:
  MMM:
    type: struct
    fields:
      FFF:
        type: string?
        example: 1
"""
                },
                out_exception_type=InvalidSpecError,
                out_exception_str=r"invalid specification: invalid property type: prop_uri='foo\.yaml#/models/MMM/fields/FFF/example' prop_type=int expected_prop_type=str",
            ),
            common.TestData(
                in_file_path_2_file_data={
                    "foo.yaml": """
models:
  MMM:
    type: struct
    fields:
      FFF:
        type: string*
        example: 1
"""
                },
                out_exception_type=InvalidSpecError,
                out_exception_str=r"invalid specification: invalid property type: prop_uri='foo\.yaml#/models/MMM/fields/FFF/example' prop_type=int expected_prop_type=list",
            ),
            common.TestData(
                in_file_path_2_file_data={
                    "foo.yaml": """
models:
  MMM:
    type: struct
    fields:
      FFF:
        type: string*
        example: [1]
"""
                },
                out_exception_type=InvalidSpecError,
                out_exception_str=r"invalid specification: invalid property type: prop_uri='foo\.yaml#/models/MMM/fields/FFF/example\[0\]' prop_type=int expected_prop_type=str",
            ),
            common.TestData(
                in_file_path_2_file_data={
                    "foo.yaml": """
models:
  MMM:
    type: struct
    fields:
      FFF:
        type: string+
        example: 1
"""
                },
                out_exception_type=InvalidSpecError,
                out_exception_str=r"invalid specification: invalid property type: prop_uri='foo\.yaml#/models/MMM/fields/FFF/example' prop_type=int expected_prop_type=list",
            ),
            common.TestData(
                in_file_path_2_file_data={
                    "foo.yaml": """
models:
  MMM:
    type: struct
    fields:
      FFF:
        type: string+
        example: [1]
"""
                },
                out_exception_type=InvalidSpecError,
                out_exception_str=r"invalid specification: invalid property type: prop_uri='foo\.yaml#/models/MMM/fields/FFF/example\[0\]' prop_type=int expected_prop_type=str",
            ),
            common.TestData(
                in_file_path_2_file_data={
                    "foo.yaml": """
models:
  MMM:
    type: struct
    fields:
      FFF:
        type: NNN
        example:
          X: 100
  NNN:
    type: struct
    fields:
      X:
        type: int32
"""
                },
                out_ignored_prop_uris={"foo.yaml#/models/MMM/fields/FFF/example"},
            ),
            common.TestData(
                in_file_path_2_file_data={
                    "foo.yaml": """
models:
  MMM:
    type: struct
    fields:
      FFF:
        type: Hello.World
        description: 1
"""
                },
                out_exception_type=InvalidSpecError,
                out_exception_str=r"invalid property type: prop_uri='foo\.yaml#/models/MMM/fields/FFF/description' prop_type=int expected_prop_type=str",
            ),
            common.TestData(
                in_file_path_2_file_data={
                    "foo.yaml": """
models:
  MMM:
    type: struct
    fields:
      FFF:
        type: string
        xyz: 1
"""
                },
                out_ignored_prop_uris={"foo.yaml#/models/MMM/fields/FFF/xyz"},
            ),
        ]
        common.test(self, test_data_list)

    def test_enum(self):
        test_data_list = [
            common.TestData(
                in_file_path_2_file_data={
                    "foo.yaml": """
models:
  MMM:
    type: enum
"""
                },
                out_exception_type=InvalidSpecError,
                out_exception_str=r"invalid specification: missing property: prop_uri='foo\.yaml#/models/MMM/underlying_type'",
            ),
            common.TestData(
                in_file_path_2_file_data={
                    "foo.yaml": """
models:
  MMM:
    type: enum
    underlying_type: vvv
"""
                },
                out_exception_type=InvalidSpecError,
                out_exception_str=r"invalid specification: invalid enum underlying type; prop_uri='foo\.yaml#/models/MMM/underlying_type' enum_underlying_type='vvv'",
            ),
            common.TestData(
                in_file_path_2_file_data={
                    "foo.yaml": """
models:
  MMM:
    type: enum
    underlying_type: int32
    constants: 1
"""
                },
                out_exception_type=InvalidSpecError,
                out_exception_str=r"invalid specification: invalid property type: prop_uri='foo\.yaml#/models/MMM/constants' prop_type=int expected_prop_type=dict",
            ),
            common.TestData(
                in_file_path_2_file_data={
                    "foo.yaml": """
models:
  MMM:
    type: enum
    underlying_type: int32
    constants: {}
"""
                },
                out_exception_type=InvalidSpecError,
                out_exception_str=r"invalid specification: property should not be empty: prop_uri='foo\.yaml#/models/MMM/constants'",
            ),
        ]
        common.test(self, test_data_list)

    def test_constant(self):
        test_data_list = [
            common.TestData(
                in_file_path_2_file_data={
                    "foo.yaml": """
models:
  MMM:
    type: enum
    underlying_type: int32
    constants:
      ccc: {}
"""
                },
                out_exception_type=InvalidSpecError,
                out_exception_str=r"invalid specification: invalid id; prop_uri='foo\.yaml#/models/MMM/constants/ccc' id='ccc'",
            ),
            common.TestData(
                in_file_path_2_file_data={
                    "foo.yaml": """
models:
  MMM:
    type: enum
    underlying_type: int32
    constants:
      CCC: {}
"""
                },
                out_exception_type=InvalidSpecError,
                out_exception_str=r"invalid specification: missing property: prop_uri='foo\.yaml#/models/MMM/constants/CCC/value'",
            ),
            common.TestData(
                in_file_path_2_file_data={
                    "foo.yaml": """
models:
  MMM:
    type: enum
    underlying_type: int32
    constants:
      CCC:
        value: 1.1
"""
                },
                out_exception_type=InvalidSpecError,
                out_exception_str=r"invalid specification: invalid property type: prop_uri='foo\.yaml#/models/MMM/constants/CCC/value' prop_type=float expected_prop_type=int",
            ),
            common.TestData(
                in_file_path_2_file_data={
                    "foo.yaml": """
models:
  MMM:
    type: enum
    underlying_type: int64
    constants:
      CCC:
        value: vvv
"""
                },
                out_exception_type=InvalidSpecError,
                out_exception_str=r"invalid specification: invalid property type: prop_uri='foo\.yaml#/models/MMM/constants/CCC/value' prop_type=str expected_prop_type=int",
            ),
            common.TestData(
                in_file_path_2_file_data={
                    "foo.yaml": """
models:
  MMM:
    type: enum
    underlying_type: string
    constants:
      CCC:
        value: 100
"""
                },
                out_exception_type=InvalidSpecError,
                out_exception_str=r"invalid specification: invalid property type: prop_uri='foo\.yaml#/models/MMM/constants/CCC/value' prop_type=int expected_prop_type=str",
            ),
            common.TestData(
                in_file_path_2_file_data={
                    "foo.yaml": """
models:
  MMM:
    type: enum
    underlying_type: string
    constants:
      CCC:
        value: foo
        description: 1
"""
                },
                out_exception_type=InvalidSpecError,
                out_exception_str=r"invalid specification: invalid property type: prop_uri='foo\.yaml#/models/MMM/constants/CCC/description' prop_type=int expected_prop_type=str",
            ),
            common.TestData(
                in_file_path_2_file_data={
                    "foo.yaml": """
models:
  MMM:
    type: enum
    underlying_type: string
    constants:
      CCC:
        value: foo
        xyz: 1
"""
                },
                out_ignored_prop_uris={"foo.yaml#/models/MMM/constants/CCC/xyz"},
            ),
        ]
        common.test(self, test_data_list)

    def test_error(self):
        test_data_list = [
            common.TestData(
                in_file_path_2_file_data={
                    "foo.yaml": """
errors:
  eee: {}
"""
                },
                out_exception_type=InvalidSpecError,
                out_exception_str=r"invalid specification: invalid id; prop_uri='foo\.yaml#/errors/eee' id='eee'",
            ),
            common.TestData(
                in_file_path_2_file_data={
                    "foo.yaml": """
errors:
  EEE: 1
"""
                },
                out_exception_type=InvalidSpecError,
                out_exception_str=r"invalid specification: invalid property type: prop_uri='foo\.yaml#/errors/EEE' prop_type=int expected_prop_type=dict",
            ),
            common.TestData(
                in_file_path_2_file_data={
                    "foo.yaml": """
errors:
  EEE: {}
"""
                },
                out_exception_type=InvalidSpecError,
                out_exception_str=r"invalid specification: missing property: prop_uri='foo\.yaml#/errors/EEE/code'",
            ),
            common.TestData(
                in_file_path_2_file_data={
                    "foo.yaml": """
errors:
  EEE:
    code: e100
"""
                },
                out_exception_type=InvalidSpecError,
                out_exception_str=r"invalid property type: prop_uri='foo\.yaml#/errors/EEE/code' prop_type=str expected_prop_type=int",
            ),
            common.TestData(
                in_file_path_2_file_data={
                    "foo.yaml": """
errors:
  EEE:
    code: 100
    xyz: 1
"""
                },
                out_ignored_prop_uris={"foo.yaml#/errors/EEE/xyz"},
            ),
        ]
        common.test(self, test_data_list)


if __name__ == "__main__":
    unittest.main()
