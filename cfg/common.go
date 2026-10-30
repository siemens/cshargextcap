// (c) Siemens AG 2023
//
// SPDX-License-Identifier: MIT

package cfg

import (
	"cmp"
	"fmt"
	"io"
	"time"

	"github.com/siemens/csharg/v2"

	"github.com/siemens/cshargextcap/extcapargs"
)

// ArgVal represents a single config argument value for the container list
// selector.
type ArgVal struct {
	Display   string // text displayed in selector
	Container string // container JSON data
}

// CompareArgVals compares two ArgVal elements by their display texts, returning
// a negative number when a sorts before b, a positive number when a sorts after
// b, and zero otherwise. It is intended for use with [slices.SortFunc] and
// friends, such as:
//
//	slices.SortFunc(argvals, cfg.CompareArgVals)
func CompareArgVals(a, b ArgVal) int {
	return cmp.Compare(a.Display, b.Display)
}

// ArgVals is a slice of ArgVal elements.
type ArgVals []ArgVal

// Dump all ArgVal elements to the specified writer, as values of the specified
// argument ID.
func (vals ArgVals) Dump(w io.Writer, argID int) {
	for _, val := range vals {
		_, _ = fmt.Fprintf(w, "value {arg=%d}%s\n",
			argID, extcapargs.Value(val.Container).Display(val.Display).Default(false))
	}
}

// CommonArgs are the configuration args common to all extcap nifs of this
// plugin, except the packetflix nif (which lacks the Standard tab common args
// completely).
var CommonArgs = func() map[int]extcapargs.ArgDef {
	args, err := extcapargs.Merge(CommonStdArgs, CommonOtherTabArgs)
	if err != nil {
		panic(err)
	}
	return args
}()

// CommonStdArgs defines the configuration args in the Standard tab, common to
// all extcap nifs except the packetflix nif.
var CommonStdArgs = map[int]extcapargs.ArgDef{
	Containers: extcapargs.SelectorArg("--container").WithReload().
		Display("Containers").Placeholder("refresh"),
	ContainerNifs: extcapargs.SelectorArg("--nif").WithReload().
		Display("Interfaces").Placeholder("refresh"),

	ShowPods: extcapargs.FlagArg("--showpods").
		Display("Show pods"),
	ShowStandaloneContainers: extcapargs.FlagArg("--showcontainers").
		Display("Show standalone containers"),
	ShowProcs: extcapargs.FlagArg("--showprocs").
		Display("Show processes").Default(false),
	ShowEmptyNetNS: extcapargs.FlagArg("--showemptynetns").
		Display("Show process-less IP stacks").Default(false),

	NoProm: extcapargs.FlagArg("--noprom").
		Display("No promiscuous mode").Default(false),
}

// CommonOtherTabArgs defines the configuration args common to all config tabs,
// except the "Standard" tab.
var CommonOtherTabArgs = map[int]extcapargs.ArgDef{
	// Advanced tab
	ServiceTimeout: extcapargs.UnsignedArg("--timeout").InGroup(AdvancedTabName).
		Default(uint(csharg.DefaultServiceTimeout / time.Second)).
		Display("Service timeout (s)").
		Tooltip("time limit in seconds for discovery and establishing capture connections"),

	// Poxy Proxy tab
	ProxyOff: extcapargs.FlagArg("--proxyoff").InGroup(ProxyTabName).Default(true).
		Display("Disable proxies").Tooltip("disables proxies completely"),
	HTTPProxy: extcapargs.StringArg("--httpproxy").InGroup(ProxyTabName).
		Display("HTTP proxy").Placeholder("http://poxyproxy.acme.com:1234"),
	HTTPSProxy: extcapargs.StringArg("--httpsproxy").InGroup(ProxyTabName).
		Display("HTTPS proxy").Placeholder("http://poxyproxy.acme.com:1234"),
	NoProxy: extcapargs.StringArg("--noproxy").InGroup(ProxyTabName).
		Display("No proxy").Placeholder("destinations where no proxy should be used"),
}

// CommonArgValues defines the initial configuration values common to all extcap
// nifs (except packetflix).
var CommonArgValues = map[int][]extcapargs.ValueDef{
	ContainerNifs: {extcapargs.Value("any").Display("all").Default(true)},
}
