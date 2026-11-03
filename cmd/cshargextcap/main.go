// (c) Siemens AG 2023
//
// SPDX-License-Identifier: MIT

package main

import (
	"errors"
	"log/slog"
	"os"

	"github.com/siemens/cshargextcap/cli"
	"github.com/siemens/cshargextcap/cli/action"
	_ "github.com/siemens/cshargextcap/cli/all"
)

func main() {
	// Ensure consistent log formatting even when cobra fails before clippy
	// gets to set up logging.
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stderr, nil)))
	// This is cobra boilerplate documentation, except for the missing call to
	// fmt.Println(err) which in the original boilerplate is just plain wrong:
	// it renders the error message twice, see also:
	// https://github.com/spf13/cobra/issues/304
	rootCmd := newRootCmd()
	rootCmd.SetArgs(cli.FixArgs(os.Args[1:]))
	if err := rootCmd.Execute(); err != nil {
		if code, ok := errors.AsType[action.ExitCodeError](err); ok {
			os.Exit(int(code))
		}
		for _, arg := range os.Args[1:] {
			slog.Warn("arg", slog.String("arg", arg))
		}
		os.Exit(1)
	}
}
