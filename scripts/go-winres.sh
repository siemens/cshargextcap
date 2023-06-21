#!/bin/bash
set -e

if ! command -v go-winres &>/dev/null; then
    export PATH="$(go env GOPATH)/bin:$PATH"
    if ! command -v go-winres &>/dev/null; then
        echo "installing go-winres..."
        go install github.com/tc-hib/go-winres@latest
    fi
fi

go-winres "$@"
