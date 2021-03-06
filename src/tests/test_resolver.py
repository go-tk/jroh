import unittest

from ..jroh.resolver import InvalidSpecError
from . import common


class TestResolver(unittest.TestCase):
    def test_duplicate_id(self):
        test_data_list = [
            common.TestData(
                in_file_path_2_file_data={
                    "foo.yaml": """
services:
  Foo:
    version: 1.1.1
""",
                    "foo2.yaml": """
services:
  Foo:
    version: 1.1.1
""",
                },
                out_exception_type=InvalidSpecError,
                out_exception_re=r"invalid spec: duplicate service id; node_uri1='foo2\.yaml#/services/Foo' node_uri2='foo\.yaml#/services/Foo'",
            ),
            common.TestData(
                in_file_path_2_file_data={
                    "foo.yaml": """
methods:
  Hello:
    service_id: World
""",
                    "foo2.yaml": """
methods:
  Hello:
    service_id: World
""",
                },
                out_exception_type=InvalidSpecError,
                out_exception_re=r"invalid spec: duplicate method id; node_uri1='foo2\.yaml#/methods/Hello' node_uri2='foo\.yaml#/methods/Hello'",
            ),
            common.TestData(
                in_file_path_2_file_data={
                    "foo.yaml": """
models:
  Test:
    type: struct
""",
                    "foo2.yaml": """
models:
  Test:
    type: enum
    underlying_type: int32
""",
                },
                out_exception_type=InvalidSpecError,
                out_exception_re=r"invalid spec: duplicate model id; node_uri1='foo2\.yaml#/models/Test' node_uri2='foo\.yaml#/models/Test'",
            ),
            common.TestData(
                in_file_path_2_file_data={
                    "foo.yaml": """
models:
  EEE:
    type: enum
    underlying_type: int32
    constants:
      Apple:
        value: 100
""",
                    "foo2.yaml": """
models:
  TTT:
    type: enum
    underlying_type: string
    constants:
      Apple:
        value: "apple"
""",
                },
                out_exception_type=InvalidSpecError,
                out_exception_re=r"invalid spec: duplicate constant id; node_uri1='foo2\.yaml#/models/TTT/constants/Apple' node_uri2='foo\.yaml#/models/EEE/constants/Apple'",
            ),
            common.TestData(
                in_file_path_2_file_data={
                    "foo.yaml": """
errors:
  Err:
    code: 1000
    status_code: 400
""",
                    "foo2.yaml": """
errors:
  Err:
    code: 2000
    status_code: 400
""",
                },
                out_exception_type=InvalidSpecError,
                out_exception_re=r"invalid spec: duplicate error id; node_uri1='foo2\.yaml#/errors/Err' node_uri2='foo\.yaml#/errors/Err'",
            ),
            common.TestData(
                in_file_path_2_file_data={
                    "foo.yaml": """
services:
  World:
    version: 1.1.1
methods:
  Hello:
    service_id: World
models:
  Test:
    type: struct
errors:
  Err:
    code: 1000
    status_code: 400
""",
                    "foo2.yaml": """
namespace: New
services:
  World:
    version: 1.1.1
methods:
  Hello:
    service_id: World
models:
  Test:
    type: struct
errors:
  Err:
    code: 1000
    status_code: 400
""",
                },
            ),
        ]
        common.test(self, test_data_list)

    def test_duplicate_error_code(self):
        test_data_list = [
            common.TestData(
                in_file_path_2_file_data={
                    "foo.yaml": """
errors:
  Foo:
    code: 1000
    status_code: 400
""",
                    "foo2.yaml": """
errors:
  Bar:
    code: 1000
    status_code: 400
""",
                },
                out_exception_type=InvalidSpecError,
                out_exception_re=r"invalid spec: duplicate error code; node_uri1='foo2\.yaml#/errors/Bar/code' node_uri2='foo\.yaml#/errors/Foo/code' error_code=1",
            ),
            common.TestData(
                in_file_path_2_file_data={
                    "foo.yaml": """
namespace: New
errors:
  Foo:
    code: 1000
    status_code: 400
""",
                    "foo2.yaml": """
errors:
  Bar:
    code: 1000
    status_code: 400
""",
                },
            ),
        ]
        common.test(self, test_data_list)

    def test_error_code_conflict(self):
        test_data_list = [
            common.TestData(
                in_file_path_2_file_data={
                    "foo.yaml": """
services:
  Hello:
    version: 1.1.1
methods:
  World:
    service_id: Hello
    error_cases:
      Foo: {}
      New.Bar: {}
errors:
  Foo:
    code: 1000
    status_code: 400
""",
                    "foo2.yaml": """
namespace: New
errors:
  Bar:
    code: 1000
    status_code: 400
""",
                },
                out_exception_type=InvalidSpecError,
                out_exception_re=r"invalid spec: error code conflict; node_uri1='foo\.yaml#/methods/World/error_cases/New\.Bar' node_uri2='foo\.yaml#/methods/World/error_cases/Foo' error_code=1000",
            ),
        ]
        common.test(self, test_data_list)

    def test_service_ref(self):
        test_data_list = [
            common.TestData(
                in_file_path_2_file_data={
                    "foo.yaml": """
methods:
  World:
    service_id: Hello1
"""
                },
                out_exception_type=InvalidSpecError,
                out_exception_re=r"invalid spec: service not found; node_uri='foo\.yaml#/methods/World/service_ids\[0\]' service_id='Hello1'",
            ),
            common.TestData(
                in_file_path_2_file_data={
                    "foo.yaml": """
services:
  Hello:
    version: 1.1.1
methods:
  World:
    service_id: Hello
"""
                },
            ),
            common.TestData(
                in_file_path_2_file_data={
                    "foo.yaml": """
services:
  Hello:
    version: 1.1.1
methods:
  World:
    service_ids:
    - Hello
    - Hello1
"""
                },
                out_exception_type=InvalidSpecError,
                out_exception_re=r"invalid spec: service not found; node_uri='foo\.yaml#/methods/World/service_ids\[1\]' service_id='Hello1'",
            ),
            common.TestData(
                in_file_path_2_file_data={
                    "foo.yaml": """
namespace: New
services:
  Hello:
    version: 1.1.1
""",
                    "bar.yaml": """
methods:
  World:
    service_id: Hello
""",
                },
                out_exception_type=InvalidSpecError,
                out_exception_re=r"invalid spec: service not found; node_uri='bar\.yaml#/methods/World/service_ids\[0\]' service_id='Hello'",
            ),
            common.TestData(
                in_file_path_2_file_data={
                    "foo.yaml": """
services:
  Hello:
    version: 1.1.1
  World:
    version: 1.1.2
""",
                    "bar.yaml": """
methods:
  World:
    service_ids:
    - Hello
    - World
""",
                },
            ),
        ]
        common.test(self, test_data_list)

    def test_error_ref(self):
        test_data_list = [
            common.TestData(
                in_file_path_2_file_data={
                    "foo.yaml": """
services:
  Hello:
    version: 1.1.1
methods:
  World:
    service_id: Hello
    error_cases:
      New.Something-Wrong: {}
"""
                },
                out_exception_type=InvalidSpecError,
                out_exception_re=r"invalid spec: error not found; node_uri='foo\.yaml#/methods/World/error_cases/New.Something-Wrong' namespace='New' error_id='Something-Wrong'",
            ),
            common.TestData(
                in_file_path_2_file_data={
                    "foo.yaml": """
services:
  Hello:
    version: 1.1.1
methods:
  World:
    service_id: Hello
    error_cases:
      Something-Wrong: {}
errors:
  Something-Wrong:
    code: 1000
    status_code: 403
"""
                },
            ),
            common.TestData(
                in_file_path_2_file_data={
                    "foo.yaml": """
namespace: New
services:
  Hello:
    version: 1.1.1
methods:
  World:
    service_id: Hello
    error_cases:
      Something-Wrong: {}
""",
                    "bar.yaml": """
errors:
  Something-Wrong:
    code: 1000
    status_code: 403
""",
                },
                out_exception_type=InvalidSpecError,
                out_exception_re=r"invalid spec: error not found; node_uri='foo\.yaml#/methods/World/error_cases/Something-Wrong' namespace='New' error_id='Something-Wrong'",
            ),
            common.TestData(
                in_file_path_2_file_data={
                    "foo.yaml": """
services:
  Hello:
    version: 1.1.1
methods:
  World:
    service_id: Hello
    error_cases:
      Something-Wrong: {}
""",
                    "bar.yaml": """
namespace: New
errors:
  Something-Wrong:
    code: 1000
    status_code: 403
""",
                },
                out_exception_type=InvalidSpecError,
                out_exception_re=r"invalid spec: error not found; node_uri='foo\.yaml#/methods/World/error_cases/Something-Wrong' namespace='Default' error_id='Something-Wrong'",
            ),
            common.TestData(
                in_file_path_2_file_data={
                    "foo.yaml": """
services:
  Hello:
    version: 1.1.1
methods:
  World:
    service_id: Hello
    error_cases:
      Something-Wrong: {}
""",
                    "bar.yaml": """
errors:
  Something-Wrong:
    code: 1000
    status_code: 403
""",
                },
            ),
            common.TestData(
                in_file_path_2_file_data={
                    "foo.yaml": """
namespace: New
services:
  Hello:
    version: 1.1.1
methods:
  World:
    service_id: Hello
    error_cases:
      Default.Something-Wrong: {}
""",
                    "bar.yaml": """
errors:
  Something-Wrong:
    code: 1000
    status_code: 403
""",
                },
            ),
            common.TestData(
                in_file_path_2_file_data={
                    "foo.yaml": """
services:
  Hello:
    version: 1.1.1
methods:
  World:
    service_id: Hello
    error_cases:
      New.Something-Wrong: {}
""",
                    "bar.yaml": """
namespace: New
errors:
  Something-Wrong:
    code: 1000
    status_code: 403
""",
                },
            ),
        ]
        common.test(self, test_data_list)

    def test_field_ref(self):
        test_data_list = [
            common.TestData(
                in_file_path_2_file_data={
                    "foo.yaml": """
services:
  Hello:
    version: 1.1.1
methods:
  World:
    service_id: Hello
    params:
      Bar:
        type: Bar
"""
                },
                out_exception_type=InvalidSpecError,
                out_exception_re=r"invalid spec: model not found; node_uri='foo\.yaml#/methods/World/params/Bar/type' namespace='Default' model_id='Bar'",
            ),
            common.TestData(
                in_file_path_2_file_data={
                    "foo.yaml": """
services:
  Hello:
    version: 1.1.1
methods:
  World:
    service_id: Hello
    params:
      Bar:
        type: Bar
models:
  Bar:
    type: struct
"""
                },
            ),
            common.TestData(
                in_file_path_2_file_data={
                    "foo.yaml": """
namespace: New
services:
  Hello:
    version: 1.1.1
methods:
  World:
    service_id: Hello
    results:
      Bar:
        type: Bar
      I:
        type: int32
""",
                    "bar.yaml": """
models:
  Bar:
    type: struct
""",
                },
                out_exception_type=InvalidSpecError,
                out_exception_re=r"invalid spec: model not found; node_uri='foo\.yaml#/methods/World/results/Bar/type' namespace='New' model_id='Bar'",
            ),
            common.TestData(
                in_file_path_2_file_data={
                    "foo.yaml": """
services:
  Hello:
    version: 1.1.1
methods:
  World:
    service_id: Hello
    params:
      Bar:
        type: Bar
""",
                    "bar.yaml": """
models:
  Bar:
    type: struct
""",
                },
            ),
            common.TestData(
                in_file_path_2_file_data={
                    "foo.yaml": """
namespace: New
services:
  Hello:
    version: 1.1.1
methods:
  World:
    service_id: Hello
    params:
      Bar:
        type: Default.Bar
""",
                    "bar.yaml": """
models:
  Bar:
    type: struct
""",
                },
            ),
            common.TestData(
                in_file_path_2_file_data={
                    "foo.yaml": """
services:
  Hello:
    version: 1.1.1
methods:
  World:
    service_id: Hello
    params:
      Bar:
        type: New.Bar
""",
                    "bar.yaml": """
namespace: New
models:
  Bar:
    type: struct
    fields:
      Baz:
        type: Baz
  Baz:
    type: struct
""",
                },
            ),
        ]
        common.test(self, test_data_list)

    def test_unused_node(self):
        test_data_list = [
            common.TestData(
                in_file_path_2_file_data={
                    "foo1.yaml": """
services:
  Hello:
    version: 1.1.1
""",
                    "foo2.yaml": """
models:
  Foo:
    type: enum
    underlying_type: int32
""",
                    "foo3.yaml": """
errors:
  Wrong:
    code: 1000
    status_code: 403
""",
                },
                out_unused_node_uris={
                    "foo1.yaml#/services/Hello",
                    "foo2.yaml#/models/Foo",
                    "foo3.yaml#/errors/Wrong",
                },
            ),
        ]
        common.test(self, test_data_list)


if __name__ == "__main__":
    unittest.main()
