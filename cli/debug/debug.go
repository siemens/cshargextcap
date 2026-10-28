// (c) Siemens AG 2023
//
// SPDX-License-Identifier: MIT

package debug

import (
	"fmt"
	"log/slog"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/thediveo/clippy/cliplugin"
	clippydebug "github.com/thediveo/clippy/debug"
	"github.com/thediveo/go-plugger/v3"

	"github.com/siemens/cshargextcap"
)

// EnvLogFile is the name of an environment variable specifying the filename
// (and path) of a log file. If the file doesn't exist, it will be created.
// Passing this environment variable to our extcap plugin works like
// --debug-file (which also implies --debug).
const EnvLogFile = "CSHARK_LOG"

// DebugFileFlagName is the name of the CLI flag defined in this package.
const DebugFileFlagName = "debug-file"

func init() {
	plugger.Group[cliplugin.SetupCLI]().Register(
		DebugSetupCLI, plugger.WithPlugin("debug-file"))
	plugger.Group[cliplugin.BeforeCommand]().Register(
		DebugFileBeforeCommand, plugger.WithPlugin("debug-file"),
		plugger.WithPlacement("<clippy/debug"))
	plugger.Group[cliplugin.BeforeCommand]().Register(
		DebugInfoBeforeCommand, plugger.WithPlugin("debug-info"),
		plugger.WithPlacement(">clippy/debug"))
}

// DebugSetupCLI adds the “--debug-file” flag; “--debug” comes from
// clippy/debug.
func DebugSetupCLI(rootCmd *cobra.Command) {
	rootCmd.PersistentFlags().String(DebugFileFlagName, "",
		"enables debug logging to file (implies --debug)")
}

// DebugFileBeforeCommand redirects logging into the file specified either by
// “--debug-file” or the CSHARK_LOG environment variable, forcing debug level.
// It must run before clippy/debug sets up the default logger. The log file
// stays open for the remaining lifetime of the process.
func DebugFileBeforeCommand(cmd *cobra.Command) error {
	filename, _ := cmd.Flags().GetString(DebugFileFlagName)
	if logfile := os.Getenv(EnvLogFile); logfile != "" && filename == "" {
		filename = logfile
	}
	if filename == "" {
		return nil
	}
	logf, err := os.OpenFile(filename, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0660)
	if err != nil {
		return fmt.Errorf("cannot open log file %q, reason: %w",
			filename, err)
	}
	clippydebug.SetWriter(cmd, logf)
	clippydebug.SetLevel(cmd, slog.LevelDebug)
	return nil
}

// DebugInfoBeforeCommand logs version, extcaps, and CLI args at debug level,
// after clippy/debug has set up the default logger.
func DebugInfoBeforeCommand(cmd *cobra.Command) error {
	slog.Debug("starting",
		slog.String("command", cmd.CalledAs()),
		slog.String("version", cshargextcap.SemVersion),
		slog.String("extcaps", strings.Join(cshargextcap.ExtcapNifNames(), ", ")),
		slog.Any("args", os.Args[1:]))
	return nil
}
