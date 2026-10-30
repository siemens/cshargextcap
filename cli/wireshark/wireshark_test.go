// (c) Siemens AG 2026
//
// SPDX-License-Identifier: MIT

package wireshark

import (
	"github.com/spf13/cobra"
	"github.com/thediveo/clippy"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Wireshark CLI flags", func() {

	var rootCmd *cobra.Command

	BeforeEach(func() {
		rootCmd = &cobra.Command{}
		SetupCLI(rootCmd)
	})

	It("defaults to empty values", func() {
		Expect(rootCmd.ParseFlags(nil)).To(Succeed())
		Expect(ReloadOption(rootCmd)).To(BeEmpty())
		Expect(Version(rootCmd)).To(BeEmpty())
		Expect(FifoPath(rootCmd)).To(BeEmpty())
	})

	It("returns the parsed values", func() {
		Expect(rootCmd.ParseFlags([]string{
			"--extcap-reload-option=container",
			"--extcap-version=4.4.0",
			"--fifo", "/tmp/wireshark-fifo",
		})).To(Succeed())
		Expect(ReloadOption(rootCmd)).To(Equal("container"))
		Expect(Version(rootCmd)).To(Equal("4.4.0"))
		Expect(FifoPath(rootCmd)).To(Equal("/tmp/wireshark-fifo"))
	})

	It("returns empty values when the flags aren't registered", func() {
		cmd := &cobra.Command{}
		Expect(ReloadOption(cmd)).To(BeEmpty())
		Expect(Version(cmd)).To(BeEmpty())
		Expect(FifoPath(cmd)).To(BeEmpty())
	})

	It("registers its flags as a CLI plugin", func() {
		cmd := &cobra.Command{}
		clippy.AddFlags(cmd)
		Expect(cmd.ParseFlags([]string{"--fifo=/tmp/foo"})).To(Succeed())
		Expect(FifoPath(cmd)).To(Equal("/tmp/foo"))
	})

})
