// (c) Siemens AG 2024
//
// SPDX-License-Identifier: MIT

package timeout

import (
	"time"

	"github.com/siemens/csharg"
	"github.com/siemens/cshargextcap/cli/cliplugin"
	"github.com/spf13/cobra"
	"github.com/thediveo/go-plugger/v3"
)

// Timeout is the HTTP discovery request timeout in seconds.
var Discovery uint

func init() {
	plugger.Group[cliplugin.SetupCLI]().Register(
		TimeoutSetupCLI, plugger.WithPlugin("timeout"))
}

func TimeoutSetupCLI(rootCmd *cobra.Command) {
	pf := rootCmd.PersistentFlags()
	pf.UintVar(&Discovery, "timeout",
		uint(csharg.DefaultServiceTimeout/time.Second), "discovery timeout in seconds")
}
