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

package stdx

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBigIntFromString_Binary(t *testing.T) {
	tests := []struct {
		str      string
		expected uint64
	}{
		{"0b0100100001101001", 18537},
		{"0B0100100001101001", 18537},
	}
	for _, tt := range tests {
		t.Run(tt.str, func(t *testing.T) {
			bigint, ok := BigIntFromString(tt.str)
			assert.True(t, ok, "Failed to parse octal string")
			assert.Equal(t, tt.expected, bigint.Uint64(), "Wrong value from octal")
		})
	}
}

func TestBigIntFromString_Decimal(t *testing.T) {
	tests := []struct {
		str      string
		expected uint64
	}{
		{"255", 255},
		{"7959", 7959},
		{"0255", 255},
		{"07959", 7959},
		{"007959", 7959},
	}
	for _, tt := range tests {
		t.Run(tt.str, func(t *testing.T) {
			bigint, ok := BigIntFromString(tt.str)
			assert.True(t, ok, "Failed to parse decimal string")
			assert.Equal(t, tt.expected, bigint.Uint64(), "Wrong value from decimal")
		})
	}
}

func TestBigIntFromString_Hexadecimal(t *testing.T) {
	tests := []struct {
		str      string
		expected uint64
	}{
		{"0x255", 597},
		{"0X255", 597},
		{"0x1f17", 7959},
		{"0X1F17", 7959},
	}
	for _, tt := range tests {
		t.Run(tt.str, func(t *testing.T) {
			bigint, ok := BigIntFromString(tt.str)
			assert.True(t, ok, "Failed to parse hexadecimal string")
			assert.Equal(t, tt.expected, bigint.Uint64(), "Wrong value from hexadecimal")
		})
	}
}

func TestBigIntFromString_Octal(t *testing.T) {
	tests := []struct {
		str      string
		expected uint64
	}{
		{"0o255", 173},
		{"0O255", 173},
		{"0o1717", 975},
		{"0O1717", 975},
	}
	for _, tt := range tests {
		t.Run(tt.str, func(t *testing.T) {
			bigint, ok := BigIntFromString(tt.str)
			assert.True(t, ok, "Failed to parse octal string")
			assert.Equal(t, tt.expected, bigint.Uint64(), "Wrong value from octal")
		})
	}
}

func TestBigIntFromString_Invalid(t *testing.T) {
	tests := []struct {
		name string
		str  string
	}{
		{"Invalid decimal character", "1234a567890"},
		{"Invalid decimal character", "12345678z90"},
		{"Invalid hexadecimal character", "0x1234x567890"},
		{"Invalid hexadecimal character", "0X12345678z90"},
		{"Invalid binary character", "0b1f17"},
		{"Invalid binary character", "0b1717"},
		{"Invalid octal character", "0o1f17"},
		{"Invalid octal character", "0oO1717"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, ok := BigIntFromString(tt.str)
			assert.False(t, ok, "Parsed number should be invalid")
		})
	}
}
