import unittest

from . import common


class TestTranslator(unittest.TestCase):
    def test_service(self):
        self.maxDiff = None
        test_data_list = [
            common.TestData(
                in_file_path_2_file_data={
                    "a.yaml": """
services:
  Greeting:
    version: 1.1.1
  Greeting-V2:
    description: Test
    version: 1.1.1
""",
                    "a2.yaml": """
namespace: XYZ
services:
  Greeting:
    version: 10.1.1
    method_path_template: 'test/{namespace}.{service_id}.{method_id}'
""",
                    "b.yaml": """
methods:
  Say-Hello:
    service_ids: [Greeting]
  Say-Hello-V2:
    service_ids: [Greeting-V2]
""",
                    "b2.yaml": """
namespace: XYZ
methods:
  Say-Hello:
    service_ids: [Greeting]
""",
                },
                out_file_path_2_file_data={
                    "greeting_api.yaml": """\
openapi: 3.0.0
info:
  title: Greeting API
  version: 1.1.1
paths:
  /rpc/Greeting.SayHello:
    post:
      operationId: sayHello
      responses:
        '200':
          description: ''
          content:
            application/json:
              schema:
                $ref: builtins.yaml#/components/schemas/rpcRespWithoutResult
""",
                    "greeting_v2_api.yaml": """\
openapi: 3.0.0
info:
  title: Greeting V2 API
  version: 1.1.1
paths:
  /rpc/GreetingV2.SayHelloV2:
    post:
      operationId: sayHelloV2
      responses:
        '200':
          description: ''
          content:
            application/json:
              schema:
                $ref: builtins.yaml#/components/schemas/rpcRespWithoutResult
""",
                    "xyz/greeting_api.yaml": """\
openapi: 3.0.0
info:
  title: Greeting API
  version: 10.1.1
paths:
  /test/XYZ.Greeting.SayHello:
    post:
      operationId: sayHello
      responses:
        '200':
          description: ''
          content:
            application/json:
              schema:
                $ref: ../builtins.yaml#/components/schemas/rpcRespWithoutResult
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
                    "a.yaml": """
services:
  Greeting:
    version: 1.1.1
    method_path_template: '{namespace}.{service_id}.{method_id}'
  Greeting-X:
    version: 1.1.1
    method_path_template: '{namespace}-x.{service_id}.{method_id}'
""",
                    "b.yaml": """
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
    result:
      Bar:
        type: string
""",
                },
                out_file_path_2_file_data={
                    "greeting_api.yaml": """\
openapi: 3.0.0
info:
  title: Greeting API
  version: 1.1.1
paths:
  /global.Greeting.SayHello:
    post:
      operationId: sayHello
      responses:
        '200':
          description: ''
          content:
            application/json:
              schema:
                $ref: builtins.yaml#/components/schemas/rpcRespWithoutResult
  /global.Greeting.SayHelloV2:
    post:
      operationId: sayHelloV2
      summary: Haha
      requestBody:
        content:
          application/json:
            schema:
              $ref: schemas.yaml#/components/schemas/sayHelloV2Params
      responses:
        '200':
          description: ''
          content:
            application/json:
              schema:
                $ref: builtins.yaml#/components/schemas/rpcRespWithoutResult
  /global.Greeting.SayHelloV3:
    post:
      operationId: sayHelloV3
      description: Test
      responses:
        '200':
          description: ''
          content:
            application/json:
              schema:
                $ref: schemas.yaml#/components/schemas/sayHelloV3Resp
""",
                    "greeting_x_api.yaml": """\
openapi: 3.0.0
info:
  title: Greeting X API
  version: 1.1.1
paths:
  /global-x.GreetingX.SayHelloV2:
    post:
      operationId: sayHelloV2
      summary: Haha
      requestBody:
        content:
          application/json:
            schema:
              $ref: schemas.yaml#/components/schemas/sayHelloV2Params
      responses:
        '200':
          description: ''
          content:
            application/json:
              schema:
                $ref: builtins.yaml#/components/schemas/rpcRespWithoutResult
""",
                    "schemas.yaml": """\
openapi: 3.0.0
info:
  title: Schemas
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
    sayHelloV3Resp:
      type: object
      properties:
        id:
          type: integer
          format: int64
          description: The RPC identifier generated by the server.
        error:
          allOf:
          - $ref: builtins.yaml#/components/schemas/rpcError
          - description: The RPC error encountered. This field is mutually exclusive
              of the `result` field.
        result:
          allOf:
          - $ref: '#/components/schemas/sayHelloV3Result'
          - description: The RPC result returned. This field is mutually exclusive
              of the `error` field.
      required:
      - id
    sayHelloV3Result:
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
                    "a.yaml": """
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
                    "b.yaml": """
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
                    "greeting_api.yaml": """\
openapi: 3.0.0
info:
  title: Greeting API
  version: 1.2.1
paths:
  /rpc/Greeting.SayHello:
    post:
      operationId: sayHello
      requestBody:
        content:
          application/json:
            schema:
              $ref: schemas.yaml#/components/schemas/sayHelloParams
      responses:
        '200':
          description: ''
          content:
            application/json:
              schema:
                $ref: builtins.yaml#/components/schemas/rpcRespWithoutResult
""",
                    "schemas.yaml": """\
openapi: 3.0.0
info:
  title: Schemas
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
                    "a.yaml": """
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
                    "b.yaml": """
errors:
  Fail:
    code: 300
  Bad-Situation:
    code: 400
""",
                    "c.yaml": """
namespace: Abc
errors:
  Xyz:
    code: 123
""",
                },
                out_file_path_2_file_data={
                    "greeting_api.yaml": """\
openapi: 3.0.0
info:
  title: Greeting API
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
            | 300 | fail | Failed |
            | 400 | bad situation |  |
            | 123 | xyz | (>_<) |
          content:
            application/json:
              schema:
                $ref: builtins.yaml#/components/schemas/rpcRespWithoutResult
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
                    "a.yaml": """
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
        type: global.S+
        description: aaa
      Fun:
        type: No-Where.Fun
        description: Have fun
""",
                    "b.yaml": """
namespace: NS1
models:
  Foo:
    type: struct
    fields:
      X:
        type: int32*
        description: ABC
      Y:
        type: global.S?
        description: CDE
""",
                    "bb.yaml": """
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
                    "c.yaml": """
models:
  S:
    type: struct
    fields:
      U:
        type: string+
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
                    "ns1/greeting_api.yaml": """\
openapi: 3.0.0
info:
  title: Greeting API
  version: 1.2.1
paths:
  /rpc/Greeting.SayHello:
    post:
      operationId: sayHello
      requestBody:
        content:
          application/json:
            schema:
              $ref: schemas.yaml#/components/schemas/sayHelloParams
      responses:
        '200':
          description: ''
          content:
            application/json:
              schema:
                $ref: ../builtins.yaml#/components/schemas/rpcRespWithoutResult
""",
                    "ns1/schemas.yaml": """\
openapi: 3.0.0
info:
  title: Schemas
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
          items:
            $ref: ../schemas.yaml#/components/schemas/s
          description: aaa
        fun:
          allOf:
          - $ref: ../no_where/schemas.yaml#/components/schemas/fun
          - description: Have fun
      required:
      - ff2
      - fun
    foo:
      type: object
      properties:
        x:
          type: array
          items:
            type: integer
            format: int32
          description: ABC
        y:
          allOf:
          - $ref: ../schemas.yaml#/components/schemas/s
          - description: CDE
""",
                    "no_where/schemas.yaml": """\
openapi: 3.0.0
info:
  title: Schemas
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
                    "schemas.yaml": """\
openapi: 3.0.0
info:
  title: Schemas
  version: ''
paths: {}
components:
  schemas:
    s:
      type: object
      properties:
        u:
          type: array
          items:
            type: string
          description: '123'
          example:
          - ufo
          - UUU
        i64:
          type: array
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
