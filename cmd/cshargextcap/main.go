// (c) Siemens AG 2023
//
// SPDX-License-Identifier: MIT

package main

import (
	"fmt"
	"log/slog"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/thediveo/clippy"

	"github.com/siemens/cshargextcap"
	"github.com/siemens/cshargextcap/cli"
	"github.com/siemens/cshargextcap/cli/action"
	_ "github.com/siemens/cshargextcap/cli/all"
)

// newRootCmd creates the root command with usage and version information, as
// well as the available CLI flags (including descriptions).
func newRootCmd() *cobra.Command {
	extcapNifNames := cshargextcap.ExtcapNifNames()
	rootCmd := &cobra.Command{
		Use: "cshargextcap",
		Short: `external capture plugin for Wireshark for remotely capturing network packets
from containers of various feathers, k8s pods, and other capture targets`,
		Version: fmt.Sprintf("%s (%s)",
			cshargextcap.SemVersion,
			strings.Join(extcapNifNames, ",")),
		Args: cobra.NoArgs,
		PersistentPreRunE: func(cmd *cobra.Command, _ []string) error {
			return clippy.BeforeCommand(cmd)
		},
		RunE: action.Call,
	}
	clippy.AddFlags(rootCmd)
	return rootCmd
}

func main() {
	cli.FixArgs()
	// This is cobra boilerplate documentation, except for the missing call to
	// fmt.Println(err) which in the original boilerplate is just plain wrong:
	// it renders the error message twice, see also:
	// https://github.com/spf13/cobra/issues/304
	if err := newRootCmd().Execute(); err != nil {
		for _, arg := range os.Args[1:] {
			slog.Warn("arg", slog.String("arg", arg))
		}
		os.Exit(1)
	}
}
