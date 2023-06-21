// (c) Siemens AG 2023
//
// SPDX-License-Identifier: MIT

package action

import (
	"errors"
	"os"

	"github.com/siemens/cshargextcap"
	"github.com/siemens/cshargextcap/cli/cliplugin"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/thediveo/go-plugger/v3"
)

// ExtcapInterfaces is true when Wireshark wants to query the list of
// available external capture network interfaces using “--extcap-interfaces”.
var ExtcapInterfaces bool

// ExtcapDlts is true when Wireshark queries the DTLs of an external capture
// network interface using “--extcap-dlts”.
var ExtcapDlts bool

// ExtcapConfig is true when Wireshark queries the configuration options for an
// external capture network interface using “--extcap-config”.
var ExtcapConfig bool

// ExtcapCapture is true when Wireshark wants to start a capture on an external
// capture network interface using “--capture”.
var ExtcapCapture bool

func init() {
	plugger.Group[cliplugin.SetupCLI]().Register(
		ActionSetupCLI, plugger.WithPlugin("action"))
}

// ActionSetupCLI registers the “generic” extcap action CLI args (that are
// mutually exclusive).
func ActionSetupCLI(rootCmd *cobra.Command) {
	pf := rootCmd.PersistentFlags()
	pf.BoolVar(&ExtcapInterfaces, "extcap-interfaces", false,
		"query list of the available external capture network interfaces")
	pf.BoolVar(&ExtcapDlts, "extcap-dlts", false,
		"query DTLs of capture network interface")
	pf.BoolVar(&ExtcapConfig, "extcap-config", false,
		"configuration options for external capture network interface")
	pf.BoolVar(&ExtcapCapture, "capture", false,
		"capture from an external capture network interface")
	rootCmd.MarkFlagsMutuallyExclusive(
		"extcap-interfaces",
		"extcap-dlts",
		"extcap-config",
		"capture",
	)
}

var osExit = os.Exit // for unit testing

// Call the particular CLI arg-specified action. Exits with a non-zero exit code
// in case an action was called that then returned a non-zero exit code.
func Call(cmd *cobra.Command, _ []string) error {
	outcome := 0
	if ExtcapInterfaces {
		outcome = cshargextcap.ExtcapInterfaces(os.Stdout)
	} else if ExtcapDlts {
		outcome = cshargextcap.ExtcapDlts(os.Stdout)
	} else if ExtcapConfig {
		outcome = cshargextcap.ExtcapConfigure(os.Stdout)
	} else if ExtcapCapture {
		outcome = cshargextcap.ExtcapCapture()
	}
	if outcome != 0 {
		log.Errorf("finished with status code %d", outcome)
		osExit(outcome)
		return errors.New("fallen off the discworld")
	}
	log.Debug("finished with status code 0")
	return nil
}
