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
    rpc_path_template: 'test/{namespace}.{service_id}.{method_id}'
""",
                    "default/b.yaml": """
methods:
  Say-Hello:
    service_ids: [Greeting]
  Say-Hello-V2:
    service_ids: [Greeting-V2]
""",
                    "default/b2.yaml": """
namespace: XYZ
methods:
  Say-Hello:
    service_ids: [Greeting]
""",
                },
                out_file_path_2_file_data={
                    "default/greeting_service.yaml": """\
openapi: 3.0.0
info:
  title: Greeting Service
  version: 1.1.1
paths:
  /rpc/Greeting.SayHello:
    post:
      operationId: sayHello
      responses:
        '200':
          description: |-
            ## Error Cases

            | Code | Message | Description |
            | - | - | - |
            | -32603 | internal error | Internal JSON-RPC error. |
          content:
            application/json:
              schema:
                $ref: models.yaml#/components/schemas/sayHelloResp
""",
                    "default/greeting_v2_service.yaml": """\
openapi: 3.0.0
info:
  title: Greeting V2 Service
  version: 1.1.1
  description: Test
paths:
  /rpc/GreetingV2.SayHelloV2:
    post:
      operationId: sayHelloV2
      responses:
        '200':
          description: |-
            ## Error Cases

            | Code | Message | Description |
            | - | - | - |
            | -32603 | internal error | Internal JSON-RPC error. |
          content:
            application/json:
              schema:
                $ref: models.yaml#/components/schemas/sayHelloV2Resp
""",
                    "default/models.yaml": """\
openapi: 3.0.0
info:
  title: Models
  version: ''
paths: {}
components:
  schemas:
    sayHelloResp:
      type: object
      properties:
        traceID:
          type: string
          description: The identifier of the trace associated with the log entry.
          example: 9m4e2mr0ui3e8a215n4g
        error:
          allOf:
          - $ref: common.yaml#/components/schemas/error
          - description: The error encountered.
      required:
      - traceID
    sayHelloV2Resp:
      type: object
      properties:
        traceID:
          type: string
          description: The identifier of the trace associated with the log entry.
          example: 9m4e2mr0ui3e8a215n4g
        error:
          allOf:
          - $ref: common.yaml#/components/schemas/error
          - description: The error encountered.
      required:
      - traceID
""",
                    "xyz/greeting_service.yaml": """\
openapi: 3.0.0
info:
  title: Greeting Service
  version: 10.1.1
paths:
  /test/XYZ.Greeting.SayHello:
    post:
      operationId: sayHello
      responses:
        '200':
          description: |-
            ## Error Cases

            | Code | Message | Description |
            | - | - | - |
            | -32603 | internal error | Internal JSON-RPC error. |
          content:
            application/json:
              schema:
                $ref: models.yaml#/components/schemas/sayHelloResp
""",
                    "xyz/models.yaml": """\
openapi: 3.0.0
info:
  title: Models
  version: ''
paths: {}
components:
  schemas:
    sayHelloResp:
      type: object
      properties:
        traceID:
          type: string
          description: The identifier of the trace associated with the log entry.
          example: 9m4e2mr0ui3e8a215n4g
        error:
          allOf:
          - $ref: ../common.yaml#/components/schemas/error
          - description: The error encountered.
      required:
      - traceID
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
    service_ids: [Greeting]
  Say-Hello-V2:
    service_ids: [Greeting, Greeting-X]
    summary: Haha
    params:
      Foo:
        type: int32
  Say-Hello-V3:
    service_ids: [Greeting]
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
      responses:
        '200':
          description: |-
            ## Error Cases

            | Code | Message | Description |
            | - | - | - |
            | -32603 | internal error | Internal JSON-RPC error. |
          content:
            application/json:
              schema:
                $ref: models.yaml#/components/schemas/sayHelloResp
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

            | Code | Message | Description |
            | - | - | - |
            | -32700 | parse error | Invalid JSON was received by the server. |
            | -32603 | internal error | Internal JSON-RPC error. |
            | -32602 | invalid params | Invalid method parameter(s). |
          content:
            application/json:
              schema:
                $ref: models.yaml#/components/schemas/sayHelloV2Resp
  /Default.Greeting.SayHelloV3:
    post:
      operationId: sayHelloV3
      description: Test
      responses:
        '200':
          description: |-
            ## Error Cases

            | Code | Message | Description |
            | - | - | - |
            | -32603 | internal error | Internal JSON-RPC error. |
          content:
            application/json:
              schema:
                $ref: models.yaml#/components/schemas/sayHelloV3Resp
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

            | Code | Message | Description |
            | - | - | - |
            | -32700 | parse error | Invalid JSON was received by the server. |
            | -32603 | internal error | Internal JSON-RPC error. |
            | -32602 | invalid params | Invalid method parameter(s). |
          content:
            application/json:
              schema:
                $ref: models.yaml#/components/schemas/sayHelloV2Resp
