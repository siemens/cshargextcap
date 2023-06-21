// (c) Siemens AG 2023
//
// SPDX-License-Identifier: MIT

// Implements the "--extcap-reload-option" action: it routes the action
// to the correct specific internally-registered extcap nif. This nif
// then checks which configuration arg needs to be reloaded (refreshed)
// and sends updated values to Wireshark.

package cshargextcap

import (
	"io"

	"github.com/siemens/cshargextcap/cli/wireshark"
	log "github.com/sirupsen/logrus"
)

// Lists the DLTs for a specific external capture network interface.
func extcapReloadOption(w io.Writer) int {
	log.Debugf("executing action --extcap-configure --extcap-reload-option \"%s\"",
		wireshark.ReloadOption)

	nif, ok := validateInterface()
	if !ok {
		return 1
	}
	nif.ReloadOption(w)
	return 0
}
