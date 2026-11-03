// (c) Siemens AG 2026
//
// SPDX-License-Identifier: MIT

package main

import (
	"encoding/json"
	"errors"
	"net/url"
	"strings"
)

// wiresharkBundleID identifies Wireshark independent of where it has been
// installed to, such as /Applications/ or /Applications/MacPorts/.
const wiresharkBundleID = "org.wireshark.Wireshark"

// openArgs returns the arguments to macOS' "open" command for starting a new
// Wireshark instance that captures from the specified packetflix URL.
func openArgs(pfURL string) ([]string, error) {
	u, err := url.Parse(pfURL)
	if err != nil {
		return nil, err
	}
	if u.Scheme != "packetflix" {
		return nil, errors.New("not a packetflix: URL")
	}
	args := []string{"-n", "-b", wiresharkBundleID, "--args", "-k", "-i", "packetflix"}
	if title := windowTitle(u); title != "" {
		args = append(args, "-o", "gui.window_title:"+title)
	}
	return append(args, "-o", "extcap.packetflix.url:"+pfURL), nil
}

// windowTitle returns a Wireshark window title in the form of "name: nif1,
// nif2, ..." derived from the capture target details in the packetflix URL's
// "container" query parameter, or "" if there are no usable details.
func windowTitle(u *url.URL) string {
	var target struct {
		Name string   `json:"name"`
		Nifs []string `json:"network-interfaces"`
	}
	if json.Unmarshal([]byte(u.Query().Get("container")), &target) != nil || target.Name == "" {
		return ""
	}
	if len(target.Nifs) == 0 {
		return target.Name
	}
	return target.Name + ": " + strings.Join(target.Nifs, ", ")
}
