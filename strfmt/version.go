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
	"fmt"
	"regexp"
	"strconv"
)

// Version represents a struct that followed the rule of semantic version.
// For more information, please refer to https://semver.org/
// For strict semantic version compliant, please keep Revision as 0.
type Version struct {
	Major    uint64
	Minor    uint64
	Patch    uint64
	Revision uint64

	PreRelease string
	BuildMeta  string
}

// Returns a new version instance.
func NewVersion(major, minor, patch, revision uint64, preRelease, buildMeta string) Version {
	return Version{
		Major:    major,
		Minor:    minor,
		Patch:    patch,
		Revision: revision,

		PreRelease: preRelease,
		BuildMeta:  buildMeta,
	}
}

// Parses a string and returns a new version instance.
func FromString(versionStr string) Version {
	semVerRegEx := regexp.MustCompile(`^(?P<major>0|[1-9]\d*)\.(?P<minor>0|[1-9]\d*)\.(?P<patch>0|[1-9]\d*)(?:\.(?P<revision>[1-9]\d*))?(?:-(?P<prerelease>(?:0|[1-9]\d*|\d*[a-zA-Z-][0-9a-zA-Z-]*)(?:\.(?:0|[1-9]\d*|\d*[a-zA-Z-][0-9a-zA-Z-]*))*))?(?:\+(?P<buildmetadata>[0-9a-zA-Z-]+(?:\.[0-9a-zA-Z-]+)*))?$`)
	matches := semVerRegEx.FindStringSubmatch(versionStr)

	if len(matches) == 0 {
		return Version{}
	}

	major, _ := strconv.ParseUint(matches[semVerRegEx.SubexpIndex("major")], 10, 0)
	minor, _ := strconv.ParseUint(matches[semVerRegEx.SubexpIndex("minor")], 10, 0)
	patch, _ := strconv.ParseUint(matches[semVerRegEx.SubexpIndex("patch")], 10, 0)
	revision, _ := strconv.ParseUint(matches[semVerRegEx.SubexpIndex("revision")], 10, 0)
	preRelease := matches[semVerRegEx.SubexpIndex("prerelease")]
	buildMeta := matches[semVerRegEx.SubexpIndex("buildmetadata")]
	return Version{
		Major:    major,
		Minor:    minor,
		Patch:    patch,
		Revision: revision,

		PreRelease: preRelease,
		BuildMeta:  buildMeta,
	}
}

// Checks current Version is compliant with semantic versioning specification.
func (v Version) IsValid() bool {
	prRegex := `^(?P<prerelease>(?:0|[1-9]\d*|\d*[a-zA-Z-][0-9a-zA-Z-]*)(?:\.(?:0|[1-9]\d*|\d*[a-zA-Z-][0-9a-zA-Z-]*))*)?$`
	bmRegex := `^(?P<buildmetadata>[0-9a-zA-Z-]+(?:\.[0-9a-zA-Z-]+)*)?$`
	r1, _ := regexp.MatchString(prRegex, v.PreRelease)
	r2, _ := regexp.MatchString(bmRegex, v.BuildMeta)
	return r1 && r2
}

// Returns a string represents semantic version.
func (v Version) String() string {
	version := ""
	if v.Revision == 0 {
		version = fmt.Sprintf("%d.%d.%d", v.Major, v.Minor, v.Patch)
	} else {
		version = fmt.Sprintf("%d.%d.%d.%d", v.Major, v.Minor, v.Patch, v.Revision)
	}
	if v.PreRelease != "" {
		version += "-" + v.PreRelease
	}
	if v.BuildMeta != "" {
		version += "+" + v.BuildMeta
	}
	return version
}
