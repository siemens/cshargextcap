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
// see when it breaks. When called, WaitTillBreak blocks until the fifo/named
// pipe finally has “broken”; that is, the reading end has been closed.
// WaitTillBreak also returns when the passed context is done.
//
// This implementation leverages [unix.Poll].
func WaitTillBreak(ctx context.Context, fifo *os.File) {
	log.Debug("constantly monitoring packet capture fifo status...")
	fifofd, err := unix.Dup(int(fifo.Fd()))
	if err != nil {
		log.Debugf("cannot duplicate packet capture fifo file descriptor, reason: %s", err.Error())
		return
	}
	defer unix.Close(fifofd)
	for {
		select {
		case <-ctx.Done():
			log.Debug("context done while monitoring packet capture fifo")
			return
		default:
		}
		fds := []unix.PollFd{
			{
				Fd:     int32(fifofd),
				Events: unix.POLLHUP, // we're interested only in POLLERR and that is ignored on input anyway.
			},
		}
		n, err := unix.Poll(fds, 100 /* ms */)
		if err != nil {
			if err == unix.EINTR {
				continue
			}
			log.Debugf("pipe polling failed, reason: %s", err.Error())
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
