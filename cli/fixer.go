// (c) Siemens AG 2023
//
// SPDX-License-Identifier: MIT

package cli

import (
	"os"
	"strings"
)

// FixArgs works around Wiresharks extcap CLI flag passing idiosyncrasies of
// especially passing bool flags in some situations as "--foo true" as opposed
// to either "--foo" or "--foo=true". We therefore simply tack on any "true"s
// and "false"s that immediately follow any "--flag" to the preceding flag.
func FixArgs() {
	args := os.Args
	fixed := []string{args[0]}
	for idx := 1; idx < len(args); idx++ {
		if strings.HasPrefix(args[idx], "--") && idx+1 < len(args) {
			switch args[idx+1] {
			case "false", "true":
				fixed = append(fixed, args[idx]+"="+args[idx+1])
				idx++
				continue
			}
		}
		fixed = append(fixed, args[idx])
	}
	os.Args = fixed
}
