// Copyright (c) 2025 T-Force I/O
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

package strfmt

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewVersion(t *testing.T) {
	tests := []struct {
		major      uint64
		minor      uint64
		patch      uint64
		revision   uint64
		prerelease string
		metadata   string
		expected   string
	}{
		{1, 0, 0, 0, "", "", "1.0.0"},
		{1, 0, 0, 1, "", "", "1.0.0.1"},
		{1, 2, 1, 0, "", "", "1.2.1"},
		{1, 2, 1, 11, "", "", "1.2.1.11"},
		{1, 3, 2, 0, "beta", "", "1.3.2-beta"},
		{1, 3, 2, 7, "beta", "", "1.3.2.7-beta"},
		{1, 4, 7, 12, "", "meta", "1.4.7.12+meta"},
		{1, 4, 7, 0, "", "meta", "1.4.7+meta"},
		{1, 4, 7, 4, "", "meta", "1.4.7.4+meta"},
		{1, 5, 9, 0, "alpha", "meta", "1.5.9-alpha+meta"},
		{1, 5, 9, 2, "alpha", "meta", "1.5.9.2-alpha+meta"},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			ver := NewVersion(tt.major, tt.minor, tt.patch, tt.revision, tt.prerelease, tt.metadata)
			assert.Equal(t, tt.expected, ver.String())
		})
	}
}

func TestFromString(t *testing.T) {
	valids := []string{
		"0.0.4",
		"1.2.3",
		"10.20.30",
		"1.1.2-prerelease+meta",
		"1.1.2+meta",
		"1.1.2+meta-valid",
		"1.0.0-alpha",
		"1.0.0-beta",
		"1.0.0-alpha.beta",
		"1.0.0-alpha.beta.1",
		"1.0.0-alpha.1",
		"1.0.0-alpha0.valid",
		"1.0.0-alpha.0valid",
		"1.0.0-alpha-a.b-c-somethinglong+build.1-aef.1-its-okay",
		"1.0.0-rc.1+build.1",
		"2.0.0-rc.1+build.123",
		"1.2.3-beta",
		"10.2.3-DEV-SNAPSHOT",
		"1.2.3-SNAPSHOT-123",
		"1.0.0",
		"2.0.0",
		"1.1.7",
		"2.0.0+build.1848",
		"2.0.1-alpha.1227",
		"1.0.0-alpha+beta",
		"1.2.3----RC-SNAPSHOT.12.9.1--.12+788",
		"1.2.3----R-S.12.9.1--.12+meta",
		"1.2.3----RC-SNAPSHOT.12.9.1--.12",
		"1.0.0+0.build.1-rc.10000aaa-kk-0.1",
		"99999999999999999.999999999999999999.99999999999999999",
		"1.0.0-0A.is.legal",

		"1.2.3.5",
	}
	for _, v := range valids {
		t.Run(v, func(t *testing.T) {
			ver := FromString(v)
			assert.True(t, ver.IsValid())
			assert.Equal(t, v, ver.String())
		})
	}

	invalids := []string{
		"1",
		"1.2",
		"1.2.3-0123",
		"1.2.3-0123.0123",
		"1.1.2+.123",
		"+invalid",
		"-invalid",
		"-invalid+invalid",
		"-invalid.01",
		"alpha",
		"alpha.beta",
		"alpha.beta.1",
		"alpha.1",
		"alpha+beta",
		"alpha_beta",
		"alpha.",
		"alpha..",
		"beta",
		"1.0.0-alpha_beta",
		"-alpha.",
		"1.0.0-alpha..",
		"1.0.0-alpha..1",
		"1.0.0-alpha...1",
		"1.0.0-alpha....1",
		"1.0.0-alpha.....1",
		"1.0.0-alpha......1",
		"1.0.0-alpha.......1",
		"01.1.1",
		"1.01.1",
		"1.1.01",
		"1.2",
		"1.2.3.DEV",
		"1.2-SNAPSHOT",
		"1.2.31.2.3----RC-SNAPSHOT.12.09.1--..12+788",
		"1.2-RC-SNAPSHOT",
		"-1.0.3-gamma+b7718",
		"+justmeta",
		"9.8.7+meta+meta",
		"9.8.7-whatever+meta+meta",
		"99999999999999999999999.999999999999999999.99999999999999999----RC-SNAPSHOT.12.09.1--------------------------------..12",

		"1.2.3.0",
	}
	for _, v := range invalids {
		t.Run(v, func(t *testing.T) {
			ver := FromString(v)
			assert.Equal(t, "0.0.0", ver.String())
		})
	}
}
