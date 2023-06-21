// (c) Siemens AG 2023
//
// SPDX-License-Identifier: MIT

// Implements the "--extcap-interfaces" action: lists all internally
// registered extcap nifs (network interfaces).

package cshargextcap

import (
	"fmt"
	"io"

	"github.com/siemens/cshargextcap/cli/target"
	log "github.com/sirupsen/logrus"
)

// Lists the available (internally registered) external capture network
// interfaces.
func ExtcapInterfaces(w io.Writer) int {
	log.Debug("executing action --extcap-interfaces")

	fmt.Fprintf(w, "extcap {version=%s}{help=%s}\n", SemVersion, HelpURL)
	for _, name := range ExtcapNifNames() {
		nif, ok := ExtcapNifByName(name)
		if ok {
			description := nif.Description()
			log.Debugf("extcap nif: \"%s\" (%s)", name, description)
			fmt.Fprintf(w, "interface {value=%s}{display=%s}\n",
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
		log.Error("missing --extcap-interface")
		return
	}
	nif, ok = ExtcapNifByName(target.Nif)
	if !ok {
		log.Errorf("unknown --extcap-interface \"%s\"", target.Nif)
	}
	return
}
