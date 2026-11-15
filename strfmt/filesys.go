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
	"os"
	"path/filepath"
	"strings"

	"github.com/tforce-io/tf-golib/multiarch"
	"github.com/tforce-io/tf-golib/stdx/stringxt"
)

var PATH_SEPARATOR = string(os.PathSeparator)

// A FileName is struct contains name and extension of a file or folder.
// Extension should have dot at the beginning.
type FileName struct {
	Prefix    string
	Name      string
	Suffix    string
	Extension string
}

// Create a new FileName from scratch
func NewFileName(name, extension string) *FileName {
	return &FileName{
		Name:      name,
		Extension: extension,
	}
}

// Parse path string to create a new FileName
func NewFileNameFromStr(path string) *FileName {
	nPath := NormalizePath(path)
	base := filepath.Base(nPath)
	if base == "." || base == PATH_SEPARATOR {
		return &FileName{
			Name:      "",
			Extension: "",
		}
	}

	ext := filepath.Ext(base)
	return &FileName{
		Name:      strings.TrimSuffix(base, ext),
		Extension: ext,
	}
}

// Make a deep copy of this Path.
func (s *FileName) Clone() *FileName {
	return &FileName{
		Prefix:    s.Prefix,
		Name:      s.Name,
		Suffix:    s.Suffix,
		Extension: s.Extension,
	}
}

// Return full name represented by this FileName.
func (s *FileName) FullName() string {
	return s.Prefix + s.Name + s.Suffix + s.Extension
}

// Check whether two FileNames are equal.
func AreEqualFileNames(x, y *FileName) bool {
	if x == nil && y == nil {
		return true
	}
	if x == nil || y == nil {
		return false
	}
	return x.Name == y.Name &&
		x.Extension == y.Extension &&
		x.Prefix == y.Prefix &&
		x.Suffix == y.Suffix
}

// A Path is a struct contains all smallest components of a path.
type Path struct {
	Parents []string
	Name    *FileName
}

// Create a new Path from scratch
func NewPath(dirs []string, name *FileName) *Path {
	return &Path{
		Parents: dirs,
		Name:    name,
	}
}

// Parse path string to create a new Path
func NewPathFromStr(path string) *Path {
	nPath := NormalizePath(path)
	dir, file := filepath.Split(nPath)
	dir = strings.TrimSuffix(dir, PATH_SEPARATOR)
	dirs := ternarySlice(stringxt.IsEmpty(dir), []string{}, strings.Split(dir, PATH_SEPARATOR))
	name := NewFileNameFromStr(file)
	return &Path{
		Parents: dirs,
		Name:    name,
	}
}

// Make a deep copy of this Path.
func (s *Path) Clone() *Path {
	directories := make([]string, len(s.Parents))
	_ = copy(s.Parents, directories)
	name := s.Name.Clone()
	return &Path{
		Parents: directories,
		Name:    name,
	}
}

// Check a Path is whether asbsolute path. Using the same rule as
// filepath.IsAbs
func (s *Path) IsAbsolute() bool {
	if isEmptySlice(s.Parents) {
		return false
	}
	fullPath := s.FullPath()
	return filepath.IsAbs(fullPath)
}

// Return full path represented by this Path.
func (s *Path) FullPath() string {
	if isEmptySlice(s.Parents) {
		return s.Name.FullName()
	}
	return s.ParentPath() + PATH_SEPARATOR + s.Name.FullName()
}

// Return parent path represented by this Path.
func (s *Path) ParentPath() string {
	return ternary(isEmptySlice(s.Parents), "", strings.Join(s.Parents, PATH_SEPARATOR))
}

// Check whether two Paths are equal.
func AreEqualPaths(x, y *Path) bool {
	if x == nil && y == nil {
		return true
	}
	if x == nil || y == nil {
		return false
	}
	return areEqualSlices(x.Parents, y.Parents) &&
		AreEqualFileNames(x.Name, y.Name)
}

// Clean path and make path consistent accross platforms. This function will perform:
// - Replace all slashes to backslashes if run on Windows.
// - Replace all backslashes to slashes if run on all UNIX-like OSes.
// - Clean the path.
func NormalizePath(path string) string {
	nPath := ternary(multiarch.IsWindows(),
		strings.ReplaceAll(path, "/", "\\"),
		strings.ReplaceAll(path, "\\", "/"),
	)
	return filepath.Clean(nPath)
}

// forked from opx package
func areEqualSlices(x, y []string) bool {
	if x == nil && y == nil {
		return true
	}
	if x == nil || y == nil {
		return false
	}
	if len(x) != len(y) {
		return false
	}
	for i := range x {
		if x[i] != y[i] {
			return false
		}
	}
	return true
}

// forked from opx package
func isEmptySlice(slice []string) bool {
	if slice == nil {
		return true
	}
	return len(slice) == 0
}

// forked from opx package
func ternary(cond bool, x, y string) string {
	if cond {
		return x
	}
	return y
}

// forked from opx package
func ternarySlice(cond bool, x, y []string) []string {
	if cond {
		return x
	}
	return y
}
