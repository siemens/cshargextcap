// (c) Siemens AG 2023
//
// SPDX-License-Identifier: MIT

package nif

import (
	"encoding/json"
	"io"
	"net/url"
	"strings"

	"github.com/siemens/csharg"
	"github.com/siemens/csharg/api"
	"github.com/siemens/cshargextcap"
	"github.com/siemens/cshargextcap/cfg"
	"github.com/siemens/cshargextcap/cli/cliplugin"
	"github.com/siemens/cshargextcap/cli/target"
	"github.com/siemens/cshargextcap/cli/wireshark"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/thediveo/go-plugger/v3"
	"golang.org/x/exp/maps"
)

// packetflixURL specifies the URL of a Packetflix capture service.
var packetflixURL string

// PacketflixNif represents an external capture network interface to a (remote)
// capture service.
type PacketflixNif struct {
	cshargextcap.ExtcapNif
}

// Automatically registers our Packetflix capture network interface.
func init() {
	plugger.Group[cshargextcap.ExtcapNifActions]().Register(&PacketflixNif{
		ExtcapNif: cshargextcap.NewExtcapNif(
			"packetflix", "packetflix:// remote cluster and container host capture")})
	plugger.Group[cliplugin.SetupCLI]().Register(
		PacketflixSetupCLI, plugger.WithPlugin("packetflix"))
}

func PacketflixSetupCLI(rootCmd *cobra.Command) {
	pf := rootCmd.PersistentFlags()
	pf.StringVar(&packetflixURL, "url", "",
		"Packetflix capture service URL")
}

// DLTs returns the data-link layer encapsulation type of this external
// capture network interface. Since we cannot know better, it's the usual
// suspect: first USER type.
func (pf *PacketflixNif) DLTs() []cshargextcap.ExtcapNifDLT {
	return []cshargextcap.ExtcapNifDLT{
		{Number: 147, Display: "remote container capture"},
	}
}

// Configure dumps the configuration arguments and their values, from
// which Wireshark then generates the configuration UI.
func (pf *PacketflixNif) Configure(w io.Writer) int {
	// Merge our specific configuration arguments with the common ones.
	opts := map[int]string{
		cfg.PacketflixURL: "{call=--url}{type=string}{required=true}" +
			"{display=capture URL}{tooltip=remote Packetflix capture service URL}" +
			"{placeholder=packetflix://...}",
		cfg.SkipVerify: "{call=--insecure}{type=boolflag}" +
			"{display=⚠INSECURE⚠ Skip server certificate validation}",
	}
	maps.Copy(opts, cfg.CommonOtherTabArgs)
	cfg.DumpConfigArgs(w, opts)

	// None here: Set the common arg values.
	// dumpConfigArgValues(w, commonArgValues)

	return 0
}

// ReloadOption refreshes values for a specific configuration argument.
func (pf *PacketflixNif) ReloadOption(w io.Writer) {
	switch wireshark.ReloadOption {
	default:
		log.Errorf("unknown option \"%s\" reload", wireshark.ReloadOption)
	}
}

// Capture captures network packets from a container, process, etc. with its
// own network namespace (virtual IP stack) on a Docker container host.
func (pf *PacketflixNif) Capture() int {
	// Turn the packetflix: URL into something the csharg Docker host
	// client can work with.
	capturl := packetflixURL
	log.Debugf("remote capture url: \"%s\"", capturl)
	// Differing from the other ClusterShark extcap nifs, we won't have any
	// sensible container reference information in the CLI args. Instead, this
	// information is part of the URL given to us. Because the csharg client
	// uses http: and https:-based Docker host URLs, we may need to map back
	// from websockets to HTTP(S).
	capturl = strings.TrimPrefix(capturl, "packetflix:")
	if !strings.HasPrefix(capturl, "ws:") && !strings.HasPrefix(capturl, "wss:") {
		capturl = "wss:" + capturl
	}
	u, err := url.Parse(capturl)
	if err != nil {
		log.Errorf("invalid packetflix URL \"%s\": %s", packetflixURL, err.Error())
		return 1
	}
	switch u.Scheme {
	case "ws":
		u.Scheme = "http"
	case "wss":
		u.Scheme = "https"
	}
	// Albeit this is slightly whacky, we simply cut of the trailing "capture"
	// path element, as csharg's container host client will tack it on again.
	u.Path = strings.TrimSuffix(u.Path, "/capture")
	if u.Path == "" {
		u.Path = "/"
	}
	// There are yet no useful capture target-related args set, as these don't
	// come in via the CLI args, but instead as part of the URL. So we need to
	// recover them in order to then make the usual csharg/extcap capture
	// mechanism work as expected.
	opts := u.Query()
	target.Container = opts.Get("container")
	target.Nifs = opts.Get("nif")
	target.CaptureFilter = opts.Get("filter")
	if _, ok := opts["chaste"]; ok {
		target.NoPromiscuous = true
	}
	// Remove anything which will trouble csharg.
	u.RawQuery = ""
	u.Fragment = ""
	// From here on, this now looks like an ordinary container host capture.
	log.Debugf("Docker host URL: %s", u.String())
	st, err := csharg.NewSharkTankOnHost(u.String(), &csharg.SharkTankOnHostOptions{
		InsecureSkipVerify: Insecure,
	})
	if err != nil {
		log.Errorf("error creating SharkTank container host client: %s", err.Error())
		return 1
	}
	// Since csharg is rather picky about seeing correct target information, we
	// need to add in the node name data, which we glance from the list of
	// targets.
	captt := target.Unpack()
	if captt == nil {
		log.Error("--container information missing or invalid")
		return 1
	}
	matches := api.Targets{}
	for _, t := range st.Targets() {
		if captt.Name == t.Name &&
			captt.Type == t.Type &&
			(captt.NodeName == "" || captt.NodeName == t.NodeName) {
			matches = append(matches, t)
		}
	}
	// Only update capture target information if exactly one match was found,
	// otherwise don't touch the capture target information, as it must be
	// sufficiently correctly supplied by the caller.
	if len(matches) == 1 {
		captt = *(&matches[0])
		if targetbytes, err := json.Marshal(captt); err == nil {
			target.Container = url.QueryEscape(string(targetbytes))
			log.Debugf("updating target information: %+v", captt)
		}
	}
	// Eh, get us the packets, finally!
	return cshargextcap.Capture(st)
}
