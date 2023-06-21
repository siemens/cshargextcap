// (c) Siemens AG 2023
//
// SPDX-License-Identifier: MIT

package cshargextcap

import (
	"io"

	"github.com/thediveo/go-plugger/v3"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

// TestNif mocks an external capture interface.
type TestNif struct {
	ExtcapNif
}

func (moo *TestNif) DLTs() []ExtcapNifDLT {
	return []ExtcapNifDLT{
		{Number: 42, Display: "testnif-" + moo.name},
	}
}
func (moo *TestNif) Configure(w io.Writer) int { return 0 }
func (moo *TestNif) ReloadOption(w io.Writer)  {}
func (moo *TestNif) Capture() int              { return 0 }

var _ = Describe("extcapnif", func() {

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
	})

	It("lists names of all registered extcap nifs", func() {
		Expect(ExtcapNifNames()).To(ConsistOf("test-moby", "test-moo"))
	})

	It("looks up extcap nifs", func() {
		for _, name := range ExtcapNifNames() {
			nif, ok := ExtcapNifByName(name)
			Expect(ok).To(BeTrue())
			Expect(nif.Name()).To(Equal(name))
			Expect(nif.DLTs()[0].Display).To(Equal("testnif-" + name))
		}
	})

	It("returns descriptions", func() {
		nif, ok := ExtcapNifByName("test-moo")
		Expect(ok).To(BeTrue())
		Expect(nif.Description()).To(Equal("captures moos"))
	})

})
