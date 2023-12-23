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
	"context"
	"os"
	"os/signal"
	"runtime"
	"strings"

	"github.com/siemens/csharg"
	"github.com/siemens/cshargextcap/cli/target"
	"github.com/siemens/cshargextcap/cli/wireshark"
	"github.com/siemens/cshargextcap/pipe"
	"golang.org/x/sys/unix"

	log "github.com/sirupsen/logrus"
)

// Capture is the workhorse: it opens the named pipe (fifo) offered by
// Wireshark, then starts a new Capture stream using the given SharkTank client
// and container target description. Then it lets csharg pump all packet Capture
// data arriving from the underlying websocket connected to the capture service
// into the Wireshark pipe.
func Capture(st csharg.SharkTank) int {
	if runtime.GOOS != "windows" {
		defer func() {
			signal.Reset(unix.SIGINT, unix.SIGTERM)
		}()
	}
	// Open packet stream pipe to Wireshark to feed it jucy packets...
	log.Debugf("opening fifo to Wireshark %s", wireshark.FifoPath)
	fifo, err := os.OpenFile(wireshark.FifoPath, os.O_WRONLY, 0)
	if err != nil {
		log.Errorf("cannot open fifo: %s", err.Error())
		return 1
	}
	defer fifo.Close()

	// Start the capture stream and wire it up to Wireshark's fifo...
	targt := target.Unpack()
	if targt == nil {
		log.Error("--container information missing or invalid")
		return 1
	}
	nifs := strings.Split(target.Nifs, "/")
	if len(nifs) == 0 || nifs[0] == "" {
		// When no network interfaces were explicitly specified, then take the
		// list from the target (container) description in the hope that
		// it will be of more use.
		nifs = targt.NetworkInterfaces
	}
	log.Debugf("capturing from: %q %q", targt.Type, targt.Name)
	log.Debugf("capturing from network interfaces: %s", strings.Join(nifs, ", "))
	cs, err := st.Capture(fifo, targt, &csharg.CaptureOptions{
		Nifs:                 nifs,
		Filter:               target.CaptureFilter,
		AvoidPromiscuousMode: target.NoPromiscuous,
	})
	if err != nil {
		log.Errorf("cannot start capture: %s", err.Error())
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
			case unix.SIGINT:
				log.Debug("received SIGINT")
			case unix.SIGTERM:
				log.Debug("received SIGTERM")
			}
			if fatal {
				// twice a signal --> immediate abort
				log.Fatal("aborting")
			}
			fatal = true
			log.Debug("shutting down capture stream")
			go func() {
				cs.Stop() // blocks, and is also idempotent.
			}()
		}
	}()
	if runtime.GOOS != "windows" {
		signal.Notify(sigs, unix.SIGINT, unix.SIGTERM)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer func() {
		cancel()
		cs.Stop() // be overly careful
	}()
	// Always keep an eye on the fifo getting closed by Wireshark: we then need
	// to stop the capture stream. This is necessary because the capture stream
	// might be idle for long times and thus we would otherwise not notice that
	// Wireshark has already stopped capturing.
	go func() {
		pipe.WaitTillBreak(ctx, fifo)
		cs.Stop()
	}()
	// ...and finally wait for the packet capture to terminate (or getting
	// ex-term-inated).
	cs.Wait()

	log.Debugf("packet capture stopped")
	// Deferred teardown will now take place in any case ... and it
	// will tear down both the fifo and the websocket ... just to be
	// safe... :)
	return 0
}
