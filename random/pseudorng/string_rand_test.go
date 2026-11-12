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

package pseudorng

import "testing"

func TestBase64(t *testing.T) {
	tests := []struct {
		group    string
		len      uint
		expected uint
	}{
		{"zero_byte", 0, 0},
		{"one_byte", 1, 4},
		{"two_bytes", 2, 4},
		{"three_bytes", 3, 4},
		{"four_bytes", 4, 8},
		{"five_bytes", 5, 8},
		{"six_bytes", 6, 8},
		{"seventy_bytes", 70, 96},
		{"seventy_one_bytes", 71, 96},
		{"seventy_two_bytes", 72, 96},
	}
	for _, tt := range tests {
		t.Run(tt.group, func(t *testing.T) {
			result := Base64(tt.len)
			if len(result) != int(tt.expected) {
				t.Errorf("expected %v actual %v underlined %v", tt.expected, len(result), result)
			}
		})
	}
}

func TestHex(t *testing.T) {
	tests := []struct {
		group    string
		len      uint
		expected uint
	}{
		{"zero_byte", 0, 0},
		{"one_byte", 1, 2},
		{"two_bytes", 2, 4},
		{"three_bytes", 3, 6},
		{"four_bytes", 4, 8},
		{"five_bytes", 5, 10},
		{"six_bytes", 6, 12},
		{"seventy_bytes", 70, 140},
		{"seventy_one_bytes", 71, 142},
		{"one_hundred_bytes", 100, 200},
	}
	for _, tt := range tests {
		t.Run(tt.group, func(t *testing.T) {
			result := Hex(tt.len)
			if len(result) != int(tt.expected) {
				t.Errorf("expected %v actual %v underlined %v", tt.expected, len(result), result)
			}
		})
	}
}
