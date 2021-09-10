override SHELL := bash
override .SHELLFLAGS := -eu$(if $(value DEBUG),x)o pipefail -c

.SILENT:
.ONESHELL:

.PHONY: all
all: format test go_vet go_test examples

.PHONY: format
format:
	autoflake --in-place --recursive --remove-all-unused-imports --remove-duplicate-keys --remove-unused-variables src
	black src
	isort src --profile=black

.PHONY: test
test:
	coverage run --include='src/jroh/*' -m unittest discover --start-directory=src/tests --top-level-directory=.
	coverage report

.PHONY: go_vet
go_vet:
	python3 -m src.jroh.compiler --go_out go/apicommon/testdata:github.com/go-tk/jroh/go/apicommon/testdata go/apicommon/testdata/*.yaml
	cd go
	go vet ./...

.PHONY: go_test
go_test:
	cd go
	go test -coverpkg=./apicommon/... ./apicommon


.PHONY: examples
examples:
	find examples -type d -path examples/output -prune -o -type f -name '*.yaml' -print |
		xargs $(if $(value DEBUG),--verbose) \
			python3 -m src.jroh.compiler \
			--oapi3_out examples/output/oapi3 \
			--go_out examples/output/go:github.com/go-tk/jroh/examples/output/go
