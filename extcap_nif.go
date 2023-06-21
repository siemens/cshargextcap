// (c) Siemens AG 2023
//
// SPDX-License-Identifier: MIT

// Provides a general handling infrastructure for external capture network
// interfaces: registration, lookup, and extcap actions interface.

package cshargextcap

import (
	"io"

	"github.com/thediveo/go-plugger/v3"
	"golang.org/x/exp/maps"
)

// ExtcapNifActions represents the actions that can be carried out on external
// capture network interfaces. These actions reflect the actions Wireshark
// signals to external capture plugins when invoking them. The only additional
// "internal" action here is returning the name of a specific external capture
// interface, which is used for looking up --ext-interface arguments.
type ExtcapNifActions interface {
	Name() string              // name of extcap network interface
	Description() string       // short description
	DLTs() []ExtcapNifDLT      // data-link layer encapsulation data type
	Configure(w io.Writer) int // dump configuration options
	ReloadOption(w io.Writer)  // dump up-to-date value(s) of specific configuration options
	Capture() int              // start packet capture
}

// ExtcapNif represents an individual, named external capture network
// interface.
type ExtcapNif struct {
	name        string // name of external capture network interface
	description string // short description of this extcap network interface
}

func NewExtcapNif(name string, description string) ExtcapNif {
	return ExtcapNif{name: name, description: description}
}

// ExtcapNifDLT describes a single DLT in terms of its number and
// description.
type ExtcapNifDLT struct {
	Number  int    // DLT number
	Display string // display string for this DLT
}

// Name returns the (guess what!) name of an external capture network
// interface. Surprise, surprise. Go for petty GoDoc rules.
func (e *ExtcapNif) Name() string {
	return e.name
}

// Description returns the short description of an external capture network
// interface.
func (e *ExtcapNif) Description() string {
	return e.description
}

// extcapNifs maps external capture network interfaces names to their specific
// implementation. To be more precise: the interfaces of the network interfaces,
// indexed by the network interface names.
var extcapNifs map[string]ExtcapNifActions

// populateDictionary ensures that the dictionary of extcap network interfaces
// is correctly populated based on the plugin registration. We can't do it in an
// init func because we cannot control order of all init's in this package, so
// that extcap plugins might not have yet allowed to register.
func populateDictionary() {
	if extcapNifs != nil {
		return
	}
	extcapNifs = map[string]ExtcapNifActions{}
	for _, nif := range plugger.Group[ExtcapNifActions]().Symbols() {
		extcapNifs[nif.Name()] = nif
	}
}

// ExtcapNifByName returns the external capture network interface with the
// specified name. Please note that capture network interfaces need to register
// their ExtcapNifActions implementation in the plugger group typed as
// ExtcapNifActions.
func ExtcapNifByName(name string) (ExtcapNifActions, bool) {
	populateDictionary()
	nif, ok := extcapNifs[name]
	return nif, ok
}

// ExtcapNifNames returns the list of names of the registered external capture
// network interfaces. The names in the returned list are in no particular
// order.
func ExtcapNifNames() []string {
	populateDictionary()
	return maps.Keys(extcapNifs)
}
