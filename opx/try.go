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

package opx

import "errors"

// Ensure that the ok is true, otherwise panic.
func Try(ok bool) {
	if !ok {
		panic(errors.New("tried and ok is false"))
	}
}

// Ensure that the ok is true then return all preceding variables, otherwise panic.
// Used to eliminate ok check for functions.
func Try1[T1 any](t1 T1, ok bool) T1 {
	if !ok {
		panic(errors.New("tried and ok is false"))
	}
	return t1
}

// Ensure that the ok is true then return all preceding variables, otherwise panic.
// Used to eliminate ok check for functions.
func Try2[T1 any, T2 any](t1 T1, t2 T2, ok bool) (T1, T2) {
	if !ok {
		panic(errors.New("tried and ok is false"))
	}
	return t1, t2
}

// Ensure that the ok is true then return all preceding variables, otherwise panic.
// Used to eliminate ok check for functions.
func Try3[T1 any, T2 any, T3 any](t1 T1, t2 T2, t3 T3, ok bool) (T1, T2, T3) {
	if !ok {
		panic(errors.New("tried and ok is false"))
	}
	return t1, t2, t3
}

// Ensure that the ok is true then return all preceding variables, otherwise panic.
// Used to eliminate ok check for functions.
func Try4[T1 any, T2 any, T3 any, T4 any](t1 T1, t2 T2, t3 T3, t4 T4, ok bool) (T1, T2, T3, T4) {
	if !ok {
		panic(errors.New("tried and ok is false"))
	}
	return t1, t2, t3, t4
}

// Ensure that the ok is true then return all preceding variables, otherwise panic.
// Used to eliminate ok check for functions.
func Try5[T1 any, T2 any, T3 any, T4 any, T5 any](t1 T1, t2 T2, t3 T3, t4 T4, t5 T5, ok bool) (T1, T2, T3, T4, T5) {
	if !ok {
		panic(errors.New("tried and ok is false"))
	}
	return t1, t2, t3, t4, t5
}
