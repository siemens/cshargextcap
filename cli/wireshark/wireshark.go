// (c) Siemens AG 2023
//
// SPDX-License-Identifier: MIT

package wireshark

import (
	"github.com/spf13/cobra"
	"github.com/thediveo/go-plugger/v3"

	"github.com/siemens/cshargextcap/cli/cliplugin"
)

var ReloadOption string
var Version string
var FifoPath string

func init() {
	plugger.Group[cliplugin.SetupCLI]().Register(
		TargetSetupCLI, plugger.WithPlugin("target"))
}

func TargetSetupCLI(rootCmd *cobra.Command) {
	pf := rootCmd.PersistentFlags()
	pf.StringVar(&ReloadOption, "extcap-reload-option", "",
		"reload a specific configuration option")
	pf.StringVar(&Version, "extcap-version", "",
		"version of Wireshark invoking the extcap plugin")
	pf.StringVar(&FifoPath, "fifo", "",
		"Wireshark fifo pathname to send captured packets to")
}
