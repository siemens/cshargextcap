// (c) Siemens AG 2023
//
// SPDX-License-Identifier: MIT

// Implements the "--extcap-dlts" action: the DLTs it sends to Wireshark
// are gathered from the internally registered extcap nifs, and then sent
// to Wireshark correctly formatted.

package cshargextcap

import (
	"fmt"
	"io"
	"log/slog"
)

// ExtcapDlts lists the DLTs for a specific external capture network interface, by
// querying the specified (internally-registered) extcap nif, and finally
// correctly formatting the answer.
func ExtcapDlts(w io.Writer) int {
	slog.Debug("executing action --extcap-dlts")
	nif, ok := validateInterface()
	if !ok {
		return 1
	}

	for _, dlt := range nif.DLTs() {
		slog.Debug("DLT", slog.Int("number", dlt.Number), slog.String("display", dlt.Display))
		_, _ = fmt.Fprintf(w, "dlt {number=%d}{name=%s}{display=%s}\n",
			dlt.Number, "CLUSTERSHARK", dlt.Display)
	}
	return 0
}
