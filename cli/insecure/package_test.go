// (c) Siemens AG 2026
//
// SPDX-License-Identifier: MIT

package insecure

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestInsecure(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "cshargextcap/cli/insecure package")
}
