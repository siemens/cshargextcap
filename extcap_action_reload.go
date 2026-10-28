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
	"log/slog"

	"github.com/spf13/cobra"

	"github.com/siemens/cshargextcap/cli/wireshark"
)

// extcapReloadOption routes the reload of a configuration option to the
// specified extcap nif.
func extcapReloadOption(cmd *cobra.Command, w io.Writer) int {
	slog.Debug("executing action --extcap-configure --extcap-reload-option",
		slog.String("option", wireshark.ReloadOption(cmd)))

	nif, ok := interfaceOk(cmd)
	if !ok {
		return 1
	}
	nif.ReloadOption(cmd, w)
	return 0
}
