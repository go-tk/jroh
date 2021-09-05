#!/usr/bin/env bash

set -euo pipefail

DIR=$(dirname "$(realpath "${0}")")
cd "${DIR}/../../compiler"
python3 -m src.jroh.compiler --go_out "${DIR}:github.com/go-tk/jroh/go/testdata" "${DIR}/"*.yaml
