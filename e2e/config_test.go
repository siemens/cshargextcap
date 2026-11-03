// (c) Siemens AG 2026
//
// SPDX-License-Identifier: MIT

package e2e

import (
	"context"
	"time"

	"github.com/onsi/gomega/gexec"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Wireshark configuring a capture", func() {

	DescribeTable("lists the containers to select from",
		func(ctx context.Context, name string, shown bool, args ...string) {
			newTarget(ctx, name, pingScript)
			display := ContainSubstring("{display=" + name + "}")
			if shown {
				Expect(string(reloadContainers(args...))).To(display)
			} else {
				Expect(string(reloadContainers(args...))).NotTo(display)
			}
		},
		Entry("when showing standalone containers", "cshargextcap-e2e-show", true, "--showcontainers"),
		Entry("when Wireshark passes a separate true", "cshargextcap-e2e-showtrue", true, "--showcontainers", "true"),
		Entry("not when hiding standalone containers", "cshargextcap-e2e-hide", false),
		Entry("not when Wireshark passes a separate false", "cshargextcap-e2e-showfalse", false, "--showcontainers", "false"),
	)

	It("lists the network interfaces of the selected container", func(ctx context.Context) {
		const name = "cshargextcap-e2e-nifs"
		newTarget(ctx, name, pingScript)
		cli := runExtcap(
			"--extcap-config", "--extcap-interface=mobyshark",
			"--extcap-reload-option=nif", "--dockerhosturl="+serviceURL(),
			"--container="+containerArg(name))
		Eventually(cli).WithTimeout(30 * time.Second).Should(gexec.Exit(0))
		Expect(string(cli.Out.Contents())).To(And(
			ContainSubstring("{value=any}{display=all}{default=true}"),
			ContainSubstring("{value=eth0}"),
			ContainSubstring("{value=lo}")))
	})

	It("lists no containers when the capture service is unreachable", func() {
		cli := runExtcap(
			"--extcap-config", "--extcap-interface=mobyshark",
			"--extcap-reload-option=container", "--dockerhosturl=http://"+nodeName+":1",
			"--showcontainers")
		Eventually(cli).WithTimeout(30 * time.Second).Should(gexec.Exit(0))
		Expect(cli.Out.Contents()).To(BeEmpty())
	})

})
