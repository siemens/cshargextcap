// (c) Siemens AG 2026
//
// SPDX-License-Identifier: MIT

package cfg

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestCfg(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "cshargextcap/cfg package")
}
