// (c) Siemens AG 2026
//
// SPDX-License-Identifier: MIT

package e2e

import (
	"context"
	"net/url"
	"os"
	"os/exec"
	"strings"
	"syscall"
	"time"

	"github.com/onsi/gomega/gexec"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	. "github.com/thediveo/success"
)

var _ = Describe("Wireshark capturing", func() {

	DescribeTable("writes a complete capture when stopped by a signal",
		func(ctx context.Context, name string, sig os.Signal) {
			newTarget(ctx, name, pingScript)
			out := newCaptureFile()
			cli := runCapture(name, out)
			awaitPackets(ctx, out)
			cli.Signal(sig)
			Eventually(cli).WithTimeout(10 * time.Second).Should(gexec.Exit(0))
			Expect(capinfos(out)).To(And(
				ContainSubstring("container-name: "+name),
				ContainSubstring("node-name: "+nodeName)))
			Expect(fields(out, "icmp", "frame.number")).NotTo(BeEmpty())
		},
		Entry("SIGINT", "cshargextcap-e2e-sigint", os.Interrupt),
		Entry("SIGTERM", "cshargextcap-e2e-sigterm", syscall.SIGTERM),
	)

	DescribeTable("captures from the selected network interfaces",
		func(ctx context.Context, name, nif string, expected ...string) {
			newTarget(ctx, name, chattyScript)
			out := newCaptureFile()
			var args []string
			if nif != "" {
				args = append(args, "--nif="+nif)
			}
			cli := runCapture(name, out, args...)
			awaitPackets(ctx, out)
			cli.Interrupt()
			Eventually(cli).WithTimeout(10 * time.Second).Should(gexec.Exit(0))
			Expect(fields(out, "", "frame.interface_name")).To(And(
				HaveEach(BeElementOf(expected)),
				ContainElements(expected)))
		},
		Entry("all by default", "cshargextcap-e2e-nif-default", "", "lo", "eth0"),
		Entry("all when Wireshark passes its default", "cshargextcap-e2e-nif-any", "any", "lo", "eth0"),
		Entry("only loopback", "cshargextcap-e2e-nif-lo", "lo", "lo"),
		Entry("only ethernet", "cshargextcap-e2e-nif-eth0", "eth0", "eth0"),
	)

	DescribeTable("applies a capture filter",
		func(ctx context.Context, name, filter string, wantUDP bool) {
			newTarget(ctx, name, chattyScript)
			out := newCaptureFile()
			// Packetflix's dumpcap applies a filter only to the last network interface.
			args := []string{"--nif=lo"}
			if filter != "" {
				args = append(args, "--extcap-capture-filter="+filter)
			}
			cli := runCapture(name, out, args...)
			awaitPackets(ctx, out)
			cli.Interrupt()
			Eventually(cli).WithTimeout(10 * time.Second).Should(gexec.Exit(0))
			Expect(fields(out, "icmp", "frame.number")).NotTo(BeEmpty())
			// ICMP port unreachable errors quote UDP, so they don't count as UDP.
			udp := fields(out, "udp && !icmp", "frame.number")
			if wantUDP {
				Expect(udp).NotTo(BeEmpty())
			} else {
				Expect(udp).To(BeEmpty())
			}
		},
		Entry("without a filter", "cshargextcap-e2e-nofilter", "", true),
		Entry("with an ICMP-only filter", "cshargextcap-e2e-icmpfilter", "icmp", false),
	)

	It("captures via a packetflix: URL with only partial target information", func(ctx context.Context) {
		const name = "cshargextcap-e2e-pfurl"
		newTarget(ctx, name, pingScript)
		out := newCaptureFile()
		cli := runExtcap("--capture", "--extcap-interface=packetflix", "--fifo="+out,
			"--url=packetflix:ws://"+serviceAddr+"/capture?nif=lo&container="+
				url.QueryEscape(`{"name":"`+name+`","type":"docker"}`))
		awaitPackets(ctx, out)
		cli.Interrupt()
		Eventually(cli).WithTimeout(10 * time.Second).Should(gexec.Exit(0))
		Expect(capinfos(out)).To(And(
			ContainSubstring("container-name: "+name),
			ContainSubstring("node-name: "+nodeName)))
		Expect(fields(out, "", "frame.interface_name")).To(HaveEach("lo"))
	})

	It("ends by itself when the target container terminates", func(ctx context.Context) {
		const name = "cshargextcap-e2e-shortlived"
		newTarget(ctx, name, "sleep 5; ping -c 3 -i 0.2 127.0.0.1")
		out := newCaptureFile()
		cli := runCapture(name, out)
		Eventually(cli).WithTimeout(30 * time.Second).Should(gexec.Exit(0))
		Expect(fields(out, "icmp", "frame.number")).NotTo(BeEmpty())
	})

	It("streams live into Wireshark's fifo", func(ctx context.Context) {
		const name = "cshargextcap-e2e-fifo"
		newTarget(ctx, name, pingScript)
		fifo := newFifo()
		tshark := Successful(gexec.Start(
			exec.Command("tshark", "-l", "-i", fifo, "-T", "fields", "-e", "frame.number"),
			GinkgoWriter, GinkgoWriter))
		DeferCleanup(func() { tshark.Kill().Wait() })
		cli := runCapture(name, fifo)
		Eventually(func() int { return len(strings.Fields(string(tshark.Out.Contents()))) }).
			WithTimeout(30 * time.Second).Should(BeNumerically(">=", 5))
		// Like Wireshark, first stop the plugin, then wait for the fifo to close.
		cli.Interrupt()
		Eventually(cli).WithTimeout(10 * time.Second).Should(gexec.Exit(0))
		Eventually(tshark).WithTimeout(10 * time.Second).Should(gexec.Exit(0))
		Expect(strings.Fields(string(tshark.Out.Contents()))[:5]).To(Equal([]string{"1", "2", "3", "4", "5"}))
	})

	It("reports an unknown capture target", func() {
		cli := runExtcap("--capture", "--extcap-interface=packetflix", "--fifo="+newCaptureFile(),
			"--url=packetflix:ws://"+serviceAddr+"/capture?container="+
				url.QueryEscape(`{"name":"cshargextcap-e2e-nonexisting","type":"docker"}`))
		Eventually(cli).WithTimeout(30 * time.Second).Should(gexec.Exit(1))
	})

})
