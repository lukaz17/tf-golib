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

/* --- INT --- */

// Return the absolute value of x for int.
func AbsInt(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

// Return the smaller of x or y for int.
func MinInt(x int, y ...int) int {
	min := x
	for _, v := range y {
		if v < min {
			min = v
		}
	}
	return min
}

// Return the larger of x or y for int.
func MaxInt(x int, y ...int) int {
	max := x
	for _, v := range y {
		if v > max {
			max = v
		}
	}
	return max
}

/* --- INT8 --- */

// Return the absolute value of x for int8.
func AbsInt8(x int8) int8 {
	if x < 0 {
		return -x
	}
	return x
}

// Return the smaller of x or y for int8.
func MinInt8(x int8, y ...int8) int8 {
	min := x
	for _, v := range y {
		if v < min {
			min = v
		}
	}
	return min
}

// Return the larger of x or y for int8.
func MaxInt8(x int8, y ...int8) int8 {
	max := x
	for _, v := range y {
		if v > max {
			max = v
		}
	}
	return max
}

/* --- INT16 --- */

// Return the absolute value of x for int16.
func AbsInt16(x int16) int16 {
	if x < 0 {
		return -x
	}
	return x
}

// Return the smaller of x or y for int16.
func MinInt16(x int16, y ...int16) int16 {
	min := x
	for _, v := range y {
		if v < min {
			min = v
		}
	}
	return min
}

// Return the larger of x or y for int16.
func MaxInt16(x int16, y ...int16) int16 {
	max := x
	for _, v := range y {
		if v > max {
			max = v
		}
	}
	return max
}

/* --- INT32 --- */

// Return the absolute value of x for int32.
func AbsInt32(x int32) int32 {
	if x < 0 {
		return -x
	}
	return x
}

// Return the smaller of x or y for int32.
func MinInt32(x int32, y ...int32) int32 {
	min := x
	for _, v := range y {
		if v < min {
			min = v
		}
	}
	return min
}

// Return the larger of x or y for int32.
func MaxInt32(x int32, y ...int32) int32 {
	max := x
	for _, v := range y {
		if v > max {
			max = v
		}
	}
	return max
}

/* --- INT64 --- */

// Return the absolute value of x for int64.
func AbsInt64(x int64) int64 {
	if x < 0 {
		return -x
	}
	return x
}

// Return the smaller of x or y for int64.
func MinInt64(x int64, y ...int64) int64 {
	min := x
	for _, v := range y {
		if v < min {
			min = v
		}
	}
	return min
}

// Return the larger of x or y for int64.
func MaxInt64(x int64, y ...int64) int64 {
	max := x
	for _, v := range y {
		if v > max {
			max = v
		}
	}
	return max
}
