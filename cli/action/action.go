// (c) Siemens AG 2023
//
// SPDX-License-Identifier: MIT

package action

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/spf13/cobra"
	"github.com/thediveo/clippy/cliplugin"
	"github.com/thediveo/go-plugger/v3"

	"github.com/siemens/cshargextcap"
)

// Names of the mutually exclusive extcap action CLI flags.
const (
	ExtcapInterfacesFlagName = "extcap-interfaces"
	ExtcapDltsFlagName       = "extcap-dlts"
	ExtcapConfigFlagName     = "extcap-config"
	CaptureFlagName          = "capture"
)

func init() {
	plugger.Group[cliplugin.SetupCLI]().Register(
		ActionSetupCLI, plugger.WithPlugin("action"))
}

// ActionSetupCLI registers the “generic” extcap action CLI args (that are
// mutually exclusive).
func ActionSetupCLI(rootCmd *cobra.Command) {
	pf := rootCmd.PersistentFlags()
	pf.Bool(ExtcapInterfacesFlagName, false,
		"query list of the available external capture network interfaces")
	pf.Bool(ExtcapDltsFlagName, false,
		"query DTLs of capture network interface")
	pf.Bool(ExtcapConfigFlagName, false,
		"configuration options for external capture network interface")
	pf.Bool(CaptureFlagName, false,
		"capture from an external capture network interface")
	rootCmd.MarkFlagsMutuallyExclusive(
		ExtcapInterfacesFlagName,
		ExtcapDltsFlagName,
		ExtcapConfigFlagName,
		CaptureFlagName,
	)
}

// ExitCodeError reports that an extcap action finished with a non-zero exit
// code that the process should exit with.
type ExitCodeError int

func (e ExitCodeError) Error() string {
	return fmt.Sprintf("extcap action failed with exit code %d", int(e))
}

// Call the particular CLI arg-specified action. Returns an [ExitCodeError] in
// case the action returned a non-zero exit code.
func Call(cmd *cobra.Command, _ []string) error {
	isSet := func(name string) bool {
		set, _ := cmd.Flags().GetBool(name)
		return set
	}
	outcome := 0
	switch {
	case isSet(ExtcapInterfacesFlagName):
		outcome = cshargextcap.ExtcapInterfaces(os.Stdout)
	case isSet(ExtcapDltsFlagName):
		outcome = cshargextcap.ExtcapDlts(cmd, os.Stdout)
	case isSet(ExtcapConfigFlagName):
		outcome = cshargextcap.ExtcapConfigure(cmd, os.Stdout)
	case isSet(CaptureFlagName):
		outcome = cshargextcap.ExtcapCapture(cmd)
	}
	if outcome != 0 {
		slog.Error("finished", slog.Int("status", outcome))
		// The action already logged what went wrong, so cobra shouldn't add
		// its own error message and usage.
		cmd.SilenceErrors = true
		cmd.SilenceUsage = true
		return ExitCodeError(outcome)
	}
	slog.Debug("finished", slog.Int("status", 0))
	return nil
}
