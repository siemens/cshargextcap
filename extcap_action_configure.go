// (c) Siemens AG 2023
//
// SPDX-License-Identifier: MIT

// Implements the "--extcap-config" action: it handles the extcap network
// interface configuration action and routes it to the appropriate extcap nif operation.

package cshargextcap

import (
	"io"

	"github.com/siemens/cshargextcap/cli/target"
	"github.com/siemens/cshargextcap/cli/wireshark"
	log "github.com/sirupsen/logrus"
)

// Wireshark wants either to know the configuration args for a specific extcap
// network interface, or it wants to update the value(s) of a specific single
// configuration arg of a network interface.
func ExtcapConfigure(w io.Writer) int {
	// Is this a configure or a reload of configuration values...? If it's
	// a reload, then we route it to a separate extcap nif operation instead.
	if wireshark.ReloadOption != "" {
		return extcapReloadOption(w)
	}

	// Proceed instead with getting the configuration args for the specific
	// extcap network interface...
	log.Debugf("executing action --extcap-configure for interface \"%s\"", target.Nif)
	if nif, ok := validateInterface(); ok {
		return nif.Configure(w)
	}
	return 1
}
