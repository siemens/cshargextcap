// (c) Siemens AG 2023
//
// SPDX-License-Identifier: MIT

package debug

import (
	"bytes"
	"context"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/thediveo/clippy"
	clippydebug "github.com/thediveo/clippy/debug"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	. "github.com/thediveo/success"
)

var _ = Describe("debug logging", func() {

	var rootCmd *cobra.Command

	BeforeEach(func() {
		DeferCleanup(slog.SetDefault, slog.Default())
		rootCmd = &cobra.Command{}
		rootCmd.SetContext(context.Background())
		clippy.AddFlags(rootCmd)
	})

	It("doesn't enable debug logging if not told to do so", func() {
		args := []string{"foo"}
		Expect(rootCmd.ParseFlags(args)).To(Succeed())
		Expect(clippy.BeforeCommand(rootCmd)).To(Succeed())
		Expect(slog.Default().Enabled(context.Background(), slog.LevelDebug)).To(BeFalse())
	})

	It("enables debug and reports version and args", func() {
		args := []string{"foo", "--debug"}
		Expect(rootCmd.ParseFlags(args)).To(Succeed())

		buff := bytes.Buffer{}
		clippydebug.SetWriter(rootCmd, &buff)

		Expect(clippy.BeforeCommand(rootCmd)).To(Succeed())
		Expect(slog.Default().Enabled(context.Background(), slog.LevelDebug)).To(BeTrue())
		Expect(buff.String()).To(MatchRegexp(`version=\S+ .*extcaps=.*args=.*-test\.timeout=`))
	})

	It("fails when logging to invalid file destination", func() {
		args := []string{"foo", "--debug-file=/nowhere/to/be/seen"}
		Expect(rootCmd.ParseFlags(args)).To(Succeed())
		Expect(clippy.BeforeCommand(rootCmd)).To(MatchError(MatchRegexp(`cannot open log file ".*", reason: .*`)))
	})

	It("picks up CSHARK_LOG env var", func() {
		logname := filepath.Join(GinkgoT().TempDir(), "cshargextcap.log")
		GinkgoT().Setenv(EnvLogFile, logname)

		args := []string{"foo"}
		Expect(rootCmd.ParseFlags(args)).To(Succeed())
		Expect(clippy.BeforeCommand(rootCmd)).To(Succeed())

		Expect(string(Successful(os.ReadFile(logname)))).To(MatchRegexp(`extcaps=.*args=`))
	})

})
