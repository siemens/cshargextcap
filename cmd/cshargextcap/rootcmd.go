// (c) Siemens AG 2023
//
// SPDX-License-Identifier: MIT

package main

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/thediveo/clippy"

	"github.com/siemens/cshargextcap"
	"github.com/siemens/cshargextcap/cli/action"
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
