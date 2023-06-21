#!/bin/bash
set -e

echo "checking pkgsite..."
go install golang.org/x/pkgsite/cmd/pkgsite@latest

echo "checking govulncheck..."
go install golang.org/x/vuln/cmd/govulncheck@latest

echo "checking gobadge..."
go install github.com/AlexBeauchemin/gobadge@latest

echo "checking go-winres..."
go install github.com/tc-hib/go-winres@latest

echo "checking goreportcard and friends..."
GOREPORTCARDTMPDIR="$(mktemp -d)"
trap 'rm -rf -- "$GOREPORTCARDTMPDIR"' EXIT
git clone https://github.com/gojp/goreportcard.git "$GOREPORTCARDTMPDIR"
(cd "$GOREPORTCARDTMPDIR" && make install && go install ./cmd/goreportcard-cli)
go install github.com/gordonklaus/ineffassign@latest
go install github.com/client9/misspell/cmd/misspell@latest

echo "checking pkgsite NPM helpers..."
(cd $HOME && npm update --silent browser-sync)
(cd $HOME && npm update --silent nodemon)
