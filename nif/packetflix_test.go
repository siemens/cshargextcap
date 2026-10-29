// (c) Siemens AG 2026
//
// SPDX-License-Identifier: MIT

package nif

import (
	"bytes"
	"log/slog"
	"net/url"

	"github.com/siemens/cshargextcap/cli/target"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("packetflix extcap nif", func() {

	var pf *PacketflixNif

	BeforeEach(func() {
		DeferCleanup(slog.SetDefault, slog.Default())
		slog.SetDefault(slog.New(slog.DiscardHandler))
		pf = &PacketflixNif{}
	})

	DescribeTable("converts packetflix: URLs into host URLs",
		func(pfurl, expected string, expectedOpts url.Values) {
			u, opts, err := hostURL(pfurl)
			Expect(err).NotTo(HaveOccurred())
			Expect(u.String()).To(Equal(expected))
			Expect(opts).To(Equal(expectedOpts))
		},
		Entry("ws", "packetflix:ws://host:5001/capture?container=foo&nif=eth0",
			"http://host:5001/", url.Values{"container": {"foo"}, "nif": {"eth0"}}),
		Entry("wss with service path", "packetflix:wss://host/svc/capture#frag",
			"https://host/svc", url.Values{}),
		Entry("schemeless", "packetflix://host:5001/capture?chaste",
			"https://host:5001/", url.Values{"chaste": {""}}),
	)

	It("rejects an invalid packetflix: URL", func() {
		_, _, err := hostURL("packetflix:ws://[::1")
		Expect(err).To(HaveOccurred())
	})

	It("dumps its configuration", func() {
		var buff bytes.Buffer
		Expect(pf.Configure(newTestCmd(), &buff)).To(BeZero())
		Expect(buff.String()).To(And(
			ContainSubstring("{call=--url}"),
			ContainSubstring("{call=--insecure}"),
			ContainSubstring("{call=--proxyoff}"),
			Not(ContainSubstring("{call=--container}"))))
	})

	It("doesn't reload any options", func() {
		var buff bytes.Buffer
		pf.ReloadOption(newTestCmd("--extcap-reload-option=container"), &buff)
		Expect(buff.Len()).To(BeZero())
	})

	It("fails capturing from an invalid URL", func() {
		Expect(pf.Capture(newTestCmd("--url=packetflix:ws://[::1"))).NotTo(BeZero())
	})

	It("fails capturing without container information, but picks up capture parameters", func() {
		cmd := newTestCmd("--url=packetflix:ws://127.0.0.1:1/capture?nif=eth0&filter=tcp&chaste")
		Expect(pf.Capture(cmd)).NotTo(BeZero())
		Expect(target.Nifs(cmd)).To(Equal("eth0"))
		Expect(target.CaptureFilter(cmd)).To(Equal("tcp"))
		Expect(target.NoPromiscuous(cmd)).To(BeTrue())
	})

})
