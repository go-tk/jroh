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
                out_exception_re=r"invalid specification: invalid node kind: node_uri='foo\.yaml#/' node_kind=sequence expected_node_kind=mapping",
            ),
            common.TestData(
                in_file_path_2_file_data={
                    "foo.yaml": """
namespace: 1
"""
                },
                out_exception_type=InvalidSpecError,
                out_exception_re=r"invalid specification: invalid node kind: node_uri='foo\.yaml#/namespace' node_kind=integer expected_node_kind=string",
            ),
            common.TestData(
                in_file_path_2_file_data={
                    "foo.yaml": """
namespace: aaa
"""
                },
                out_exception_type=InvalidSpecError,
                out_exception_re=r"invalid specification: invalid id; node_uri='foo\.yaml#/namespace' id='aaa'",
            ),
            common.TestData(
                in_file_path_2_file_data={
                    "foo.yaml": """
services: 1
"""
                },
                out_exception_type=InvalidSpecError,
                out_exception_re=r"invalid specification: invalid node kind: node_uri='foo\.yaml#/services' node_kind=integer expected_node_kind=mapping",
            ),
            common.TestData(
                in_file_path_2_file_data={
                    "foo.yaml": """
services: {}
"""
                },
                out_exception_type=InvalidSpecError,
                out_exception_re=r"invalid specification: non-empty mapping required: node_uri='foo\.yaml#/services'",
            ),
            common.TestData(
                in_file_path_2_file_data={
                    "foo.yaml": """
methods: 1
"""
                },
                out_exception_type=InvalidSpecError,
                out_exception_re=r"invalid specification: invalid node kind: node_uri='foo\.yaml#/methods' node_kind=integer expected_node_kind=mapping",
            ),
            common.TestData(
                in_file_path_2_file_data={
                    "foo.yaml": """
methods: {}
"""
                },
                out_exception_type=InvalidSpecError,
                out_exception_re=r"invalid specification: non-empty mapping required: node_uri='foo\.yaml#/methods'",
            ),
            common.TestData(
                in_file_path_2_file_data={
                    "foo.yaml": """
models: 1
"""
                },
                out_exception_type=InvalidSpecError,
                out_exception_re=r"invalid specification: invalid node kind: node_uri='foo\.yaml#/models' node_kind=integer expected_node_kind=mapping",
            ),
            common.TestData(
                in_file_path_2_file_data={
                    "foo.yaml": """
models: {}
"""
                },
                out_exception_type=InvalidSpecError,
                out_exception_re=r"invalid specification: non-empty mapping required: node_uri='foo\.yaml#/models'",
            ),
            common.TestData(
                in_file_path_2_file_data={
                    "foo.yaml": """
errors: 1
"""
                },
                out_exception_type=InvalidSpecError,
                out_exception_re=r"invalid specification: invalid node kind: node_uri='foo\.yaml#/errors' node_kind=integer expected_node_kind=mapping",
            ),
            common.TestData(
                in_file_path_2_file_data={
                    "foo.yaml": """
errors: {}
"""
                },
                out_exception_type=InvalidSpecError,
                out_exception_re=r"invalid specification: non-empty mapping required: node_uri='foo\.yaml#/errors'",
            ),
            common.TestData(
                in_file_path_2_file_data={
                    "foo.yaml": """
xyz: 1
"""
                },
                out_ignored_node_uris={"foo.yaml#/xyz"},
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
                out_exception_re=r"invalid specification: invalid id; node_uri='foo\.yaml#/services/aaa' id='aaa'",
            ),
            common.TestData(
                in_file_path_2_file_data={
                    "foo.yaml": """
services:
  AAA: 1
"""
                },
                out_exception_type=InvalidSpecError,
                out_exception_re=r"invalid specification: invalid node kind: node_uri='foo\.yaml#/services/AAA' node_kind=integer expected_node_kind=mapping",
            ),
            common.TestData(
                in_file_path_2_file_data={
                    "foo.yaml": """
services:
  AAA: {}
"""
                },
                out_exception_type=InvalidSpecError,
                out_exception_re=r"invalid specification: missing node: node_uri='foo\.yaml#/services/AAA/version'",
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
                out_exception_re=r"invalid specification: invalid node kind: node_uri='foo\.yaml#/services/AAA/version' node_kind=integer expected_node_kind=string",
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
                out_exception_re=r"invalid node kind: node_uri='foo\.yaml#/services/AAA/description' node_kind=integer expected_node_kind=string",
            ),
            common.TestData(
                in_file_path_2_file_data={
                    "foo.yaml": """
services:
  AAA:
    version: 1.0.0
    rpc_path_template: 1
"""
                },
                out_exception_type=InvalidSpecError,
                out_exception_re=r"invalid node kind: node_uri='foo\.yaml#/services/AAA/rpc_path_template' node_kind=integer expected_node_kind=string",
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
                out_ignored_node_uris={"foo.yaml#/services/AAA/xyz"},
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
                out_exception_re=r"invalid specification: invalid id; node_uri='foo\.yaml#/methods/aaa' id='aaa'",
            ),
            common.TestData(
                in_file_path_2_file_data={
                    "foo.yaml": """
methods:
  AAA: 1
"""
                },
                out_exception_type=InvalidSpecError,
                out_exception_re=r"invalid specification: invalid node kind: node_uri='foo\.yaml#/methods/AAA' node_kind=integer expected_node_kind=mapping",
            ),
            common.TestData(
                in_file_path_2_file_data={
                    "foo.yaml": """
methods:
  AAA: {}
"""
                },
                out_exception_type=InvalidSpecError,
                out_exception_re=r"invalid specification: missing node: node_uri='foo\.yaml#/methods/AAA/service_id'",
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
                out_exception_re=r"invalid specification: invalid node kind: node_uri='foo\.yaml#/methods/AAA/service_id' node_kind=integer expected_node_kind=string",
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
                out_exception_re=r"invalid specification: invalid id; node_uri='foo\.yaml#/methods/AAA/service_id' id='foo'",
            ),
            common.TestData(
                in_file_path_2_file_data={
                    "foo.yaml": """
methods:
  AAA:
    service_ids: 1
"""
                },
                out_exception_type=InvalidSpecError,
                out_exception_re=r"invalid specification: invalid node kind: node_uri='foo\.yaml#/methods/AAA/service_ids' node_kind=integer expected_node_kind=sequence",
            ),
            common.TestData(
                in_file_path_2_file_data={
                    "foo.yaml": """
methods:
  AAA:
    service_ids:
    - 1
"""
                },
                out_exception_type=InvalidSpecError,
                out_exception_re=r"invalid specification: invalid node kind: node_uri='foo\.yaml#/methods/AAA/service_ids\[0\]' node_kind=integer expected_node_kind=string",
            ),
            common.TestData(
                in_file_path_2_file_data={
                    "foo.yaml": """
methods:
  AAA:
    service_ids:
    - Abc
    - foo
"""
                },
                out_exception_type=InvalidSpecError,
                out_exception_re=r"invalid specification: invalid id; node_uri='foo\.yaml#/methods/AAA/service_ids\[1\]' id='foo'",
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
                out_exception_re=r"invalid node kind: node_uri='foo\.yaml#/methods/AAA/summary' node_kind=integer expected_node_kind=string",
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
                out_exception_re=r"invalid node kind: node_uri='foo\.yaml#/methods/AAA/description' node_kind=integer expected_node_kind=string",
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
                out_exception_re=r"invalid node kind: node_uri='foo\.yaml#/methods/AAA/params' node_kind=integer expected_node_kind=mapping",
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
                out_exception_re=r"invalid specification: non-empty mapping required: node_uri='foo\.yaml#/methods/AAA/params'",
            ),
            common.TestData(
                in_file_path_2_file_data={
                    "foo.yaml": """
methods:
  AAA:
    service_id: Foo
    results: 1
"""
                },
                out_exception_type=InvalidSpecError,
                out_exception_re=r"invalid node kind: node_uri='foo\.yaml#/methods/AAA/results' node_kind=integer expected_node_kind=mapping",
            ),
            common.TestData(
                in_file_path_2_file_data={
                    "foo.yaml": """
methods:
  AAA:
    service_id: Foo
    results: {}
"""
                },
                out_exception_type=InvalidSpecError,
                out_exception_re=r"invalid specification: non-empty mapping required: node_uri='foo\.yaml#/methods/AAA/results'",
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
                out_exception_re=r"invalid node kind: node_uri='foo\.yaml#/methods/AAA/error_cases' node_kind=integer expected_node_kind=mapping",
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
                out_exception_re=r"invalid specification: non-empty mapping required: node_uri='foo\.yaml#/methods/AAA/error_cases'",
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
                out_ignored_node_uris={"foo.yaml#/methods/AAA/xyz"},
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
                out_exception_re=r"invalid specification: invalid id; node_uri='foo\.yaml#/methods/AAA/params/aaa' id='aaa'",
            ),
        ]
        common.test(self, test_data_list)

    def test_results(self):
        test_data_list = [
            common.TestData(
                in_file_path_2_file_data={
                    "foo.yaml": """
methods:
  AAA:
    service_id: Foo
    results:
      aaa: {}
"""
                },
                out_exception_type=InvalidSpecError,
                out_exception_re=r"invalid specification: invalid id; node_uri='foo\.yaml#/methods/AAA/results/aaa' id='aaa'",
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
                out_exception_re=r"invalid specification: invalid ref; node_uri='foo\.yaml#/methods/AAA/error_cases/eee' ref='eee'",
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
                out_exception_re=r"invalid specification: invalid node kind: node_uri='foo\.yaml#/methods/AAA/error_cases/EEE/description' node_kind=integer expected_node_kind=string",
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
    code: 1230
    status_code: 400
"""
                },
                out_ignored_node_uris={"foo.yaml#/methods/AAA/error_cases/EEE/xyz"},
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
                out_exception_re=r"invalid specification: invalid id; node_uri='foo\.yaml#/models/mmm' id='mmm'",
            ),
            common.TestData(
                in_file_path_2_file_data={
                    "foo.yaml": """
models:
  MMM: {}
"""
                },
                out_exception_type=InvalidSpecError,
                out_exception_re=r"invalid specification: missing node: node_uri='foo\.yaml#/models/MMM/type'",
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
                out_exception_re=r"invalid specification: invalid node kind: node_uri='foo\.yaml#/models/MMM/type' node_kind=integer expected_node_kind=string",
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
                out_exception_re=r"invalid specification: invalid model type; node_uri='foo\.yaml#/models/MMM/type' model_type='ttt'",
            ),
            common.TestData(
                in_file_path_2_file_data={
                    "foo.yaml": """
models:
  MMM:
    type: enum
    underlying_type: int32
    description: 1
"""
                },
                out_exception_type=InvalidSpecError,
                out_exception_re=r"invalid specification: invalid node kind: node_uri='foo\.yaml#/models/MMM/description' node_kind=integer expected_node_kind=string",
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
                out_ignored_node_uris={"foo.yaml#/models/MMM/xyz"},
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
                out_exception_re=r"invalid specification: invalid node kind: node_uri='foo\.yaml#/models/MMM/fields' node_kind=integer expected_node_kind=mapping",
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
                out_exception_re=r"invalid specification: non-empty mapping required: node_uri='foo\.yaml#/models/MMM/fields'",
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
                out_exception_re=r"invalid specification: invalid id; node_uri='foo\.yaml#/models/MMM/fields/fff' id='fff'",
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
                out_exception_re=r"invalid specification: missing node: node_uri='foo\.yaml#/models/MMM/fields/FFF/type'",
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
                out_exception_re=r"invalid specification: invalid field type; node_uri='foo\.yaml#/models/MMM/fields/FFF/type' field_type='ttt'",
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
                out_exception_re=r"invalid specification: invalid node kind: node_uri='foo\.yaml#/models/MMM/fields/FFF/example' node_kind=floating-point expected_node_kind=boolean",
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
                out_exception_re=r"invalid specification: invalid node kind: node_uri='foo\.yaml#/models/MMM/fields/FFF/example' node_kind=floating-point expected_node_kind=integer",
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
                out_exception_re=r"invalid specification: invalid node kind: node_uri='foo\.yaml#/models/MMM/fields/FFF/example' node_kind=floating-point expected_node_kind=integer",
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
                out_exception_re=r"invalid specification: invalid node kind: node_uri='foo\.yaml#/models/MMM/fields/FFF/example' node_kind=string expected_node_kind=floating-point",
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
                out_exception_re=r"invalid specification: invalid node kind: node_uri='foo\.yaml#/models/MMM/fields/FFF/example' node_kind=string expected_node_kind=floating-point",
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
                out_exception_re=r"invalid specification: invalid node kind: node_uri='foo\.yaml#/models/MMM/fields/FFF/example' node_kind=integer expected_node_kind=string",
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
        is_optional: true
        example: 1
"""
                },
                out_exception_type=InvalidSpecError,
                out_exception_re=r"invalid specification: invalid node kind: node_uri='foo\.yaml#/models/MMM/fields/FFF/example' node_kind=integer expected_node_kind=string",
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
        is_repeated: true
        example: 1
"""
                },
                out_exception_type=InvalidSpecError,
                out_exception_re=r"invalid specification: invalid node kind: node_uri='foo\.yaml#/models/MMM/fields/FFF/example' node_kind=integer expected_node_kind=sequence",
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
        is_repeated: true
        example: [1, 2, "3"]
"""
                },
                out_exception_type=InvalidSpecError,
                out_exception_re=r"invalid specification: invalid node kind: node_uri='foo\.yaml#/models/MMM/fields/FFF/example\[2\]' node_kind=string expected_node_kind=integer",
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
        is_repeated: true
        min_count: -1
