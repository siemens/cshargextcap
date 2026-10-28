/*
Package cli provides [FixArgs] to work around Wireshark's idiosyncratic way of
passing boolean extcap args, which needs to be applied to the CLI args before
parsing them.

The CLI args themselves are implemented in sub-packages as CLI plugins for
[clippy], each registering its flags and optionally a check or setup to run
before the selected action:

  - [github.com/siemens/cshargextcap/cli/action]: the mutually exclusive
    extcap action flags.
  - [github.com/siemens/cshargextcap/cli/target]: what to capture from.
  - [github.com/siemens/cshargextcap/cli/show]: which kinds of capture
    targets to offer for selection.
  - [github.com/siemens/cshargextcap/cli/wireshark]: further Wireshark extcap
    args, such as the fifo to write the packet capture stream into.
  - [github.com/siemens/cshargextcap/cli/proxy]: HTTP(S) proxy settings.
  - [github.com/siemens/cshargextcap/cli/insecure]: skipping server
    certificate verification.
  - [github.com/siemens/cshargextcap/cli/debug]: debug logging to a file.

Importing [github.com/siemens/cshargextcap/cli/all] pulls in all CLI plugins
as well as the extcap network interface implementations.

[clippy]: https://github.com/thediveo/clippy
*/
package cli
