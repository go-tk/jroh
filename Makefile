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
	coverage run --omit=src/tests/*,.venv/* -m unittest discover --start-directory=src/tests --top-level-directory=.
	coverage report