"""
                },
                out_exception_type=InvalidSpecError,
                out_exception_re=r"invalid specification: number too small: node_uri='foo\.yaml#/models/MMM/fields/FFF/min_count' number=-1 min_number=0",
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
        is_repeated: true
        max_count: 0
"""
                },
                out_exception_type=InvalidSpecError,
                out_exception_re=r"invalid specification: number too small: node_uri='foo\.yaml#/models/MMM/fields/FFF/max_count' number=0 min_number=1",
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
        is_repeated: true
        min_count: 2
        max_count: 1
"""
                },
                out_exception_type=InvalidSpecError,
                out_exception_re=r"invalid specification: invalid field, min_count > max_count: node_uri='foo\.yaml#/models/MMM/fields/FFF' min_count=2 max_count=1",
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
        is_repeated: true
        max_count: 2
        example: [1, 2, 3]
"""
                },
                out_exception_type=InvalidSpecError,
                out_exception_re=r"invalid specification: sequence too long: node_uri='foo\.yaml#/models/MMM/fields/FFF/example' sequence_length=3 max_sequence_length=2",
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
        is_repeated: true
        min_count: 2
        example: [1]
"""
                },
                out_exception_type=InvalidSpecError,
                out_exception_re=r"invalid specification: sequence too short: node_uri='foo\.yaml#/models/MMM/fields/FFF/example' sequence_length=1 min_sequence_length=2",
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
                out_ignored_node_uris={"foo.yaml#/models/MMM/fields/FFF/example"},
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
                out_exception_re=r"invalid node kind: node_uri='foo\.yaml#/models/MMM/fields/FFF/description' node_kind=integer expected_node_kind=string",
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
                out_ignored_node_uris={"foo.yaml#/models/MMM/fields/FFF/xyz"},
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
                out_exception_re=r"invalid specification: missing node: node_uri='foo\.yaml#/models/MMM/underlying_type'",
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
                out_exception_re=r"invalid specification: invalid enum underlying type; node_uri='foo\.yaml#/models/MMM/underlying_type' enum_underlying_type='vvv'",
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
                out_exception_re=r"invalid specification: invalid node kind: node_uri='foo\.yaml#/models/MMM/constants' node_kind=integer expected_node_kind=mapping",
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
                out_exception_re=r"invalid specification: non-empty mapping required: node_uri='foo\.yaml#/models/MMM/constants'",
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
                out_exception_re=r"invalid specification: invalid id; node_uri='foo\.yaml#/models/MMM/constants/ccc' id='ccc'",
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
                out_exception_re=r"invalid specification: missing node: node_uri='foo\.yaml#/models/MMM/constants/CCC/value'",
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
                out_exception_re=r"invalid specification: invalid node kind: node_uri='foo\.yaml#/models/MMM/constants/CCC/value' node_kind=floating-point expected_node_kind=integer",
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
                out_exception_re=r"invalid specification: invalid node kind: node_uri='foo\.yaml#/models/MMM/constants/CCC/value' node_kind=string expected_node_kind=integer",
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
                out_exception_re=r"invalid specification: invalid node kind: node_uri='foo\.yaml#/models/MMM/constants/CCC/value' node_kind=integer expected_node_kind=string",
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
                out_exception_re=r"invalid specification: invalid node kind: node_uri='foo\.yaml#/models/MMM/constants/CCC/description' node_kind=integer expected_node_kind=string",
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
                out_ignored_node_uris={"foo.yaml#/models/MMM/constants/CCC/xyz"},
            ),
        ]
        common.test(self, test_data_list)

    def test_xprimit(self):
        test_data_list = [
            common.TestData(
                in_file_path_2_file_data={
                    "foo.yaml": """
