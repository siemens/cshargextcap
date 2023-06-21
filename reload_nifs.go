// (c) Siemens AG 2023
//
// SPDX-License-Identifier: MIT

package cshargextcap

import (
	"fmt"
	"io"
	"sort"
	"strings"

	"github.com/siemens/cshargextcap/cfg"
	"github.com/siemens/cshargextcap/cli/target"
	log "github.com/sirupsen/logrus"
)

// ReloadContainerNifs reloads the list of network interfaces of a specific
// container. That is, it refreshes the list of the interface selector config
// arg. This helper is exported so additional extcaps can make use of it.
func ReloadContainerNifs(w io.Writer) {
	// Check the container information, from which we'll later extract
	// the list of network interfaces: if the container information is
	// invalid, then silently return only the catch-all "all" network
	// interface.
	target := target.Unpack()
	fmt.Fprintf(w, "value {arg=%d}{value=any}{display=all}{default=true}\n", cfg.ContainerNifs)
	if target == nil {
		log.Error("--container information missing or invalid")
		return
	}

	// Get the list of network interface names from the container, and sort
	// it.
	log.Debugf("target: \"%s\"", target.Name)
	nifs := target.NetworkInterfaces
	sort.Strings(nifs)
	log.Debugf("container network interfaces: %s", strings.Join(nifs, ", "))

	// Finally, we can generate the list of network interfaces for this specific
	// container, including the catch-all "all" (which we already printed out
	// above).
	for _, nif := range nifs {
		fmt.Fprintf(w, "value {arg=%d}{value=%s}{display=%s}{default=false}\n",
			cfg.ContainerNifs, nif, nif)
	}
}
