// (c) Siemens AG 2023
//
// SPDX-License-Identifier: MIT

package cshargextcap

import (
	"log/slog"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestContainerSharkExtCap(t *testing.T) {
	slog.SetLogLoggerLevel(slog.LevelDebug)

	RegisterFailHandler(Fail)
	RunSpecs(t, "cshargextcap package")
}
