// (c) Siemens AG 2023
//
// SPDX-License-Identifier: MIT

package target

import (
	"encoding/json"
	"log/slog"
	"net/url"

	"github.com/siemens/csharg/v2/api"
	"github.com/spf13/cobra"
	"github.com/thediveo/clippy/cliplugin"
	"github.com/thediveo/go-plugger/v3"
	"github.com/thediveo/nonstd/xslog"
)

// Names of the CLI flags defined and used in this package.
const (
	NifFlagName           = "extcap-interface"
	ContainerFlagName     = "container"
	NifsFlagName          = "nif"
	CaptureFilterFlagName = "extcap-capture-filter"
	NoPromFlagName        = "noprom"
)

func init() {
	plugger.Group[cliplugin.SetupCLI]().Register(
		TargetSetupCLI, plugger.WithPlugin("target"))
}

func TargetSetupCLI(rootCmd *cobra.Command) {
	pf := rootCmd.PersistentFlags()
	pf.String(NifFlagName, "",
		"the extcap network interface to select")
	pf.String(ContainerFlagName, "",
		"pod, container, process reference details to capture from")
	pf.String(NifsFlagName, "",
		"network interface name(s) to capture from, separated by slashes")
	pf.String(CaptureFilterFlagName, "",
		"capture filter expression")
	pf.Bool(NoPromFlagName, false,
		"don't put the network interface(s) into promiscuous mode")
}

// Nif returns the name of the selected extcap network interface.
func Nif(cmd *cobra.Command) string {
	nif, _ := cmd.Flags().GetString(NifFlagName)
	return nif
}

// Container returns the (URL-escaped JSON) capture target reference details.
func Container(cmd *cobra.Command) string {
	c, _ := cmd.Flags().GetString(ContainerFlagName)
	return c
}

// Nifs returns the slash-separated network interface names to capture from.
func Nifs(cmd *cobra.Command) string {
	nifs, _ := cmd.Flags().GetString(NifsFlagName)
	return nifs
}

// CaptureFilter returns the capture filter expression.
func CaptureFilter(cmd *cobra.Command) string {
	filter, _ := cmd.Flags().GetString(CaptureFilterFlagName)
	return filter
}

// NoPromiscuous returns true if the network interfaces must not be put into
// promiscuous mode.
func NoPromiscuous(cmd *cobra.Command) bool {
	noprom, _ := cmd.Flags().GetBool(NoPromFlagName)
	return noprom
}

// Unpack the “--container” CLI arg into a target object, returning nil if
// there is anything going wrong.
func Unpack(cmd *cobra.Command) *api.Target {
	container := Container(cmd)
	if container == "" {
		return nil
	}
	cstr, err := url.QueryUnescape(container)
	if err != nil {
		slog.Error("invalid --container arg", xslog.Error(err))
		return nil
	}
	target := &api.Target{}
	if err = json.Unmarshal([]byte(cstr), target); err != nil {
		slog.Error("invalid --container arg", xslog.Error(err))
		return nil
	}
	return target
}
