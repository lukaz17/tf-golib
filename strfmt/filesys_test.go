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
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tforce-io/tf-golib/multiarch"
)

func TestNewFileName(t *testing.T) {
	tests := []struct {
		group    string
		name     string
		ext      string
		expected *FileName
	}{
		{"empty", "", "", &FileName{}},
		{"name_only", "main", "", &FileName{Name: "main"}},
		{"name_and_ext", "main", ".go", &FileName{Name: "main", Extension: ".go"}},
	}
	for _, tt := range tests {
		t.Run(tt.group, func(t *testing.T) {
			fileName := NewFileName(tt.name, tt.ext)
			assert.Equal(t, tt.expected, fileName)
		})
	}
}

func TestNewFileNameFromStr(t *testing.T) {
	type Test struct {
		group    string
		path     string
		expected *FileName
	}
	var tests []Test
	if multiarch.IsWindows() {
		tests = []Test{
			{"empty", "", &FileName{}},
			{"root", "/", &FileName{}},
			{"root", `\`, &FileName{}},
			{"root", "//", &FileName{}},
			{"root", `\\`, &FileName{}},
			{"drive", `c:`, &FileName{}},
			{"drive", `c:\`, &FileName{}},
			{"drive", `c:::`, &FileName{Name: "::"}},
			{"drive", `c:abc`, &FileName{Name: "abc"}},
			{"trailing_slash", "/usr/local/tforce-io/tf-golib/", &FileName{Name: "tf-golib"}},
			{"trailing_slash", "///usr/local/tforce-io/tf-golib/", &FileName{Name: "tf-golib"}},
			{"name_only", "main", &FileName{Name: "main"}},
			{"name_and_ext", "main.go", &FileName{Name: "main", Extension: ".go"}},
			{"relative_path_and_filename", "tf-golib/main", &FileName{Name: "main"}},
			{"relative_path_and_filename", "tf-golib/main.go", &FileName{Name: "main", Extension: ".go"}},
			{"absolute_path_and_filename", "/usr/local/tforce-io/tf-golib/main", &FileName{Name: "main"}},
			{"absolute_path_and_filename", "/usr/local/tforce-io/tf-golib/main.go", &FileName{Name: "main", Extension: ".go"}},
		}
	} else {
		tests = []Test{
			{"empty", "", &FileName{}},
			{"root", "/", &FileName{}},
			{"root", `\`, &FileName{}},
			{"root", "//", &FileName{}},
			{"root", `\\`, &FileName{}},
			{"trailing_slash", "/usr/local/tforce-io/tf-golib/", &FileName{Name: "tf-golib"}},
			{"trailing_slash", "///usr/local/tforce-io/tf-golib/", &FileName{Name: "tf-golib"}},
			{"name_only", "main", &FileName{Name: "main"}},
			{"name_and_ext", "main.go", &FileName{Name: "main", Extension: ".go"}},
			{"relative_path_and_filename", "tf-golib/main", &FileName{Name: "main"}},
			{"relative_path_and_filename", "tf-golib/main.go", &FileName{Name: "main", Extension: ".go"}},
			{"absolute_path_and_filename", "/usr/local/tforce-io/tf-golib/main", &FileName{Name: "main"}},
			{"absolute_path_and_filename", "/usr/local/tforce-io/tf-golib/main.go", &FileName{Name: "main", Extension: ".go"}},
		}
	}
	for _, tt := range tests {
		t.Run(tt.group, func(t *testing.T) {
			result := NewFileNameFromStr(tt.path)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestFileName_FullName(t *testing.T) {
	tests := []struct {
		group    string
		name     string
		ext      string
		prefix   string
		suffix   string
		expected string
	}{
		{"name_only", "main", "", "", "", "main"},
		{"name_with_ext", "main", ".go", "", "", "main.go"},
		{"name_with_prefix", "main", "", "copy_of_", "", "copy_of_main"},
		{"name_with_prefix", "main", ".go", "copy_of_", "", "copy_of_main.go"},
		{"name_with_suffix", "main", "", "", "_original", "main_original"},
		{"name_with_suffix", "main", ".go", "", "_original", "main_original.go"},
	}
	for _, tt := range tests {
		t.Run(tt.group, func(t *testing.T) {
			fileName := NewFileName(tt.name, tt.ext)
			fileName.Prefix = tt.prefix
			fileName.Suffix = tt.suffix
			fullName := fileName.FullName()
			assert.Equal(t, tt.expected, fullName)
		})
	}
}

func TestAreEqualFileNames(t *testing.T) {
	tests := []struct {
		group    string
		x        *FileName
		y        *FileName
		expected bool
	}{
		// all_nils
		{"all_nils", nil, nil, true},
		// one_nil
		{"one_nil", &FileName{}, nil, false},
		// deep_compare
		{"deep_compare", &FileName{}, &FileName{}, true},
		{"deep_compare", &FileName{Name: "main"}, &FileName{Name: "main"}, true},
		{"deep_compare", &FileName{Name: "main", Extension: ".go"}, &FileName{Name: "main", Extension: ".go"}, true},
	}
	for _, tt := range tests {
		t.Run(tt.group, func(t *testing.T) {
			result := AreEqualFileNames(tt.x, tt.y)
			assert.Equal(t, tt.expected, result, "left %v right %v", tt.x, tt.y)
			result = AreEqualFileNames(tt.y, tt.x)
			assert.Equal(t, tt.expected, result, "left %v right %v", tt.y, tt.x)
		})
	}
}

func TestNewPath(t *testing.T) {
	type Test struct {
		group string
		dirs  []string
		file  *FileName
	}
	var tests []Test
	if multiarch.IsWindows() {
		tests = []Test{
			{"empty", []string{}, NewFileNameFromStr("")},
			{"file_only", []string{}, NewFileNameFromStr("main")},
			{"file_only", []string{}, NewFileNameFromStr("main.go")},
			{"relative_path", []string{"tf-golib"}, NewFileNameFromStr("main.go")},
			{"absolute_path", []string{"d:", "tforce-io", "tf-golib"}, NewFileNameFromStr("main.go")},
		}
	} else {
		tests = []Test{
			{"empty", []string{}, NewFileNameFromStr("")},
			{"file_only", []string{}, NewFileNameFromStr("main")},
			{"file_only", []string{}, NewFileNameFromStr("main.go")},
			{"relative_path", []string{"tf-golib"}, NewFileNameFromStr("main.go")},
			{"absolute_path", []string{"d:", "tforce-io", "tf-golib"}, NewFileNameFromStr("main.go")},
		}
	}
	for _, tt := range tests {
		t.Run(tt.group, func(t *testing.T) {
			path := NewPath(tt.dirs, tt.file)
			assert.Equal(t, tt.dirs, path.Parents)
			assert.Equal(t, tt.file.FullName(), path.Name.FullName())
		})
	}
}

func TestNewPathFromStr(t *testing.T) {
	type Test struct {
		group   string
		path    string
		parents []string
		name    string
	}
	var tests []Test
	if multiarch.IsWindows() {
		tests = []Test{
			{"file_only", `main`, []string{}, "main"},
			{"file_only", `main.go`, []string{}, "main.go"},
			{"relative_path", `tf-golib\main.go`, []string{"tf-golib"}, "main.go"},
			{"relative_path", `tforce-io\tf-golib\main.go`, []string{"tforce-io", "tf-golib"}, "main.go"},
			{"relative_path", `.\tforce-io\tf-golib\main.go`, []string{"tforce-io", "tf-golib"}, "main.go"},
			{"relative_path", `..\tforce-io\tf-golib\main.go`, []string{"..", "tforce-io", "tf-golib"}, "main.go"},
			{"absolute_path", `d:\repositories\tforce-io\tf-golib\main.go`, []string{"d:", "repositories", "tforce-io", "tf-golib"}, "main.go"},
			{"environment_variable", `%DOCUMENTS%\main.go`, []string{"%DOCUMENTS%"}, "main.go"},
			{"trailing_slash", `main.go\`, []string{}, "main.go"},
			{"trailing_slash", `main.go\\`, []string{}, "main.go"},
			{"trailing_slash", `tf-golib\\main.go`, []string{"tf-golib"}, "main.go"},
			{"trailing_slash", `tf-golib\\main.go\\`, []string{"tf-golib"}, "main.go"},
		}
	} else {
		tests = []Test{
			{"file_only", "main", []string{}, "main"},
			{"file_only", "main.go", []string{}, "main.go"},
			{"relative_path", "tf-golib/main.go", []string{"tf-golib"}, "main.go"},
			{"relative_path", "tforce-io/tf-golib/main.go", []string{"tforce-io", "tf-golib"}, "main.go"},
			{"relative_path", "./tforce-io/tf-golib/main.go", []string{"tforce-io", "tf-golib"}, "main.go"},
			{"relative_path", "../tforce-io/tf-golib/main.go", []string{"..", "tforce-io", "tf-golib"}, "main.go"},
			{"absolute_path", "/repositories/tforce-io/tf-golib/main.go", []string{"", "repositories", "tforce-io", "tf-golib"}, "main.go"},
			{"environment_variable", "${pwd}/main.go", []string{"${pwd}"}, "main.go"},
			{"trailing_slash", "main.go/", []string{}, "main.go"},
			{"trailing_slash", "main.go//", []string{}, "main.go"},
			{"trailing_slash", "tf-golib//main.go", []string{"tf-golib"}, "main.go"},
			{"trailing_slash", "tf-golib//main.go//", []string{"tf-golib"}, "main.go"},
		}
	}
	for _, tt := range tests {
		t.Run(tt.group, func(t *testing.T) {
			path := NewPathFromStr(tt.path)
			assert.Equal(t, tt.parents, path.Parents)
			assert.Equal(t, tt.name, path.Name.FullName())
		})
	}
}

func TestPath_Clone(t *testing.T) {
	tests := []struct {
		name     string
		path     *Path
		expected *Path
	}{
		{"keeps single parent",
			&Path{Parents: []string{"tf-golib"}, Name: &FileName{Name: "main", Extension: ".go"}},
			&Path{Parents: []string{"tf-golib"}, Name: &FileName{Name: "main", Extension: ".go"}}},
		{"keeps nested parents",
			&Path{Parents: []string{"d:", "tforce-io", "tf-golib"}, Name: &FileName{Name: "main", Extension: ".go"}},
			&Path{Parents: []string{"d:", "tforce-io", "tf-golib"}, Name: &FileName{Name: "main", Extension: ".go"}}},
		{"keeps empty parents",
			&Path{Parents: []string{}, Name: &FileName{Name: "main", Extension: ".go"}},
			&Path{Parents: []string{}, Name: &FileName{Name: "main", Extension: ".go"}}},
		{"keeps nil parents",
			&Path{Name: &FileName{Name: "main", Extension: ".go"}},
			&Path{Name: &FileName{Name: "main", Extension: ".go"}}},
		{"keeps nil name",
			&Path{Parents: []string{"tf-golib"}},
			&Path{Parents: []string{"tf-golib"}}},
		{"keeps zero value", &Path{}, &Path{}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			clone := tt.path.Clone()

			assert.Equal(t, tt.expected, clone)
		})
	}
}

func TestPath_Clone_LeavesOriginalUnchanged(t *testing.T) {
	path := NewPath([]string{"tforce-io", "tf-golib"}, NewFileName("main", ".go"))

	_ = path.Clone()

	assert.Equal(t, []string{"tforce-io", "tf-golib"}, path.Parents)
}

func TestPath_Clone_IsIndependentOfOriginal(t *testing.T) {
	path := NewPath([]string{"tforce-io", "tf-golib"}, NewFileName("main", ".go"))
	clone := path.Clone()

	clone.Parents[0] = "other"
	clone.Name.Name = "other"

	assert.Equal(t, []string{"tforce-io", "tf-golib"}, path.Parents)
	assert.Equal(t, "main", path.Name.Name)
}

func TestPath_Clone_IsSafeForConcurrentUse(t *testing.T) {
	path := NewPath([]string{"tforce-io", "tf-golib"}, NewFileName("main", ".go"))
	const cloners = 8
	clones := make([]*Path, cloners)
	var wg sync.WaitGroup

	for i := 0; i < cloners; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			clones[i] = path.Clone()
		}(i)
	}
	wg.Wait()

	for _, clone := range clones {
		assert.Equal(t, []string{"tforce-io", "tf-golib"}, clone.Parents)
	}
}

func TestPath_IsAbsolute(t *testing.T) {
	type Test struct {
		group    string
		dirs     []string
		file     *FileName
		expected bool
	}
	var tests []Test
	if multiarch.IsWindows() {
		tests = []Test{
			{"empty", []string{}, NewFileNameFromStr(""), false},
			{"file_only", []string{}, NewFileNameFromStr("main"), false},
			{"file_only", []string{}, NewFileNameFromStr("main.go"), false},
			{"relative_path", []string{"tf-golib"}, NewFileNameFromStr("main.go"), false},
			{"absolute_path", []string{"d:", "tforce-io", "tf-golib"}, NewFileNameFromStr("main.go"), true},
		}
	} else {
		tests = []Test{
			{"empty", []string{}, NewFileNameFromStr(""), false},
			{"file_only", []string{}, NewFileNameFromStr("main"), false},
			{"file_only", []string{}, NewFileNameFromStr("main.go"), false},
			{"relative_path", []string{"tf-golib"}, NewFileNameFromStr("main.go"), false},
			{"absolute_path", []string{"", "tforce-io", "tf-golib"}, NewFileNameFromStr("main.go"), true},
		}
	}
	for _, tt := range tests {
		t.Run(tt.group, func(t *testing.T) {
			path := NewPath(tt.dirs, tt.file)
			result := path.IsAbsolute()
			assert.Equal(t, tt.expected, result, "intermediate %v", path.FullPath())
		})
	}
}

func TestPath_FullPath(t *testing.T) {
	type Test struct {
		group    string
		dirs     []string
		file     *FileName
		expected string
	}
	var tests []Test
	if multiarch.IsWindows() {
		tests = []Test{
			{"empty", []string{}, NewFileNameFromStr(""), ""},
			{"file_only", []string{}, NewFileNameFromStr("main"), "main"},
			{"file_only", []string{}, NewFileNameFromStr("main.go"), "main.go"},
			{"relative_path", []string{"tf-golib"}, NewFileNameFromStr("main.go"), "tf-golib\\main.go"},
			{"absolute_path", []string{"d:", "tforce-io", "tf-golib"}, NewFileNameFromStr("main.go"), "d:\\tforce-io\\tf-golib\\main.go"},
		}
	} else {
		tests = []Test{
			{"empty", []string{}, NewFileNameFromStr(""), ""},
			{"file_only", []string{}, NewFileNameFromStr("main"), "main"},
			{"file_only", []string{}, NewFileNameFromStr("main.go"), "main.go"},
			{"relative_path", []string{"tf-golib"}, NewFileNameFromStr("main.go"), "tf-golib/main.go"},
			{"absolute_path", []string{"", "tforce-io", "tf-golib"}, NewFileNameFromStr("main.go"), "/tforce-io/tf-golib/main.go"},
		}
	}
	for _, tt := range tests {
		t.Run(tt.group, func(t *testing.T) {
			path := NewPath(tt.dirs, tt.file)
			result := path.FullPath()
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestPath_ParentPath(t *testing.T) {
	type Test struct {
		group    string
		dirs     []string
		file     *FileName
		expected string
	}
	var tests []Test
	if multiarch.IsWindows() {
		tests = []Test{
			{"empty", []string{}, NewFileNameFromStr(""), ""},
			{"file_only", []string{}, NewFileNameFromStr("main"), ""},
			{"file_only", []string{}, NewFileNameFromStr("main.go"), ""},
			{"relative_path", []string{"tf-golib"}, NewFileNameFromStr("main.go"), "tf-golib"},
			{"absolute_path", []string{"d:", "tforce-io", "tf-golib"}, NewFileNameFromStr("main.go"), "d:\\tforce-io\\tf-golib"},
			{"absolute_path", []string{"d:"}, NewFileNameFromStr("main.go"), "d:"},
		}
	} else {
		tests = []Test{
			{"empty", []string{}, NewFileNameFromStr(""), ""},
			{"file_only", []string{}, NewFileNameFromStr("main"), ""},
			{"file_only", []string{}, NewFileNameFromStr("main.go"), ""},
			{"relative_path", []string{"tf-golib"}, NewFileNameFromStr("main.go"), "tf-golib"},
			{"absolute_path", []string{"", "tforce-io", "tf-golib"}, NewFileNameFromStr("main.go"), "/tforce-io/tf-golib"},
			{"absolute_path", []string{""}, NewFileNameFromStr("main.go"), ""},
		}
	}
	for _, tt := range tests {
		t.Run(tt.group, func(t *testing.T) {
			path := NewPath(tt.dirs, tt.file)
			result := path.ParentPath()
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestAreEqualPaths(t *testing.T) {
	tests := []struct {
		group    string
		x        *Path
		y        *Path
		expected bool
	}{
		// all_nils
		{"all_nils", nil, nil, true},
		// one_nil
		{"one_nil", &Path{}, nil, false},
		{"one_nil", nil, &Path{}, false},
		// deep_compare
		{"deep_compare: empty", &Path{}, &Path{}, true},
		{"deep_compare: parents_only", &Path{Parents: []string{}}, &Path{}, false},
		{"deep_compare: parents_only", &Path{Parents: []string{"tf-golib"}}, &Path{Parents: []string{}}, false},
		{"deep_compare: name_only", &Path{Name: &FileName{Name: "main"}}, &Path{Name: &FileName{Name: "main"}}, true},
		{"deep_compare: name_only", &Path{Name: &FileName{Name: "main", Extension: ".go"}}, &Path{Name: &FileName{Name: "main", Extension: ".go"}}, true},
		{"deep_compare", &Path{Parents: []string{"tf-golib"}, Name: &FileName{Name: "main", Extension: ".go"}}, &Path{Parents: []string{"tf-golib"}, Name: &FileName{Name: "main", Extension: ".go"}}, true},
	}
	for _, tt := range tests {
		t.Run(tt.group, func(t *testing.T) {
			result := AreEqualPaths(tt.x, tt.y)
			assert.Equal(t, tt.expected, result, "left %v right %v", tt.x, tt.y)
			result = AreEqualPaths(tt.y, tt.x)
			assert.Equal(t, tt.expected, result, "left %v right %v", tt.y, tt.x)
		})
	}
}
