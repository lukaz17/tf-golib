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

package diag

import (
	"math"
	"sync/atomic"
)

// Counter is a concurrency-safe struct that represents a single numerical value that can go up
// and go down. This implementation is optimized for a fast Inc/Dec/Add/Sub methods.
//
// A Counter is typically used to count requests served, tasks completed, errors occurred, etc.
//
// For Set-heavy use case. Please consider using Gauge instead.
type Counter struct {
	// valBits contains the bits of the represented float64 value, while
	// valInt stores values that are exact integers. Both have to go first
	// in the struct to guarantee alignment for atomic operations.
	// http://golang.org/pkg/sync/atomic/#pkg-note-BUG
	valBits uint64
	valInt  int64
}

// Return new Counter with the given initial value.
func NewCounter(initVal float64) *Counter {
	c := &Counter{}
	c.Add(initVal)
	return c
}

// Add the given value to the counter.
func (c *Counter) Add(v float64) {
	ival := int64(v)
	if float64(ival) == v {
		atomic.AddInt64(&c.valInt, ival)
		return
	}

	for {
		oldBits := atomic.LoadUint64(&c.valBits)
		newBits := math.Float64bits(math.Float64frombits(oldBits) + v)
		if atomic.CompareAndSwapUint64(&c.valBits, oldBits, newBits) {
			return
		}
	}
}

// Subtract the given value to the counter.
func (c *Counter) Sub(v float64) {
	c.Add(v * -1)
}

// Increment the counter by 1.
func (c *Counter) Inc() {
	c.Add(1)
}

// Decrement the counter by 1.
func (c *Counter) Dec() {
	c.Add(-1)
}

// Set an arbitary value to the counter.
func (c *Counter) Set(v float64) {
	for {
		oldBits := atomic.LoadUint64(&c.valBits)
		ival := atomic.LoadInt64(&c.valInt)
		newBits := math.Float64bits(v - float64(ival))
		if atomic.CompareAndSwapUint64(&c.valBits, oldBits, newBits) {
			return
		}
	}
}

// Return the current value stored in the counter.
func (c *Counter) Value() float64 {
	fval := math.Float64frombits(atomic.LoadUint64(&c.valBits))
	ival := atomic.LoadInt64(&c.valInt)
	return fval + float64(ival)
}
