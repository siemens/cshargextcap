// (c) Siemens AG 2023
//
// SPDX-License-Identifier: MIT

package cfg

import (
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"

	"github.com/siemens/csharg"
	"golang.org/x/exp/maps"
	"golang.org/x/exp/slices"
)

// DumpConfigArgs dumps all configuration "args" to a writer, sorted by their
// arg numbers. This allows us to assemble the configuration args from common
// and additional elements without having to ensure that they are already in the
// correct order. Instead, we just now bring our configuration args into order,
// sorted by their arg numbers. This is necessary, as Wireshark has some quirks
// in the way it generates the extcap configuration Qt UI dialog. In particular,
// while it sorts the args by their numbers, it places the initial focus on the
// first config arg seen in a group, not the first one by its number.
func DumpConfigArgs(w io.Writer, options map[int]string) {
	keys := maps.Keys(options)
	slices.Sort(keys)
	for _, key := range keys {
		fmt.Fprintf(w, "arg {number=%d}%s\n", key, options[key])
	}
}

// ArgVal represents a single config argument value for the container list
// selector.
type ArgVal struct {
	Display   string // text displayed in selector
	Container string // container JSON data
}

// ArgVals is a (sortable) slice of ArgVal elements.
type ArgVals []ArgVal

// Sort sorts a slice of ArgVal elements in place.
func (a ArgVals) Sort() {
	slices.SortFunc(a, func(a, b ArgVal) int { return strings.Compare(a.Display, b.Display) })
}

// Dump all ArgVal elements to the specified writer, as values of the specified
// argument ID.
func (a ArgVals) Dump(w io.Writer, argid int) {
	for _, argval := range a {
		fmt.Fprintf(w, "value {arg=%d}{value=%s}{display=%s}{default=false}\n",
			argid, argval.Container, argval.Display)
	}
}

// DumpConfigArgValues dumps the values for configuration arguments. Here, order
// does not matter on the level of the config args, but only for the value set
// for a single config arg.
func DumpConfigArgValues(w io.Writer, values map[int][]string) {
	for key, perargvals := range values {
		for _, val := range perargvals {
			fmt.Fprintf(w, "value {arg=%d}%s\n", key, val)
		}
	}
}

// CommonArgs are the configuration args common to all extcap nifs of this
// plugin, except the packetflix nif (which lacks the Standard tab common args
// completely).
var CommonArgs = map[int]string{}

// Merge the common Standard tab args as well as the other tabs common args into
// a single set of command args.
func init() {
	maps.Copy(CommonArgs, CommonStdArgs)
	maps.Copy(CommonArgs, CommonOtherTabArgs)
}

// CommonStdArgs defines the configuration args in the Standard tab, common to
// all extcap nifs except the packetflix nif.
var CommonStdArgs = map[int]string{
	// Standard tab
	Containers: "{call=--container}{type=selector}{reload=true}" +
		"{display=Containers}{placeholder=refresh}",
	ContainerNifs: "{call=--nif}{type=selector}{reload=true}" +
		"{display=Interfaces}{placeholder=refresh}",

	ShowPods: "{call=--showpods}{type=boolflag}" +
		"{display=Show pods}",
	ShowStandaloneContainers: "{call=--showcontainers}{type=boolflag}" +
		"{display=Show standalone containers}",
	ShowProcs: "{call=--showprocs}{type=boolflag}" +
		"{display=Show processes}{default=false}",
	ShowEmptyNetNS: "{call=--showemptynetns}{type=boolflag}" +
		"{display=Show process-less IP stacks}{default=false}",

	NoProm: "{call=--noprom}{type=boolflag}" +
		"{display=No promiscuous mode}{default=false}",
}

// CommonOtherTabArgs defines the configuration args common to all config tabs,
// except the "Standard" tab.
var CommonOtherTabArgs = map[int]string{
	// Advanced tab (currently no shared config arguments)
	DiscoveryTimeout: "{group=" + AdvancedTabName + "}" +
		"{call=--timeout}{type=unsigned}{default=" + strconv.FormatUint(uint64(csharg.DefaultServiceTimeout/time.Second), 10) + "}" +
		"{display=Discovery timeout}{tooltip=discovery timeout in s}",

	// Poxy Proxy tab
	ProxyOff: "{group=" + ProxyTabName + "}{call=--proxyoff}{type=boolflag}{default=yes}" +
		"{display=Disable proxies}{tooltip=disables proxies completely}",
	HTTPProxy: "{group=" + ProxyTabName + "}{call=--httpproxy}{type=string}" +
		"{display=HTTP proxy}{placeholder=http://poxyproxy.acme.com:1234}",
	HTTPSProxy: "{group=" + ProxyTabName + "}{call=--httpsproxy}{type=string}" +
		"{display=HTTPS proxy}{placeholder=http://poxyproxy.acme.com:1234}",
	NoProxy: "{group=" + ProxyTabName + "}{call=--noproxy}{type=string}" +
		"{display=No proxy}{placeholder=destinations where no proxy should be used}",
}

// CommonArgValues defines the initial configuration values common to all extcap
// nifs (except packetflix)
var CommonArgValues = map[int][]string{
	ContainerNifs: {"{value=any}{display=all}{default=true}"},
}

var CommonOtherArgValues = map[int][]string{
	DiscoveryTimeout: {"value=}{default=+ strconv.FormatUint(uint64(csharg.DefaultServiceTimeout/time.Second), 10) +}"},
}
