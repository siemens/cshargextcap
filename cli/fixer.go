// (c) Siemens AG 2023
//
// SPDX-License-Identifier: MIT

package cli

import (
	"strings"
)

// FixArgs works around Wireshark's extcap CLI flag passing idiosyncrasies of
// especially passing bool flags in some situations as "--foo true" as opposed
// to either "--foo" or "--foo=true". It returns a copy of the specified args
// with any "true" or "false" immediately following a "--flag" merged into
// that flag.
func FixArgs(args []string) []string {
	fixed := make([]string, 0, len(args))
	for _, arg := range args {
		if n := len(fixed); n > 0 && (arg == "true" || arg == "false") {
			if prev := fixed[n-1]; strings.HasPrefix(prev, "--") && !strings.Contains(prev, "=") {
				fixed[n-1] = prev + "=" + arg
				continue
			}
		}
		fixed = append(fixed, arg)
	}
	return fixed
}
