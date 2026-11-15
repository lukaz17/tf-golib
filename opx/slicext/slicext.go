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

// Compare two slices x and y of comparable types that is deep equals.
func AreEqual[S ~[]T, T comparable](x, y S) bool {
	if x == nil && y == nil {
		return true
	}
	if x == nil || y == nil {
		return false
	}
	if len(x) != len(y) {
		return false
	}
	for i := range x {
		if x[i] != y[i] {
			return false
		}
	}
	return true
}

// Compare two slices x and y of any type that is deep equals.
func AreEqualFunc[S ~[]T, T any](x, y S, equalFunc func(x, y T) bool) bool {
	if x == nil && y == nil {
		return true
	}
	if x == nil || y == nil {
		return false
	}
	if len(x) != len(y) {
		return false
	}
	for i := range x {
		if !equalFunc(x[i], y[i]) {
			return false
		}
	}
	return true
}

// Check whether a value v is a member of slice s.
func Contains[S ~[]T, T comparable](s S, v T) bool {
	for i := range s {
		if v == s[i] {
			return true
		}
	}
	return false
}

// Check whether a value v is a member of slice s.
func ContainsFunc[S ~[]T, T comparable](s S, v T, equalFunc func(x, y T) bool) bool {
	for i := range s {
		if equalFunc(v, s[i]) {
			return true
		}
	}
	return false
}

// Check whether slice s is Nil nor zero in length.
func IsEmpty[S ~[]T, T any](s S) bool {
	if s == nil {
		return true
	}
	return len(s) == 0
}

// Map a slice contains items of type T to a new slice contains equivalent number items of type N.
// Nil handling and special conversion will be taken care of in mapFunc.
func Map[S ~[]T, M ~[]N, T, N any](s S, mapFunc func(t T) N) M {
	m := make([]N, len(s))
	for i, t := range s {
		m[i] = mapFunc(t)
	}
	return m
}

// Map a slice contains items of type T to a Key-Value map.
// keyFunc must returns extra bool to indicate whether slice item T will be included in the resulted map.
func MapKV[S ~[]T, K comparable, V, T any](s S, keyFunc func(t T) (K, bool), valFunc func(t T) V) map[K]V {
	m := make(map[K]V)
	for _, t := range s {
		if k, ok := keyFunc(t); ok {
			m[k] = valFunc(t)
		}
	}
	return m
}
