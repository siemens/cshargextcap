// (c) Siemens AG 2023
//
// SPDX-License-Identifier: MIT

//go:build windows

package pipe

import (
	"os"
	"syscall"
	"time"

	log "github.com/sirupsen/logrus"
)

// WaitTillBreak continuously checks a fifo/pipe to see when it breaks. When
// called, PipeChecker blocks until the fifo/pipe finally has broken.
//
// As the Windows platform lacks a generally useful [syscall.Select]
// implementation that can also handle pipes. Instead, we will try to write 0
// octets at regular intervals to see if the pipe is broken. Usually,
// unsynchronized concurrent writes are a really bad idea, but in this case
// we're not really writing anything, but just poking things to see if they're
// dead already.
func WaitTillBreak(fifo *os.File) {
	log.Debug("constantly monitoring packet capture fifo status...")
	nothing := []byte{}
	ticker := time.NewTicker(1 * time.Second)
	for {
		select {
		case <-ticker.C:
			// Avoid the usual higher level writes, because of their
			// optimizations. While at this time the Windows writer
			// seems to write even zero-length data, we cannot be sure
			// this will hold for all future. So dive down into the
			// syscall basement to have full control.
			n, err := syscall.Write(syscall.Handle(fifo.Fd()), nothing)
			if n != 0 || err != nil {
				// Either the pipe was broken by Wireshark, or we
				// did break it on purpose in the piping process.
				// Anyway, we're done.
				log.Debug("capture fifo broken, stopped monitoring.")
				return
			}
		}
	}
}
