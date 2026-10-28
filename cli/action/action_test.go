// (c) Siemens AG 2023
//
// SPDX-License-Identifier: MIT

package action

import (
	"errors"
	"log/slog"
	"os"

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

		BeforeEach(func() {
			oldLogger := slog.Default()
			oldStdout := os.Stdout
			oldStderr := os.Stderr
			os.Stdout = Successful(os.Create(stdoutfile))
			os.Stderr = os.Stdout
			slog.SetDefault(slog.New(slog.DiscardHandler))
			DeferCleanup(func() {
				_ = os.Stdout.Close()
				os.Stdout = oldStdout
				os.Stderr = oldStderr
				_ = os.Remove(stdoutfile)
				slog.SetDefault(oldLogger)
			})
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
		})

		DescribeTable("returns a non-zero exit code on failure",
			func(flag string) {
				Expect(rootCmd.ParseFlags([]string{"foo", flag})).To(Succeed())
				err := Call(rootCmd, nil)
				code, ok := errors.AsType[ExitCodeError](err)
				Expect(ok).To(BeTrue(), "expected ExitCodeError, got: %v", err)
				Expect(int(code)).To(BeNumerically(">", 0))
				Expect(rootCmd.SilenceErrors).To(BeTrue())
				Expect(rootCmd.SilenceUsage).To(BeTrue())
			},
			Entry("DLTs", "--extcap-dlts"),
			Entry("configure", "--extcap-config"),
			Entry("capture", "--capture"),
		)

	})

})
