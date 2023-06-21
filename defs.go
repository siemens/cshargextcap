// (c) Siemens AG 2023
//
// SPDX-License-Identifier: MIT

package cshargextcap

const (
	// HelpURL points to the help page for this external capture plugin;
	// when a Wireshark user clicks on the "help" button in the configuration
	// dialog for this plugin, then this web page will be navigated to.
	HelpURL = "https://github.com/siemens/cshargextcap"

	// ServiceDefaultPort is the default port of the Packetflix+GhostWire
	// capture service for streaming packet captures snatched up in containers.
	ServiceDefaultPort = int32(5001)
)
