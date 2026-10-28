// (c) Siemens AG 2023
//
// SPDX-License-Identifier: MIT

package wireshark

import (
	"github.com/spf13/cobra"
	"github.com/thediveo/clippy/cliplugin"
	"github.com/thediveo/go-plugger/v3"
)

// Names of the CLI flags defined and used in this package.
const (
	ReloadOptionFlagName = "extcap-reload-option"
	VersionFlagName      = "extcap-version"
	FifoFlagName         = "fifo"
)

func init() {
	plugger.Group[cliplugin.SetupCLI]().Register(
		SetupCLI, plugger.WithPlugin("wireshark"))
}

// SetupCLI registers the Wireshark-specific extcap CLI flags.
func SetupCLI(rootCmd *cobra.Command) {
	pf := rootCmd.PersistentFlags()
	pf.String(ReloadOptionFlagName, "",
		"reload a specific configuration option")
	pf.String(VersionFlagName, "",
		"version of Wireshark invoking the extcap plugin")
	pf.String(FifoFlagName, "",
		"Wireshark fifo pathname to send captured packets to")
}

// ReloadOption returns the name of the configuration arg to reload values
// for, or "" if this isn't a reload.
func ReloadOption(cmd *cobra.Command) string {
	opt, _ := cmd.Flags().GetString(ReloadOptionFlagName)
	return opt
}

// Version returns the version of Wireshark invoking this plugin.
func Version(cmd *cobra.Command) string {
	v, _ := cmd.Flags().GetString(VersionFlagName)
	return v
}

// FifoPath returns the path of Wireshark's named pipe to write captured
// packets to.
func FifoPath(cmd *cobra.Command) string {
	path, _ := cmd.Flags().GetString(FifoFlagName)
	return path
}
