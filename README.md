# JROH

Framework for **J**SON **R**PC **o**ver **H**TTP

## Why not OpenAPI?

OpenAPI addresses the definition of RESTful APIs, when it comes to JSON RPCs, some important elements
for defining RPCs is missing, e.g., service name, method name and RPC error code.

## The user story

1. Users shall define JSON RPCs in **YAML**.
2. Users can compile the YAML files into stub (client-side) code and skeleton (server-side) code in **Go**.
3. Users can compile the YAML files into OpenAPI 3.0 specifications so that it's able to leverage
**Swagger UI** as viewer for JSON RPCs.

## Integration

- [Zerolog](./go/middleware/zerologmw)
- [OpenTelemetry](./go/middleware/opentelemetrymw)
- [Prometheus](./go/middleware/prometheusmw)

## Installation

No installation required, we run the toolchain in Docker containers.

## Getting started

### 1. Prepare

Create a temporary directory as workspace:

```sh
mkdir temp && cd temp
```

### 2. Define JSON RPC(s) with YAML

Create new file `./jroh/hello_world/greeter_service.yaml` by the following command:

```sh
mkdir -p ./jroh/hello_world && cat >./jroh/hello_world/greeter_service.yaml <<EOF
######### BEGIN greeter_service.yaml ##########
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
      Greeting:
        type: string
        min_length: 1
        description: The greeting returned.
        example: Hi, Roy!
    error_cases:
      User-Not-Allowed: {}

errors:
  User-Not-Allowed:
    code: 1000
    status_code: 403
    description: The user is in the blacklist.
########## END greeter_service.yaml ##########
EOF
```

### 3. Generate Go code and OpenAPI 3 specifications

Transform the current directory into Go module `pkg.go.test`:

```sh
go mod init pkg.go.test
```

Generate files:

```sh
docker run --rm \
  --volume="${PWD}:/workspace" \
  --workdir=/workspace \
  ghcr.io/go-tk/jrohc:v0.11.0 \
    --go_out=./api:pkg.go.test/api \
    --oapi3_out=./oapi3 \
    ./jroh/hello_world/greeter_service.yaml
```

Check generated files:

```sh
ls -R ./api ./oapi3
```

```sh
# Output:
#
# ./api:
# helloworldapi
#
# ./api/helloworldapi:
# errors_generated.go  greeteractor_generated.go  greeterclient_generated.go  misc_generated.go  models_generated.go
#
# ./oapi3:
# common.yaml  hello_world
#
# ./oapi3/hello_world:
# greeter_service.yaml  models.yaml
```

### 4. Start up an RPC server

Create new file `./server/server.go` by the following command:

```sh
mkdir -p ./server && cat >./server/server.go <<EOF
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
        actor := helloworldapi.GreeterActorFuncs{
                SayHelloFunc: func(
                        ctx context.Context,
                        params *helloworldapi.SayHelloParams,
                        results *helloworldapi.SayHelloResults,
                ) error {
                        log.Printf("Received: %v", params.Name)

                        if params.Name == "God" {
                                return helloworldapi.NewUserNotAllowedError()
                        }
                        results.Greeting = fmt.Sprintf("Hi, %v!", params.Name)
                        return nil
                },
        }
        router := apicommon.NewRouter()
        helloworldapi.RegisterGreeterActor(&actor, router, apicommon.ActorOptions{})
        log.Printf("route infos: %#v", router.RouteInfos())

        apicommon.DebugMode = true
        err := http.ListenAndServe(":2220", router)
        log.Fatal(err)
}
////////// END server.go //////////
EOF
```

Run the program:

```sh
go mod tidy && go run ./server/server.go
```

### 5.a. Invoke RPCs with `curl`

Call `HelloWorld.Greeter.SayHello` with good parameters:

```sh
curl -XPOST -d'{"name": "Roy"}' -D- http://127.0.0.1:2220/rpc/HelloWorld.Greeter.SayHello
```

```sh
# Output:
#
# HTTP/1.1 200 OK
# Content-Type: application/json; charset=utf-8
# Jroh-Trace-Id: Uv38ByGCZU8WP18PmmIdcg
# Date: Sun, 02 Jan 2022 07:34:15 GMT
# Content-Length: 28
#
# {
#   "greeting": "Hi, Roy!"
# }
```

Call `HelloWorld.Greeter.SayHello` with bad parameters:

```sh
curl -XPOST -d'{"name": "God"}' -D- http://127.0.0.1:2220/rpc/HelloWorld.Greeter.SayHello
```

```sh
# Output:
#
# HTTP/1.1 403 Forbidden
# Content-Type: application/json; charset=utf-8
# Jroh-Error-Code: 1000
# Jroh-Trace-Id: lWbHTRADfE17uwQH0eLGSQ
# Date: Sun, 02 Jan 2022 07:35:01 GMT
# Content-Length: 36
#
# {
#   "message": "user not allowed"
# }
```

### 5.b. Invoke RPCs with Go client

Create new file `./client/client.go` by the following command:

```sh
mkdir -p ./client && cat >./client/client.go <<EOF
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
        results, err := client.SayHello(context.Background(), &helloworldapi.SayHelloParams{
                Name: "Roy",
        })
        if err != nil {
                log.Fatal(err)
        }
        fmt.Printf("1 - %#v\n", *results)

        _, err = client.SayHello(context.Background(), &helloworldapi.SayHelloParams{
                Name: "God",
        })
        if error := (*apicommon.Error)(nil); errors.As(err, &error) && error.Code == helloworldapi.ErrorUserNotAllowed {
                fmt.Printf("2 - %#v\n", error)
        }
}

////////// END client.go //////////
EOF
```

Run the program:

```sh
go mod tidy && go run ./client/client.go
```

```sh
# Output:
#
# 1 - &helloworldapi.SayHelloResults{Greeting:"Hi, Roy!"}
# 2 - &apicommon.Error{Code:1000, StatusCode:403, Message:"user not allowed", Details:"", Data:apicommon.ErrorData(nil)}
```

### 6. View the definition of JSON RPC(s)

Run an instance of [Swagger UI](https://swagger.io/tools/swagger-ui/) as viewer:

```sh
docker run --rm \
  --volume="${PWD}:/usr/share/nginx/html/data" \
  --env=SWAGGER_JSON_URL=./data/oapi3/hello_world/greeter_service.yaml \
  --publish=2333:8080 \
  swaggerapi/swagger-ui:latest
```

Open http://127.0.0.1:2333 in the browser:

![screenshot](https://user-images.githubusercontent.com/6377788/148944012-143461b2-c399-46eb-b649-6183d485cd3b.png)

## Advanced examples

- [Petstore](examples/2-petstore)
