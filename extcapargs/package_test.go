// (c) Siemens AG 2026
//
// SPDX-License-Identifier: MIT

package extcapargs

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestExtcapArgs(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "cshargextcap/extcapargs package")
}