""",
                    "default/models.yaml": """\
openapi: 3.0.0
info:
  title: Models
  version: ''
paths: {}
components:
  schemas:
    sayHelloResp:
      type: object
      properties:
        traceID:
          type: string
          description: The identifier of the trace associated with the log entry.
          example: 9m4e2mr0ui3e8a215n4g
        error:
          allOf:
          - $ref: common.yaml#/components/schemas/error
          - description: The error encountered.
      required:
      - traceID
    sayHelloV2Params:
      type: object
      properties:
        foo:
          type: integer
          format: int32
      required:
      - foo
    sayHelloV2Resp:
      type: object
      properties:
        traceID:
          type: string
          description: The identifier of the trace associated with the log entry.
          example: 9m4e2mr0ui3e8a215n4g
        error:
          allOf:
          - $ref: common.yaml#/components/schemas/error
          - description: The error encountered.
      required:
      - traceID
    sayHelloV3Resp:
      type: object
      properties:
        traceID:
          type: string
          description: The identifier of the trace associated with the log entry.
          example: 9m4e2mr0ui3e8a215n4g
        error:
          allOf:
          - $ref: common.yaml#/components/schemas/error
          - description: The error encountered. This field is mutually exclusive of
              the `results` field.
        results:
          allOf:
          - $ref: '#/components/schemas/sayHelloV3Results'
          - description: The results returned. This field is mutually exclusive of
              the `error` field.
      required:
      - traceID
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
    service_ids: [Greeting]
    params:
      Color:
        type: Color?
      Color2:
        type: Color
        description: None
      Fruits1:
        type: Fruit*
      Fruits2:
        type: Fruit+
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
  /rpc/Greeting.SayHello:
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

            | Code | Message | Description |
            | - | - | - |
            | -32700 | parse error | Invalid JSON was received by the server. |
            | -32603 | internal error | Internal JSON-RPC error. |
            | -32602 | invalid params | Invalid method parameter(s). |
          content:
            application/json:
              schema:
                $ref: models.yaml#/components/schemas/sayHelloResp
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
            $ref: '#/components/schemas/fruit'
          description: |-
            Test

            Constants:

            - APPLE("Ap"): An apple
            - BANANA("Bn"): A banana
      required:
      - color2
      - fruits2
    sayHelloResp:
      type: object
      properties:
        traceID:
          type: string
          description: The identifier of the trace associated with the log entry.
          example: 9m4e2mr0ui3e8a215n4g
        error:
          allOf:
          - $ref: common.yaml#/components/schemas/error
          - description: The error encountered.
      required:
      - traceID
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
    service_ids: [Greeting]
    error_cases:
      Fail:
        description: Failed
      Bad-Situation: {}
      Abc.Xyz:
        description: (>_<)
""",
                    "default/b.yaml": """
errors:
  Fail:
    code: 300
  Bad-Situation:
    code: 400
""",
                    "default/c.yaml": """
namespace: Abc
errors:
  Xyz:
    code: 123
""",
                },
                out_file_path_2_file_data={
                    "default/greeting_service.yaml": """\
openapi: 3.0.0
info:
  title: Greeting Service
  version: 1.2.1
