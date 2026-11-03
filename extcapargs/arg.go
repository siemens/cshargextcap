// (c) Siemens AG 2026
//
// SPDX-License-Identifier: MIT

package extcapargs

// ArgDef defines an extcap configuration argument, rendered as the elements
// of an “arg” sentence (without its number). Create ArgDefs using the
// “…Arg” constructors, such as [SelectorArg], one for each argument type.
type ArgDef struct {
	attrs attrs
}

// newArg returns the definition of an extcap configuration argument of the
// specified type.
func newArg(call, typ string) ArgDef {
	return ArgDef{attrs: attrs{{key: "call", value: call}, {key: "type", value: typ}}}
}

// IntegerArg returns an arg for entering an integer value.
func IntegerArg(call string) ArgDef { return newArg(call, "integer") }

// UnsignedArg returns an arg for entering an unsigned integer value.
func UnsignedArg(call string) ArgDef { return newArg(call, "unsigned") }

// LongArg returns an arg for entering a long integer value.
func LongArg(call string) ArgDef { return newArg(call, "long") }

// DoubleArg returns an arg for entering a floating point value.
func DoubleArg(call string) ArgDef { return newArg(call, "double") }

// StringArg returns an arg for entering a text value.
func StringArg(call string) ArgDef { return newArg(call, "string") }

// PasswordArg returns an arg for entering a masked text value that doesn't
// get saved with the other capture settings.
func PasswordArg(call string) ArgDef { return newArg(call, "password") }

// BoolValueArg returns a checkbox arg that is always passed on the command
// line together with its value, as in “--foo true” or “--foo false”.
func BoolValueArg(call string) ArgDef { return newArg(call, "boolean") }

// FlagArg returns a checkbox arg that is passed as a flag on the command line
// only when checked.
func FlagArg(call string) ArgDef { return newArg(call, "boolflag") }

// FileSelectArg returns an arg for selecting a file path.
func FileSelectArg(call string) ArgDef { return newArg(call, "fileselect") }

// FolderSelectArg returns an arg for selecting a folder path.
func FolderSelectArg(call string) ArgDef { return newArg(call, "folderselect") }

// SelectorArg returns an arg for selecting a single value from a list of
// values, which are defined using [ValueDef]s.
func SelectorArg(call string) ArgDef { return newArg(call, "selector") }

// EditSelectorArg returns a [SelectorArg] whose selected value can be edited.
func EditSelectorArg(call string) ArgDef { return newArg(call, "editselector") }

// RadioArg returns an arg for selecting a single value from a group of radio
// buttons, which are defined using [ValueDef]s.
func RadioArg(call string) ArgDef { return newArg(call, "radio") }

// MultiCheckArg returns an arg for selecting multiple values, which are
// defined using [ValueDef]s.
func MultiCheckArg(call string) ArgDef { return newArg(call, "multicheck") }

// TableArg returns an arg for a table of values populated by the user.
func TableArg(call string) ArgDef { return newArg(call, "table") }

// TimestampArg returns an arg for entering a date and time.
func TimestampArg(call string) ArgDef { return newArg(call, "timestamp") }

// Display sets the name of the argument shown in the configuration dialog.
func (a ArgDef) Display(name string) ArgDef { return ArgDef{a.attrs.with("display", name)} }

// Tooltip sets the tooltip of the argument in the configuration dialog.
func (a ArgDef) Tooltip(text string) ArgDef { return ArgDef{a.attrs.with("tooltip", text)} }

// Placeholder sets the text shown in an empty entry field, or on the reload
// button of a reloadable selector.
func (a ArgDef) Placeholder(text string) ArgDef {
	return ArgDef{a.attrs.with("placeholder", text)}
}

// Default sets the default value of the argument.
func (a ArgDef) Default(value any) ArgDef { return ArgDef{a.attrs.with("default", value)} }

// IsRequired marks the argument as needing a value before starting a capture.
func (a ArgDef) IsRequired() ArgDef { return ArgDef{a.attrs.with("required", true)} }

// WithReload marks a selector argument as reloadable, so that Wireshark
// queries its values using “--extcap-reload-option”.
func (a ArgDef) WithReload() ArgDef { return ArgDef{a.attrs.with("reload", true)} }

// InGroup places the argument on the named tab of the configuration dialog.
func (a ArgDef) InGroup(name string) ArgDef { return ArgDef{a.attrs.with("group", name)} }

// Call returns the CLI flag Wireshark passes the argument with.
func (a ArgDef) Call() string { return a.attrs.get("call") }

// String returns the “{key=value}” elements of the argument definition.
func (a ArgDef) String() string { return a.attrs.String() }
