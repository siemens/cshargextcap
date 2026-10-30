/*
Packetflix-handler is the macOS handler for packetflix: URLs, starting a new
Wireshark instance that captures from the capture service given in the URL.

It is packaged as packetflix-handler.app by packaging/macos/build.sh, which
also works on Linux. macOS hands over URLs as Apple events, not command line
arguments, so packetflix-handler talks to Cocoa using purego instead of cgo
in order to keep cross-compiling.
*/
package main