models:
  PPP:
    type: int32
    min: aaa
"""
                },
                out_exception_type=InvalidSpecError,
                out_exception_re=r"invalid specification: invalid node kind: node_uri='foo\.yaml#/models/PPP/min' node_kind=string expected_node_kind=integer",
            ),
            common.TestData(
                in_file_path_2_file_data={
                    "foo.yaml": """
models:
  PPP:
    type: int64
    max: aaa
"""
                },
                out_exception_type=InvalidSpecError,
                out_exception_re=r"invalid specification: invalid node kind: node_uri='foo\.yaml#/models/PPP/max' node_kind=string expected_node_kind=integer",
            ),
            common.TestData(
                in_file_path_2_file_data={
                    "foo.yaml": f"""
models:
  PPP:
    type: int32
    min: {-(1<<31)-1}
"""
                },
                out_exception_type=InvalidSpecError,
                out_exception_re=r"invalid specification: number too small: node_uri='foo\.yaml#/models/PPP/min' number=-2147483649 min_number=-2147483648",
            ),
            common.TestData(
                in_file_path_2_file_data={
                    "foo.yaml": f"""
models:
  PPP:
    type: int32
    min: {1<<31}
"""
                },
                out_exception_type=InvalidSpecError,
                out_exception_re=r"invalid specification: number too large: node_uri='foo\.yaml#/models/PPP/min' number=2147483648 max_number=2147483647",
            ),
            common.TestData(
                in_file_path_2_file_data={
                    "foo.yaml": f"""
