SHELL:=/bin/bash
GOGEN:=go generate .
BUILDTAGS:="osusergo,netgo"

.PHONY: help clean dist pkgsite report run vuln

help: ## list available targets
	@# Derived from Gomega's Makefile (github.com/onsi/gomega) under MIT License
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-16s\033[0m %s\n", $$1, $$2}'

dist: ## build snapshot cshargextcap binary packages+archives in dist/
# gorelease will run go generate anyway
	@scripts/goreleaser.sh --snapshot --clean
	@ls -lh dist/cshargextcap_*
	@echo "🏁  done"

clean: ## cleans up build and testing artefacts
	rm -rf dist
	find . -name __debug_bin -delete
	rm -f coverage.html coverage.out coverage.txt

pkgsite: ## serves Go documentation on port 6060
	@echo "navigate to: http://localhost:6060/github.com/siemens/cshargextcap"
	@scripts/pkgsite.sh

report: ## runs goreportcard
	@scripts/goreportcard.sh

test: ## runs all tests
	go test -v -p=1 -count=1 ./...

vuln: ## run go vulnerabilities check
	@scripts/vuln.sh
