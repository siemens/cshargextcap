// (c) Siemens AG 2026
//
// SPDX-License-Identifier: MIT

package extcapargs

import (
	"fmt"
	"slices"
	"strings"
)

// attr is a single “{key=value}” element of an extcap sentence.
type attr struct {
	key, value string
}

// attrs is an ordered list of extcap sentence elements.
type attrs []attr

// with returns a copy with the element for key set to value, keeping the
// position of an existing element and otherwise appending it.
func (as attrs) with(key string, value any) attrs {
	as = slices.Clone(as)
	v := fmt.Sprint(value)
	if i := slices.IndexFunc(as, func(a attr) bool { return a.key == key }); i >= 0 {
		as[i].value = v
		return as
	}
	return append(as, attr{key: key, value: v})
}

// get returns the value of the element for key, or "" if not present.
func (as attrs) get(key string) string {
	if i := slices.IndexFunc(as, func(a attr) bool { return a.key == key }); i >= 0 {
		return as[i].value
	}
	return ""
}

func (as attrs) String() string {
	var b strings.Builder
	for _, a := range as {
		b.WriteString("{")
		b.WriteString(a.key)
		b.WriteString("=")
		b.WriteString(a.value)
		b.WriteString("}")
	}
	return b.String()
}
