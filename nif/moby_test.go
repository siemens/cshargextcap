// (c) Siemens AG 2026
//
// SPDX-License-Identifier: MIT

package nif

import (
	"bytes"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gleak"
)

var _ = Describe("moby extcap nif", func() {

	var moby *MobyNif

	BeforeEach(func() {
		goodgos := Goroutines()
		DeferCleanup(func() {
			Eventually(Goroutines).ShouldNot(HaveLeaked(goodgos))
		})
		DeferCleanup(slog.SetDefault, slog.Default())
		slog.SetDefault(slog.New(slog.DiscardHandler))
		moby = &MobyNif{}
	})

	It("has a DLT", func() {
		Expect(moby.DLTs()).To(HaveLen(1))
	})

	It("dumps its configuration", func() {
		var buff bytes.Buffer
		Expect(moby.Configure(newTestCmd(), &buff)).To(BeZero())
		Expect(buff.String()).To(And(
			ContainSubstring("{call=--dockerhosturl}"),
			ContainSubstring("{call=--insecure}"),
			ContainSubstring("{call=--container}"),
			MatchRegexp(`\{call=--showpods\}.*\{default=false\}`),
			MatchRegexp(`\{call=--timeout\}.*\{default=30\}`),
			ContainSubstring("{group=Advanced}"),
			ContainSubstring("value {arg=20}{value=any}")))
	})

	It("reloads the container nifs", func() {
		var buff bytes.Buffer
		moby.ReloadOption(newTestCmd("--extcap-reload-option=nif"), &buff)
		Expect(buff.String()).To(Equal("value {arg=20}{value=any}{display=all}{default=true}\n"))
	})

	It("ignores unknown reload options", func() {
		var buff bytes.Buffer
		moby.ReloadOption(newTestCmd("--extcap-reload-option=foobar"), &buff)
		Expect(buff.Len()).To(BeZero())
	})

	It("returns no containers when the host isn't reachable", func() {
		var buff bytes.Buffer
		moby.ReloadOption(newTestCmd(
			"--extcap-reload-option=container",
			"--dockerhosturl=http://127.0.0.1:1"), &buff)
		Expect(buff.Len()).To(BeZero())
	})

	It("gives up discovering containers after the configured timeout", func() {
		release := make(chan struct{})
		srv := httptest.NewServer(http.HandlerFunc(func(http.ResponseWriter, *http.Request) {
			<-release
		}))
		DeferCleanup(srv.Close)
		DeferCleanup(func() { close(release) })

		var buff bytes.Buffer
		start := time.Now()
		moby.ReloadOption(newTestCmd(
			"--extcap-reload-option=container",
			"--dockerhosturl="+srv.URL,
			"--timeout=1"), &buff)
		Expect(time.Since(start)).To(BeNumerically("<", 5*time.Second))
		Expect(buff.Len()).To(BeZero())
	})

	It("fails capturing from an invalid host URL", func() {
		Expect(moby.Capture(newTestCmd("--dockerhosturl=http://127.0.0.1:1/?foo=bar"))).NotTo(BeZero())
	})

	It("fails capturing without Wireshark fifo", func() {
		Expect(moby.Capture(newTestCmd(
			"--dockerhosturl=http://127.0.0.1:1",
			"--fifo=/nowhere/to/be/found"))).NotTo(BeZero())
	})

})
