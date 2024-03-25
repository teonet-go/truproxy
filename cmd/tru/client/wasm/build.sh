#!/bin/bash

set -e

GOOS=js GOARCH=wasm go build -o main.wasm
# cp main.wasm ../public

# Test run
# $(go env GOROOT)/misc/wasm/go_js_wasm_exec ./main.wasm