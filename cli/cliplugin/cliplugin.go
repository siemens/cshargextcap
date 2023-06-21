// (c) Siemens AG 2023
//
// SPDX-License-Identifier: MIT

/*
Package cliplugin defines the plugin group types for setting up the extcap
plugin CLI args in a modular way. This way, we can separate the individual
concerns into different [plugin groups] and at the same time provide extension
points for additional capture-service related extcap network interface types.

[plugin groups]: https://github.com/thediveo/go-plugger
*/
package cliplugin

import "github.com/spf13/cobra"

// SetupCLI defines an exposed plugin symbol type for adding things to a cobra
// root command.
type SetupCLI func(*cobra.Command)

// BeforeCommand defines an exposed plugin symbol type for running checks after
// the command line args have been processed and before running the command.
type BeforeCommand func(*cobra.Command) error
