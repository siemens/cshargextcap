// (c) Siemens AG 2023
//
// SPDX-License-Identifier: MIT

package cfg

// Shared "Standard" tab controls, except for the packetflix extcap network
// interface. These IDs do not start at 1 (or 0 fwiw) in order to give head
// space to extcap interface-specific configuration controls.
const (
	Containers    = 10
	ContainerNifs = 20

	ShowPods                 = 30
	ShowStandaloneContainers = 31
	ShowProcs                = 32
	ShowEmptyNetNS           = 33

	NoProm = 40
)

// Configuration UI controls identifiers for the Docker capture "Standard" tab
// controls; please note that these go first in the rendered dialog, as these
// are the most important ones.
const (
	DockerHostURL = 1
	SkipVerify    = 2
)

// Packetflix configuration "Standard" tab controls
const (
	PacketflixURL = 10
)

// Kubernetes capture "Standard" tab controls
const (
	ClusterContext = 1
)
