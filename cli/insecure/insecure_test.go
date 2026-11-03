// (c) Siemens AG 2026
//
// SPDX-License-Identifier: MIT

package insecure

import (
	"github.com/spf13/cobra"
	"github.com/thediveo/clippy"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("--insecure CLI flag", func() {

	DescribeTable("skips verification only when told so",
		func(args []string, expected bool) {
			cmd := &cobra.Command{}
			clippy.AddFlags(cmd)
			Expect(cmd.ParseFlags(args)).To(Succeed())
			Expect(SkipVerify(cmd)).To(Equal(expected))
		},
		Entry("default", []string{}, false),
		Entry("--insecure", []string{"--insecure"}, true),
	)

	It("doesn't skip verification when the flag isn't registered", func() {
		Expect(SkipVerify(&cobra.Command{})).To(BeFalse())
	})

})
