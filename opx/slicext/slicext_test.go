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

package slicext

import (
	"testing"
)

func TestAreEqual(t *testing.T) {
	tests := []struct {
		group    string
		x        []int
		y        []int
		expected bool
	}{
		// all nils
		{"all nils", nil, nil, true},
		// one nils
		{"one nil", []int{}, nil, false},
		// different length
		{"different length", []int{1}, []int{1, 2}, false},
		{"different length", []int{1}, []int{1, 2, 3}, false},
		{"different length", []int{1, 2}, []int{1, 2, 3}, false},
		{"different length", []int{1, 2}, []int{1, 2, 3, 4}, false},
		// deep compare
		{"deep compare", []int{}, []int{}, true},
		{"deep compare", []int{1}, []int{1}, true},
		{"deep compare", []int{1, 3, 5}, []int{1, 3, 5}, true},
		{"deep compare", []int{5, 3, 1}, []int{1, 3, 5}, false},
	}
	for _, tt := range tests {
		t.Run(tt.group, func(t *testing.T) {
			result := AreEqual(tt.x, tt.y)
			if result != tt.expected {
				t.Errorf("expected %v actual left %v actual right %v ", tt.expected, tt.x, tt.y)
			}
			result = AreEqual(tt.y, tt.x)
			if result != tt.expected {
				t.Errorf("expected %v actual left %v actual right %v ", tt.expected, tt.y, tt.x)
			}
		})
	}
}

func TestContains(t *testing.T) {
	tests := []struct {
		group    string
		items    []int
		value    int
		expected bool
	}{
		{"true", []int{3, 5, 7, 2, 4, 6}, 7, true},
		{"true", []int{3, 5, 7, 2, 4, 6}, 4, true},
		{"true", []int{}, 7, false},
		{"true", []int{3, 5, 2, 4, 6}, 7, false},
	}
	for _, tt := range tests {
		t.Run(tt.group, func(t *testing.T) {
			result := Contains(tt.items, tt.value)
			if result != tt.expected {
				t.Errorf("expected %v actual %v intermediate %v %v", tt.expected, result, tt.items, tt.value)
			}
		})
	}
}

func TestIsEmpty(t *testing.T) {
	tests := []struct {
		group    string
		slice    []int
		expected bool
	}{
		// nil
		{"nil", nil, true},
		// empty
		{"empty", []int{}, true},
		// one or more
		{"one or more", []int{1}, false},
		{"one or more", []int{1, 2, 3}, false},
		{"one or more", []int{1, 2, 3, 4}, false},
		{"one or more", []int{1, 2, 3, 4, 5, 6, 7}, false},
	}
	for _, tt := range tests {
		t.Run(tt.group, func(t *testing.T) {
			result := IsEmpty(tt.slice)
			if result != tt.expected {
				t.Errorf("expected %v actual %v intermediate %v", tt.expected, result, tt.slice)
			}
		})
	}
}
