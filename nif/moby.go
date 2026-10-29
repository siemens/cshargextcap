// (c) Siemens AG 2023
//
// SPDX-License-Identifier: MIT

package nif

import (
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/url"
	"slices"

	"github.com/siemens/csharg/v2"
	"github.com/spf13/cobra"
	"github.com/thediveo/clippy/cliplugin"
	"github.com/thediveo/go-plugger/v3"
	"github.com/thediveo/nonstd/xslog"

	"github.com/siemens/cshargextcap"
	"github.com/siemens/cshargextcap/cfg"
	"github.com/siemens/cshargextcap/cli/insecure"
	"github.com/siemens/cshargextcap/cli/show"
	"github.com/siemens/cshargextcap/cli/wireshark"
	"github.com/siemens/cshargextcap/extcapargs"
)

// DockerHostURLFlagName is the name of the CLI flag specifying the Docker host
// URL to capture from.
const DockerHostURLFlagName = "dockerhosturl"

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
	rootCmd.PersistentFlags().String(DockerHostURLFlagName, "",
		"Docker host URL to capture from")
}

// DockerHostURL returns the URL of the Docker host to capture from.
func DockerHostURL(cmd *cobra.Command) string {
	u, _ := cmd.Flags().GetString(DockerHostURLFlagName)
	return u
}

// newMobySharkTank returns a csharg client for the Docker host specified on
// the command line.
func newMobySharkTank(cmd *cobra.Command) (csharg.SharkTank, error) {
	return csharg.NewSharkTankOnHost(DockerHostURL(cmd), &csharg.SharkTankOnHostOptions{
		InsecureSkipVerify: insecure.SkipVerify(cmd),
	})
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
func (moby *MobyNif) Configure(_ *cobra.Command, w io.Writer) int {
	// Merge our specific configuration arguments with the common ones.
	opts, err := extcapargs.Merge(map[int]extcapargs.ArgDef{
		cfg.DockerHostURL: extcapargs.StringArg("--dockerhosturl").IsRequired().
			Display("Docker host URL").Tooltip("DNS name or IP address of Docker host").
			Placeholder("http://host-ip-or-dns:5001"),
		cfg.SkipVerify: extcapargs.FlagArg("--insecure").
			Display("⚠INSECURE⚠ Skip server certificate validation"),
	}, cfg.CommonArgs)
	if err != nil {
		slog.Error("invalid configuration args", xslog.Error(err))
		return 1
	}

	// Set sensible defaults for the pod and standalone container filters.
	opts[cfg.ShowPods] = opts[cfg.ShowPods].Default(false)
	opts[cfg.ShowStandaloneContainers] = opts[cfg.ShowStandaloneContainers].Default(true)
	extcapargs.Dump(w, opts)

	// Set the common arg values.
	extcapargs.DumpValues(w, cfg.CommonArgValues)

	return 0
}

// ReloadOption refreshes values for a specific configuration argument.
func (moby *MobyNif) ReloadOption(cmd *cobra.Command, w io.Writer) {
	switch opt := wireshark.ReloadOption(cmd); opt {
	case "container":
		reloadMobyContainers(cmd, w)
	case "nif":
		cshargextcap.ReloadContainerNifs(cmd, w)
	default:
		slog.Error("unknown option reload", slog.String("option", opt))
	}
}

// reloadMobyContainers retrieves a fresh list of capture targets from the
// Docker host's capture service and then returns an updated (container) list in
// the form Wireshark expects from extcap plugins.
func reloadMobyContainers(cmd *cobra.Command, w io.Writer) {
	slog.Debug("reloading containers", slog.String("url", DockerHostURL(cmd)))
	st, err := newMobySharkTank(cmd)
	if err != nil {
		slog.Error("error creating SharkTank container host client", xslog.Error(err))
		return
	}
	slog.Debug("discovering targets...")
	targets, err := st.Targets()
	if err != nil {
		slog.Error("cannot discover targets", xslog.Error(err))
		return
	}
	argvals := cfg.ArgVals{}
	for _, t := range targets {
		slog.Debug("discovered target",
			slog.String("type", t.Type),
			slog.String("name", t.Name),
			slog.Any("netns", t.NetNS))
		display := t.Name
		if len(t.Prefix) > 0 {
			display = fmt.Sprintf("%s:%s", t.Prefix, display)
		}
		if show.Target(cmd, t) {
			value, _ := json.Marshal(t)
			argvals = append(argvals, cfg.ArgVal{
				Display:   display,
				Container: url.QueryEscape(string(value)),
			})
		}
	}
	slices.SortFunc(argvals, cfg.CompareArgVals)
	argvals.Dump(w, cfg.Containers)
}

// Capture captures network packets from a container, process, etc. with its
// own network namespace (virtual IP stack) on a Docker container host.
func (moby *MobyNif) Capture(cmd *cobra.Command) int {
	slog.Debug("capturing", slog.String("url", DockerHostURL(cmd)))
	st, err := newMobySharkTank(cmd)
	if err != nil {
		slog.Error("error creating SharkTank container host client", xslog.Error(err))
		return 1
	}
	return cshargextcap.Capture(cmd, st)
}
