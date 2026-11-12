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

func TestAbsFloat32(t *testing.T) {
	tests := []struct {
		name     string
		x        float32
		expected float32
	}{
		{"positive_value", 5.5, 5.5},
		{"negative_value", -5.5, 5.5},
		{"zero_value", 0, 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := AbsFloat32(tt.x)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestMinFloat32(t *testing.T) {
	tests := []struct {
		name     string
		x        float32
		y        []float32
		expected float32
	}{
		{"single_value", 5.5, nil, 5.5},
		{"multiple_values", 10.5, []float32{5.5, 20.5, 3.5}, 3.5},
		{"all_equal", 10.5, []float32{10.5, 10.5, 10.5}, 10.5},
		{"negative_values", -5.5, []float32{-10.5, -3.5}, -10.5},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := MinFloat32(tt.x, tt.y...)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestMaxFloat32(t *testing.T) {
	tests := []struct {
		name     string
		x        float32
		y        []float32
		expected float32
	}{
		{"single_value", 10.5, nil, 10.5},
		{"multiple_values", 10.5, []float32{5.5, 15.5, 25.5}, 25.5},
		{"all_equal", 10.5, []float32{10.5, 10.5, 10.5}, 10.5},
		{"negative_values", -5.5, []float32{-10.5, -3.5}, -3.5},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := MaxFloat32(tt.x, tt.y...)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestAbsFloat64(t *testing.T) {
	tests := []struct {
		name     string
		x        float64
		expected float64
	}{
		{"positive_value", 5.5, 5.5},
		{"negative_value", -5.5, 5.5},
		{"zero_value", 0, 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := AbsFloat64(tt.x)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestMinFloat64(t *testing.T) {
	tests := []struct {
		name     string
		x        float64
		y        []float64
		expected float64
	}{
		{"single_value", 5.5, nil, 5.5},
		{"multiple_values", 10.5, []float64{5.5, 20.5, 3.5}, 3.5},
		{"all_equal", 10.5, []float64{10.5, 10.5, 10.5}, 10.5},
		{"negative_values", -5.5, []float64{-10.5, -3.5}, -10.5},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := MinFloat64(tt.x, tt.y...)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestMaxFloat64(t *testing.T) {
	tests := []struct {
		name     string
		x        float64
		y        []float64
		expected float64
	}{
		{"single_value", 10.5, nil, 10.5},
		{"multiple_values", 10.5, []float64{5.5, 15.5, 25.5}, 25.5},
		{"all_equal", 10.5, []float64{10.5, 10.5, 10.5}, 10.5},
		{"negative_values", -5.5, []float64{-10.5, -3.5}, -3.5},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := MaxFloat64(tt.x, tt.y...)
			assert.Equal(t, tt.expected, result)
		})
	}
}
