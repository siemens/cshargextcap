// (c) Siemens AG 2026
//
// SPDX-License-Identifier: MIT

package extcapargs

// ValueDef defines a value of a selector-like extcap configuration argument,
// rendered as the elements of a “value” sentence (without its arg number).
type ValueDef struct {
	attrs attrs
}

// Value returns the definition of an argument value.
func Value(value string) ValueDef {
	return ValueDef{attrs: attrs{{key: "value", value: value}}}
}

// Display sets the text shown for this value.
func (v ValueDef) Display(text string) ValueDef { return ValueDef{v.attrs.with("display", text)} }

// Default sets whether this value is selected by default.
func (v ValueDef) Default(selected bool) ValueDef {
	return ValueDef{v.attrs.with("default", selected)}
}

// String returns the “{key=value}” elements of the value definition.
func (v ValueDef) String() string { return v.attrs.String() }
