// (c) Siemens AG 2023
//
// SPDX-License-Identifier: MIT

//go:build !windows

package pipe

import (
	"os"

	"golang.org/x/sys/unix"

	log "github.com/sirupsen/logrus"
)

// WaitTillBreak continuously checks a fifo/pipe to see when it breaks. When
// called, WaitTillBreak blocks until the fifo/pipe finally has broken.
//
// This implementation leverages [syscall.Select].
func WaitTillBreak(fifo *os.File) {
	log.Debug("constantly monitoring packet capture fifo status...")
	fds := []unix.PollFd{
		{
			Fd:     int32(fifo.Fd()),
			Events: unix.POLLIN + unix.POLLERR,
		},
	}
	for {
		// Check the fifo becomming readable, which signals that it has been
		// closed. In this case, ex-termi-nate ;) Oh, and remember to correctly
		// initialize the fdset each time before calling select() ... well, just
		// because that's a good idea to do. :(
		n, err := unix.Poll(fds, 1000 /*ms*/)
		if err != nil {
			if err == unix.EINTR {
				continue
			}
			log.Debugf("capture fifo broken, reason: %s", err.Error())
			return
		}
		if n <= 0 {
			continue
		}
		log.Debugf("poll: %+v", fds)
		if fds[0].Revents&unix.POLLERR != 0 {
			// Either the pipe was broken by Wireshark, or we did break it on
			// purpose in the piping process. Anyway, we're done.
			log.Debug("capture fifo broken, stopped monitoring.")
			return
		}
	}
}
