// (c) Siemens AG 2023
//
// SPDX-License-Identifier: MIT

package cshargextcap

import (
	"github.com/onsi/gomega/gbytes"
	"github.com/siemens/cshargextcap/cli/target"
	"github.com/thediveo/go-plugger/v3"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("extcap_dlts", func() {

	var out *gbytes.Buffer

	BeforeEach(func() {
		extcapNifsSaved := extcapNifs
		DeferCleanup(func() {
			extcapNifs = extcapNifsSaved
		})
		// For testing, wipe out any automatic registrations, and only
		// test with our small test set.
		extcapNifs = nil
		plugger.Group[ExtcapNifActions]().Register(&TestNif{
			ExtcapNif{name: "test-moo", description: "captures moos"}})
		plugger.Group[ExtcapNifActions]().Register(&TestNif{
			ExtcapNif{name: "test-moby", description: "Captn Ahab gives +1"}})

		out = gbytes.NewBuffer()
		DeferCleanup(func() { target.Nif = "" })
	})

	It("lists the DLT for an external capture interface", func() {
		target.Nif = ExtcapNifNames()[0]
		Expect(ExtcapDlts(out)).To(BeZero())
		Expect(out).To(gbytes.Say("dlt {number=.*}{name=.*}{display=.*}\n"))
	})

})
