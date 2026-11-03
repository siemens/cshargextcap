// (c) Siemens AG 2026
//
// SPDX-License-Identifier: MIT

package extcapargs

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("fluent arg and value definitions", func() {

	DescribeTable("constructs args of a specific type",
		func(ctor func(string) ArgDef, typ string) {
			Expect(ctor("--foo").String()).To(Equal("{call=--foo}{type=" + typ + "}"))
		},
		Entry(nil, IntegerArg, "integer"),
		Entry(nil, UnsignedArg, "unsigned"),
		Entry(nil, LongArg, "long"),
		Entry(nil, DoubleArg, "double"),
		Entry(nil, StringArg, "string"),
		Entry(nil, PasswordArg, "password"),
		Entry(nil, BoolValueArg, "boolean"),
		Entry(nil, FlagArg, "boolflag"),
		Entry(nil, FileSelectArg, "fileselect"),
		Entry(nil, FolderSelectArg, "folderselect"),
		Entry(nil, SelectorArg, "selector"),
		Entry(nil, EditSelectorArg, "editselector"),
		Entry(nil, RadioArg, "radio"),
		Entry(nil, MultiCheckArg, "multicheck"),
		Entry(nil, TableArg, "table"),
		Entry(nil, TimestampArg, "timestamp"),
	)

	It("renders an arg definition in the order of the fluent calls", func() {
		Expect(SelectorArg("--foo").WithReload().IsRequired().
			Display("Foo").Tooltip("all about foo").Placeholder("bar").
			InGroup("Tab").Default(42).String()).To(Equal(
			"{call=--foo}{type=selector}{reload=true}{required=true}" +
				"{display=Foo}{tooltip=all about foo}{placeholder=bar}" +
				"{group=Tab}{default=42}"))
	})

	It("replaces existing elements in place", func() {
		Expect(StringArg("--foo").Default("a").Display("Foo").Default("b").String()).
			To(Equal("{call=--foo}{type=string}{default=b}{display=Foo}"))
	})

	It("returns the CLI flag", func() {
		Expect(StringArg("--foo").Display("Foo").Call()).To(Equal("--foo"))
	})

	It("doesn't modify the arg definition derived from", func() {
		base := FlagArg("--foo").Display("Foo")
		derived := base.Default(true)
		Expect(base.String()).To(Equal("{call=--foo}{type=boolflag}{display=Foo}"))
		Expect(derived.String()).To(Equal("{call=--foo}{type=boolflag}{display=Foo}{default=true}"))
		Expect(base.Display("Bar").String()).To(Equal("{call=--foo}{type=boolflag}{display=Bar}"))
		Expect(derived.String()).To(Equal("{call=--foo}{type=boolflag}{display=Foo}{default=true}"))
	})

	It("renders a value definition", func() {
		base := Value("any").Display("all")
		Expect(base.Default(true).String()).To(Equal("{value=any}{display=all}{default=true}"))
		Expect(base.String()).To(Equal("{value=any}{display=all}"))
	})

})