models:
  PPP:
    type: int32
    max: {-(1<<31)-1}
"""
                },
                out_exception_type=InvalidSpecError,
                out_exception_re=r"invalid specification: number too small: node_uri='foo\.yaml#/models/PPP/max' number=-2147483649 min_number=-2147483648",
            ),
            common.TestData(
                in_file_path_2_file_data={
                    "foo.yaml": f"""
models:
  PPP:
    type: int32
    max: {1<<31}
"""
                },
                out_exception_type=InvalidSpecError,
                out_exception_re=r"invalid specification: number too large: node_uri='foo\.yaml#/models/PPP/max' number=2147483648 max_number=2147483647",
            ),
            common.TestData(
                in_file_path_2_file_data={
                    "foo.yaml": f"""
models:
  PPP:
    type: int64
    min: {-(1<<63)-1}
"""
                },
                out_exception_type=InvalidSpecError,
                out_exception_re=r"invalid specification: number too small: node_uri='foo\.yaml#/models/PPP/min' number=-9223372036854775809 min_number=-9223372036854775808",
            ),
            common.TestData(
                in_file_path_2_file_data={
                    "foo.yaml": f"""
models:
  PPP:
    type: int64
    min: {1<<63}
"""
                },
                out_exception_type=InvalidSpecError,
                out_exception_re=r"invalid specification: number too large: node_uri='foo\.yaml#/models/PPP/min' number=9223372036854775808 max_number=9223372036854775807",
            ),
            common.TestData(
                in_file_path_2_file_data={
                    "foo.yaml": f"""
