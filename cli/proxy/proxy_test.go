// (c) Siemens AG 2023
//
// SPDX-License-Identifier: MIT

package proxy

import (
	"os"

	"github.com/spf13/cobra"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("proxy setting", func() {

	var rootCmd *cobra.Command

	BeforeEach(func() {
		oldsettings := map[string]string{}
		for _, envvarname := range crappyEnvVars {
			val, ok := os.LookupEnv(envvarname)
			if !ok {
				continue
			}
			oldsettings[envvarname] = val
		}
		DeferCleanup(func() {
			for name, val := range oldsettings {
				os.Setenv(name, val)
			}
		})

		rootCmd = &cobra.Command{}
		ProxySetupCLI(rootCmd)
	})

	It("sets proxy-related env variables", func() {
		args := []string{"foo", "--noproxy", "nope", "--httpproxy", "http://example.org", "--httpsproxy", "https://example.org"}
		Expect(rootCmd.ParseFlags(args)).To(Succeed())
		Expect(ProxyBeforeCommand(rootCmd)).To(Succeed())
		Expect(os.Getenv("no_proxy")).To(Equal("nope"))
		Expect(os.Getenv("NO_PROXY")).To(Equal("nope"))
		Expect(os.Getenv("http_proxy")).To(Equal("http://example.org"))
		Expect(os.Getenv("HTTP_PROXY")).To(Equal("http://example.org"))
		Expect(os.Getenv("https_proxy")).To(Equal("https://example.org"))
		Expect(os.Getenv("HTTPS_PROXY")).To(Equal("https://example.org"))
	})

	It("clears proxy-related env variables", func() {
		for _, envvar := range crappyEnvVars {
			os.Setenv(envvar, "free the Internet!")
		}
		args := []string{"foo", "--proxyoff"}
		Expect(rootCmd.ParseFlags(args)).To(Succeed())
		Expect(ProxyBeforeCommand(rootCmd)).To(Succeed())
		for _, envvar := range crappyEnvVars {
			_, ok := os.LookupEnv(envvar)
			Expect(ok).To(BeFalse(), "didn't expect the Spanish Inquisition: %s", envvar)
		}
	})

})
