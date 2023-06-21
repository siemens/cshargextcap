// (c) Siemens AG 2023
//
// SPDX-License-Identifier: MIT

package cshargextcap

import (
	"github.com/siemens/cshargextcap/cli/target"
	log "github.com/sirupsen/logrus"
)

// Implements the "--capture" action: it routes the capture action to
// the specified extcap network interface.
func ExtcapCapture() int {
	log.Debugf("executing action --capture for interface \"%s\"", target.Nif)
	nif, ok := validateInterface()
	if !ok {
		return 1
	}
	return nif.Capture()
}
