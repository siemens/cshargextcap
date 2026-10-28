// (c) Siemens AG 2026
//
// SPDX-License-Identifier: MIT

package cshargextcap

import (
	"bytes"
	"encoding/json"
	"io"
	"log/slog"
	"net/url"

	"github.com/siemens/csharg/v2/api"
	"github.com/spf13/cobra"
	"github.com/thediveo/go-plugger/v3"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	. "github.com/thediveo/success"
)

type reloadingNif struct {
	TestNif
	reloaded bool
}

func (r *reloadingNif) ReloadOption(cmd *cobra.Command, w io.Writer) { r.reloaded = true }

var _ = Describe("reloading", func() {

	BeforeEach(func() {
		DeferCleanup(slog.SetDefault, slog.Default())
		slog.SetDefault(slog.New(slog.DiscardHandler))
	})

	It("reloads container nifs in sorted order", func() {
		container := url.QueryEscape(string(Successful(json.Marshal(&api.Target{
			Name:              "foo",
			Type:              "docker.com",
			NetworkInterfaces: []string{"eth1", "eth0"},
		}))))
		var buff bytes.Buffer
		ReloadContainerNifs(newTestCmd("--container="+container), &buff)
		Expect(buff.String()).To(Equal(
			"value {arg=20}{value=any}{display=all}{default=true}\n" +
				"value {arg=20}{value=eth0}{display=eth0}{default=false}\n" +
				"value {arg=20}{value=eth1}{display=eth1}{default=false}\n"))
	})

	It("reloads only the catch-all nif without container information", func() {
		var buff bytes.Buffer
		ReloadContainerNifs(newTestCmd(), &buff)
		Expect(buff.String()).To(Equal(
			"value {arg=20}{value=any}{display=all}{default=true}\n"))
	})

	Context("routing reload options", func() {

		var nif *reloadingNif

		BeforeEach(func() {
			extcapNifsSaved := extcapNifs
			DeferCleanup(func() { extcapNifs = extcapNifsSaved })
			nifs := plugger.Group[ExtcapNifActions]()
			DeferCleanup(nifs.Restore, nifs.Backup())
			extcapNifs = nil
			nif = &reloadingNif{TestNif: TestNif{ExtcapNif{name: "test-reload"}}}
			plugger.Group[ExtcapNifActions]().Register(nif)
		})

		It("routes a reload to the specified extcap nif", func() {
			cmd := newTestCmd("--extcap-interface=test-reload", "--extcap-reload-option=foo")
			Expect(ExtcapConfigure(cmd, &bytes.Buffer{})).To(BeZero())
			Expect(nif.reloaded).To(BeTrue())
		})

		It("rejects a reload for an unknown extcap nif", func() {
			cmd := newTestCmd("--extcap-interface=foobar", "--extcap-reload-option=foo")
			Expect(ExtcapConfigure(cmd, &bytes.Buffer{})).NotTo(BeZero())
			Expect(nif.reloaded).To(BeFalse())
		})

	})

})
