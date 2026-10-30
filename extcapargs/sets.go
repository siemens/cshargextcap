// (c) Siemens AG 2026
//
// SPDX-License-Identifier: MIT

package extcapargs

import (
	"fmt"
	"io"
	"maps"
	"slices"
)

// Merge returns a new set of args merged from the specified sets, keyed by
// their arg numbers. It fails if an arg number occurs more than once, as
// Wireshark requires unique arg numbers, or if a CLI flag is used by more than
// one arg.
func Merge(sets ...map[int]ArgDef) (map[int]ArgDef, error) {
	merged := map[int]ArgDef{}
	calls := map[string]int{}
	for _, set := range sets {
		for _, num := range slices.Sorted(maps.Keys(set)) {
			if _, ok := merged[num]; ok {
				return nil, fmt.Errorf("duplicate config arg number %d", num)
			}
			call := set[num].Call()
			if other, ok := calls[call]; ok {
				return nil, fmt.Errorf("duplicate config arg call %q for numbers %d and %d",
					call, other, num)
			}
			merged[num] = set[num]
			calls[call] = num
		}
	}
	return merged, nil
}

// Dump writes the “arg” sentences for the specified args to w, sorted by their
// arg numbers. While Wireshark sorts the args by their numbers in its
// configuration dialog, it places the initial focus on the first arg seen in a
// group, not the first one by number.
func Dump(w io.Writer, args map[int]ArgDef) {
	for _, num := range slices.Sorted(maps.Keys(args)) {
		_, _ = fmt.Fprintf(w, "arg {number=%d}%s\n", num, args[num])
	}
}

// DumpValues writes the “value” sentences for the specified arg values to w,
// sorted by their arg numbers while keeping the order of the values of each
// arg.
func DumpValues(w io.Writer, values map[int][]ValueDef) {
	for _, num := range slices.Sorted(maps.Keys(values)) {
		for _, val := range values[num] {
			_, _ = fmt.Fprintf(w, "value {arg=%d}%s\n", num, val)
		}
	}
}
