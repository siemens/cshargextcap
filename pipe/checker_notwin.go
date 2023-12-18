// (c) Siemens AG 2023
//
// SPDX-License-Identifier: MIT

//go:build !windows

package pipe

import (
	"context"
	"os"

	"golang.org/x/sys/unix"

	log "github.com/sirupsen/logrus"
)

// WaitTillBreak continuously checks a fifo/pipe's producer end (writing end) to
// see when it breaks. When called, WaitTillBreak blocks until the fifo/pipe
// finally has broken. It also returns when the passed context is done.
//
// This implementation leverages [unix.Poll].
func WaitTillBreak(ctx context.Context, fifo *os.File) {
	log.Debug("constantly monitoring packet capture fifo status...")
	fds := []unix.PollFd{
		{
			Fd:     int32(fifo.Fd()),
			Events: 0, // we're interested only in POLLERR and that is ignored here anyway.
		},
	}
	for {
		// Check the fifo becomming readable, which signals that it has been
		// closed. In this case, ex-termi-nate ;) Oh, and remember to correctly
		// initialize the fdset each time before calling select() ... well, just
		// because that's a good idea to do. :(
		n, err := unix.Poll(fds, 100 /* ms */)
		select {
		case <-ctx.Done():
			log.Debug("context done while monitoring packet capture fifo")
			return
		default:
		}
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
		if fds[0].Revents&unix.POLLERR != 0 {
			// Either the pipe was broken by Wireshark, or we did break it on
			// purpose in the piping process. Anyway, we're done.
			log.Debug("capture fifo broken, stopped monitoring.")
			return
		}
	}
}
