// (c) Siemens AG 2023
//
// SPDX-License-Identifier: MIT

// Implements the "--extcap-interfaces" action: lists all internally
// registered extcap nifs (network interfaces).

package cshargextcap

import (
	"fmt"
	"io"
	"log/slog"

	"github.com/siemens/cshargextcap/cli/target"
)

// ExtcapInterfaces lists the available (internally registered) external capture network
// interfaces.
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

// Convenience method to check the --extcap-interface argument; returns
// the extcap network interface object and an ok signal. If there is
// something wrong, then an additional error gets logged.
func validateInterface() (nif ExtcapNifActions, ok bool) {
	if target.Nif == "" {
		slog.Error("missing --extcap-interface")
		return
	}
	nif, ok = ExtcapNifByName(target.Nif)
	if !ok {
		slog.Error("unknown --extcap-interface", slog.String("interface", target.Nif))
	}
	return
}
