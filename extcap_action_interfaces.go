// (c) Siemens AG 2023
//
// SPDX-License-Identifier: MIT

package cshargextcap

import (
	"fmt"
	"io"
	"log/slog"

	"github.com/spf13/cobra"

	"github.com/siemens/cshargextcap/cli/target"
)

// ExtcapInterfaces implements the "--extcap-interfaces" action: it lists the
// available (internally registered) external capture network interfaces, see
// “Query For Available Interfaces” in the [extcap guide].
//
// [extcap guide]: https://www.wireshark.org/docs/wsdg_html_chunked/ChCaptureExtcap.html
func ExtcapInterfaces(w io.Writer) int {
	slog.Debug("executing action --extcap-interfaces")

	_, _ = fmt.Fprintf(w, "extcap {version=%s}{help=%s}\n", SemVersion, HelpURL)
	for _, name := range ExtcapNifNames() {
		nif, ok := ExtcapNifByName(name)
		if ok {
			description := nif.Description()
			slog.Debug("extcap nif", slog.String("name", name), slog.String("description", description))
			_, _ = fmt.Fprintf(w, "interface {value=%s}{display=%s}\n",
				name,
				fmt.Sprintf("%s (%s)", description, SemVersion))
		}
	}
	return 0
}

// interfaceOk checks the --extcap-interface argument; if ok, it returns the
// extcap network interface object and true. If there is something wrong, then
// an additional error gets logged and false is returned instead.
func interfaceOk(cmd *cobra.Command) (nif ExtcapNifActions, ok bool) {
	name := target.Nif(cmd)
	if name == "" {
		slog.Error("missing --extcap-interface")
		return
	}
	nif, ok = ExtcapNifByName(name)
	if !ok {
		slog.Error("unknown --extcap-interface", slog.String("interface", name))
	}
	return
}