models:
  PPP:
    type: int64
    max: {-(1<<63)-1}
"""
                },
                out_exception_type=InvalidSpecError,
                out_exception_re=r"invalid specification: number too small: node_uri='foo\.yaml#/models/PPP/max' number=-9223372036854775809 min_number=-9223372036854775808",
            ),
            common.TestData(
                in_file_path_2_file_data={
                    "foo.yaml": f"""
models:
  PPP:
    type: int64
    max: {1<<63}
"""
                },
                out_exception_type=InvalidSpecError,
                out_exception_re=r"invalid specification: number too large: node_uri='foo\.yaml#/models/PPP/max' number=9223372036854775808 max_number=9223372036854775807",
            ),
            common.TestData(
                in_file_path_2_file_data={
                    "foo.yaml": """
models:
  PPP:
    type: int64
    max: 1
    min: 2
"""
                },
                out_exception_type=InvalidSpecError,
                out_exception_re=r"invalid specification: invalid primitive constraints, min > max: node_uri='foo\.yaml#/models/PPP' min=2 max=1",
            ),
            common.TestData(
                in_file_path_2_file_data={
                    "foo.yaml": """
models:
  PPP:
    type: float32
    min: sss
"""
                },
                out_exception_type=InvalidSpecError,
                out_exception_re=r"invalid specification: invalid node kind: node_uri='foo\.yaml#/models/PPP/min' node_kind=string expected_node_kind=floating-point",
            ),
            common.TestData(
                in_file_path_2_file_data={
                    "foo.yaml": """
models:
  PPP:
    type: float64
    max: sss
"""
                },
                out_exception_type=InvalidSpecError,
                out_exception_re=r"invalid specification: invalid node kind: node_uri='foo\.yaml#/models/PPP/max' node_kind=string expected_node_kind=floating-point",
            ),
            common.TestData(
                in_file_path_2_file_data={
                    "foo.yaml": """
models:
  PPP:
    type: float32
    min: 1.0
    min_is_exclusive: 1
"""
                },
                out_exception_type=InvalidSpecError,
                out_exception_re=r"invalid specification: invalid node kind: node_uri='foo\.yaml#/models/PPP/min_is_exclusive' node_kind=integer expected_node_kind=boolean",
            ),
            common.TestData(
                in_file_path_2_file_data={
                    "foo.yaml": """
models:
  PPP:
    type: float32
    max: 1.0
    max_is_exclusive: 1
"""
                },
                out_exception_type=InvalidSpecError,
                out_exception_re=r"invalid specification: invalid node kind: node_uri='foo\.yaml#/models/PPP/max_is_exclusive' node_kind=integer expected_node_kind=boolean",
            ),
            common.TestData(
                in_file_path_2_file_data={
                    "foo.yaml": """
models:
  PPP:
    type: float64
    min: 1.0
    min_is_exclusive: true
    max: 1.0
"""
                },
                out_exception_type=InvalidSpecError,
                out_exception_re=r"invalid specification: invalid primitive constraints, min > max: node_uri='foo\.yaml#/models/PPP' exclusive_min=1.0 max=1.0",
            ),
            common.TestData(
                in_file_path_2_file_data={
                    "foo.yaml": """
models:
  PPP:
    type: float64
    min: 1.0
    max: 1.0
    max_is_exclusive: true
"""
                },
                out_exception_type=InvalidSpecError,
                out_exception_re=r"invalid specification: invalid primitive constraints, min > max: node_uri='foo\.yaml#/models/PPP' min=1.0 exclusive_max=1.0",
            ),
            common.TestData(
                in_file_path_2_file_data={
                    "foo.yaml": """
