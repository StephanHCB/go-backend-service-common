# go-backend-service-common

A library that bundles common building blocks for golang microservices.

## Features

This library provides:

- read and validate **configuration** from environment variables (and from a file on localhost)
- json **logging** (and human readable plaintext on localhost)
- a **vault** client
- a **health** controller
- a controller for serving a bundled **swagger ui** and an openapi v3 spec
- **middlewares** for
  - cors headers
  - distributed tracing (request id headers)
  - incoming request logging
  - incoming request metrics
  - incoming request timeouts
  - panic recovery
  - apm tracing

It aims to be compatible with a typical Spring microservice:

- Request IDs and the Request ID header match
- Request metrics names are customized to match
  [Spring Boot Default Metrics](https://tomgregory.com/spring-boot-default-metrics/)
- JSON logging format has been adapted to match ECS logging schema

## It's not Beans it's ... Acorns

Our ready-made singleton components implement the 
[Acorn interface](https://github.com/StephanHCB/go-autumn-acorn-registry/blob/main/api/interface.go#L102)
from [StephanHCB/go-autumn-acorn-registry](https://github.com/StephanHCB/go-autumn-acorn-registry).

You **can** wire up your application manually instead, but if you decide to play along, and build the rest of your
application from Acorns (=singleton classes that implement 
[Acorn](https://github.com/StephanHCB/go-autumn-acorn-registry/blob/main/api/interface.go#L102)),
it will give you **dependency injection** and automatic **setup in the correct order**, but without the usual 
drawbacks:
- zero additional dependencies
- no reflection
- no hidden auto magic
- you are in full control of the entire component lifecycle

## Development

### initial setup

clone this outside of your GOPATH (on linux, defaults to ~/go)

_Tip: On Windows, maybe don't have the GOPATH in your profile, because a large body of source code goes there._

### build

`go build ./...`

### run tests

Run all tests, collecting coverage across all directories

`go test -coverpkg=./... -v ./...`

In IntelliJ/GoLand, if you want to check code coverage, you must set Go Tool Arguments to `-coverpkg=./...`, 
so cross-package coverage from acceptance tests is considered. You may wish to set this under 
Edit Configuration - Edit Configuration Templates, so it will be set on all new test run configurations.

### Goland terminal configuration

Goland has the annoying habit of limiting line width on the output terminal to 80 characters no matter how wide the window is.
You can fix this. Menu: Help -> Find Action... -> search for "Registry"

Uncheck `go.run.processes.with.pty`.

### List dependency tree

`go mod graph > deps.txt`
