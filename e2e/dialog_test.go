// (c) Siemens AG 2026
//
// SPDX-License-Identifier: MIT

package e2e

import (
	"context"
	"os"
	"os/exec"
	"time"

	"github.com/siemens/cshargextcap/cfg"

	"github.com/onsi/gomega/gexec"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	. "github.com/thediveo/success"
)

// startWireshark starts Wireshark on the desktop with only the cshargextcap
// plugin as its extcap plugin.
func (d *desktop) startWireshark() {
	GinkgoHelper()
	// List only extcap interfaces, so ours don't end up scrolled out of view.
	cmd := exec.Command("wireshark", "-o", "gui.interfaces_hidden_types:0,1,2,3,4,5,6,7,9")
	cmd.Env = append(os.Environ(), d.env()...)
	ws := Successful(gexec.Start(cmd, GinkgoWriter, GinkgoWriter))
	DeferCleanup(func() { ws.Kill().Wait(5 * time.Second) })
}

// openExtcapOptions opens the options dialog of the named extcap interface by
// clicking its gear in Wireshark's interface list, returning the dialog.
func (d *desktop) openExtcapOptions(nif string) accessible {
	GinkgoHelper()
	var item accessible
	Eventually(func() (ok bool) {
		item, ok = d.root().find("tree item", HaveSuffix(": "+nif))
		return ok
	}).Within(30 * time.Second).ProbeEvery(500 * time.Millisecond).Should(BeTrue())
	// The gear sits in the cell left of the interface name.
	d.clickOn(item.parent().child(item.indexInParent() - 1))

	var dialog accessible
	Eventually(func() (ok bool) {
		dialog, ok = d.root().find("dialog", ContainSubstring("Interface Options: "))
		return ok
	}).Within(10 * time.Second).ProbeEvery(250 * time.Millisecond).Should(BeTrue())
	Expect(dialog.name()).To(HaveSuffix(": " + nif))
	return dialog
}

// control is a labelling or input element as rendered by Wireshark.
type control struct {
	Role    string
	Name    string
	Checked bool
}

// shownControls returns the labels, check boxes and combo boxes currently
// showing in the subtree of a, in rendering order.
func shownControls(a accessible) []control {
	var controls []control
	if a.state(stateShowing) {
		switch role := a.role(); role {
		case "label", "check box", "combo box":
			controls = append(controls, control{Role: role, Name: a.name(), Checked: a.state(stateChecked)})
		}
	}
	for _, child := range a.children() {
		controls = append(controls, shownControls(child)...)
	}
	return controls
}

var _ = Describe("Wireshark's extcap options dialog", func() {

	It("renders the args of all tabs", func() {
		tabs := []struct {
			name     string
			controls []control
		}{
			{"Default", []control{
				{Role: "label", Name: "Docker host URL"},
				{Role: "check box", Name: "⚠INSECURE⚠ Skip server certificate validation"},
				{Role: "label", Name: "Containers"},
				{Role: "combo box"},
				{Role: "label", Name: "Interfaces"},
				{Role: "combo box", Name: "all"},
				{Role: "check box", Name: "Show pods"},
				{Role: "check box", Name: "Show standalone containers", Checked: true},
				{Role: "check box", Name: "Show processes"},
				{Role: "check box", Name: "Show process-less IP stacks"},
				{Role: "check box", Name: "No promiscuous mode"},
			}},
			{cfg.ProxyTabName, []control{
				{Role: "check box", Name: "Disable proxies", Checked: true},
				{Role: "label", Name: "HTTP proxy"},
				{Role: "label", Name: "HTTPS proxy"},
				{Role: "label", Name: "No proxy"},
			}},
			{cfg.AdvancedTabName, []control{
				{Role: "label", Name: "Service timeout (s)"},
			}},
		}

		d := newDesktop()
		d.startWireshark()
		dialog := d.openExtcapOptions("mobyshark")
		pages := dialog.mustFind("layered pane", BeEmpty())
		var tabNames []string
		for _, tab := range tabs {
			tabNames = append(tabNames, tab.name)
		}
		Expect(names(dialog.findAll("page tab"))).To(Equal(tabNames))

		for _, tab := range tabs {
			By("switching to the " + tab.name + " tab")
			pageTab := dialog.mustFind("page tab", Equal(tab.name))
			Expect(pageTab.press()).To(Succeed())
			Eventually(func() bool { return pageTab.state(stateSelected) }).Should(BeTrue())
			Eventually(shownControls).WithArguments(pages).
				Within(5*time.Second).ProbeEvery(250*time.Millisecond).
				Should(Equal(tab.controls), "controls on the %s tab", tab.name)
		}
	})

	It("lists the containers and the network interfaces of the selected one", func(ctx context.Context) {
		const name = "cshargextcap-e2e-gui"
		newTarget(ctx, name, pingScript)

		d := newDesktop()
		d.startWireshark()
		dialog := d.openExtcapOptions("mobyshark")
		Expect(dialog.mustFind("label", Equal("Docker host URL")).next().setText(serviceURL())).To(Succeed())
		containers := dialog.mustFind("label", Equal("Containers")).next()
		Expect(containers.mustFind("push button", Equal("refresh")).press()).To(Succeed())
		Eventually(func() []string { return names(containers.findAll("list item")) }).
			Within(30 * time.Second).ProbeEvery(500 * time.Millisecond).
			Should(ContainElement(name))

		By("selecting the container")
		// Qt's AT-SPI selection doesn't change a combo box's current item, so click it.
		combo := containers.mustFind("combo box", Not(BeNil()))
		Expect(combo.press()).To(Succeed())
		d.clickOn(containers.mustFind("list item", Equal(name)))
		Eventually(combo.name).Within(5 * time.Second).Should(Equal(name))

		By("refreshing the container's network interfaces")
		interfaces := dialog.mustFind("label", Equal("Interfaces")).next()
		Expect(interfaces.mustFind("push button", Equal("refresh")).press()).To(Succeed())
		Eventually(func() []string { return names(interfaces.findAll("list item")) }).
			Within(30 * time.Second).ProbeEvery(500 * time.Millisecond).
			Should(Equal([]string{"all", "eth0", "lo"}))
		Expect(interfaces.mustFind("combo box", Not(BeNil())).name()).To(Equal("all"))
	})

})
