override SHELL := bash
override .SHELLFLAGS := -eu$(if $(value DEBUG),x)o pipefail -c

.SILENT:
.ONESHELL:

default:
.PHONY: default

testdata:
	PYTHONPATH=.. python3 -m src.jroh.compiler --go_out apicommon/testdata:github.com/go-tk/jroh/go/apicommon/testdata apicommon/testdata/*.yaml
.PHONY: testdata
default: testdata

fmt:
	go fmt ./...
.PHONY: fmt
default: fmt

vet:
	go vet ./...
.PHONY: vet
default: vet

test:
	go test -coverpkg=./apicommon/... ./apicommon
	go test -cover ./middleware/...
.PHONY: test
default: test