paths:
  /rpc/Greeting.SayHello:
    post:
      operationId: sayHello
      responses:
        '200':
          description: |-
            ## Error Cases

            | Code | Message | Description |
            | - | - | - |
            | -32603 | internal error | Internal JSON-RPC error. |
            | 123 | xyz | (>_<) |
            | 300 | fail | Failed |
            | 400 | bad situation |  |
          content:
            application/json:
              schema:
                $ref: models.yaml#/components/schemas/sayHelloResp
""",
                    "default/models.yaml": """\
openapi: 3.0.0
info:
  title: Models
  version: ''
paths: {}
components:
  schemas:
    sayHelloResp:
      type: object
      properties:
        traceID:
          type: string
          description: The identifier of the trace associated with the log entry.
          example: 9m4e2mr0ui3e8a215n4g
        error:
          allOf:
          - $ref: common.yaml#/components/schemas/error
          - description: The error encountered.
      required:
      - traceID
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
    service_ids: [Greeting]
    params:
      F:
        type: Foo?
      FF2:
        type: Default.S+
        description: aaa
      FF22:
        type: Default.S{1,}
        description: aaa
      FF33:
        type: Default.S{2,}
        description: aaa
      FF44:
        type: Default.S{2,2}
        description: aaa
      FF55:
        type: Default.S{0,2}
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
        type: int32*
        min: 32
        description: ABC
      X2:
        type: int32{0,}
        max: -64
        description: ABC
      Y:
        type: Default.S?
        description: CDE
      Y2:
        type: Default.S{0,1}
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
      AA2:
        type: AA{1,1}
  AA:
    type: struct
""",
                    "default/c.yaml": """
models:
  S:
    type: struct
    fields:
      U:
        type: string+
        min_length: 3
        max_length: 10
        description: '123'
        example:
        - ufo
        - UUU
      I64:
        type: int64+
        example: [4, 5, 6]
      F32:
        type: float32?
        example: 4.0
      F64:
        type: float64*
        example: [1.0, 2.0, 3.0]
      B:
        type: bool
""",
                },
                out_file_path_2_file_data={
                    "ns1/greeting_service.yaml": """\
openapi: 3.0.0
info:
  title: Greeting Service
  version: 1.2.1
paths:
  /rpc/Greeting.SayHello:
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

            | Code | Message | Description |
            | - | - | - |
            | -32700 | parse error | Invalid JSON was received by the server. |
            | -32603 | internal error | Internal JSON-RPC error. |
            | -32602 | invalid params | Invalid method parameter(s). |
          content:
            application/json:
              schema:
                $ref: models.yaml#/components/schemas/sayHelloResp
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
            $ref: ../default/models.yaml#/components/schemas/s
          description: aaa
        ff22:
          type: array
          minItems: 1
          items:
            $ref: ../default/models.yaml#/components/schemas/s
          description: aaa
        ff33:
          type: array
          minItems: 2
          items:
            $ref: ../default/models.yaml#/components/schemas/s
          description: aaa
        ff44:
          type: array
          minItems: 2
          maxItems: 2
          items:
            $ref: ../default/models.yaml#/components/schemas/s
          description: aaa
        ff55:
          type: array
          maxItems: 2
          items:
            $ref: ../default/models.yaml#/components/schemas/s
          description: aaa
        fun:
          allOf:
          - $ref: ../no_where/models.yaml#/components/schemas/fun
          - description: Have fun
      required:
      - ff2
      - ff22
      - ff33
      - ff44
      - fun
    sayHelloResp:
      type: object
      properties:
        traceID:
          type: string
          description: The identifier of the trace associated with the log entry.
          example: 9m4e2mr0ui3e8a215n4g
        error:
          allOf:
          - $ref: ../common.yaml#/components/schemas/error
          - description: The error encountered.
      required:
      - traceID
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
        y2:
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
        aa2:
          $ref: '#/components/schemas/aa'
      required:
      - aa
      - aa2
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
          example: 4.0
        f64:
          type: array
          items:
            type: number
            format: double
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
