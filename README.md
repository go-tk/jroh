# JROH

Solution & Framework for **J**SON-**R**PC **o**ver **H**TTP

## Why not OpenAPI?

OpenAPI addresses the definition of RESTful APIs, when it comes to JSON-RPCs, some important elements
for defining RPCs is missing, e.g., service name, method name and RPC error code.

## The user story

1. Users shall define JSON-RPCs in **YAML**.
2. Users can compile the YAML files into stub (client-side) code and skeleton (server-side) code in **Go**.
3. Users can compile the YAML files into OpenAPI 3.0 specifications so that it's able to leverage
**Swagger UI** as a browser for JSON-RPCs.

## Features

TODO

## Installation

No installation required, we run the toolchain in Docker containers.

## Getting started

**1. Create a temporary directory as workspace**

```sh
$ mkdir temp && cd temp
```

---

**2. Define JSON-RPC(s)**

```sh
$ mkdir -p ./jroh/hello_world

$ cat >./jroh/hello_world/greeter_service.yaml <<EOF
########## BEGIN greeter_service.yaml ##########
namespace: Hello-World

services:
  Greeter:
    description: The service of greetings.
    version: 1.0.0
    rpc_path_template: /rpc/{namespace}.{service_id}.{method_id}

methods:
  Say-Hello:
    service_id: Greeter
    summary: Sends a greeting
    description: Sends a greeting on behalf of someone.
    params:
      Name:
        type: string
        min_length: 1
        description: The user's name.
        example: Roy
    results:
      Message:
        type: string
        min_length: 1
        description: The greeting returned.
        example: Hi, Roy!
    error_cases:
      User-Not-Allowed: {}

errors:
  User-Not-Allowed:
    code: 1001
    description: The user is in the blacklist.
########## END greeter_service.yaml ##########
EOF
```

---

**3. Generate Go code and OpenAPI 3.0 specification(s)**

```sh
$ go mod init pkg.go.test

$ docker run --rm \
  --volume="${PWD}:/workspace" \
  --workdir=/workspace \
  ghcr.io/go-tk/jrohc:v0.4.2 \
    --go_out=./api:pkg.go.test/api \
    --oapi3_out=./oapi3 \
    ./jroh/hello_world/greeter_service.yaml

$ ls --recursive -l ./api ./oapi3
```

---

**4. Start up an RPC server**

```sh
$ mkdir ./server

$ cat >./server/server.go <<EOF
////////// BEGIN server.go //////////
package main

import (
    "context"
    "fmt"
    "log"
    "net/http"

    "github.com/go-tk/jroh/go/apicommon"
    "pkg.go.test/api/helloworldapi"
)

func main() {
    server := helloworldapi.GreeterServerFuncs{
        SayHelloFunc: func(
            ctx context.Context,
            params *helloworldapi.SayHelloParams,
            results *helloworldapi.SayHelloResults,
        ) error {
            log.Printf("Received: %v", params.Name)

            if params.Name == "God" {
                return helloworldapi.ErrUserNotAllowed
            }
            results.Message = fmt.Sprintf("Hi, %v!", params.Name)
            return nil
        },
    }
    rpcRouter := apicommon.NewRPCRouter(nil)
    helloworldapi.RegisterGreeterServer(&server, rpcRouter, apicommon.ServerOptions{})
    log.Printf("rpc route infos: %#v", rpcRouter.RPCRouteInfos())

    apicommon.DebugMode = true
    err := http.ListenAndServe(":2220", rpcRouter.ServeMux())
    log.Fatal(err)
}
////////// END server.go //////////
EOF

$ go run ./server/server.go
```

---

**5.a. Send RPC requests with `curl`**

```sh
$ curl --data '{"name": "Roy"}' http://127.0.0.1:2220/rpc/HelloWorld.Greeter.SayHello
# Output:
# {
#   "traceID": "Uv38ByGCZU8WP18PmmIdcg",
#   "results": {
#     "message": "Hi, Roy!"
#   }
# }

$ curl --data '{"name": "God"}' http://127.0.0.1:2220/rpc/HelloWorld.Greeter.SayHello
# Output:
# {
#   "traceID": "lWbHTRADfE17uwQH0eLGSQ",
#   "error": {
#     "code": 1001,
#     "message": "user not allowed"
#   }
# }
```

---

**5.b. Send RPC requests with Go client**

```sh
$ mkdir ./client

$ cat >./client/client.go <<EOF
////////// BEGIN client.go //////////
package main

import (
    "context"
    "errors"
    "fmt"
    "log"

    "github.com/go-tk/jroh/go/apicommon"
    "pkg.go.test/api/helloworldapi"
)

func main() {
    client := helloworldapi.NewGreeterClient("http://127.0.0.1:2220", apicommon.ClientOptions{})
    results1, err := client.SayHello(context.Background(), &helloworldapi.SayHelloParams{
        Name: "Roy",
    })
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("1 - %#v\n", *results1)

    _, err = client.SayHello(context.Background(), &helloworldapi.SayHelloParams{
        Name: "God",
    })
    if !errors.Is(err, helloworldapi.ErrUserNotAllowed) {
        log.Fatal(err)
    }
    fmt.Printf("2 - %v\n", err)
}
////////// END client.go //////////
EOF

$ go run ./client/client.go
# Output:
# 1 - helloworldapi.SayHelloResults{Message:"Hi, Roy!"}
# 2 - apicommon: rpc failed; namespace="HelloWorld" serviceName="Greeter" methodName="SayHello" traceID="lWbHTRADfE17uwQH0eLGSQ": api: user not allowed (1001)
```

---

**6. Browse  JSON-RPC(s) with Swagger UI**

```sh
$ docker run --rm \
  --volume="${PWD}:/usr/share/nginx/html/data" \
  --env=SWAGGER_JSON_URL=./data/oapi3/hello_world/greeter_service.yaml \
  --publish=2333:8080 \
  swaggerapi/swagger-ui:latest
```

Open http://127.0.0.1:2333 in the browser.

## Advanced examples

- [Petstore](examples/2-petstore)
