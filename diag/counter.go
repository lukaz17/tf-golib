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
	"sync"
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
	// valInt stores values that are exact integers. valBitsSeq is odd while
	// valBits is being written and changes with every write. All three have
	// to go first in the struct to guarantee alignment for atomic operations.
	// http://golang.org/pkg/sync/atomic/#pkg-note-BUG
	valBits    uint64
	valInt     int64
	valBitsSeq uint64
	valBitsMu  sync.Mutex
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

	c.updateValBits(func(fval float64) float64 { return fval + v })
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
	c.updateValBits(func(float64) float64 { return v - float64(atomic.LoadInt64(&c.valInt)) })
}

// Return the current value stored in the counter.
func (c *Counter) Value() float64 {
	seq := atomic.LoadUint64(&c.valBitsSeq)
	fval := math.Float64frombits(atomic.LoadUint64(&c.valBits))
	ival := atomic.LoadInt64(&c.valInt)
	if seq%2 == 0 && atomic.LoadUint64(&c.valBitsSeq) == seq {
		return fval + float64(ival)
	}

	// A write to valBits overlapped the reads: waiting for it instead of retrying
	// keeps Value bounded under a stream of writes.
	c.valBitsMu.Lock()
	defer c.valBitsMu.Unlock()
	return math.Float64frombits(atomic.LoadUint64(&c.valBits)) + float64(atomic.LoadInt64(&c.valInt))
}

// Replace the float part of the counter with what update returns for it, in one write
// that Value never observes half-done. Read valInt inside update, not before calling:
// Value may run between an earlier read and the write.
func (c *Counter) updateValBits(update func(fval float64) float64) {
	c.valBitsMu.Lock()
	defer c.valBitsMu.Unlock()
	atomic.AddUint64(&c.valBitsSeq, 1)
	fval := math.Float64frombits(atomic.LoadUint64(&c.valBits))
	atomic.StoreUint64(&c.valBits, math.Float64bits(update(fval)))
	atomic.AddUint64(&c.valBitsSeq, 1)
}
