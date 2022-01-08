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

## Integration

- [Zerolog](./go/middleware/zerologmw)
- [OpenTelemetry](./go/middleware/opentelemetrymw)
- [Prometheus](./go/middleware/prometheusmw)

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
# Create file ./jroh/hello_world/greeter_service.yaml
$ install -Dm 644 /dev/stdin ./jroh/hello_world/greeter_service.yaml <<EOF
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

---

**3. Generate Go code and OpenAPI 3.0 specification(s)**

```sh
$ go mod init pkg.go.test

$ docker run --rm \
  --volume="${PWD}:/workspace" \
  --workdir=/workspace \
  ghcr.io/go-tk/jrohc:v0.9.0 \
    --go_out=./api:pkg.go.test/api \
    --oapi3_out=./oapi3 \
    ./jroh/hello_world/greeter_service.yaml

$ ls -R ./api ./oapi3
# Output:
#
# ./api:
# helloworldapi
#
# ./api/helloworldapi:
# errors.go  greeteractor.go  greeterclient.go  misc.go  models.go
#
# ./oapi3:
# common.yaml  hello_world
#
# ./oapi3/hello_world:
# greeter_service.yaml  models.yaml
```

---

**4. Start up an RPC server**

```sh
# Create file ./server/server.go
$ install -Dm 644 /dev/stdin ./server/server.go <<EOF
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

$ go run -v ./server/server.go
```

---

**5.a. Send RPC requests with `curl`**

```sh
$ curl -XPOST -d'{"name": "Roy"}' -D- http://127.0.0.1:2220/rpc/HelloWorld.Greeter.SayHello
# Output:
#
# HTTP/1.1 200 OK
# Content-Type: application/json; charset=utf-8
# X-Jroh-Trace-Id: Uv38ByGCZU8WP18PmmIdcg
# Date: Sun, 02 Jan 2022 07:34:15 GMT
# Content-Length: 28
#
# {
#   "greeting": "Hi, Roy!"
# }

$ curl -XPOST -d'{"name": "God"}' -D- http://127.0.0.1:2220/rpc/HelloWorld.Greeter.SayHello
# Output:
#
# HTTP/1.1 403 Forbidden
# Content-Type: application/json; charset=utf-8
# X-Jroh-Error-Code: 1000
# X-Jroh-Trace-Id: lWbHTRADfE17uwQH0eLGSQ
# Date: Sun, 02 Jan 2022 07:35:01 GMT
# Content-Length: 36
#
# {
#   "message": "user not allowed"
# }
```

---

**5.b. Send RPC requests with Go client**

```sh
# Create file ./client/client.go
$ install -Dm 644 /dev/stdin ./client/client.go <<EOF
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

$ go run -v ./client/client.go
# Output:
#
# 1 - &helloworldapi.SayHelloResults{Greeting:"Hi, Roy!"}
# 2 - &apicommon.Error{Code:1000, StatusCode:403, Message:"user not allowed", Details:"", Data:apicommon.ErrorData(nil)}
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

![screenshot](https://user-images.githubusercontent.com/6377788/148325351-d57e6dd1-0646-4b66-ae82-370eedcdd16f.png)

## Advanced examples

- [Petstore](examples/2-petstore)
