// (c) Siemens AG 2023
//
// SPDX-License-Identifier: MIT

package pipe

import (
	"context"
	"io"
	"os"
	"time"

	"golang.org/x/sys/unix"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gleak"
	. "github.com/thediveo/success"
)

var _ = Describe("pipes", func() {

	BeforeEach(func() {
		goodgos := Goroutines()
		DeferCleanup(func() {
			Eventually(Goroutines).Within(2 * time.Second).ProbeEvery(100 * time.Millisecond).
				ShouldNot(HaveLeaked(goodgos))
		})
	})

	It("detects on the write end when a pipe breaks", func(ctx context.Context) {
		// As Wireshark uses a named pipe it passes an extcap its name (path)
		// and then expects the extcap to open this named pipe for writing
		// packet capture data into it. For this test we simulate Wireshark
		// closing its reading end and we must properly detect this situation on
		// our writing end of the pipe.
		By("creating a temporary named pipe/fifo and opening its ends")
		tmpfifodir := Successful(os.MkdirTemp("", "test-fifo-*"))
		defer os.RemoveAll(tmpfifodir)

		fifoname := tmpfifodir + "/fifo"
		Expect(unix.Mkfifo(fifoname, 0600)).To(Succeed())

		// Open both ends of the named pipe, once for reading and once for
		// writing. As this is a rendevouz operation, we run the two open
		// operations concurrently and proceed after we've succeeded on both
		// ends.
		rch := make(chan *os.File)
		go func() {
			defer GinkgoRecover()
			rch <- Successful(os.OpenFile(fifoname, os.O_RDONLY, 0))
		}()

		wch := make(chan *os.File)
		go func() {
			defer GinkgoRecover()
			wch <- Successful(os.OpenFile(fifoname, os.O_WRONLY, 0))
		}()

		var r, w *os.File
		Eventually(rch).Should(Receive(&r))
		Eventually(wch).Should(Receive(&w))
		defer w.Close()

		go func() {
			defer GinkgoRecover()
			By("continously draining the read end of the pipe into /dev/null...")
			null := Successful(os.OpenFile("/dev/null", os.O_WRONLY, 0))
			defer null.Close()
			io.Copy(null, r)
			By("...pipe draining done")
		}()

		go func() {
			defer GinkgoRecover()
			time.Sleep(2 * time.Second)
			By("closing read end of pipe")
			Expect(r.Close()).To(Succeed())
		}()

		go func() {
			defer GinkgoRecover()
			time.Sleep(500 * time.Microsecond)
			By("writing some data into the pipe")
			w.WriteString("Wireshark rulez")
		}()

		By("waiting for pipe to break")
		ctx, cancel := context.WithTimeout(ctx, 4*time.Second)
		defer cancel()
		start := time.Now()
		WaitTillBreak(ctx, w)
		Expect(ctx.Err()).To(BeNil(), "break detection failed")
		Expect(time.Since(start).Milliseconds()).To(
			BeNumerically(">", 1900), "false positive: pipe wasn't broken yet")
	})

})
