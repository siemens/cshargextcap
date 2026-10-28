/*
Package extcapargs provides a fluent API for defining the configuration
arguments of Wireshark extcap plugins, as well as the values of selector-like
arguments. It then emits them as the “arg” and “value” sentences Wireshark
expects in response to “--extcap-config” and “--extcap-reload-option”. The
sentences are documented in “Extcap Arguments” of the [Wireshark Developer's
Guide] and in the [extcap(4)] manual page.

Arguments are created using the constructor for the particular argument type,
such as [SelectorArg] or [FlagArg], passing the CLI flag Wireshark will use
to hand the argument's value to the extcap plugin. The remaining elements of an
argument are then set fluently, keeping the order of the calls in the emitted
sentence:

	args := map[int]extcapargs.ArgDef{
		10: extcapargs.SelectorArg("--container").WithReload().
			Display("Containers").Placeholder("refresh"),
		20: extcapargs.FlagArg("--noprom").
			Display("No promiscuous mode").Default(false),
	}

Values of selector-like arguments are created using [Value]:

	values := map[int][]extcapargs.ValueDef{
		10: {extcapargs.Value("any").Display("all").Default(true)},
	}

[ArgDef] and [ValueDef] are immutable: their methods return modified copies, so
shared definitions can be safely derived from. Setting the same element again
replaces its value in place.

Sets of arguments are keyed by their arg numbers, which also control the order
of the arguments in Wireshark's configuration dialog. Use [Merge] to combine
sets while rejecting duplicate arg numbers and CLI flags, and [Dump] as well as
[DumpValues] to finally emit the sentences.

[Wireshark Developer's Guide]: https://www.wireshark.org/docs/wsdg_html_chunked/ChCaptureExtcap.html
[extcap(4)]: https://www.wireshark.org/docs/man-pages/extcap.html
*/
package extcapargs
