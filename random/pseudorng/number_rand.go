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

import "math/rand"

// Int returns a non-negative pseudo-random int from the default [Source].
func Int() int {
	return rand.Int()
}

// Intn returns, as an int, a non-negative pseudo-random number in the half-open interval [0,n)
// from the default [Source].
// It panics if n <= 0.
func Intn(n int) int {
	return rand.Intn(n)
}

// Intr returns, as an int, a non-negative pseudo-random number in the half-open interval [i,n)
// from the default [Source].
// It panics if n <= 0 or n <= i.
func Intr(i, n int) int {
	if n <= 0 || n <= i {
		panic("invalid argument to Intr")
	}
	return i + rand.Intn(n-i)
}

// Int31 returns a non-negative pseudo-random 31-bit integer as an int32
// from the default [Source].
func Int31() int32 {
	return rand.Int31()
}

// Int31n returns, as an int32, a non-negative pseudo-random number in the half-open interval [0,n)
// from the default [Source].
// It panics if n <= 0.
func Int31n(n int32) int32 {
	return rand.Int31n(n)
}

// Int31r returns, as an int32, a non-negative pseudo-random number in the half-open interval [i,n)
// from the default [Source].
// It panics if n <= 0 or n <= i.
func Int31r(i, n int32) int32 {
	if n <= 0 || n <= i {
		panic("invalid argument to Int31r")
	}
	return i + rand.Int31n(n-i)
}

// Int63 returns a non-negative pseudo-random 63-bit integer as an int64
// from the default [Source].
func Int63() int64 {
	return rand.Int63()
}

// Int63n returns, as an int64, a non-negative pseudo-random number in the half-open interval [0,n)
// from the default [Source].
// It panics if n <= 0.
func Int63n(n int64) int64 {
	return rand.Int63n(n)
}

// Int63r returns, as an int64, a non-negative pseudo-random number in the half-open interval [i,n)
// from the default [Source].
// It panics if n <= 0 or n <= i.
func Int63r(i, n int64) int64 {
	if n <= 0 || n <= i {
		panic("invalid argument to Int63r")
	}
	return i + rand.Int63n(n-i)
}

// Uint32 returns a pseudo-random 32-bit value as a uint32
// from the default [Source].
func Uint32() uint32 {
	return rand.Uint32()
}

// Uint32n returns, as an uint32, a pseudo-random number in the half-open interval [0,n)
// from the default [Source].
// It panics if n <= 0.
func Uint32n(n uint32) uint32 {
	if n <= 0 {
		panic("invalid argument to Uint32n")
	}
	uint32Max := ^uint32(0)
	if n&(n-1) == 0 { // n is power of two, can mask
		return rand.Uint32() & (n - 1)
	}
	max := uint32Max - uint32Max%n
	v := rand.Uint32()
	for v > max {
		v = rand.Uint32()
	}
	return v % n
}

// Uint32r returns, as an uint32, a pseudo-random number in the half-open interval [i,n)
// from the default [Source].
// It panics if n <= 0 or n <= i.
func Uint32r(i, n uint32) uint32 {
	if n <= 0 || n <= i {
		panic("invalid argument to Uint32r")
	}
	return i + Uint32n(n-i)
}

// Uint64 returns a pseudo-random 64-bit value as a uint64
// from the default [Source].
func Uint64() uint64 {
	return rand.Uint64()
}

// Uint64n returns, as an uint64, a pseudo-random number in the half-open interval [0,n)
// from the default [Source].
// It panics if n <= 0.
func Uint64n(n uint64) uint64 {
	if n <= 0 {
		panic("invalid argument to Uint64n")
	}
	uint64Max := ^uint64(0)
	if n&(n-1) == 0 { // n is power of two, can mask
		return rand.Uint64() & (n - 1)
	}
	max := uint64Max - uint64Max%n
	v := rand.Uint64()
	for v > max {
		v = rand.Uint64()
	}
	return v % n
}

// Uint64r returns, as an uint64, a pseudo-random number in the half-open interval [i,n)
// from the default [Source].
// It panics if n <= 0 or n <= i.
func Uint64r(i, n uint64) uint64 {
	if n <= 0 || n <= i {
		panic("invalid argument to Uint64r")
	}
	return i + Uint64n(n-i)
}

// Float32 returns, as a float32, a pseudo-random number in the half-open interval [0.0,1.0)
// from the default [Source].
func Float32() float32 {
	return rand.Float32()
}

// Float64 returns, as a float64, a pseudo-random number in the half-open interval [0.0,1.0)
// from the default [Source].
func Float64() float64 {
	return rand.Float64()
}

// NormFloat64 returns a normally distributed float64 in the range
// [-[math.MaxFloat64], +[math.MaxFloat64]] with
// standard normal distribution (mean = 0, stddev = 1)
// from the default [Source].
// To produce a different normal distribution, callers can
// adjust the output using:
//
//	sample = NormFloat64() * desiredStdDev + desiredMean
func NormFloat64() float64 {
	return rand.NormFloat64()
}

// NormFloat64n returns a float64 in the range
// [-[math.MaxFloat64], +[math.MaxFloat64]] with custom distribution (stddev, mean)
// from the default [Source].
// The custom distribution method is taken from Golang std library:
//
//	result = NormFloat64() * stddev + mean
func NormFloat64n(stddev, mean float64) float64 {
	return rand.NormFloat64()*stddev + mean
}

// ExpFloat64 returns an exponentially distributed float64 in the range
// (0, +[math.MaxFloat64]] with an exponential distribution whose rate parameter
// (lambda) is 1 and whose mean is 1/lambda (1) from the default [Source].
// To produce a distribution with a different rate parameter,
// callers can adjust the output using:
//
//	sample = ExpFloat64() / desiredRateParameter
func ExpFloat64() float64 {
	return rand.ExpFloat64()
}

// ExpFloat64 returns an exponentially distributed float64 in the range
// (0, +[math.MaxFloat64]] with an exponential distribution with custom
// rate parameter lambda from the default [Source].
// The custom distribution method is taken from Golang std library:
//
//	result = ExpFloat64() / lamda
func ExpFloat64n(lamda float64) float64 {
	return rand.ExpFloat64() / lamda
}
