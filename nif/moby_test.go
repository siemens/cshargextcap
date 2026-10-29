// (c) Siemens AG 2026
//
// SPDX-License-Identifier: MIT

package nif

import (
	"bytes"
	"log/slog"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("moby extcap nif", func() {

	var moby *MobyNif

	BeforeEach(func() {
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

	It("fails capturing from an invalid host URL", func() {
		Expect(moby.Capture(newTestCmd("--dockerhosturl=http://127.0.0.1:1/?foo=bar"))).NotTo(BeZero())
	})

	It("fails capturing without Wireshark fifo", func() {
		Expect(moby.Capture(newTestCmd(
			"--dockerhosturl=http://127.0.0.1:1",
			"--fifo=/nowhere/to/be/found"))).NotTo(BeZero())
	})

})
