// (c) Siemens AG 2026
//
// SPDX-License-Identifier: MIT

package extcapargs

import (
	"bytes"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("arg sets", func() {

	It("dumps args sorted by their numbers", func() {
		var buff bytes.Buffer
		Dump(&buff, map[int]ArgDef{20: StringArg("--b"), 3: StringArg("--c"), 10: StringArg("--a")})
		Expect(buff.String()).To(Equal(
			"arg {number=3}{call=--c}{type=string}\n" +
				"arg {number=10}{call=--a}{type=string}\n" +
				"arg {number=20}{call=--b}{type=string}\n"))
	})

	It("dumps arg values sorted by arg numbers, keeping value order", func() {
		var buff bytes.Buffer
		DumpValues(&buff, map[int][]ValueDef{
			7: {Value("x"), Value("a")},
			3: {Value("z")},
		})
		Expect(buff.String()).To(Equal(
			"value {arg=3}{value=z}\nvalue {arg=7}{value=x}\nvalue {arg=7}{value=a}\n"))
	})

	It("merges args without modifying the sets", func() {
		a := map[int]ArgDef{1: StringArg("--a")}
		b := map[int]ArgDef{2: StringArg("--b")}
		merged, err := Merge(a, b)
		Expect(err).NotTo(HaveOccurred())
		Expect(merged).To(HaveLen(2))
		Expect(a).To(HaveLen(1))
		Expect(b).To(HaveLen(1))
	})

	It("rejects duplicate arg numbers", func() {
		_, err := Merge(map[int]ArgDef{1: StringArg("--a")}, map[int]ArgDef{1: StringArg("--b")})
		Expect(err).To(MatchError("duplicate config arg number 1"))
	})

	It("rejects duplicate arg calls", func() {
		_, err := Merge(map[int]ArgDef{1: StringArg("--a"), 2: FlagArg("--a")})
		Expect(err).To(MatchError(`duplicate config arg call "--a" for numbers 1 and 2`))
	})

})
