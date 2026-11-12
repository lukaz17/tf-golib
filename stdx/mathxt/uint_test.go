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

package mathxt

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMinUint(t *testing.T) {
	tests := []struct {
		name     string
		x        uint
		y        []uint
		expected uint
	}{
		{"single_value", 5, nil, 5},
		{"multiple_values", 10, []uint{5, 20, 3}, 3},
		{"all_equal", 10, []uint{10, 10, 10}, 10},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := MinUint(tt.x, tt.y...)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestMaxUint(t *testing.T) {
	tests := []struct {
		name     string
		x        uint
		y        []uint
		expected uint
	}{
		{"single_value", 10, nil, 10},
		{"multiple_values", 10, []uint{5, 15, 25}, 25},
		{"all_equal", 10, []uint{10, 10, 10}, 10},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := MaxUint(tt.x, tt.y...)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestMinUint8(t *testing.T) {
	tests := []struct {
		name     string
		x        uint8
		y        []uint8
		expected uint8
	}{
		{"single_value", 5, nil, 5},
		{"multiple_values", 10, []uint8{5, 20, 3}, 3},
		{"all_equal", 10, []uint8{10, 10, 10}, 10},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := MinUint8(tt.x, tt.y...)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestMaxUint8(t *testing.T) {
	tests := []struct {
		name     string
		x        uint8
		y        []uint8
		expected uint8
	}{
		{"single_value", 10, nil, 10},
		{"multiple_values", 10, []uint8{5, 15, 25}, 25},
		{"all_equal", 10, []uint8{10, 10, 10}, 10},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := MaxUint8(tt.x, tt.y...)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestMinUint16(t *testing.T) {
	tests := []struct {
		name     string
		x        uint16
		y        []uint16
		expected uint16
	}{
		{"single_value", 5, nil, 5},
		{"multiple_values", 10, []uint16{5, 20, 3}, 3},
		{"all_equal", 10, []uint16{10, 10, 10}, 10},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := MinUint16(tt.x, tt.y...)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestMaxUint16(t *testing.T) {
	tests := []struct {
		name     string
		x        uint16
		y        []uint16
		expected uint16
	}{
		{"single value", 10, nil, 10},
		{"two values", 10, []uint16{20}, 20},
		{"multiple values", 10, []uint16{5, 15, 25}, 25},
		{"all equal values", 10, []uint16{10, 10, 10}, 10},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := MaxUint16(tt.x, tt.y...)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestMinUint32(t *testing.T) {
	tests := []struct {
		name     string
		x        uint32
		y        []uint32
		expected uint32
	}{
		{"single_value", 5, nil, 5},
		{"multiple_values", 10, []uint32{5, 20, 3}, 3},
		{"all_equal", 10, []uint32{10, 10, 10}, 10},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := MinUint32(tt.x, tt.y...)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestMaxUint32(t *testing.T) {
	tests := []struct {
		name     string
		x        uint32
		y        []uint32
		expected uint32
	}{
		{"single_value", 10, nil, 10},
		{"multiple_values", 10, []uint32{5, 15, 25}, 25},
		{"all_equal", 10, []uint32{10, 10, 10}, 10},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := MaxUint32(tt.x, tt.y...)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestMinUint64(t *testing.T) {
	tests := []struct {
		name     string
		x        uint64
		y        []uint64
		expected uint64
	}{
		{"single_value", 5, nil, 5},
		{"multiple_values", 10, []uint64{5, 20, 3}, 3},
		{"all_equal", 10, []uint64{10, 10, 10}, 10},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := MinUint64(tt.x, tt.y...)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestMaxUint64(t *testing.T) {
	tests := []struct {
		name     string
		x        uint64
		y        []uint64
		expected uint64
	}{
		{"single_value", 10, nil, 10},
		{"multiple_values", 10, []uint64{5, 15, 25}, 25},
		{"all_equal", 10, []uint64{10, 10, 10}, 10},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := MaxUint64(tt.x, tt.y...)
			assert.Equal(t, tt.expected, result)
		})
	}
}
