// (c) Siemens AG 2023
//
// SPDX-License-Identifier: MIT

package nif

import (
	"encoding/json"
	"fmt"
	"io"
	"net/url"
	"time"

	"github.com/siemens/csharg"
	"github.com/siemens/cshargextcap"
	"github.com/siemens/cshargextcap/cfg"
	"github.com/siemens/cshargextcap/cli/cliplugin"
	"github.com/siemens/cshargextcap/cli/show"
	"github.com/siemens/cshargextcap/cli/timeout"
	"github.com/siemens/cshargextcap/cli/wireshark"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/thediveo/go-plugger/v3"
	"golang.org/x/exp/maps"
)

var DockerHostURL string
var Insecure bool

// MobyNif represents an external capture network interface connecting to the
// packetflix service of a Docker (“Moby”) container host.
type MobyNif struct {
	cshargextcap.ExtcapNif
}

// Automatically register our Docker/Moby host capture network interface.
func init() {
	plugger.Group[cshargextcap.ExtcapNifActions]().Register(&MobyNif{
		ExtcapNif: cshargextcap.NewExtcapNif("mobyshark", "Docker host capture")})
	plugger.Group[cliplugin.SetupCLI]().Register(
		MobySetupCLI, plugger.WithPlugin("mobyshark"))
}

func MobySetupCLI(rootCmd *cobra.Command) {
	pf := rootCmd.PersistentFlags()
	pf.StringVar(&DockerHostURL, "dockerhosturl", "",
		"Docker host URL to capture from")
	pf.BoolVar(&Insecure, "insecure", false,
		"skip server certificate verification")
}

// DLTs returns the data-link layer encapsulation type of this external
// capture network interface. Since we cannot know better, it's the usual
// suspect: first USER type.
func (moby *MobyNif) DLTs() []cshargextcap.ExtcapNifDLT {
	return []cshargextcap.ExtcapNifDLT{
		{Number: 147, Display: "remote container capture"},
	}
}

// Configure dumps the configuration arguments and their values, from
// which Wireshark then generates the configuration UI.
func (moby *MobyNif) Configure(w io.Writer) int {
	// Merge our specific configuration arguments with the common ones.
	opts := map[int]string{
		// Standard tab
		cfg.DockerHostURL: "{call=--dockerhosturl}{type=string}{required=true}" +
			"{display=Docker host URL}{tooltip=DNS name or IP address of Docker host}" +
			"{placeholder=http://host-ip-or-dns:5001}",
		cfg.SkipVerify: "{call=--insecure}{type=boolflag}" +
			"{display=⚠INSECURE⚠ Skip server certificate validation}",
	}
	maps.Copy(opts, cfg.CommonArgs)

	// Set sensible defaults for the pod and standalone container filters.
	opts[cfg.ShowPods] += "{default=false}"
	opts[cfg.ShowStandaloneContainers] += "{default=true}"
	cfg.DumpConfigArgs(w, opts)

	// Set the common arg values and other values.
	cfg.DumpConfigArgValues(w, cfg.CommonArgValues)
	cfg.DumpConfigArgValues(w, cfg.CommonOtherArgValues)

	return 0
}

// ReloadOption refreshes values for a specific configuration argument.
func (moby *MobyNif) ReloadOption(w io.Writer) {
	switch wireshark.ReloadOption {
	case "container":
		reloadMobyContainers(w)
	case "nif":
		cshargextcap.ReloadContainerNifs(w)
	default:
		log.Errorf("unknown option \"%s\" reload", wireshark.ReloadOption)
	}
}

// reloadMobyContainers retrieves a fresh list of capture targets from the
// Docker host's capture service and then returns an updated (container) list in
// the form Wireshark expects from extcap plugins.
func reloadMobyContainers(w io.Writer) {
	log.Debugf("Docker host URL: %s", DockerHostURL)
	st, err := csharg.NewSharkTankOnHost(DockerHostURL, &csharg.SharkTankOnHostOptions{
		InsecureSkipVerify: Insecure,
	})
	if err != nil {
		log.Errorf("error creating SharkTank container host client: %s", err.Error())
		return
	}
	log.Debug("discovering targets...")
	argvals := cfg.ArgVals{}
	for _, t := range st.Targets() {
		log.Debugf("%s: \"%s\" netns:[%d] ",
			t.Type, t.Name, t.NetNS)
		display := t.Name
		if len(t.Prefix) > 0 {
			display = fmt.Sprintf("%s:%s", t.Prefix, display)
		}
		if show.Target(t) {
			value, _ := json.Marshal(t)
			argvals = append(argvals, cfg.ArgVal{
				Display:   display,
				Container: url.QueryEscape(string(value)),
			})
		}
	}
	argvals.Sort()
	argvals.Dump(w, cfg.Containers)
}

// Capture captures network packets from a container, process, etc. with its
// own network namespace (virtual IP stack) on a Docker container host.
func (moby *MobyNif) Capture() int {
	log.Debugf("Docker host URL: %s", DockerHostURL)
	st, err := csharg.NewSharkTankOnHost(DockerHostURL, &csharg.SharkTankOnHostOptions{
		CommonClientOptions: csharg.CommonClientOptions{
			Timeout: time.Duration(timeout.Discovery) * time.Second,
		},
		InsecureSkipVerify: Insecure,
	})
	if err != nil {
		log.Errorf("error creating SharkTank container host client: %s", err.Error())
		return 1
	}
	return cshargextcap.Capture(st)
}
