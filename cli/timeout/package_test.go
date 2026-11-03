// (c) Siemens AG 2026
//
// SPDX-License-Identifier: MIT

package timeout

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestTimeout(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "cshargextcap/cli/timeout package")
}
