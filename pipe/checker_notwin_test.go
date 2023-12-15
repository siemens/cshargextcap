// (c) Siemens AG 2023
//
// SPDX-License-Identifier: MIT

package pipe

import (
	"os"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	. "github.com/thediveo/success"
)

var _ = Describe("pipes", func() {

	It("detects on the write end when a pipe breaks", func() {
		// As Wireshark uses a named pipe it passes an extcap its name (path)
		// and then expects the extcap to open this named pipe for writing
		// packet capture data into it. For this test we simulate Wireshark
		// closing its reading end and we must properly detect this situation on
		// our writing end of the pipe.
		r, w := Successful2R(os.Pipe())
		defer w.Close()
		go func() {
			GinkgoHelper()
			time.Sleep(2 * time.Second)
			Expect(r.Close()).To(Succeed())
		}()
		start := time.Now()
		WaitTillBreak(w)
		Expect(time.Since(start).Milliseconds()).To(BeNumerically(">", 1900))
	})

})
