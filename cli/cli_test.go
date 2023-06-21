// (c) Siemens AG 2023
//
// SPDX-License-Identifier: MIT

package cli

import (
	"errors"

	"github.com/siemens/cshargextcap/cli/cliplugin"
	"github.com/spf13/cobra"
	"github.com/thediveo/go-plugger/v3"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func init() {
	plugger.Group[cliplugin.SetupCLI]().Register(UnittestSetupCLI)
	plugger.Group[cliplugin.BeforeCommand]().Register(UnittestBeforeRun)
}

var setupCLI = 0

func UnittestSetupCLI(rootCmd *cobra.Command) {
	setupCLI++
}

var beforeRun = 0
var beforeRunErr = error(nil)

func UnittestBeforeRun(*cobra.Command) error {
	beforeRun++
	return beforeRunErr
}

var _ = Describe("CLI cmd plugins", func() {

	It("calls AddFlags plugin method", func() {
		rootCmd := cobra.Command{}
		rootCmd.SetArgs([]string{})
		setupCLI = 0
		AddFlags(&rootCmd)
		Expect(setupCLI).To(Equal(1))
	})

	It("calls BeforeCommand plugin method", func() {
		beforeRun = 0
		beforeRunErr = errors.New("fooerror")
		Expect(BeforeCommand(nil)).To(HaveOccurred())
		Expect(beforeRun).To(Equal(1))

		beforeRunErr = nil
		Expect(BeforeCommand(nil)).ToNot(HaveOccurred())
		Expect(beforeRun).To(Equal(2))
	})

})
