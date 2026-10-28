// (c) Siemens AG 2023
//
// SPDX-License-Identifier: MIT

package debug

import (
	"fmt"
	"io"
	"log/slog"
	"os"
	"sort"
	"strings"

	"github.com/spf13/cobra"
	"github.com/thediveo/go-plugger/v3"

	"github.com/siemens/cshargextcap"
	"github.com/siemens/cshargextcap/cli/cliplugin"
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

var output io.Writer = os.Stderr // allow redirecting under test

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
	out := output
	if filename != "" {
		enabled = true
		var err error
		logf, err = os.OpenFile(filename, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0660)
		if err != nil {
			return fmt.Errorf("cannot open log file %q, reason: %w",
				filename, err)
		}
		out = logf
	}

	// If debugging has been enabled, then lower the logging level to
	// debug and then log our command line args.
	if enabled {
		slog.SetDefault(slog.New(slog.NewTextHandler(out, &slog.HandlerOptions{
			Level: slog.LevelDebug,
			ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
				if a.Key == slog.TimeKey && len(groups) == 0 {
					return slog.String(a.Key, a.Value.Time().Format("15:04:05"))
				}
				return a
			},
		})))

		extcapNames := cshargextcap.ExtcapNifNames()
		sort.Strings(extcapNames)
		slog.Debug("starting",
			slog.String("command", rootCmd.CalledAs()),
			slog.String("version", cshargextcap.SemVersion),
			slog.String("extcaps", strings.Join(extcapNames, ", ")),
			slog.Any("args", os.Args[1:]))
	}

	return nil
}
