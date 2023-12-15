// (c) Siemens AG 2023
//
// SPDX-License-Identifier: MIT

package pipe

import (
	"io"
	"os"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	. "github.com/thediveo/success"
	"golang.org/x/sys/unix"
)

var _ = Describe("pipes", func() {

	It("detects on the write end when a pipe breaks", func() {
		// As Wireshark uses a named pipe it passes an extcap its name (path)
		// and then expects the extcap to open this named pipe for writing
		// packet capture data into it. For this test we simulate Wireshark
		// closing its reading end and we must properly detect this situation on
		// our writing end of the pipe.
		By("creating a temporary named pipe/fifo and opening its ends")
		tmpfifodir := Successful(os.MkdirTemp("", "test-fifo-*"))
		defer os.RemoveAll(tmpfifodir)

		fifoname := tmpfifodir + "/fifo"
		unix.Mkfifo(fifoname, 0660)
		wch := make(chan *os.File)
		go func() {
			defer GinkgoRecover()
			wch <- Successful(os.OpenFile(fifoname, os.O_WRONLY, os.ModeNamedPipe))
		}()

		rch := make(chan *os.File)
		go func() {
			defer GinkgoRecover()
			rch <- Successful(os.OpenFile(fifoname, os.O_RDONLY, os.ModeNamedPipe))
		}()

		var r, w *os.File
		Eventually(rch).Should(Receive(&r))
		Eventually(wch).Should(Receive(&w))
		defer w.Close()

		go func() {
			defer GinkgoRecover()
			By("continously draining the read end of the pipe into /dev/null")
			null := Successful(os.OpenFile("/dev/null", os.O_WRONLY, 0))
			defer null.Close()
			io.Copy(null, r)
			By("pipe draining done")
		}()

		go func() {
			defer GinkgoRecover()
			time.Sleep(2 * time.Second)
			By("closing read end of pipe")
			Expect(r.Close()).To(Succeed())
		}()

		go func() {
			defer GinkgoRecover()
			time.Sleep(300 * time.Microsecond)
			By("writing some data into the pipe")
			w.WriteString("Wireshark rulez")
		}()

		start := time.Now()
		WaitTillBreak(w)
		Expect(time.Since(start).Milliseconds()).To(BeNumerically(">", 1900))
	})

})
