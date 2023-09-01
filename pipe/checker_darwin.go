// (c) Siemens AG 2023
//
// SPDX-License-Identifier: MIT

//go:build darwin

package pipe

import (
	"os"
	"syscall"

	log "github.com/sirupsen/logrus"
)

// WaitTillBreak continuously checks a fifo/pipe to see when it breaks. When
// called, WaitTillBreak blocks until the fifo/pipe finally has broken.
//
// This implementation leverages [syscall.Select].
func WaitTillBreak(fifo *os.File) {
	log.Debug("constantly monitoring packet capture fifo status...")
	fds := syscall.FdSet{}
	for {
		// Check the fifo becomming readable, which signals that it has been
		// closed. In this case, ex-termi-nate ;) Oh, and remember to correctly
		// initialize the fdset each time before calling select() ... well, just
		// because that's a good idea to do. :(
		fds.Bits[fifo.Fd()>>3] |= 0x01 << (fifo.Fd() & 7)
		n := syscall.Select(
			int(fifo.Fd())+1, // highest fd is our file descriptor.
			&fds, nil, nil,   // only watch readable.
			nil, // no timeout, ever.
		)
		if n != nil {
			// Either the pipe was broken by Wireshark, or we did break it on
			// purpose in the piping process. Anyway, we're done.
			log.Debug("capture fifo broken, stopped monitoring.")
			return
		}
	}
}
