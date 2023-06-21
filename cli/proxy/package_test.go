// (c) Siemens AG 2023
//
// SPDX-License-Identifier: MIT

package proxy

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestAction(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "cshargextcap/cli/proxy package")
}
