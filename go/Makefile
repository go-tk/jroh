override SHELL := bash
override .SHELLFLAGS := -eu$(if $(value DEBUG),x)o pipefail -c

.SILENT:
.ONESHELL:

.PHONY: all
all: testdata fmt vet test

.PHONY: testdata
testdata:
	PYTHONPATH=.. python3 -m src.jroh.compiler --go_out apicommon/testdata:github.com/go-tk/jroh/go/apicommon/testdata apicommon/testdata/*.yaml

.PHONY: fmt
fmt:
	go fmt ./...

.PHONY: vet
vet:
	go vet ./...

.PHONY: test
test:
	go test -coverpkg=./apicommon/... ./apicommon
	go test -cover ./middleware/...
