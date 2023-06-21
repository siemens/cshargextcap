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
	"os"
	"strings"

	"github.com/siemens/csharg"
	"github.com/siemens/cshargextcap/cli/target"
	"github.com/siemens/cshargextcap/cli/wireshark"
	"github.com/siemens/cshargextcap/pipe"

	log "github.com/sirupsen/logrus"
)

// Capture is the workhorse: it opens the pipe (fifo) offered by Wireshark, then
// starts a new Capture stream using the given SharkTank client and container
// target description. Then it lets csharg pump all packet Capture data arriving
// from the underlying websocket connected to the Capture service into the
// Wireshark pipe.
func Capture(st csharg.SharkTank) int {
	// Open packet stream pipe to Wireshark to feed it jucy packets...
	log.Debugf("fifo to Wireshark %s", wireshark.FifoPath)
	fifo, err := os.OpenFile(wireshark.FifoPath, os.O_WRONLY, os.ModeNamedPipe)
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
	defer cs.Stop() // be overly careful

	// Always keep an eye on the fifo getting closed by Wireshark: we then need
	// to stop the capture stream. This is necessary because the capture stream
	// might be idle for long times and thus we would otherwise not notice that
	// Wireshark has already stopped capturing.
	go func() {
		pipe.WaitTillBreak(fifo)
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
