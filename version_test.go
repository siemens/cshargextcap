// (c) Siemens AG 2026
//
// SPDX-License-Identifier: MIT

package cshargextcap

import (
	"runtime/debug"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("module version", func() {

	It("has a version", func() {
		Expect(SemVersion).NotTo(BeEmpty())
	})

	DescribeTable("from build information",
		func(bi *debug.BuildInfo, expected string) {
			DeferCleanup(func(old func() (*debug.BuildInfo, bool)) { readBuildInfo = old }, readBuildInfo)
			readBuildInfo = func() (*debug.BuildInfo, bool) { return bi, bi != nil }
			Expect(buildInfoVersion()).To(Equal(expected))
		},
		Entry("without build information", nil, "(devel)"),
		Entry("as main module", &debug.BuildInfo{
			Main: debug.Module{Path: modulePath, Version: "v0.11.0"},
		}, "0.11.0"),
		Entry("as development main module", &debug.BuildInfo{
			Main: debug.Module{Path: modulePath, Version: "(devel)"},
		}, "(devel)"),
		Entry("as dependency", &debug.BuildInfo{
			Main: debug.Module{Path: "example.org/foo", Version: "v1.0.0"},
			Deps: []*debug.Module{
				{Path: "example.org/bar", Version: "v0.1.0"},
				{Path: modulePath, Version: "v0.10.7"},
			},
		}, "0.10.7"),
		Entry("as replaced dependency", &debug.BuildInfo{
			Main: debug.Module{Path: "example.org/foo"},
			Deps: []*debug.Module{
				{Path: modulePath, Version: "v0.10.7", Replace: &debug.Module{Path: "example.org/fork", Version: "v0.10.8"}},
			},
		}, "0.10.8"),
		Entry("as dependency replaced by local directory", &debug.BuildInfo{
			Main: debug.Module{Path: "example.org/foo"},
			Deps: []*debug.Module{
				{Path: modulePath, Version: "v0.10.7", Replace: &debug.Module{Path: "../cshargextcap"}},
			},
		}, "(devel)"),
		Entry("when missing", &debug.BuildInfo{
			Main: debug.Module{Path: "example.org/foo", Version: "v1.0.0"},
		}, "(devel)"),
	)

})
