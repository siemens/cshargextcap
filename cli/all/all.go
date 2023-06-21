// (c) Siemens AG 2023
//
// SPDX-License-Identifier: MIT

/*
Package all ensures to pull in the required packages for (obscure) CLI args as
well as the set of extcap implementations
*/
package all

import (
	_ "github.com/siemens/cshargextcap/cli/proxy" // pull in proxy CLI args
	_ "github.com/siemens/cshargextcap/nif"       // pull in extcap nif implementations
)
