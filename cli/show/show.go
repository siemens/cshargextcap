// (c) Siemens AG 2023
//
// SPDX-License-Identifier: MIT

package show

import (
	"github.com/siemens/csharg/v2/api"
	"github.com/spf13/cobra"
	"github.com/thediveo/clippy/cliplugin"
	"github.com/thediveo/go-plugger/v3"
)

// Names of the CLI flags defined and used in this package.
const (
	PodsFlagName                 = "showpods"
	StandaloneContainersFlagName = "showcontainers"
	ProcsFlagName                = "showprocs"
	EmptyNetNSFlagName           = "showemptynetns"
)

func init() {
	plugger.Group[cliplugin.SetupCLI]().Register(
		ShowSetupCLI, plugger.WithPlugin("show"))
}

func ShowSetupCLI(rootCmd *cobra.Command) {
	pf := rootCmd.PersistentFlags()
	pf.Bool(PodsFlagName, false,
		"discover pods")
	pf.Bool(StandaloneContainersFlagName, false,
		"discover stand-alone containers")
	pf.Bool(ProcsFlagName, false,
		"discover non-containers with virtual IP stacks")
	pf.Bool(EmptyNetNSFlagName, false,
		"discover IP stacks without any attached processes")
}

// Target returns true if a capture target should be shown according to the
// filter settings, or not.
func Target(cmd *cobra.Command, target *api.Target) bool {
	var flagName string
	switch target.Type {
	case "pod":
		flagName = PodsFlagName
	case "proc":
		flagName = ProcsFlagName
	case "bindmount":
		flagName = EmptyNetNSFlagName
	default:
		flagName = StandaloneContainersFlagName
	}
	show, _ := cmd.Flags().GetBool(flagName)
	return show
}
