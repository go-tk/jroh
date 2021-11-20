override SHELL := bash
override .SHELLFLAGS := -eu$(if $(value DEBUG),x)o pipefail -c

.SILENT:
.ONESHELL:

default:
.PHONY: default

format:
	autoflake --in-place --recursive --remove-all-unused-imports --remove-duplicate-keys --remove-unused-variables src
	black src
	isort src --profile=black
.PHONY: format
default: format

test:
	coverage run --include='src/jroh/*' -m unittest discover --start-directory=src/tests --top-level-directory=.
	coverage report
.PHONY: test
default: test

go:
	$(MAKE) --directory=go
.PHONY: go
default: go

examples:
	find examples -type d -path examples/output -prune -o -type f -name '*.yaml' -print |
		xargs $(if $(value DEBUG),--verbose) \
			python3 -m src.jroh.compiler \
			--oapi3_out examples/output/oapi3 \
			--go_out examples/output/go:github.com/go-tk/jroh/examples/output/go
.PHONY: examples
default: examples

envrc:
	python3 -m venv .venv
	.venv/bin/pip3 install .
	.venv/bin/pip3 install --requirement requirements-dev.txt
	cat >.envrc <<EOF
	source .venv/bin/activate
	EOF
.PHONY: envrc
