// (c) Siemens AG 2023
//
// SPDX-License-Identifier: MIT

package cshargextcap

import (
	"log/slog"

	"github.com/spf13/cobra"

	"github.com/siemens/cshargextcap/cli/target"
)

// ExtcapCapture implements the "--capture" action: it routes the capture action to
// the specified extcap network interface; see “The Extcap Capture Process” in
// the [extcap guide].
//
// [extcap guide]: https://www.wireshark.org/docs/wsdg_html_chunked/ChCaptureExtcap.html
func ExtcapCapture(cmd *cobra.Command) int {
	slog.Debug("executing action --capture", slog.String("interface", target.Nif(cmd)))
	nif, ok := interfaceOk(cmd)
	if !ok {
		return 1
	}
	return nif.Capture(cmd)
}
