// (c) Siemens AG 2023
//
// SPDX-License-Identifier: MIT

package cshargextcap

import (
	"fmt"
	"io"
	"log/slog"

	"github.com/spf13/cobra"
)

// ExtcapDlts implements the "--extcap-dlts" action: it lists the DLTs for a
// specific external capture network interface, by querying the specified
// (internally-registered) extcap nif, and finally correctly formatting the
// answer; see “Ask For DLTs For Each Interface” in the [extcap guide].
//
// [extcap guide]: https://www.wireshark.org/docs/wsdg_html_chunked/ChCaptureExtcap.html
func ExtcapDlts(cmd *cobra.Command, w io.Writer) int {
	slog.Debug("executing action --extcap-dlts")
	nif, ok := interfaceOk(cmd)
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
