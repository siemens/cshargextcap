// (c) Siemens AG 2026
//
// SPDX-License-Identifier: MIT

package insecure

import (
	"github.com/spf13/cobra"
	"github.com/thediveo/clippy/cliplugin"
	"github.com/thediveo/go-plugger/v3"
)

// FlagName is the name of the CLI flag defined in this package.
const FlagName = "insecure"

func init() {
	plugger.Group[cliplugin.SetupCLI]().Register(
		SetupCLI, plugger.WithPlugin("insecure"))
}

// SetupCLI registers the “--insecure” CLI flag.
func SetupCLI(rootCmd *cobra.Command) {
	rootCmd.PersistentFlags().Bool(FlagName, false,
		"skip server certificate verification")
}

// SkipVerify returns true if server certificate verification should be
// skipped.
func SkipVerify(cmd *cobra.Command) bool {
	skip, _ := cmd.Flags().GetBool(FlagName)
	return skip
}
