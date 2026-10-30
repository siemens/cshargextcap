// (c) Siemens AG 2023
//
// SPDX-License-Identifier: MIT

package cshargextcap

import (
	"io"
	"log/slog"

	"github.com/spf13/cobra"

	"github.com/siemens/cshargextcap/cli/target"
	"github.com/siemens/cshargextcap/cli/wireshark"
)

// ExtcapConfigure handles the "--extcap-config" configuration action: Wireshark
// wants either to know the configuration args for a specific extcap network
// interface, or it wants to update the value(s) of a specific single
// configuration arg of a network interface. See “The Extcap Configuration
// Interface” and “Reload A Selector” in the [extcap guide].
//
// [extcap guide]: https://www.wireshark.org/docs/wsdg_html_chunked/ChCaptureExtcap.html
func ExtcapConfigure(cmd *cobra.Command, w io.Writer) int {
	// Is this a configure or a reload of configuration values...? If it's
	// a reload, then we route it to a separate extcap nif operation instead.
	if wireshark.ReloadOption(cmd) != "" {
		return extcapReloadOption(cmd, w)
	}

	// Proceed instead with getting the configuration args for the specific
	// extcap network interface...
	slog.Debug("executing action --extcap-configure", slog.String("interface", target.Nif(cmd)))
	if nif, ok := interfaceOk(cmd); ok {
		return nif.Configure(cmd, w)
	}
	return 1
}
