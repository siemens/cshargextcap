// (c) Siemens AG 2026
//
// SPDX-License-Identifier: MIT

package target

import (
	"encoding/json"
	"log/slog"
	"net/url"

	"github.com/siemens/csharg/v2/api"
	"github.com/spf13/cobra"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	. "github.com/thediveo/success"
)

var _ = Describe("capture target", func() {

	var rootCmd *cobra.Command

	BeforeEach(func() {
		DeferCleanup(slog.SetDefault, slog.Default())
		slog.SetDefault(slog.New(slog.DiscardHandler))
		rootCmd = &cobra.Command{}
		TargetSetupCLI(rootCmd)
	})

	It("registers its CLI flags", func() {
		Expect(rootCmd.ParseFlags([]string{
			"--extcap-interface=foo", "--nif=eth0/eth1",
			"--extcap-capture-filter=tcp", "--noprom",
		})).To(Succeed())
		Expect(Nif(rootCmd)).To(Equal("foo"))
		Expect(Nifs(rootCmd)).To(Equal("eth0/eth1"))
		Expect(CaptureFilter(rootCmd)).To(Equal("tcp"))
		Expect(NoPromiscuous(rootCmd)).To(BeTrue())
	})

	DescribeTable("rejects missing or invalid --container",
		func(container string) {
			Expect(rootCmd.ParseFlags([]string{"--container=" + container})).To(Succeed())
			Expect(Unpack(rootCmd)).To(BeNil())
		},
		Entry("missing", ""),
		Entry("invalid escaping", "%zz"),
		Entry("invalid JSON", url.QueryEscape("{foo")),
	)

	It("unpacks a valid --container", func() {
		container := url.QueryEscape(string(Successful(json.Marshal(&api.Target{
			Name: "foo",
			Type: "docker.com",
		}))))
		Expect(rootCmd.ParseFlags([]string{"--container=" + container})).To(Succeed())
		t := Unpack(rootCmd)
		Expect(t).NotTo(BeNil())
		Expect(t.Name).To(Equal("foo"))
		Expect(t.Type).To(Equal("docker.com"))
	})

})
