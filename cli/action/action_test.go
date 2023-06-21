// (c) Siemens AG 2023
//
// SPDX-License-Identifier: MIT

package action

import (
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	. "github.com/thediveo/success"
)

const stdoutfile = "/tmp/csharg-test-stdout"

var _ = Describe("extcap actions", func() {

	var rootCmd *cobra.Command

	BeforeEach(func() {
		rootCmd = &cobra.Command{
			// we need to provide a run fn, as otherwise cmd.Execute will return
			// a nil error, urgh.
			RunE: func(cmd *cobra.Command, args []string) error { return nil },
		}
		ActionSetupCLI(rootCmd)
	})

	DescribeTable("accepts a single action",
		func(action string) {
			args := []string{"--" + action}
			rootCmd.SetArgs(args)
			Expect(rootCmd.Execute()).To(Succeed())
		},
		Entry(nil, "capture"),
		Entry(nil, "extcap-interfaces"),
		Entry(nil, "extcap-dlts"),
		Entry(nil, "extcap-config"),
	)

	It("accepts no action", func() {
		args := []string{""}
		rootCmd.SetArgs(args)
		Expect(rootCmd.Execute()).To(Succeed())
	})

	Context("call to action", func() {

		var exitCode int

		BeforeEach(func() {
			oldOsExit := osExit
			oldLevel := log.GetLevel()
			oldStdout := os.Stdout
			oldStderr := os.Stderr
			os.Stdout = Successful(os.Create(stdoutfile))
			os.Stderr = os.Stdout
			log.SetLevel(log.FatalLevel)
			DeferCleanup(func() {
				os.Stdout.Close()
				os.Stdout = oldStdout
				os.Stderr = oldStderr
				os.Remove(stdoutfile)
				osExit = oldOsExit
				log.SetLevel(oldLevel)
			})

			exitCode = -1
			osExit = func(code int) { exitCode = code }
		})

		It("rejects mutually exclusive actions", func() {
			args := []string{"--capture", "--extcap-interfaces"}
			rootCmd.SetArgs(args)
			Expect(rootCmd.Execute()).To(MatchError(ContainSubstring("if any flags in the group")))
		})

		It("returns without action", func() {
			args := []string{"foo"}
			Expect(rootCmd.ParseFlags(args)).To(Succeed())
			Expect(Call(rootCmd, nil)).To(Succeed())
		})

		It("lists no DLTs", func() {
			args := []string{"foo", "--extcap-interfaces"}
			Expect(rootCmd.ParseFlags(args)).To(Succeed())
			Expect(Call(rootCmd, nil)).To(Succeed())
			Expect(exitCode).To(Equal(-1))
		})

		It("fails to list DLTs", func() {
			args := []string{"foo", "--extcap-dlts"}
			Expect(rootCmd.ParseFlags(args)).To(Succeed())
			Expect(Call(rootCmd, nil)).NotTo(Succeed())
			Expect(exitCode).To(BeNumerically(">", 0))
		})

		It("fails to configure", func() {
			args := []string{"foo", "--extcap-config"}
			Expect(rootCmd.ParseFlags(args)).To(Succeed())
			Expect(Call(rootCmd, nil)).NotTo(Succeed())
			Expect(exitCode).To(BeNumerically(">", 0))
		})

		It("fails to capture", func() {
			args := []string{"foo", "--capture"}
			Expect(rootCmd.ParseFlags(args)).To(Succeed())
			Expect(Call(rootCmd, nil)).NotTo(Succeed())
			Expect(exitCode).To(BeNumerically(">", 0))
		})

	})

})