models:
  PPP:
    type: float64
    min: 2.0
    max: 1.0
"""
                },
                out_exception_type=InvalidSpecError,
                out_exception_re=r"invalid specification: invalid primitive constraints, min > max: node_uri='foo\.yaml#/models/PPP' min=2.0 max=1.0",
            ),
            common.TestData(
                in_file_path_2_file_data={
                    "foo.yaml": """
models:
  PPP:
    type: string
    min_length: aaa
"""
                },
                out_exception_type=InvalidSpecError,
                out_exception_re=r"invalid specification: invalid node kind: node_uri='foo\.yaml#/models/PPP/min_length' node_kind=string expected_node_kind=integer",
            ),
            common.TestData(
                in_file_path_2_file_data={
                    "foo.yaml": """
models:
  PPP:
    type: string
    min_length: -1
"""
                },
                out_exception_type=InvalidSpecError,
                out_exception_re=r"invalid specification: number too small: node_uri='foo\.yaml#/models/PPP/min_length' number=-1 min_number=0",
            ),
            common.TestData(
                in_file_path_2_file_data={
                    "foo.yaml": """
models:
  PPP:
    type: string
    max_length: aaa
"""
                },
                out_exception_type=InvalidSpecError,
                out_exception_re=r"invalid specification: invalid node kind: node_uri='foo\.yaml#/models/PPP/max_length' node_kind=string expected_node_kind=integer",
            ),
            common.TestData(
                in_file_path_2_file_data={
                    "foo.yaml": """
models:
  PPP:
    type: string
    max_length: 0
"""
                },
                out_exception_type=InvalidSpecError,
                out_exception_re=r"invalid specification: number too small: node_uri='foo\.yaml#/models/PPP/max_length' number=0 min_number=1",
            ),
            common.TestData(
                in_file_path_2_file_data={
                    "foo.yaml": """
models:
  PPP:
    type: string
    min_length: 2
    max_length: 1
"""
                },
                out_exception_type=InvalidSpecError,
                out_exception_re=r"invalid specification: invalid primitive constraints, min_length > max_length: node_uri='foo\.yaml#/models/PPP' min_length=2 max_length=1",
            ),
            common.TestData(
                in_file_path_2_file_data={
                    "foo.yaml": """
models:
  PPP:
    type: bool
"""
                },
                out_exception_type=InvalidSpecError,
                out_exception_re=r"invalid specification: invalid model type; node_uri='foo\.yaml#/models/PPP/type' model_type='bool' expected_pattern='struct\|enum\|int32\|int64\|float32\|float64\|string'",
            ),
            common.TestData(
                in_file_path_2_file_data={
                    "foo.yaml": """
models:
  PPP:
    type: int32
    example: aaa
"""
                },
                out_exception_type=InvalidSpecError,
                out_exception_re=r"invalid specification: invalid node kind: node_uri='foo\.yaml#/models/PPP/example' node_kind=string expected_node_kind=integer",
            ),
            common.TestData(
                in_file_path_2_file_data={
                    "foo.yaml": f"""
models:
  PPP:
    type: int32
    example: {-(1<<31)-1}
"""
                },
                out_exception_type=InvalidSpecError,
                out_exception_re=r"invalid specification: number too small: node_uri='foo\.yaml#/models/PPP/example' number=-2147483649 min_number=-2147483648",
            ),
            common.TestData(
                in_file_path_2_file_data={
                    "foo.yaml": """
models:
  PPP:
    type: int32
    min: 2
    example: 1
"""
                },
                out_exception_type=InvalidSpecError,
                out_exception_re=r"invalid specification: number too small: node_uri='foo\.yaml#/models/PPP/example' number=1 min_number=2",
            ),
            common.TestData(
                in_file_path_2_file_data={
                    "foo.yaml": f"""
models:
  PPP:
    type: int32
    example: {1<<31}
"""
                },
                out_exception_type=InvalidSpecError,
                out_exception_re=r"invalid specification: number too large: node_uri='foo\.yaml#/models/PPP/example' number=2147483648 max_number=2147483647",
            ),
            common.TestData(
                in_file_path_2_file_data={
                    "foo.yaml": """
models:
  PPP:
    type: int32
    max: 0
    example: 1
"""
                },
                out_exception_type=InvalidSpecError,
                out_exception_re=r"invalid specification: number too large: node_uri='foo\.yaml#/models/PPP/example' number=1 max_number=0",
            ),
            common.TestData(
                in_file_path_2_file_data={
                    "foo.yaml": """
models:
  PPP:
    type: int64
    example: aaa
