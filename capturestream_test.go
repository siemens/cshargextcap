// (c) Siemens AG 2026
//
// SPDX-License-Identifier: MIT

package cshargextcap

import (
	"log/slog"
	"os"
	"path/filepath"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gleak"
)

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

})
