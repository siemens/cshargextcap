// (c) Siemens AG 2023
//
// SPDX-License-Identifier: MIT

package debug

import (
	"bytes"
	"context"
	"log/slog"
	"os"

	"github.com/spf13/cobra"
	"github.com/thediveo/once"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	. "github.com/thediveo/success"
)

const logtmpname = "/tmp/csharg-test-logfile"

var _ = Describe("debug logging", func() {

	var rootCmd *cobra.Command

	BeforeEach(func() {
		DeferCleanup(slog.SetDefault, slog.Default())
		rootCmd = &cobra.Command{}
		DebugSetupCLI(rootCmd)
	})

	It("doesn't enable debug logging if not told to do so", func() {
		args := []string{"foo"}
		Expect(rootCmd.ParseFlags(args)).To(Succeed())
		Expect(DebugBeforeCommand(rootCmd)).To(Succeed())
		Expect(slog.Default().Enabled(context.Background(), slog.LevelDebug)).To(BeFalse())
	})

	It("enables debug and reports version and args", func() {
		args := []string{"foo", "--debug"}
		Expect(rootCmd.ParseFlags(args)).To(Succeed())

		buff := bytes.Buffer{}
		oldOutput := output
		DeferCleanup(func() { output = oldOutput })
		output = &buff

		Expect(DebugBeforeCommand(rootCmd)).To(Succeed())
		Expect(slog.Default().Enabled(context.Background(), slog.LevelDebug)).To(BeTrue())
		Expect(buff.String()).To(MatchRegexp(`version=\d+\.\d+\.\d+.*extcaps=.*args=.*-test\.timeout=`))
	})

	It("fails when logging to invalid file destination", func() {
		args := []string{"foo", "--debug-file=/nowhere/to/be/seen"}
		Expect(rootCmd.ParseFlags(args)).To(Succeed())
		Expect(DebugBeforeCommand(rootCmd)).To(MatchError(MatchRegexp(`cannot open log file ".*", reason: .*`)))
	})

	It("picks up CSHARK_LOG env var", func() {
		_ = os.Remove(logtmpname)
		closeLog := once.Once(func() { _ = logf.Close() })
		defer func() {
			closeLog.Do()
			_ = os.Unsetenv(EnvLogFile)
			_ = os.Remove(logtmpname)
		}()
		Expect(os.Setenv(EnvLogFile, logtmpname)).To(Succeed())

		args := []string{"foo"}
		Expect(rootCmd.ParseFlags(args)).To(Succeed())
		Expect(DebugBeforeCommand(rootCmd)).To(Succeed())

		closeLog.Do()
		Expect(string(Successful(os.ReadFile(logtmpname)))).To(MatchRegexp(`extcaps=.*args=`))
	})

})
