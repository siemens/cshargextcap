// (c) Siemens AG 2026
//
// SPDX-License-Identifier: MIT

package main

import (
	"net/url"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	. "github.com/thediveo/success"
)

var _ = Describe("packetflix URL handler", func() {

	DescribeTable("deriving window titles",
		func(container, expected string) {
			u := Successful(url.Parse("packetflix:ws://host:5001/capture?container=" + url.QueryEscape(container)))
			Expect(windowTitle(u)).To(Equal(expected))
		},
		Entry("name and nifs", `{"name":"foo","type":"docker","network-interfaces":["eth0","lo"]}`, "foo: eth0, lo"),
		Entry("name only", `{"name":"foo"}`, "foo"),
		Entry("no name", `{"network-interfaces":["eth0"]}`, ""),
		Entry("no JSON", `foo`, ""),
		Entry("missing", ``, ""),
	)

	It("starts a new Wireshark instance located by its bundle ID", func() {
		pfURL := "packetflix:ws://host:5001/capture?container=" +
			url.QueryEscape(`{"name":"foo","network-interfaces":["eth0"]}`) + "&nif=eth0"
		Expect(openArgs(pfURL)).To(Equal([]string{
			"-n", "-b", "org.wireshark.Wireshark", "--args",
			"-k", "-i", "packetflix",
			"-o", "gui.window_title:foo: eth0",
			"-o", "extcap.packetflix.url:" + pfURL,
		}))
	})

	It("passes shell metacharacters verbatim as a single argument", func() {
		pfURL := `packetflix:ws://host/$(touch$IFS/tmp/pwned)"';|&`
		args := Successful(openArgs(pfURL))
		Expect(args[len(args)-1]).To(Equal("extcap.packetflix.url:" + pfURL))
	})

	DescribeTable("rejecting non-packetflix URLs",
		func(pfURL string) {
			Expect(openArgs(pfURL)).Error().To(HaveOccurred())
		},
		Entry("other scheme", "https://example.org"),
		Entry("no scheme", "//host:5001/capture"),
		Entry("invalid", "packetflix:ws://host:5001/\x7f"),
	)

})
