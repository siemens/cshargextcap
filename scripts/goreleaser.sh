#!/bin/bash
set -e

if ! command -v goreleaser &>/dev/null; then
    export PATH="$(go env GOPATH)/bin:$PATH"
    if ! command -v goreleaser &>/dev/null; then
        echo "installing goreleaser..."
        go install github.com/goreleaser/goreleaser@latest
    fi
fi
if [[ $(find "$(go env GOPATH)/bin/goreleaser" -mtime +1 -print) ]]; then
    echo "updating goreleaser to @latest..."
    go install github.com/goreleaser/goreleaser@latest
fi

goreleaser "$@"
