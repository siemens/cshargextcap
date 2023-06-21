#!/bin/bash
set -e

if ! command -v goreportcard-cli &>/dev/null; then
    export PATH="$(go env GOPATH)/bin:$PATH"
    if ! command -v goreportcard-cli &>/dev/null; then
        rm -rf /tmp/goreportcard || true
        git clone https://github.com/gojp/goreportcard.git /tmp/goreportcard
        (cd /tmp/goreportcard && make install && go install ./cmd/goreportcard-cli)
        rm -rf /tmp/goreportcard || true
        # Install a somewhat recent ineffassign over the totally outdated one
        # that goreportcard still insists of infecting the system with.
        go install github.com/gordonklaus/ineffassign@latest
        # Install the missing misspell, oh well...
        go install github.com/client9/misspell/cmd/misspell@latest
    fi
fi

goreportcard-cli -v ./...
