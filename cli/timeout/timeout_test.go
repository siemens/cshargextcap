// (c) Siemens AG 2026
//
// SPDX-License-Identifier: MIT

package timeout

import (
	"time"

	"github.com/siemens/csharg/v2"
	"github.com/spf13/cobra"
	"github.com/thediveo/clippy"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("--timeout CLI flag", func() {

	DescribeTable("returns the service timeout",
		func(args []string, expected time.Duration) {
			cmd := &cobra.Command{}
			clippy.AddFlags(cmd)
			Expect(cmd.ParseFlags(args)).To(Succeed())
			Expect(Service(cmd)).To(Equal(expected))
		},
		Entry("default", []string{}, csharg.DefaultServiceTimeout),
		Entry("--timeout=5", []string{"--timeout=5"}, 5*time.Second),
		Entry("--timeout=0", []string{"--timeout=0"}, csharg.DefaultServiceTimeout),
	)

	It("rejects negative timeouts", func() {
		cmd := &cobra.Command{}
		clippy.AddFlags(cmd)
		Expect(cmd.ParseFlags([]string{"--timeout=-1"})).NotTo(Succeed())
	})

	It("returns the default timeout when the flag isn't registered", func() {
		Expect(Service(&cobra.Command{})).To(Equal(csharg.DefaultServiceTimeout))
	})

})
