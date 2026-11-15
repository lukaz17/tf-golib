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

package stringxt

import "testing"

func TestIsEmpty(t *testing.T) {
	tests := []struct {
		group    string
		str      string
		expected bool
	}{
		// empty
		{"empty", "", true},
		// ascii
		{"ascii", "a", false},
		{"ascii", "hello", false},
		{"ascii", "hello world", false},
	}
	for _, tt := range tests {
		t.Run(tt.group, func(t *testing.T) {
			result := IsEmpty(tt.str)
			if result != tt.expected {
				t.Errorf("expected %v actual %v intermediate %v", tt.expected, result, tt.str)
			}
		})
	}
}

func TestIsEmptyOrWhitespace(t *testing.T) {
	tests := []struct {
		group    string
		str      string
		expected bool
	}{
		// empty
		{"empty", "", true},
		// all whitespace
		{"all whitespace", " ", true},
		{"all whitespace", "  ", true},
		{"all whitespace", "       ", true},
		// all tabs
		{"all tabs", "\t", true},
		{"all tabs", "\t\t", true},
		{"all tabs", "\t\t\t\t\t\t\t", true},
		// mixed whitespaces and tabs
		{"mixed whitespaces and tabs", " \t", true},
		{"mixed whitespaces and tabs", " \t ", true},
		{"mixed whitespaces and tabs", " \t \t", true},
		{"mixed whitespaces and tabs", " \t \t\t \t", true},
		{"mixed whitespaces and tabs", "\t ", true},
		{"mixed whitespaces and tabs", "\t \t", true},
		{"mixed whitespaces and tabs", "\t \t ", true},
		{"mixed whitespaces and tabs", "\t \t  \t ", true},
		// ascii
		{"ascii", "a", false},
		{"ascii", "hello", false},
		{"ascii", "hello world", false},
	}
	for _, tt := range tests {
		t.Run(tt.group, func(t *testing.T) {
			result := IsEmptyOrWhitespace(tt.str)
			if result != tt.expected {
				t.Errorf("expected %v actual %v intermediate %v", tt.expected, result, tt.str)
			}
		})
	}
}
