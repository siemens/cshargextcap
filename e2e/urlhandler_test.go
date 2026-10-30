// (c) Siemens AG 2026
//
// SPDX-License-Identifier: MIT

package e2e

import (
	"context"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/onsi/gomega/gexec"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	. "github.com/thediveo/success"
)

const urlHandlerDesktopFile = "../packaging/linux/com.github.siemens.packetflix.desktop"

// installURLHandler installs the packaged packetflix: URL handler for the
// desktop's user, the same way the Linux packages do system-wide, returning the
// handler's name.
func (d *desktop) installURLHandler() string {
	GinkgoHelper()
	entry := Successful(os.ReadFile(urlHandlerDesktopFile))
	apps := filepath.Join(d.home, ".local", "share", "applications")
	Expect(os.MkdirAll(apps, 0o755)).To(Succeed())
	Expect(os.WriteFile(filepath.Join(apps, filepath.Base(urlHandlerDesktopFile)), entry, 0o644)).
		To(Succeed())
	Expect(exec.Command("update-desktop-database", apps).Run()).To(Succeed())
	name := regexp.MustCompile(`(?m)^Name=(.+)$`).FindSubmatch(entry)
	Expect(name).NotTo(BeNil(), "URL handler without name")
	return string(name[1])
}

// firefoxPrefs skip Firefox's first-run onboarding and start pages, which would
// otherwise cover or hide the page under test.
const firefoxPrefs = `
user_pref("browser.shell.checkDefaultBrowser", false);
user_pref("browser.aboutwelcome.enabled", false);
user_pref("browser.preonboarding.enabled", false);
user_pref("termsofuse.bypassNotification", true);
user_pref("datareporting.policy.dataSubmissionEnabled", false);
user_pref("datareporting.policy.dataSubmissionPolicyBypassNotification", true);
user_pref("toolkit.telemetry.reportingpolicy.firstRun", false);
user_pref("trailhead.firstrun.didSeeAboutWelcome", true);
user_pref("startup.homepage_welcome_url", "");
user_pref("startup.homepage_welcome_url.additional", "");
user_pref("browser.startup.homepage_override.mstone", "ignore");
user_pref("browser.startup.firstrunSkipsHomepage", true);
user_pref("browser.startup.page", 0);
user_pref("browser.startup.homepage", "about:blank");
user_pref("browser.newtabpage.enabled", false);
`

// startFirefox starts Firefox with a fresh profile on the desktop, showing the
// specified URL.
func (d *desktop) startFirefox(url string) {
	GinkgoHelper()
	profile := tempDir()
	Expect(os.WriteFile(filepath.Join(profile, "user.js"), []byte(firefoxPrefs), 0o600)).To(Succeed())
	cmd := exec.Command("firefox", "--no-remote", "--new-instance", "--profile", profile, url)
	cmd.Env = append(append(os.Environ(), d.env()...), "MOZ_ENABLE_WAYLAND=0")
	// Handlers launched by Firefox inherit its output, so leave killing them all
	// to the desktop's cleanup.
	_ = Successful(gexec.Start(cmd, GinkgoWriter, GinkgoWriter))
}

var _ = Describe("Edgeshark capture links", func() {

	It("start capturing in Wireshark", func(ctx context.Context) {
		const name = "cshargextcap-e2e-url"
		newTarget(ctx, name, pingScript)
		d := newDesktop()
		handler := d.installURLHandler()

		By("starting a capture from the container's details view in Edgeshark")
		d.startFirefox(serviceURL() + "/n/" + name)
		var refresh accessible
		Eventually(func() (ok bool) {
			refresh, ok = d.root().find("button", Equal("refresh"))
			return ok
		}).Within(30 * time.Second).ProbeEvery(500 * time.Millisecond).Should(BeTrue())
		// The details view only knows the container after a discovery.
		d.clickOn(refresh)
		var link accessible
		Eventually(func() string {
			for _, l := range d.root().findAll("link") {
				if uri := l.uri(); strings.HasPrefix(uri, "packetflix:") {
					link = l
					return uri
				}
			}
			return ""
		}).Within(30 * time.Second).ProbeEvery(500 * time.Millisecond).
			Should(ContainSubstring(url.QueryEscape(`"name":"` + name + `"`)))
		d.clickOn(link)

		By("confirming to the browser to open the capture link")
		Eventually(func() (ok bool) {
			_, ok = d.root().find("label", Equal("Allow this site to open the packetflix link with "+handler+"?"))
			return ok
		}).Within(30 * time.Second).ProbeEvery(500 * time.Millisecond).Should(BeTrue())
		open := d.root().mustFind("button", Equal("Open Link"))
		// Firefox enables dialog buttons only after a delay against clickjacking.
		Eventually(func() bool { return open.state(stateEnabled) }).Within(10 * time.Second).Should(BeTrue())
		Expect(open.press()).To(Succeed())

		By("capturing in Wireshark")
		var capturing accessible
		Eventually(func() (ok bool) {
			wireshark, ok := d.app("Wireshark")
			if ok {
				capturing, ok = wireshark.find("frame", HavePrefix("Capturing from packetflix"))
			}
			return ok
		}).Within(60 * time.Second).ProbeEvery(time.Second).Should(BeTrue())
		packets := capturing.mustFind("tree", Equal("Packet list"))
		// Probe only the first packets, as incoming packets invalidate later cells.
		Eventually(func() []string {
			var texts []string
			for row := range int32(10) {
				for column := range packets.columns() {
					texts = append(texts, packets.cell(row, column).name())
				}
			}
			return texts
		}).Within(30 * time.Second).ProbeEvery(time.Second).
			Should(ContainElements("ICMP", HavePrefix("Echo (ping) request")))
	}, SpecTimeout(3*time.Minute))

})