"""
                },
                out_exception_type=InvalidSpecError,
                out_exception_re=r"invalid specification: invalid node kind: node_uri='foo\.yaml#/models/PPP/example' node_kind=string expected_node_kind=integer",
            ),
            common.TestData(
                in_file_path_2_file_data={
                    "foo.yaml": f"""
models:
  PPP:
    type: int64
    example: {-(1<<63)-1}
"""
                },
                out_exception_type=InvalidSpecError,
                out_exception_re=r"invalid specification: number too small: node_uri='foo\.yaml#/models/PPP/example' number=-9223372036854775809 min_number=-9223372036854775808",
            ),
            common.TestData(
                in_file_path_2_file_data={
                    "foo.yaml": """
models:
  PPP:
    type: int64
    min: 2
    example: 1
"""
                },
                out_exception_type=InvalidSpecError,
                out_exception_re=r"invalid specification: number too small: node_uri='foo\.yaml#/models/PPP/example' number=1 min_number=2",
            ),
            common.TestData(
                in_file_path_2_file_data={
                    "foo.yaml": f"""
models:
  PPP:
    type: int64
    example: {1<<63}
"""
                },
                out_exception_type=InvalidSpecError,
                out_exception_re=r"invalid specification: number too large: node_uri='foo\.yaml#/models/PPP/example' number=9223372036854775808 max_number=9223372036854775807",
            ),
            common.TestData(
                in_file_path_2_file_data={
                    "foo.yaml": """
models:
  PPP:
    type: int64
    max: 0
    example: 1
"""
                },
                out_exception_type=InvalidSpecError,
                out_exception_re=r"invalid specification: number too large: node_uri='foo\.yaml#/models/PPP/example' number=1 max_number=0",
            ),
            common.TestData(
                in_file_path_2_file_data={
                    "foo.yaml": """
models:
  PPP:
    type: float32
    example: aaa
"""
                },
                out_exception_type=InvalidSpecError,
                out_exception_re=r"invalid specification: invalid node kind: node_uri='foo\.yaml#/models/PPP/example' node_kind=string expected_node_kind=floating-point",
            ),
            common.TestData(
                in_file_path_2_file_data={
                    "foo.yaml": """
models:
  PPP:
    type: float32
    min: 2.0
    example: 1.0
"""
                },
                out_exception_type=InvalidSpecError,
                out_exception_re=r"invalid specification: number too small: node_uri='foo\.yaml#/models/PPP/example' number=1.0 min_number=2.0",
            ),
            common.TestData(
                in_file_path_2_file_data={
                    "foo.yaml": """
models:
  PPP:
    type: float32
    min: 1.0
    min_is_exclusive: true
    example: 1.0
"""
                },
                out_exception_type=InvalidSpecError,
                out_exception_re=r"invalid specification: number too small: node_uri='foo\.yaml#/models/PPP/example' number=1.0 exclusive_min_number=1.0",
            ),
            common.TestData(
                in_file_path_2_file_data={
                    "foo.yaml": """
models:
  PPP:
    type: float64
    example: aaa
"""
                },
                out_exception_type=InvalidSpecError,
                out_exception_re=r"invalid specification: invalid node kind: node_uri='foo\.yaml#/models/PPP/example' node_kind=string expected_node_kind=floating-point",
            ),
            common.TestData(
                in_file_path_2_file_data={
                    "foo.yaml": """
models:
  PPP:
    type: float64
    max: 0.0
    example: 1.0
"""
                },
                out_exception_type=InvalidSpecError,
                out_exception_re=r"invalid specification: number too large: node_uri='foo\.yaml#/models/PPP/example' number=1.0 max_number=0.0",
            ),
            common.TestData(
                in_file_path_2_file_data={
                    "foo.yaml": """
models:
  PPP:
    type: float64
    max: 1.0
    max_is_exclusive: true
    example: 1.0
"""
                },
                out_exception_type=InvalidSpecError,
                out_exception_re=r"invalid specification: number too large: node_uri='foo\.yaml#/models/PPP/example' number=1.0 exclusive_max_number=1.0",
            ),
            common.TestData(
                in_file_path_2_file_data={
                    "foo.yaml": """
models:
  PPP:
    type: string
    example: 11
"""
                },
                out_exception_type=InvalidSpecError,
                out_exception_re=r"invalid specification: invalid node kind: node_uri='foo\.yaml#/models/PPP/example' node_kind=integer expected_node_kind=string",
            ),
            common.TestData(
                in_file_path_2_file_data={
                    "foo.yaml": """
models:
  PPP:
    type: string
    min_length: 2
    example: a
"""
                },
                out_exception_type=InvalidSpecError,
                out_exception_re=r"invalid specification: string too short: node_uri='foo\.yaml#/models/PPP/example' string_length=1 min_string_length=2",
            ),
            common.TestData(
                in_file_path_2_file_data={
                    "foo.yaml": """
models:
  PPP:
    type: string
    max_length: 2
    example: abc
