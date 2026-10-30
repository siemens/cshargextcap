// (c) Siemens AG 2026
//
// SPDX-License-Identifier: MIT

package cli

import (
	"slices"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("fixing Wireshark's CLI args", func() {

	DescribeTable("merges bool values into their preceding flags",
		func(args []string, expected []string) {
			orig := slices.Clone(args)
			Expect(FixArgs(args)).To(Equal(expected))
			Expect(args).To(Equal(orig), "must not modify the args passed in")
		},
		Entry("no args", []string{}, []string{}),
		Entry("true", []string{"--foo", "true", "--bar"}, []string{"--foo=true", "--bar"}),
		Entry("false", []string{"--foo", "false"}, []string{"--foo=false"}),
		Entry("non-bool value", []string{"--foo", "bar"}, []string{"--foo", "bar"}),
		Entry("trailing flag", []string{"--foo"}, []string{"--foo"}),
		Entry("lone bool", []string{"true", "--foo", "true"}, []string{"true", "--foo=true"}),
		Entry("flag with value", []string{"--foo=bar", "true"}, []string{"--foo=bar", "true"}),
	)

})
