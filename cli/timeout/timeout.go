// (c) Siemens AG 2026
//
// SPDX-License-Identifier: MIT

package timeout

import (
	"time"

	"github.com/siemens/csharg/v2"
	"github.com/spf13/cobra"
	"github.com/thediveo/clippy/cliplugin"
	"github.com/thediveo/go-plugger/v3"
)

// FlagName is the name of the CLI flag defined in this package.
const FlagName = "timeout"

// DefaultSeconds is the default service timeout in seconds.
const DefaultSeconds = uint(csharg.DefaultServiceTimeout / time.Second)

func init() {
	plugger.Group[cliplugin.SetupCLI]().Register(
		SetupCLI, plugger.WithPlugin("timeout"))
}

// SetupCLI registers the “--timeout” CLI flag.
func SetupCLI(rootCmd *cobra.Command) {
	rootCmd.PersistentFlags().Uint(FlagName, DefaultSeconds,
		"time limit in seconds for discovery and establishing capture connections")
}

// Service returns the time limit for discovery requests and for establishing
// capture connections. A zero timeout or a missing flag result in the default
// timeout, as csharg would otherwise not limit requests at all.
func Service(cmd *cobra.Command) time.Duration {
	secs, err := cmd.Flags().GetUint(FlagName)
	if err != nil || secs == 0 {
		return csharg.DefaultServiceTimeout
	}
	return time.Duration(secs) * time.Second
}
