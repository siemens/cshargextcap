// (c) Siemens AG 2026
//
// SPDX-License-Identifier: MIT

package cshargextcap

import (
	"runtime/debug"
	"strings"
)

const modulePath = "github.com/siemens/cshargextcap"

// SemVersion is the semantic version string of the cshargextcap module.
// Release builds set it using “-ldflags -X”, otherwise it is taken from the Go
// build information.
var SemVersion string

// readBuildInfo is replaceable so tests can supply fake build information.
var readBuildInfo = debug.ReadBuildInfo

func init() {
	if SemVersion == "" {
		SemVersion = buildInfoVersion()
	}
}

// buildInfoVersion returns the cshargextcap module version from the binary's
// build information, regardless of cshargextcap being the main module or a
// dependency.
func buildInfoVersion() string {
	const devel = "(devel)"
	bi, ok := readBuildInfo()
	if !ok {
		return devel
	}
	for _, mod := range append([]*debug.Module{&bi.Main}, bi.Deps...) {
		if mod.Path != modulePath {
			continue
		}
		if mod.Replace != nil {
			mod = mod.Replace
		}
		if mod.Version == "" || mod.Version == devel {
			return devel
		}
		return strings.TrimPrefix(mod.Version, "v")
	}
	return devel
}
