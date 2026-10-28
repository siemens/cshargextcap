// (c) Siemens AG 2023
//
// SPDX-License-Identifier: MIT

// Starts a live capture session by opening a websocket stream to the capture
// service, and then piping the packet capture stream into the pipe (fifo) made
// available by Wireshark. This module has basically outsourced most of its
// functionality to the csharg standalone package for better re-use.
//
// Additionally, also edits the beginning of the packet capture stream in order
// to insert meta data about the container capture origin, such as cluster
// identity, container identity, et cetera.

package cshargextcap

import (
	"log/slog"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/siemens/csharg/v2"
	"github.com/spf13/cobra"
	"github.com/thediveo/nonstd/xslog"

	"github.com/siemens/cshargextcap/cli/target"
	"github.com/siemens/cshargextcap/cli/wireshark"
)

// Capture is the workhorse: it opens the named pipe (fifo) offered by
// Wireshark, then starts a new Capture stream using the given SharkTank client
// and container target description. Then it lets csharg pump all packet Capture
// data arriving from the underlying websocket connected to the capture service
// into the Wireshark pipe.
func Capture(cmd *cobra.Command, st csharg.SharkTank) int {
	// Open packet stream pipe to Wireshark to feed it jucy packets...
	fifoPath := wireshark.FifoPath(cmd)
	slog.Debug("opening fifo to Wireshark", slog.String("path", fifoPath))
	fifo, err := os.OpenFile(fifoPath, os.O_WRONLY, 0)
	if err != nil {
		slog.Error("cannot open fifo", xslog.Error(err))
		return 1
	}
	defer func() { _ = fifo.Close() }()

	// Start the capture stream and wire it up to Wireshark's fifo...
	targt := target.Unpack(cmd)
	if targt == nil {
		slog.Error("--container information missing or invalid")
		return 1
	}
	nifs := strings.Split(target.Nifs(cmd), "/")
	if len(nifs) == 0 || nifs[0] == "" {
		// When no network interfaces were explicitly specified, then take the
		// list from the target (container) description in the hope that
		// it will be of more use.
		nifs = targt.NetworkInterfaces
	}
	slog.Debug("capturing",
		slog.Group("target",
			slog.String("type", targt.Type),
			slog.String("name", targt.Name)),
		slog.String("interfaces", strings.Join(nifs, ", ")))
	cs, err := st.Capture(fifo, targt, &csharg.CaptureOptions{
		Nifs:                 nifs,
		Filter:               target.CaptureFilter(cmd),
		AvoidPromiscuousMode: target.NoPromiscuous(cmd),
	})
	if err != nil {
		slog.Error("cannot start capture", xslog.Error(err))
		return 1
	}

	// Wireshark on unix systems sends SIGINT upon stopping a capture and
	// SIGTERM upon wanting to quit. We here use Debug logs as otherwise
	// Wireshark will report the logging as errors to the user. We only accept
	// that in case of a fatal abort when catching one of the signals twice or
	// one after the other.
	sigs := make(chan os.Signal, 1)
	go func() {
		fatal := false
		for sig := range sigs {
			switch sig {
			case syscall.SIGINT:
				slog.Debug("received SIGINT")
			case syscall.SIGTERM:
				slog.Debug("received SIGTERM")
			}
			if fatal {
				// twice a signal --> immediate abort
				slog.Error("aborting")
				os.Exit(1)
			}
			fatal = true
			slog.Debug("shutting down capture stream")
			go func() {
				_ = cs.Stop() // blocks, and is also idempotent.
			}()
		}
	}()
	// As mentioned above, we unconditionally handle SIGINT and SIGTERM on all
	// platforms. While this is currently not needed on Windows, some day it
	// might become alive.
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	defer func() {
		signal.Stop(sigs)
		close(sigs)
	}()

	defer func() { _ = cs.Stop() }() // be overly careful
	// ...and finally wait for the packet capture to terminate (or getting
	// ex-term-inated).
	if err := cs.Wait(); err != nil {
		slog.Error("packet capture failed", xslog.Error(err))
		return 1
	}

	slog.Debug("packet capture stopped")
	// Deferred teardown will now take place in any case ... and it
	// will tear down both the fifo and the websocket ... just to be
	// safe... :)
	return 0
}
