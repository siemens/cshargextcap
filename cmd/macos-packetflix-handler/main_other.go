// (c) Siemens AG 2026
//
// SPDX-License-Identifier: MIT

//go:build !darwin

package main

import (
	"fmt"
	"os"
)

func main() {
	fmt.Fprintln(os.Stderr, "packetflix-handler is macOS-only")
	os.Exit(1)
}
