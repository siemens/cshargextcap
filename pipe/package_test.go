// (c) Siemens AG 2023
//
// SPDX-License-Identifier: MIT

package pipe

import (
	"testing"

	log "github.com/sirupsen/logrus"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestContainerSharkExtCap(t *testing.T) {
	log.SetLevel(log.DebugLevel)

	RegisterFailHandler(Fail)
	RunSpecs(t, "cshargextcap/pipe")
}
