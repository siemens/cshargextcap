// (c) Siemens AG 2023
//
// SPDX-License-Identifier: MIT

package show

import (
	"github.com/siemens/csharg/api"
	"github.com/siemens/cshargextcap/cli/cliplugin"
	"github.com/spf13/cobra"
	"github.com/thediveo/go-plugger/v3"
)

// Pods enables or disables discovery of k8s pods.
var Pods bool

// StandaloneContainers enables or disables discovery of non-pod containers.
var StandaloneContainers bool

// Procs enables or disables discovery of targets with only stand-alone
// processes.
var Procs bool

// EmptyNetNS enables or disables discovery of targets with no processes at all.
var EmptyNetNS bool

func init() {
	plugger.Group[cliplugin.SetupCLI]().Register(
		ShowSetupCLI, plugger.WithPlugin("show"))
}

func ShowSetupCLI(rootCmd *cobra.Command) {
	pf := rootCmd.PersistentFlags()
	pf.BoolVar(&Pods, "showpods", false,
		"discover pods")
	pf.BoolVar(&StandaloneContainers, "showcontainers", false,
		"discover stand-alone containers")
	pf.BoolVar(&Procs, "showprocs", false,
		"discover non-containers with virtual IP stacks")
	pf.BoolVar(&EmptyNetNS, "showemptynetns", false,
		"discover IP stacks without any attached processes")
}

// Target returns true if a capture target should be shown according to the
// filter settings, or not.
func Target(target *api.Target) bool {
	switch target.Type {
	case "pod":
		return Pods
	case "proc":
		return Procs
	case "bindmount":
		return EmptyNetNS
	default:
		return StandaloneContainers
	}
}
