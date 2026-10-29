SHELL:=/bin/bash
GOGEN:=go generate .
BUILDTAGS:="osusergo,netgo"

PKGTESTS:=alpine debian-12 debian-13 ubuntu-24.04 ubuntu-26.04 fedora-43 fedora-44 almalinux-9 almalinux-10
PKGTEST_IMAGE_alpine:=alpine:latest
PKGTEST_IMAGE_debian-12:=debian:12-slim
PKGTEST_IMAGE_debian-13:=debian:13-slim
PKGTEST_IMAGE_ubuntu-24.04:=ubuntu:24.04
PKGTEST_IMAGE_ubuntu-26.04:=ubuntu:26.04
# Fedora: keep in sync with the two current GA releases (each lives ~13 months).
PKGTEST_IMAGE_fedora-43:=fedora:43
PKGTEST_IMAGE_fedora-44:=fedora:44
PKGTEST_IMAGE_almalinux-9:=almalinux:9
PKGTEST_IMAGE_almalinux-10:=almalinux:10
# space-separated list, such as "linux/amd64 linux/arm64"; foreign platforms need QEMU binfmt.
PKGTEST_PLATFORMS?=linux/$(shell go env GOARCH)

# $(call report,STATUS,COLOR,TEXT) prints a status line, colorized only on terminals.
report = if [[ -t 1 && -z "$${NO_COLOR:-}" ]]; then echo -e "\033[1;$(2)m$(1)\033[0m $(3)"; else echo "$(1) $(3)"; fi

.PHONY: help clean cshargextcap dist dist-build dist-test $(addprefix pkgtest-,$(PKGTESTS)) e2e

help: ## list available targets
	@# Derived from Gomega's Makefile (github.com/onsi/gomega) under MIT License
	@grep -E '^[a-zA-Z0-9_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-16s\033[0m %s\n", $$1, $$2}'

dist: ## build snapshot packages+archives in dist/ and test the Linux packages
	$(MAKE) dist-build
	$(MAKE) dist-test
	@echo "🏁  done"

dist-build: ## build snapshot packages+archives in dist/ (without testing)
	@if goreleaser release --snapshot --clean && ls -lh dist/cshargextcap_*; then \
		$(call report,SUCCESS,32,dist-build); \
	else \
		$(call report,FAIL,31,dist-build); exit 1; \
	fi

dist-test: $(addprefix pkgtest-,$(PKGTESTS)) ## test the Linux packages in dist/ on all distros (use -j for parallel)

$(addprefix pkgtest-,$(PKGTESTS)): pkgtest-%:
	@for platform in $(PKGTEST_PLATFORMS); do \
		echo "== test package installation on $* ($${platform}) ==" \
		&& docker run --rm --pull=always --platform "$${platform}" \
			-v "$(CURDIR)/dist:/dist:ro" \
			-v "$(CURDIR)/packaging/linux/test/pkgtest.sh:/pkgtest.sh:ro" \
			$(PKGTEST_IMAGE_$*) /bin/sh /pkgtest.sh \
		|| { $(call report,FAIL,31,$* ($${platform})); exit 1; }; \
	done

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
