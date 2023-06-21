// (c) Siemens AG 2023
//
// SPDX-License-Identifier: MIT

package proxy

import (
	"os"

	"github.com/siemens/cshargextcap/cli/cliplugin"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/thediveo/go-plugger/v3"
)

// Off is true if any HTTP(S) proxy setting(s) should be categorically ignored.
var off bool

// httpProxy specifies the URL, if any, to an httpProxy proxy.
var httpProxy string

// httpsProxy specifies the URL, if any, to an httpsProxy proxy.
var httpsProxy string

// No specifies the destination(s) to which noProxy proxy should be used at all.
var noProxy string

// crappyEnvVars contain the names of proxy-related environment variables to
// sanitize.
var crappyEnvVars = []string{
	"http_proxy", "HTTP_PROXY",
	"https_proxy", "HTTPS_PROXY",
	"no_proxy", "NO_PROXY",
}

func init() {
	plugger.Group[cliplugin.SetupCLI]().Register(
		ProxySetupCLI, plugger.WithPlugin("proxy"))
	plugger.Group[cliplugin.BeforeCommand]().Register(
		ProxyBeforeCommand, plugger.WithPlugin("proxy"))
}

func ProxySetupCLI(rootCmd *cobra.Command) {
	pf := rootCmd.PersistentFlags()
	pf.BoolVar(&off, "proxyoff", false,
		"never use any HTTP/HTTPS proxy")
	pf.StringVar(&httpProxy, "httpproxy", "",
		"HTTP proxy URL")
	pf.StringVar(&httpsProxy, "httpsproxy", "",
		"HTTPS proxy URL")
	pf.StringVar(&noProxy, "noproxy", "",
		"destination(s) to which no proxy should be used")
}

func ProxyBeforeCommand(rootCmd *cobra.Command) error {
	// Carry out some dirty deeds before we can correctly kick into
	// action...
	if off {
		for _, crap := range crappyEnvVars {
			os.Unsetenv(crap)
		}
	} else {
		if httpProxy != "" {
			os.Setenv("http_proxy", httpProxy)
			os.Setenv("HTTP_PROXY", httpProxy)
		}
		if httpsProxy != "" {
			os.Setenv("https_proxy", httpsProxy)
			os.Setenv("HTTPS_PROXY", httpsProxy)
		}
		if noProxy != "" {
			os.Setenv("no_proxy", noProxy)
			os.Setenv("NO_PROXY", noProxy)
		}
	}

	for _, crap := range crappyEnvVars {
		log.Debugf("proxy setting: %s=%s", crap, os.Getenv(crap))
	}

	return nil
}
