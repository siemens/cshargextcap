// (c) Siemens AG 2023
//
// SPDX-License-Identifier: MIT

package cshargextcap

import (
	"strings"

	"github.com/onsi/gomega/gbytes"
	"github.com/siemens/cshargextcap/cli/target"
	"github.com/thediveo/go-plugger/v3"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("extcap_interfaces", func() {

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
		target.Nif = ""
	})

	It("lists its external capture network interfaces and details", func() {
		Expect(ExtcapInterfaces(out)).To(BeZero())
		Expect(out).Should(gbytes.Say("extcap {version=.*}{help=.*}\n"))
		nifs := ExtcapNifNames()
		regex := "(" + strings.Join(nifs, "|") + ")"
		for i := 1; i <= len(nifs); i++ {
			Expect(out).Should(gbytes.Say("interface {value=" + regex + "}{display=.*}\n"))
		}
	})

	It("validates existing --extcap-interface name", func() {
		target.Nif = ExtcapNifNames()[0]
		nif, ok := validateInterface()
		Expect(ok).To(BeTrue())
		Expect(nif.Name()).To(Equal(target.Nif))
	})

	It("invalidates missing --extcap-interface", func() {
		_, ok := validateInterface()
		Expect(ok).To(BeFalse())
	})

	It("invalidates unknown --extcap-interface name", func() {
		target.Nif = "foobar"
		_, ok := validateInterface()
		Expect(ok).To(BeFalse())
	})

})
