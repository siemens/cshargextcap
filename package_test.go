// (c) Siemens AG 2023
//
// SPDX-License-Identifier: MIT

package cshargextcap

import (
	"log/slog"
	"testing"

	"github.com/spf13/cobra"

	"github.com/siemens/cshargextcap/cli/target"
	"github.com/siemens/cshargextcap/cli/wireshark"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestContainerSharkExtCap(t *testing.T) {
	slog.SetLogLoggerLevel(slog.LevelDebug)

	RegisterFailHandler(Fail)
	RunSpecs(t, "cshargextcap package")
}

// newTestCmd returns a command with the target and Wireshark CLI flags parsed
// from the specified args.
func newTestCmd(args ...string) *cobra.Command {
	GinkgoHelper()
	cmd := &cobra.Command{}
	target.TargetSetupCLI(cmd)
	wireshark.SetupCLI(cmd)
	Expect(cmd.ParseFlags(args)).To(Succeed())
	return cmd
}
