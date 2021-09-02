override SHELL := bash
override .SHELLFLAGS := -eu$(if $(value DEBUG),x)o pipefail -c

.SILENT:
.ONESHELL:

.PHONY: all
all: format test

.PHONY: format
format:
	autoflake --in-place --recursive --remove-all-unused-imports --remove-duplicate-keys --remove-unused-variables src
	black src
	isort src --profile=black

.PHONY: test
test:
	coverage run --include='src/jroh/*' -m unittest discover --start-directory=src/tests --top-level-directory=.
	coverage report


.PHONY: example
example:
	find ../example -type f -name '*.yaml' | xargs python3 -m src.jroh.compiler --out ../example_out --go_out ../example_go_out:example/api
