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
	"github.com/thediveo/go-plugger/v3"
	"github.com/thediveo/nonstd/xslog"

	"github.com/siemens/cshargextcap/cli/cliplugin"
)

var Nif string
var Container string
var Nifs string
var CaptureFilter string
var NoPromiscuous bool

func init() {
	plugger.Group[cliplugin.SetupCLI]().Register(
		TargetSetupCLI, plugger.WithPlugin("target"))
}

func TargetSetupCLI(rootCmd *cobra.Command) {
	pf := rootCmd.PersistentFlags()
	pf.StringVar(&Nif, "extcap-interface", "",
		"the extcap network interface to select")
	pf.StringVar(&Container, "container", "",
		"pod, container, process reference details to capture from")
	pf.StringVar(&Nifs, "nif", "",
		"network interface name(s) to capture from, separated by slashes")
	pf.StringVar(&CaptureFilter, "extcap-capture-filter", "",
		"capture filter expression")
	pf.BoolVar(&NoPromiscuous, "noprom", false,
		"don't put the network interface(s) into promiscuous mode")
}

// Unpack the “--container” CLI arg into a target object, returning nil if
// there is anything going wrong.
func Unpack() *api.Target {
	if Container == "" {
		return nil
	}
	cstr, err := url.QueryUnescape(Container)
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
