// (c) Siemens AG 2026
//
// SPDX-License-Identifier: MIT

package nif

import (
	"testing"

	"github.com/spf13/cobra"
	"github.com/thediveo/clippy"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestNif(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "cshargextcap/nif package")
}

// newTestCmd returns a command with all registered CLI flags parsed from the
// specified args.
func newTestCmd(args ...string) *cobra.Command {
	GinkgoHelper()
	cmd := &cobra.Command{}
	clippy.AddFlags(cmd)
	Expect(cmd.ParseFlags(args)).To(Succeed())
	return cmd
}