"""
                },
                out_exception_type=InvalidSpecError,
                out_exception_re=r"invalid specification: string too long: node_uri='foo\.yaml#/models/PPP/example' string_length=3 max_string_length=2",
            ),
            common.TestData(
                in_file_path_2_file_data={
                    "foo.yaml": """
models:
  PPP:
    type: string
    pattern: ''
"""
                },
                out_exception_type=InvalidSpecError,
                out_exception_re=r"invalid specification: string too short: node_uri='foo\.yaml#/models/PPP/pattern' string_length=0 min_string_length=1",
            ),
            common.TestData(
                in_file_path_2_file_data={
                    "foo.yaml": """
models:
  PPP:
    type: string
    pattern: (
"""
                },
                out_exception_type=InvalidSpecError,
                out_exception_re=r"invalid specification: invalid regular expression \(re2\): node_uri='foo\.yaml#/models/PPP' reg_exp='\('",
            ),
            common.TestData(
                in_file_path_2_file_data={
                    "foo.yaml": """
models:
  PPP:
    type: string
    pattern: '[0-9]+'
    example: aaa
"""
                },
                out_exception_type=InvalidSpecError,
                out_exception_re=r"invalid specification: unexpected string: node_uri='foo\.yaml#/models/PPP/example' string='aaa' expected_pattern='\[0-9\]\+'",
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
                out_exception_re=r"invalid specification: invalid id; node_uri='foo\.yaml#/errors/eee' id='eee'",
            ),
            common.TestData(
                in_file_path_2_file_data={
                    "foo.yaml": """
errors:
  EEE: 1
"""
                },
                out_exception_type=InvalidSpecError,
                out_exception_re=r"invalid specification: invalid node kind: node_uri='foo\.yaml#/errors/EEE' node_kind=integer expected_node_kind=mapping",
            ),
            common.TestData(
                in_file_path_2_file_data={
                    "foo.yaml": """
errors:
  EEE: {}
"""
                },
                out_exception_type=InvalidSpecError,
                out_exception_re=r"invalid specification: missing node: node_uri='foo\.yaml#/errors/EEE/code'",
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
                out_exception_re=r"invalid node kind: node_uri='foo\.yaml#/errors/EEE/code' node_kind=string expected_node_kind=integer",
            ),
            common.TestData(
                in_file_path_2_file_data={
                    "foo.yaml": """
errors:
  EEE:
    code: 1000
    status_code: 400
    description: 1
"""
                },
                out_exception_type=InvalidSpecError,
                out_exception_re=r"invalid specification: invalid node kind: node_uri='foo\.yaml#/errors/EEE/description' node_kind=integer expected_node_kind=string",
            ),
            common.TestData(
                in_file_path_2_file_data={
                    "foo.yaml": """
errors:
  EEE:
    code: 999
"""
                },
                out_exception_type=InvalidSpecError,
                out_exception_re=r"invalid specification: number too small: node_uri='foo\.yaml#/errors/EEE/code' number=999 min_number=1",
            ),
            common.TestData(
                in_file_path_2_file_data={
                    "foo.yaml": f"""
errors:
  EEE:
    code: {1<<31}
"""
                },
                out_exception_type=InvalidSpecError,
                out_exception_re=r"invalid specification: number too large: node_uri='foo\.yaml#/errors/EEE/code' number=2147483648 max_number=2147483647",
            ),
            common.TestData(
                in_file_path_2_file_data={
                    "foo.yaml": f"""
errors:
  EEE:
    code: 1000
    status_code: s
"""
                },
                out_exception_type=InvalidSpecError,
                out_exception_re=r"invalid specification: invalid node kind: node_uri='foo\.yaml#/errors/EEE/status_code' node_kind=string expected_node_kind=integer",
            ),
            common.TestData(
                in_file_path_2_file_data={
                    "foo.yaml": f"""
errors:
  EEE:
    code: 1000
    status_code: 99
"""
                },
                out_exception_type=InvalidSpecError,
                out_exception_re=r"invalid specification: number too small: node_uri='foo\.yaml#/errors/EEE/status_code' number=99 min_number=100",
            ),
            common.TestData(
                in_file_path_2_file_data={
                    "foo.yaml": f"""
errors:
  EEE:
    code: 1000
    status_code: 600
"""
                },
                out_exception_type=InvalidSpecError,
                out_exception_re=r"invalid specification: number too large: node_uri='foo\.yaml#/errors/EEE/status_code' number=600 max_number=599",
            ),
            common.TestData(
                in_file_path_2_file_data={
                    "foo.yaml": """
errors:
  EEE:
    code: 1000
    status_code: 400
    xyz: 1
"""
                },
                out_ignored_node_uris={"foo.yaml#/errors/EEE/xyz"},
            ),
        ]
        common.test(self, test_data_list)


if __name__ == "__main__":
    unittest.main()
