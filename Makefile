override SHELL := bash
override .SHELLFLAGS := -eu$(if $(value DEBUG),x)o pipefail -c

.SILENT:
.ONESHELL:

.PHONY: all
all: format test go examples

.PHONY: envrc
envrc:
	python3 -m venv .venv
	.venv/bin/pip3 install .
	.venv/bin/pip3 install --requirement requirements-dev.txt
	cat >.envrc <<EOF
	source .venv/bin/activate
	EOF

.PHONY: format
format:
	autoflake --in-place --recursive --remove-all-unused-imports --remove-duplicate-keys --remove-unused-variables src
	black src
	isort src --profile=black

.PHONY: test
test:
	coverage run --include='src/jroh/*' -m unittest discover --start-directory=src/tests --top-level-directory=.
	coverage report

.PHONY: go
go:
	$(MAKE) --directory=go


.PHONY: examples
examples:
	find examples -type d -path examples/output -prune -o -type f -name '*.yaml' -print |
		xargs $(if $(value DEBUG),--verbose) \
			python3 -m src.jroh.compiler \
			--oapi3_out examples/output/oapi3 \
			--go_out examples/output/go:$$(cd examples/output/go; go list -m)
