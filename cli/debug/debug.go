// (c) Siemens AG 2023
//
// SPDX-License-Identifier: MIT

package debug

import (
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/siemens/cshargextcap"
	"github.com/siemens/cshargextcap/cli/cliplugin"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/thediveo/go-plugger/v3"
	prefixed "github.com/x-cray/logrus-prefixed-formatter"
)

// EnvLogFile is the name of an environment variable specifying the filename
// (and path) of a log file. If the file doesn't exist, it will be created.
// Passing this environment variable to our extcap plugin works like
// --debug-file (which also implies --debug).
const EnvLogFile = "CSHARK_LOG"

// enabled is true if debug logging is enabled.
var enabled bool

// filename specifies the name of a file to write logging information to.
var filename string

func init() {
	plugger.Group[cliplugin.SetupCLI]().Register(
		DebugSetupCLI, plugger.WithPlugin("debug"))
	plugger.Group[cliplugin.BeforeCommand]().Register(
		DebugBeforeCommand, plugger.WithPlugin("debug"))
}

func DebugSetupCLI(rootCmd *cobra.Command) {
	pf := rootCmd.PersistentFlags()
	pf.BoolVar(&enabled, "debug", false,
		"enables debug logging")
	pf.StringVar(&filename, "debug-file", "",
		"enabled debug logging to file (implies --debug)")
}

var logf *os.File // allow closing under test

func DebugBeforeCommand(rootCmd *cobra.Command) error {
	// If a environment variable has been set to point to a log file, then use
	// this as "--debug-file" if the latter hasn't been specified as a CLI arg.
	// Otherwise, ignore it.
	if logfile := os.Getenv(EnvLogFile); logfile != "" && filename == "" {
		filename = logfile
	}

	// If "--debug-file" has been specified, append to this log file it it
	// already exists, otherwise create it. Specifying a debug log file also
	// includes "--debug".
	if filename != "" {
		enabled = true
		var err error
		logf, err = os.OpenFile(filename, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0660)
		if err != nil {
			return fmt.Errorf("cannot open log file %q, reason: %w",
				filename, err)
		}
		log.SetOutput(logf)
	}

	// If debugging has been enabled, then lower the logging level to
	// debug and then log our command line args.
	if enabled {
		log.SetLevel(log.DebugLevel)
		f := new(prefixed.TextFormatter)
		f.DisableColors = true
		f.ForceFormatting = true
		f.FullTimestamp = true
		f.TimestampFormat = "15:04:05"
		log.SetFormatter(f)

		log.Debugf("%s version %s", rootCmd.CalledAs(), cshargextcap.SemVersion)
		extcapNames := cshargextcap.ExtcapNifNames()
		sort.Strings(extcapNames)
		log.Debugf("extcaps: %s", strings.Join(extcapNames, ", "))
		for _, arg := range os.Args[1:] {
			log.Debugf("arg: \"%s\"", arg)
		}
	}

	return nil
}
