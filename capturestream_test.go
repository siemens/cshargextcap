// (c) Siemens AG 2026
//
// SPDX-License-Identifier: MIT

package cshargextcap

import (
	"errors"
	"io"
	"log/slog"
	"net/url"
	"os"
	"path/filepath"
	"sync"
	"syscall"

	"github.com/siemens/csharg/v2"
	"github.com/siemens/csharg/v2/api"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gleak"
)

// fakeSharkTank records the parameters of its last Capture call.
type fakeSharkTank struct {
	csharg.SharkTank
	stream *fakeStream
	err    error
	w      io.Writer
	target *api.Target
	opts   *csharg.CaptureOptions
}

func (st *fakeSharkTank) Capture(w io.Writer, t *api.Target, opts *csharg.CaptureOptions) (csharg.CaptureStreamer, error) {
	st.w, st.target, st.opts = w, t, opts
	if st.err != nil {
		return nil, st.err
	}
	return st.stream, nil
}

// fakeStream's Wait blocks until the stream gets stopped or finished.
type fakeStream struct {
	csharg.CaptureStreamer
	done chan struct{}
	once sync.Once
	err  error
}

func newFakeStream(err error) *fakeStream {
	return &fakeStream{done: make(chan struct{}), err: err}
}

func (s *fakeStream) finish() { s.once.Do(func() { close(s.done) }) }

func (s *fakeStream) Stop() error { s.finish(); return nil }

func (s *fakeStream) Wait() error { <-s.done; return s.err }

var _ = Describe("capture stream", func() {

	BeforeEach(func() {
		goodgos := Goroutines()
		DeferCleanup(func() {
			Eventually(Goroutines).ShouldNot(HaveLeaked(goodgos))
		})
		DeferCleanup(slog.SetDefault, slog.Default())
		slog.SetDefault(slog.New(slog.DiscardHandler))
	})

	It("fails without Wireshark fifo", func() {
		Expect(Capture(newTestCmd("--fifo=/nowhere/to/be/found"), nil)).NotTo(BeZero())
	})

	It("fails without container information", func() {
		fifo := filepath.Join(GinkgoT().TempDir(), "fifo")
		Expect(os.WriteFile(fifo, nil, 0o600)).To(Succeed())
		Expect(Capture(newTestCmd("--fifo="+fifo), nil)).NotTo(BeZero())
	})

	Context("with a capture service", func() {

		container := "--container=" + url.QueryEscape(
			`{"name":"foo","type":"docker","network-interfaces":["eth0","lo"]}`)

		var fifo string
		var sigs chan chan<- os.Signal

		BeforeEach(func() {
			fifo = filepath.Join(GinkgoT().TempDir(), "fifo")
			Expect(os.WriteFile(fifo, nil, 0o600)).To(Succeed())

			sigs = make(chan chan<- os.Signal, 1)
			oldNotify := notifySignals
			DeferCleanup(func() { notifySignals = oldNotify })
			notifySignals = func(c chan<- os.Signal, _ ...os.Signal) { sigs <- c }
		})

		It("fails when the capture cannot be started", func() {
			st := &fakeSharkTank{err: errors.New("no capture for you")}
			Expect(Capture(newTestCmd("--fifo="+fifo, container), st)).NotTo(BeZero())
		})

		It("passes the fifo, target, and capture options", func() {
			stream := newFakeStream(nil)
			stream.finish()
			st := &fakeSharkTank{stream: stream}
			Expect(Capture(newTestCmd("--fifo="+fifo, container,
				"--nif=eth1/eth2", "--extcap-capture-filter=tcp", "--noprom"), st)).To(BeZero())

			Expect(st.w).To(BeAssignableToTypeOf(&os.File{}))
			Expect(st.w.(*os.File).Name()).To(Equal(fifo))
			Expect(st.target).To(HaveField("Name", "foo"))
			Expect(st.opts.Nifs).To(ConsistOf("eth1", "eth2"))
			Expect(st.opts.Filter).To(Equal("tcp"))
			Expect(st.opts.AvoidPromiscuousMode).To(BeTrue())
		})

		It("falls back to the target's network interfaces", func() {
			stream := newFakeStream(nil)
			stream.finish()
			st := &fakeSharkTank{stream: stream}
			Expect(Capture(newTestCmd("--fifo="+fifo, container), st)).To(BeZero())
			Expect(st.opts.Nifs).To(ConsistOf("eth0", "lo"))
		})

		It("fails when the capture fails", func() {
			stream := newFakeStream(errors.New("broken stream"))
			stream.finish()
			st := &fakeSharkTank{stream: stream}
			Expect(Capture(newTestCmd("--fifo="+fifo, container), st)).NotTo(BeZero())
		})

		It("stops capturing upon a signal", func() {
			st := &fakeSharkTank{stream: newFakeStream(nil)}
			cmd := newTestCmd("--fifo="+fifo, container)
			result := make(chan int)
			go func() {
				defer GinkgoRecover()
				result <- Capture(cmd, st)
			}()

			var sig chan<- os.Signal
			Eventually(sigs).Should(Receive(&sig))
			Consistently(result, "100ms").ShouldNot(Receive())
			sig <- syscall.SIGINT
			Eventually(result).Should(Receive(BeZero()))
		})

	})

})
