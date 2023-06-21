// (c) Siemens AG 2023
//
// SPDX-License-Identifier: MIT

package show

import (
	"github.com/siemens/csharg/api"
	"github.com/spf13/cobra"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("showing or hiding targets", func() {

	var rootCmd *cobra.Command

	BeforeEach(func() {
		rootCmd = &cobra.Command{}
		ShowSetupCLI(rootCmd)
	})

	DescribeTable("the show must go on",
		func(flag string, ttype string) {
			Expect(Target(&api.Target{Type: ttype})).To(BeFalse())
			args := []string{"foo", flag}
			Expect(rootCmd.ParseFlags(args)).To(Succeed())
			Expect(Target(&api.Target{Type: ttype})).To(BeTrue())
		},
		Entry(nil, "--showpods", "pod"),
		Entry(nil, "--showprocs", "proc"),
		Entry(nil, "--showcontainers", "docker.com"),
		Entry(nil, "--showemptynetns", "bindmount"),
	)

})
