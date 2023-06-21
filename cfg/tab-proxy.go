// (c) Siemens AG 2023
//
// SPDX-License-Identifier: MIT

package cfg

// ProxyTabName specifies the tab name for proxy configuration options.
const ProxyTabName = "Proxy"

// Shared "Proxy" tab controls; please node that the control IDs were randomly
// chosen, see also	https://xkcd.com/221/ for the algorithm used.
const (
	ProxyOff   = 600
	HTTPProxy  = 666
	HTTPSProxy = 667
	NoProxy    = 668
)
