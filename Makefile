SHELL:=/bin/bash
GOGEN:=go generate .
BUILDTAGS:="osusergo,netgo"

.PHONY: help clean cshargextcap dist e2e test vuln

help: ## list available targets
	@# Derived from Gomega's Makefile (github.com/onsi/gomega) under MIT License
	@grep -E '^[a-zA-Z0-9_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-16s\033[0m %s\n", $$1, $$2}'

dist: ## build snapshot cshargextcap binary packages+archives in dist/ and test them
# gorelease will run go generate anyway
	@goreleaser release --snapshot --clean
	for distro in alpine debian fedora ubuntu; do \
		( \
			echo "== test package installation on $${distro} ==" \
			&& cd packaging/linux/test/$${distro} \
			&& ./test.sh \
		) || exit 1; \
	done
	@ls -lh dist/cshargextcap_*
	@echo "🏁  done"

cshargextcap: ## build local extcap binary
	go build \
		-tags netgo,osusergo \
		-ldflags "-s -w" \
		./cmd/cshargextcap

clean: ## cleans up build and testing artefacts
	rm -rf dist
	find . -name __debug_bin -delete
	rm -f coverage.html coverage.out coverage.txt

e2e: ## runs the opt-in end-to-end tests
	go test -v -race -p 1 -count=1 ./e2e/... -ginkgo.label-filter=e2e
