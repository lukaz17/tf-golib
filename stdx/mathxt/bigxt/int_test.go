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

package bigxt

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tforce-io/tf-golib/stdx"
)

func TestAreEqualInt(t *testing.T) {
	tests := []struct {
		name     string
		a        string
		b        string
		expected bool
	}{
		{"Equal positive numbers", "100", "100", true},
		{"Equal negative numbers", "-100", "-100", true},
		{"Equal zero", "0", "0", true},
		{"Different positive numbers", "100", "200", false},
		{"Different negative numbers", "-100", "-200", false},
		{"Positive and negative", "100", "-100", false},
		{"Zero and positive", "0", "1", false},
		{"Zero and negative", "0", "-1", false},
		{"Large equal numbers", "999999999999999999999999999999", "999999999999999999999999999999", true},
		{"Large different numbers", "999999999999999999999999999999", "999999999999999999999999999998", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a, _ := stdx.BigIntFromString(tt.a)
			b, _ := stdx.BigIntFromString(tt.b)
			result := AreEqualInt(a, b)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestGreaterThanInt(t *testing.T) {
	tests := []struct {
		name     string
		a        string
		b        string
		expected bool
	}{
		{"Greater positive numbers", "200", "100", true},
		{"Greater negative numbers", "-100", "-200", true},
		{"Positive greater than negative", "1", "-1", true},
		{"Positive greater than zero", "1", "0", true},
		{"Zero greater than negative", "0", "-1", true},
		{"Not greater - equal", "100", "100", false},
		{"Not greater - less than", "100", "200", false},
		{"Negative not greater than positive", "-1", "1", false},
		{"Negative not greater than zero", "-1", "0", false},
		{"Large number greater", "1000000000000000000000000000000", "999999999999999999999999999999", true},
		{"Large number not greater", "999999999999999999999999999999", "1000000000000000000000000000000", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a, _ := stdx.BigIntFromString(tt.a)
			b, _ := stdx.BigIntFromString(tt.b)
			result := IsGreaterThanInt(a, b)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestBigIntGreatorThanEqual(t *testing.T) {
	tests := []struct {
		name     string
		a        string
		b        string
		expected bool
	}{
		{"Greater positive numbers", "200", "100", true},
		{"Greater negative numbers", "-100", "-200", true},
		{"Equal positive numbers", "100", "100", true},
		{"Equal negative numbers", "-100", "-100", true},
		{"Equal zero", "0", "0", true},
		{"Positive greater than negative", "1", "-1", true},
		{"Positive greater than zero", "1", "0", true},
		{"Zero greater than negative", "0", "-1", true},
		{"Not greater or equal - less than", "100", "200", false},
		{"Negative not greater or equal to positive", "-1", "1", false},
		{"Negative not greater or equal to zero", "-1", "0", false},
		{"Large number greater or equal (greater)", "1000000000000000000000000000000", "999999999999999999999999999999", true},
		{"Large number greater or equal (equal)", "999999999999999999999999999999", "999999999999999999999999999999", true},
		{"Large number not greater or equal", "999999999999999999999999999999", "1000000000000000000000000000000", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a, _ := stdx.BigIntFromString(tt.a)
			b, _ := stdx.BigIntFromString(tt.b)
			result := IsGreaterThanOrEqualInt(a, b)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestBigIntLessThan(t *testing.T) {
	tests := []struct {
		name     string
		a        string
		b        string
		expected bool
	}{
		{"Less than positive numbers", "100", "200", true},
		{"Less than negative numbers", "-200", "-100", true},
		{"Negative less than positive", "-1", "1", true},
		{"Negative less than zero", "-1", "0", true},
		{"Zero less than positive", "0", "1", true},
		{"Not less than - equal", "100", "100", false},
		{"Not less than - greater", "200", "100", false},
		{"Positive not less than negative", "1", "-1", false},
		{"Zero not less than negative", "0", "-1", false},
		{"Large number less than", "999999999999999999999999999999", "1000000000000000000000000000000", true},
		{"Large number not less than", "1000000000000000000000000000000", "999999999999999999999999999999", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a, _ := stdx.BigIntFromString(tt.a)
			b, _ := stdx.BigIntFromString(tt.b)
			result := IsLessThanInt(a, b)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestBigIntLessThanEqual(t *testing.T) {
	tests := []struct {
		name     string
		a        string
		b        string
		expected bool
	}{
		{"Less than positive numbers", "100", "200", true},
		{"Less than negative numbers", "-200", "-100", true},
		{"Equal positive numbers", "100", "100", true},
		{"Equal negative numbers", "-100", "-100", true},
		{"Equal zero", "0", "0", true},
		{"Negative less than positive", "-1", "1", true},
		{"Negative less than zero", "-1", "0", true},
		{"Zero less than positive", "0", "1", true},
		{"Not less than or equal - greater", "200", "100", false},
		{"Positive not less than or equal to negative", "1", "-1", false},
		{"Zero not less than or equal to negative", "0", "-1", false},
		{"Large number less than or equal (less)", "999999999999999999999999999999", "1000000000000000000000000000000", true},
		{"Large number less than or equal (equal)", "999999999999999999999999999999", "999999999999999999999999999999", true},
		{"Large number not less than or equal", "1000000000000000000000000000000", "999999999999999999999999999999", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a, _ := stdx.BigIntFromString(tt.a)
			b, _ := stdx.BigIntFromString(tt.b)
			result := IsLessThanOrEqualInt(a, b)
			assert.Equal(t, tt.expected, result)
		})
	}
}
