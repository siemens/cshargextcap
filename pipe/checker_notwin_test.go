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

	It("detects when a pipe breaks", func() {
		r, w := Successful2R(os.Pipe())
		defer r.Close()
		go func() {
			GinkgoHelper()
			time.Sleep(2 * time.Second)
			Expect(w.Close()).To(Succeed())
		}()
		start := time.Now()
		WaitTillBreak(r)
		Expect(time.Since(start).Milliseconds()).To(BeNumerically(">", 1500))
	})

})
