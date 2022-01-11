import unittest

from . import common


class TestTranslator(unittest.TestCase):
    def test_service(self):
        self.maxDiff = None
        test_data_list = [
            common.TestData(
                in_file_path_2_file_data={
                    "default/a.yaml": """
services:
  Greeting:
    version: 1.1.1
  Greeting-V2:
    description: Test
    version: 1.1.1
""",
                    "default/a2.yaml": """
namespace: XYZ
services:
  Greeting:
    version: 10.1.1
    rpc_path_template: 'test/{namespace}-{service_id}-{method_id}'
""",
                    "default/b.yaml": """
methods:
  Say-Hello:
    service_id: Greeting
  Say-Hello-V2:
    service_id: Greeting-V2
""",
                    "default/b2.yaml": """
namespace: XYZ
methods:
  Say-Hello:
    service_id: Greeting
""",
                },
                out_file_path_2_file_data={
                    "default/greeting_service.yaml": """\
openapi: 3.0.0
info:
  title: Greeting Service
  version: 1.1.1
paths:
  /rpc/Default.Greeting.SayHello:
    post:
      operationId: sayHello
      requestBody:
        content:
          application/json:
            schema:
              type: object
      responses:
        '200':
          description: |-
            ## Error Cases

            | Error Code | Status Code | Message | Description |
            | - | - | - | - |
            | -1 | ... | ... | Low-level error. |
            | 1 | 501 | not implemented | The method is not implemented. |
          headers:
            Jroh-Trace-Id:
              description: The trace identifier.
              schema:
                type: string
            Jroh-Error-Code:
              description: The error code. This header is present only if error occurs.
              schema:
                type: int32
          content:
            application/json:
              schema:
                oneOf:
                - type: object
                - $ref: ../common.yaml#/components/schemas/error
""",
                    "default/greeting_v2_service.yaml": """\
openapi: 3.0.0
info:
  title: Greeting V2 Service
  version: 1.1.1
  description: Test
paths:
  /rpc/Default.GreetingV2.SayHelloV2:
    post:
      operationId: sayHelloV2
      requestBody:
        content:
          application/json:
            schema:
              type: object
      responses:
        '200':
          description: |-
            ## Error Cases

            | Error Code | Status Code | Message | Description |
            | - | - | - | - |
            | -1 | ... | ... | Low-level error. |
            | 1 | 501 | not implemented | The method is not implemented. |
          headers:
            Jroh-Trace-Id:
              description: The trace identifier.
              schema:
                type: string
            Jroh-Error-Code:
              description: The error code. This header is present only if error occurs.
              schema:
                type: int32
          content:
            application/json:
              schema:
                oneOf:
                - type: object
                - $ref: ../common.yaml#/components/schemas/error
""",
                    "xyz/greeting_service.yaml": """\
openapi: 3.0.0
info:
  title: Greeting Service
  version: 10.1.1
paths:
  /test/XYZ-Greeting-SayHello:
    post:
      operationId: sayHello
      requestBody:
        content:
          application/json:
            schema:
              type: object
      responses:
        '200':
          description: |-
            ## Error Cases

            | Error Code | Status Code | Message | Description |
            | - | - | - | - |
            | -1 | ... | ... | Low-level error. |
            | 1 | 501 | not implemented | The method is not implemented. |
          headers:
            Jroh-Trace-Id:
              description: The trace identifier.
              schema:
                type: string
            Jroh-Error-Code:
              description: The error code. This header is present only if error occurs.
              schema:
                type: int32
          content:
            application/json:
              schema:
                oneOf:
                - type: object
                - $ref: ../common.yaml#/components/schemas/error
""",
                },
            ),
        ]
        common.test(self, test_data_list)

    def test_method(self):
        self.maxDiff = None
        test_data_list = [
            common.TestData(
                in_file_path_2_file_data={
                    "default/a.yaml": """
services:
  Greeting:
    version: 1.1.1
    rpc_path_template: '{namespace}.{service_id}.{method_id}'
  Greeting-X:
    version: 1.1.1
    rpc_path_template: '{namespace}-x.{service_id}.{method_id}'
""",
                    "default/b.yaml": """
methods:
  Say-Hello:
    service_id: Greeting
  Say-Hello-V2:
    service_ids:
    - Greeting
    - Greeting-X
    summary: Haha
    params:
      Foo:
        type: int32
  Say-Hello-V3:
    service_id: Greeting
    description: Test
    results:
      Bar:
        type: string
""",
                },
                out_file_path_2_file_data={
                    "default/greeting_service.yaml": """\
openapi: 3.0.0
info:
  title: Greeting Service
  version: 1.1.1
paths:
  /Default.Greeting.SayHello:
    post:
      operationId: sayHello
      requestBody:
        content:
          application/json:
            schema:
              type: object
      responses:
        '200':
          description: |-
            ## Error Cases

            | Error Code | Status Code | Message | Description |
            | - | - | - | - |
            | -1 | ... | ... | Low-level error. |
            | 1 | 501 | not implemented | The method is not implemented. |
          headers:
            Jroh-Trace-Id:
              description: The trace identifier.
              schema:
                type: string
            Jroh-Error-Code:
              description: The error code. This header is present only if error occurs.
              schema:
                type: int32
          content:
            application/json:
              schema:
                oneOf:
                - type: object
                - $ref: ../common.yaml#/components/schemas/error
  /Default.Greeting.SayHelloV2:
    post:
      operationId: sayHelloV2
      summary: Haha
      requestBody:
        content:
          application/json:
            schema:
              $ref: models.yaml#/components/schemas/sayHelloV2Params
      responses:
        '200':
          description: |-
            ## Error Cases

            | Error Code | Status Code | Message | Description |
            | - | - | - | - |
            | -1 | ... | ... | Low-level error. |
            | 1 | 501 | not implemented | The method is not implemented. |
            | 2 | 422 | invalid params | Invalid method parameter(s). |
          headers:
            Jroh-Trace-Id:
              description: The trace identifier.
              schema:
                type: string
            Jroh-Error-Code:
              description: The error code. This header is present only if error occurs.
              schema:
                type: int32
          content:
            application/json:
              schema:
                oneOf:
                - type: object
                - $ref: ../common.yaml#/components/schemas/error
  /Default.Greeting.SayHelloV3:
    post:
      operationId: sayHelloV3
      description: Test
      requestBody:
        content:
          application/json:
            schema:
              type: object
      responses:
        '200':
          description: |-
            ## Error Cases

            | Error Code | Status Code | Message | Description |
            | - | - | - | - |
            | -1 | ... | ... | Low-level error. |
            | 1 | 501 | not implemented | The method is not implemented. |
          headers:
            Jroh-Trace-Id:
              description: The trace identifier.
              schema:
                type: string
            Jroh-Error-Code:
              description: The error code. This header is present only if error occurs.
              schema:
                type: int32
          content:
            application/json:
              schema:
                oneOf:
                - $ref: models.yaml#/components/schemas/sayHelloV3Results
                - $ref: ../common.yaml#/components/schemas/error
""",
                    "default/greeting_x_service.yaml": """\
openapi: 3.0.0
info:
  title: Greeting X Service
  version: 1.1.1
paths:
  /Default-x.GreetingX.SayHelloV2:
    post:
      operationId: sayHelloV2
      summary: Haha
      requestBody:
        content:
          application/json:
            schema:
              $ref: models.yaml#/components/schemas/sayHelloV2Params
      responses:
        '200':
          description: |-
            ## Error Cases

            | Error Code | Status Code | Message | Description |
            | - | - | - | - |
            | -1 | ... | ... | Low-level error. |
            | 1 | 501 | not implemented | The method is not implemented. |
            | 2 | 422 | invalid params | Invalid method parameter(s). |
          headers:
            Jroh-Trace-Id:
              description: The trace identifier.
              schema:
                type: string
            Jroh-Error-Code:
              description: The error code. This header is present only if error occurs.
              schema:
                type: int32
          content:
            application/json:
              schema:
                oneOf:
                - type: object
                - $ref: ../common.yaml#/components/schemas/error
""",
                    "default/models.yaml": """\
openapi: 3.0.0
info:
  title: Models
  version: ''
paths: {}
components:
  schemas:
    sayHelloV2Params:
      type: object
      properties:
        foo:
          type: integer
          format: int32
      required:
      - foo
    sayHelloV3Results:
      type: object
      properties:
        bar:
          type: string
      required:
      - bar
""",
                },
            ),
        ]
        common.test(self, test_data_list)

    def test_enum(self):
        self.maxDiff = None
        test_data_list = [
            common.TestData(
                in_file_path_2_file_data={
                    "default/a.yaml": """
services:
  Greeting:
    version: 1.2.1
methods:
  Say-Hello:
    service_id: Greeting
    params:
      Color:
        type: Color
        is_optional: true
      Color2:
        type: Color
        description: None
      Fruits1:
        type: Fruit
        is_repeated: true
      Fruits2:
        type: Fruit
        is_repeated: true
        min_count: 1
        description: Test
""",
                    "default/b.yaml": """
models:
  Color:
    type: enum
    underlying_type: int32
    constants:
      Red:
        value: 1
      Black:
        value: 22
        description: Black
  Fruit:
    type: enum
    description: An fruit
    underlying_type: string
    constants:
      Apple:
        value: Ap
        description: An apple
      Banana:
        value: Bn
        description: A banana
""",
                },
                out_file_path_2_file_data={
                    "default/greeting_service.yaml": """\
openapi: 3.0.0
info:
  title: Greeting Service
  version: 1.2.1
paths:
  /rpc/Default.Greeting.SayHello:
    post:
      operationId: sayHello
      requestBody:
        content:
          application/json:
            schema:
              $ref: models.yaml#/components/schemas/sayHelloParams
      responses:
        '200':
          description: |-
            ## Error Cases

            | Error Code | Status Code | Message | Description |
            | - | - | - | - |
            | -1 | ... | ... | Low-level error. |
            | 1 | 501 | not implemented | The method is not implemented. |
            | 2 | 422 | invalid params | Invalid method parameter(s). |
          headers:
            Jroh-Trace-Id:
              description: The trace identifier.
              schema:
                type: string
            Jroh-Error-Code:
              description: The error code. This header is present only if error occurs.
              schema:
                type: int32
          content:
            application/json:
              schema:
                oneOf:
                - type: object
                - $ref: ../common.yaml#/components/schemas/error
""",
                    "default/models.yaml": """\
openapi: 3.0.0
info:
  title: Models
  version: ''
paths: {}
components:
  schemas:
    sayHelloParams:
      type: object
      properties:
        color:
          allOf:
          - $ref: '#/components/schemas/color'
          - description: |-
              Constants:

              - RED(1)
              - BLACK(22): Black
        color2:
          allOf:
          - $ref: '#/components/schemas/color'
          - description: |-
              None

              Constants:

              - RED(1)
              - BLACK(22): Black
        fruits1:
          type: array
          items:
            allOf:
            - $ref: '#/components/schemas/fruit'
            - description: |-
                An fruit

                Constants:

                - APPLE("Ap"): An apple
                - BANANA("Bn"): A banana
        fruits2:
          type: array
          minItems: 1
          items:
            allOf:
            - $ref: '#/components/schemas/fruit'
            - description: |-
                Test

                Constants:

                - APPLE("Ap"): An apple
                - BANANA("Bn"): A banana
      required:
      - color2
      - fruits2
    color:
      type: integer
      format: int32
      enum:
      - 1
      - 22
    fruit:
      type: string
      enum:
      - Ap
      - Bn
""",
                },
            ),
        ]
        common.test(self, test_data_list)

    def test_xprimit(self):
        self.maxDiff = None
        test_data_list = [
            common.TestData(
                in_file_path_2_file_data={
                    "default/a.yaml": """
services:
  Greeting:
    version: 1.2.1
methods:
  Say-Hello:
    service_id: Greeting
    params:
      Nickname:
        type: Nickname
        description: The Niiickname
      Score:
        type: Score
      Weight:
        type: Weight
        is_optional: true
        description: +, it's optional
      Heights:
        type: Height
        is_repeated: true
        description: None
      Age:
        type: Age
""",
                    "default/b.yaml": """
models:
  Nickname:
    type: string
    min_length: 1
    max_length: 100
    pattern: '[a-z]+'
    description: The Nickname
    example: tommy
  Score:
    type: float32
    min: 0.0
    max: 100.0
    description: The Score
    example: 99.0
  Weight:
    type: float64
    min: 0.0
    min_is_exclusive: true
    description: The Weight
    example: 1.0
  Height:
    type: int32
    min: 1
    max: 300
    description: The Height
    example: 200
  Age:
    type: int64
    min: 1
    max: 200
    description: The Age
    example: 100
""",
                },
                out_file_path_2_file_data={
                    "default/greeting_service.yaml": """\
openapi: 3.0.0
info:
  title: Greeting Service
  version: 1.2.1
paths:
  /rpc/Default.Greeting.SayHello:
    post:
      operationId: sayHello
      requestBody:
        content:
          application/json:
            schema:
              $ref: models.yaml#/components/schemas/sayHelloParams
      responses:
        '200':
          description: |-
            ## Error Cases

            | Error Code | Status Code | Message | Description |
            | - | - | - | - |
            | -1 | ... | ... | Low-level error. |
            | 1 | 501 | not implemented | The method is not implemented. |
            | 2 | 422 | invalid params | Invalid method parameter(s). |
          headers:
            Jroh-Trace-Id:
              description: The trace identifier.
              schema:
                type: string
            Jroh-Error-Code:
              description: The error code. This header is present only if error occurs.
              schema:
                type: int32
          content:
            application/json:
              schema:
                oneOf:
                - type: object
                - $ref: ../common.yaml#/components/schemas/error
""",
                    "default/models.yaml": """\
openapi: 3.0.0
info:
  title: Models
  version: ''
paths: {}
components:
  schemas:
    sayHelloParams:
      type: object
      properties:
        nickname:
          allOf:
          - $ref: '#/components/schemas/nickname'
          - description: The Niiickname
        score:
          allOf:
          - $ref: '#/components/schemas/score'
          - description: The Score
        weight:
          allOf:
          - $ref: '#/components/schemas/weight'
          - description: The Weight, it's optional
        heights:
          type: array
          items:
            allOf:
            - $ref: '#/components/schemas/height'
            - description: None
        age:
          allOf:
          - $ref: '#/components/schemas/age'
          - description: The Age
      required:
      - nickname
      - score
      - age
    nickname:
      type: string
      minLength: 1
      maxLength: 100
      pattern: '[a-z]+'
      example: tommy
    score:
      type: number
      format: float
      minimum: 0.0
      maximum: 100.0
      example: 99.0
    weight:
      type: number
      format: double
      minimum: 0.0
      exclusiveMinimum: true
      example: 1.0
    height:
      type: integer
      format: int32
      minimum: 1
      maximum: 300
      example: 200
    age:
      type: integer
      format: int64
      minimum: 1
      maximum: 200
      example: 100
""",
                },
            ),
        ]
        common.test(self, test_data_list)

    def test_error(self):
        self.maxDiff = None
        test_data_list = [
            common.TestData(
                in_file_path_2_file_data={
                    "default/a.yaml": """
services:
  Greeting:
    version: 1.2.1
methods:
  Say-Hello:
    service_id: Greeting
    error_cases:
      Fail:
        description: Failed
      Bad-Situation: {}
      Abc.Xyz:
        description: + (>_<)
""",
                    "default/b.yaml": """
errors:
  Fail:
    code: 3000
    status_code: 400
  Bad-Situation:
    code: 4000
    status_code: 500
    description: None
""",
                    "default/c.yaml": """
namespace: Abc
errors:
  Xyz:
    code: 1230
    status_code: 400
    description: Too Bad!
""",
                },
                out_file_path_2_file_data={
                    "default/greeting_service.yaml": """\
openapi: 3.0.0
info:
  title: Greeting Service
  version: 1.2.1
paths:
  /rpc/Default.Greeting.SayHello:
    post:
      operationId: sayHello
      requestBody:
        content:
          application/json:
            schema:
              type: object
      responses:
        '200':
          description: |-
            ## Error Cases

            | Error Code | Status Code | Message | Description |
            | - | - | - | - |
            | -1 | ... | ... | Low-level error. |
            | 1 | 501 | not implemented | The method is not implemented. |
            | 1230 | 400 | xyz | Too Bad! (>_<) |
            | 3000 | 400 | fail | Failed |
            | 4000 | 500 | bad situation | None |
          headers:
            Jroh-Trace-Id:
              description: The trace identifier.
              schema:
                type: string
            Jroh-Error-Code:
              description: The error code. This header is present only if error occurs.
              schema:
                type: int32
          content:
            application/json:
              schema:
                oneOf:
                - type: object
                - $ref: ../common.yaml#/components/schemas/error
""",
                },
            ),
        ]
        common.test(self, test_data_list)

    def test_struct(self):
        self.maxDiff = None
        test_data_list = [
            common.TestData(
                in_file_path_2_file_data={
                    "default/a.yaml": """
namespace: NS1
services:
  Greeting:
    version: 1.2.1
methods:
  Say-Hello:
    service_id: Greeting
    params:
      F:
        type: Foo
        is_optional: true
      FF2:
        type: Default.S
        is_repeated: true
        min_count: 1
        description: + aaa
      FF44:
        type: Default.S
        is_repeated: true
        min_count: 2
        max_count: 2
        description: aaa
      FF55:
        type: Default.S
        is_repeated: true
        max_count: 2
        description: aaa
      Fun:
        type: No-Where.Fun
        description: Have fun
""",
                    "default/b.yaml": """
namespace: NS1
models:
  Foo:
    type: struct
    fields:
      X:
        type: int32
        min: 32
        is_repeated: true
        description: ABC
      X2:
        type: int32
        is_repeated: true
        max: -64
        description: ABC
      Y:
        type: Default.S
        is_optional: true
        description: CDE
""",
                    "default/bb.yaml": """
namespace: No-Where
models:
  Fun:
    type: struct
    fields:
      AA:
        type: AA
  AA:
    type: struct
""",
                    "default/c.yaml": """
models:
  S:
    type: struct
    fields:
      U:
        type: string
        min_length: 3
        max_length: 10
        is_repeated: true
        min_count: 1
        description: '123'
        example:
        - ufo
        - UUU
      I64:
        type: int64
        is_repeated: true
        min_count: 1
        example: [4, 5, 6]
      F32:
        type: float32
        min: 1.0
        min_is_exclusive: true
        is_optional: true
        example: 4.0
      F64:
        type: float64
        max: 4.0
        max_is_exclusive: true
        is_repeated: true
        example: [1.0, 2.0, 3.0]
      B:
        type: bool
    description: struct S
""",
                },
                out_file_path_2_file_data={
                    "ns1/greeting_service.yaml": """\
openapi: 3.0.0
info:
  title: Greeting Service
  version: 1.2.1
paths:
  /rpc/NS1.Greeting.SayHello:
    post:
      operationId: sayHello
      requestBody:
        content:
          application/json:
            schema:
              $ref: models.yaml#/components/schemas/sayHelloParams
      responses:
        '200':
          description: |-
            ## Error Cases

            | Error Code | Status Code | Message | Description |
            | - | - | - | - |
            | -1 | ... | ... | Low-level error. |
            | 1 | 501 | not implemented | The method is not implemented. |
            | 2 | 422 | invalid params | Invalid method parameter(s). |
          headers:
            Jroh-Trace-Id:
              description: The trace identifier.
              schema:
                type: string
            Jroh-Error-Code:
              description: The error code. This header is present only if error occurs.
              schema:
                type: int32
          content:
            application/json:
              schema:
                oneOf:
                - type: object
                - $ref: ../common.yaml#/components/schemas/error
""",
                    "ns1/models.yaml": """\
openapi: 3.0.0
info:
  title: Models
  version: ''
paths: {}
components:
  schemas:
    sayHelloParams:
      type: object
      properties:
        f:
          $ref: '#/components/schemas/foo'
        ff2:
          type: array
          minItems: 1
          items:
            allOf:
            - $ref: ../default/models.yaml#/components/schemas/s
            - description: struct S aaa
        ff44:
          type: array
          minItems: 2
          maxItems: 2
          items:
            allOf:
            - $ref: ../default/models.yaml#/components/schemas/s
            - description: aaa
        ff55:
          type: array
          maxItems: 2
          items:
            allOf:
            - $ref: ../default/models.yaml#/components/schemas/s
            - description: aaa
        fun:
          allOf:
          - $ref: ../no_where/models.yaml#/components/schemas/fun
          - description: Have fun
      required:
      - ff2
      - ff44
      - fun
    foo:
      type: object
      properties:
        x:
          type: array
          items:
            type: integer
            format: int32
            minimum: 32
            description: ABC
        x2:
          type: array
          items:
            type: integer
            format: int32
            maximum: -64
            description: ABC
        y:
          allOf:
          - $ref: ../default/models.yaml#/components/schemas/s
          - description: CDE
""",
                    "no_where/models.yaml": """\
openapi: 3.0.0
info:
  title: Models
  version: ''
paths: {}
components:
  schemas:
    fun:
      type: object
      properties:
        aa:
          $ref: '#/components/schemas/aa'
      required:
      - aa
    aa:
      type: object
""",
                    "default/models.yaml": """\
openapi: 3.0.0
info:
  title: Models
  version: ''
paths: {}
components:
  schemas:
    s:
      type: object
      properties:
        u:
          type: array
          minItems: 1
          items:
            type: string
            minLength: 3
            maxLength: 10
            description: '123'
          example:
          - ufo
          - UUU
        i64:
          type: array
          minItems: 1
          items:
            type: integer
            format: int64
          example:
          - 4
          - 5
          - 6
        f32:
          type: number
          format: float
          minimum: 1.0
          exclusiveMinimum: true
          example: 4.0
        f64:
          type: array
          items:
            type: number
            format: double
            maximum: 4.0
            exclusiveMaximum: true
          example:
          - 1.0
          - 2.0
          - 3.0
        b:
          type: boolean
      required:
      - u
      - i64
      - b
""",
                },
            ),
        ]
        common.test(self, test_data_list)


if __name__ == "__main__":
    unittest.main()
