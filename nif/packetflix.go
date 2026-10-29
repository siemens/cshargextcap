// (c) Siemens AG 2023
//
// SPDX-License-Identifier: MIT

package nif

import (
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/url"
	"strings"

	"github.com/siemens/csharg/v2"
	"github.com/siemens/csharg/v2/api"
	"github.com/spf13/cobra"
	"github.com/thediveo/clippy/cliplugin"
	"github.com/thediveo/go-plugger/v3"
	"github.com/thediveo/nonstd/xslog"

	"github.com/siemens/cshargextcap"
	"github.com/siemens/cshargextcap/cfg"
	"github.com/siemens/cshargextcap/cli/insecure"
	"github.com/siemens/cshargextcap/cli/target"
	"github.com/siemens/cshargextcap/cli/wireshark"
	"github.com/siemens/cshargextcap/extcapargs"
)

// PacketflixURLFlagName is the name of the CLI flag specifying the URL of a
// Packetflix capture service.
const PacketflixURLFlagName = "url"

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
	rootCmd.PersistentFlags().String(PacketflixURLFlagName, "",
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
func (pf *PacketflixNif) Configure(_ *cobra.Command, w io.Writer) int {
	// Merge our specific configuration arguments with the common ones.
	opts, err := extcapargs.Merge(map[int]extcapargs.ArgDef{
		cfg.PacketflixURL: extcapargs.StringArg("--url").IsRequired().
			Display("capture URL").Tooltip("remote Packetflix capture service URL").
			Placeholder("packetflix://..."),
		cfg.SkipVerify: extcapargs.FlagArg("--insecure").
			Display("⚠INSECURE⚠ Skip server certificate validation"),
	}, cfg.CommonOtherTabArgs)
	if err != nil {
		slog.Error("invalid configuration args", xslog.Error(err))
		return 1
	}
	extcapargs.Dump(w, opts)

	// None here: Set the common arg values.
	// dumpConfigArgValues(w, commonArgValues)

	return 0
}

// ReloadOption refreshes values for a specific configuration argument.
func (pf *PacketflixNif) ReloadOption(cmd *cobra.Command, w io.Writer) {
	slog.Error("unknown option reload", slog.String("option", wireshark.ReloadOption(cmd)))
}

// Capture captures network packets from a container, process, etc. with its
// own network namespace (virtual IP stack) on a Docker container host.
func (pf *PacketflixNif) Capture(cmd *cobra.Command) int {
	packetflixURL, _ := cmd.Flags().GetString(PacketflixURLFlagName)
	slog.Debug("remote capture", slog.String("url", packetflixURL))
	u, opts, err := hostURL(packetflixURL)
	if err != nil {
		slog.Error("invalid packetflix URL", slog.String("url", packetflixURL), xslog.Error(err))
		return 1
	}
	// There are yet no useful capture target-related args set, as these don't
	// come in via the CLI args, but instead as part of the URL. So we need to
	// recover them in order to then make the usual csharg/extcap capture
	// mechanism work as expected.
	flags := cmd.Flags()
	err = errors.Join(
		flags.Set(target.ContainerFlagName, opts.Get("container")),
		flags.Set(target.NifsFlagName, opts.Get("nif")),
		flags.Set(target.CaptureFilterFlagName, opts.Get("filter")))
	if _, ok := opts["chaste"]; ok {
		err = errors.Join(err, flags.Set(target.NoPromFlagName, "true"))
	}
	if err != nil {
		slog.Error("cannot set capture parameters", xslog.Error(err))
		return 1
	}
	// From here on, this now looks like an ordinary container host capture.
	slog.Debug("capturing", slog.String("url", u.String()))
	st, err := csharg.NewSharkTankOnHost(u.String(), &csharg.SharkTankOnHostOptions{
		InsecureSkipVerify: insecure.SkipVerify(cmd),
	})
	if err != nil {
		slog.Error("error creating SharkTank container host client", xslog.Error(err))
		return 1
	}
	// Since csharg is rather picky about seeing correct target information, we
	// need to add in the node name data, which we glance from the list of
	// targets.
	captt := target.Unpack(cmd)
	if captt == nil {
		slog.Error("--container information missing or invalid")
		return 1
	}
	targets, err := st.Targets()
	if err != nil {
		slog.Warn("cannot discover targets", xslog.Error(err))
	}
	matches := api.Targets{}
	for _, t := range targets {
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
		captt = matches[0]
		if targetbytes, err := json.Marshal(captt); err == nil &&
			flags.Set(target.ContainerFlagName, url.QueryEscape(string(targetbytes))) == nil {
			slog.Debug("updating target information", slog.Any("target", captt))
		}
	}
	// Eh, get us the packets, finally!
	return cshargextcap.Capture(cmd, st)
}

// hostURL converts a packetflix: URL into the http(s) base URL of the capture
// service as expected by csharg's host client, and returns the capture
// parameters passed in the URL's query separately.
func hostURL(packetflixURL string) (*url.URL, url.Values, error) {
	capturl := strings.TrimPrefix(packetflixURL, "packetflix:")
	if !strings.HasPrefix(capturl, "ws:") && !strings.HasPrefix(capturl, "wss:") {
		capturl = "wss:" + capturl
	}
	u, err := url.Parse(capturl)
	if err != nil {
		return nil, nil, err
	}
	switch u.Scheme {
	case "ws":
		u.Scheme = "http"
	case "wss":
		u.Scheme = "https"
	}
	// csharg's host client tacks the "/capture" path element on again.
	u.Path = strings.TrimSuffix(u.Path, "/capture")
	if u.Path == "" {
		u.Path = "/"
	}
	opts := u.Query()
	u.RawQuery = ""
	u.Fragment = ""
	return u, opts, nil
}
