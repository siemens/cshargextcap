// (c) Siemens AG 2026
//
// SPDX-License-Identifier: MIT

package wireshark

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestWireshark(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "cshargextcap/cli/wireshark package")
}
