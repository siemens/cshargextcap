#!/bin/bash
set -e

if ! command -v govulncheck &>/dev/null; then
    export PATH="$(go env GOPATH)/bin:$PATH"
    if ! command -v govulncheck &>/dev/null; then
        echo "installing govulncheck..."
        go install golang.org/x/vuln/cmd/govulncheck@latest
    fi
fi
if [[ $(find "$(go env GOPATH)/bin/govulncheck" -mtime +1 -print) ]]; then
    echo "updating govulncheck to @latest..."
    go install golang.org/x/vuln/cmd/govulncheck@latest
fi

govulncheck ./...
