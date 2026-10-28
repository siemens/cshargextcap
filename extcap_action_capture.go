// (c) Siemens AG 2023
//
// SPDX-License-Identifier: MIT

package cshargextcap

import (
	"log/slog"

	"github.com/siemens/cshargextcap/cli/target"
)

// ExtcapCapture implements the "--capture" action: it routes the capture action to
// the specified extcap network interface.
func ExtcapCapture() int {
	slog.Debug("executing action --capture", slog.String("interface", target.Nif))
	nif, ok := validateInterface()
	if !ok {
		return 1
	}
	return nif.Capture()
}
