// (c) Siemens AG 2023
//
// SPDX-License-Identifier: MIT

package debug

import (
	"bytes"
	"os"

	log "github.com/sirupsen/logrus"
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
		rootCmd = &cobra.Command{}
		DebugSetupCLI(rootCmd)
	})

	AfterEach(func() {
		log.SetLevel(log.InfoLevel)
	})

	It("doesn't enable debug logging if not told to do so", func() {
		args := []string{"foo"}
		Expect(rootCmd.ParseFlags(args)).To(Succeed())
		Expect(DebugBeforeCommand(rootCmd)).To(Succeed())
		Expect(log.GetLevel()).To(Equal(log.InfoLevel))
	})

	It("enables debug and reports version and args", func() {
		args := []string{"foo", "--debug"}
		Expect(rootCmd.ParseFlags(args)).To(Succeed())

		defer func() { log.SetOutput(os.Stdout) }()
		buff := bytes.Buffer{}
		log.SetOutput(&buff)

		Expect(DebugBeforeCommand(rootCmd)).To(Succeed())
		Expect(log.GetLevel()).To(Equal(log.DebugLevel))
		Expect(string(buff.Bytes())).To(MatchRegexp(`(?s)version \d+.\d+.\d+.*extcaps:.*arg: "-test.timeout=.*"`))
	})

	It("fails when logging to invalid file destination", func() {
		args := []string{"foo", "--debug-file=/nowhere/to/be/seen"}
		Expect(rootCmd.ParseFlags(args)).To(Succeed())
		Expect(DebugBeforeCommand(rootCmd)).To(MatchError(MatchRegexp(`cannot open log file ".*", reason: .*`)))
	})

	It("picks up CSHARK_LOG env var", func() {
		os.Remove(logtmpname)
		closeLog := once.Once(func() { logf.Close() })
		defer func() {
			log.SetOutput(os.Stdout)
			closeLog.Do()
			os.Unsetenv(EnvLogFile)
			os.Remove(logtmpname)
		}()
		os.Setenv(EnvLogFile, logtmpname)

		args := []string{"foo"}
		Expect(rootCmd.ParseFlags(args)).To(Succeed())
		Expect(DebugBeforeCommand(rootCmd)).To(Succeed())

		closeLog.Do()
		Expect(string(Successful(os.ReadFile(logtmpname)))).To(MatchRegexp(`(?s)extcaps:.*arg:`))
	})

})
