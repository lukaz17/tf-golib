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

package multiarch

import (
	"runtime"
)

var linuxFamily = []string{"linux", "freebsd", "netbsd", "openbsd", "aix", "illumos", "plan9", "solaris"}
var wasmFamily = []string{"wasip1"}

// Check whether current platform is Android.
func IsAndroid() bool {
	return runtime.GOOS == "android"
}

// Check whether current platform is iOS.
func IsIos() bool {
	return runtime.GOOS == "ios"
}

// Check whether current platform is Linux.
// "linux", "freebsd", "netbsd", "openbsd", "aix", "illumos", "plan9", "solaris" are considered Linux.
func IsLinux() bool {
	for _, os := range linuxFamily {
		if os == runtime.GOOS {
			return true
		}
	}
	return false
}

// Check whether current platform is MacOSX.
func IsMacintosh() bool {
	return runtime.GOOS == "darwin"
}

// Check whether current platform is Web Assembly.
func IsWebAssembly() bool {
	for _, os := range wasmFamily {
		if os == runtime.GOOS {
			return true
		}
	}
	return false
}

// Check whether current platform is Windows.
func IsWindows() bool {
	return runtime.GOOS == "windows"
}
