// (c) Siemens AG 2026
//
// SPDX-License-Identifier: MIT

package cfg

import (
	"bytes"
	"maps"
	"slices"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("configuration args", func() {

	It("compares arg values by their display texts", func() {
		vals := ArgVals{
			{Display: "zoo", Container: "A"},
			{Display: "abc", Container: "Z"},
		}
		slices.SortFunc(vals, CompareArgVals)
		Expect(vals).To(HaveExactElements(
			HaveField("Display", "abc"),
			HaveField("Display", "zoo")))
		Expect(CompareArgVals(vals[0], vals[0])).To(BeZero())
	})

	It("dumps arg values", func() {
		var buff bytes.Buffer
		ArgVals{
			{Display: "zoo", Container: "Z"},
			{Display: "abc", Container: "A"},
		}.Dump(&buff, 42)
		Expect(buff.String()).To(Equal(
			"value {arg=42}{value=Z}{display=zoo}{default=false}\n" +
				"value {arg=42}{value=A}{display=abc}{default=false}\n"))
	})

	It("merges the common args", func() {
		Expect(slices.Collect(maps.Keys(CommonArgs))).To(ConsistOf(
			slices.Concat(
				slices.Collect(maps.Keys(CommonStdArgs)),
				slices.Collect(maps.Keys(CommonOtherTabArgs)))))
	})

})
