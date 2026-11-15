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

import "math"

/* --- FLOAT32 --- */

// Return the absolute value of x for float32.
func AbsFloat32(x float32) float32 {
	return float32(math.Abs(float64(x)))
}

// Return the smaller of x or y for float32.
func MinFloat32(x float32, y ...float32) float32 {
	min := x
	for _, v := range y {
		if v < min {
			min = v
		}
	}
	return min
}

// Return the larger of x or y for float32.
func MaxFloat32(x float32, y ...float32) float32 {
	max := x
	for _, v := range y {
		if v > max {
			max = v
		}
	}
	return max
}

/* --- FLOAT64 --- */

// Return the absolute value of x for float64.
func AbsFloat64(x float64) float64 {
	return math.Abs(x)
}

// Return the smaller of x or y for float64.
func MinFloat64(x float64, y ...float64) float64 {
	min := x
	for _, v := range y {
		if v < min {
			min = v
		}
	}
	return min
}

// Return the larger of x or y for float64.
func MaxFloat64(x float64, y ...float64) float64 {
	max := x
	for _, v := range y {
		if v > max {
			max = v
		}
	}
	return max
}
