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

// Gauge is a concurrency-safe struct that represents a single numerical value that can go up
// and go down. This implementation is optimized for a fast Set method.
//
// A Gauge is typically used for measured values like temperatures or current memory usage,
// but also "counts" that can go up and down, like the number of running goroutines.
//
// For Inc-heavy use case. Please consider using Counter instead.
type Gauge struct {
	// valBits contains the bits of the represented float64 value. It has
	// to go first in the struct to guarantee alignment for atomic
	// operations.
	// http://golang.org/pkg/sync/atomic/#pkg-note-BUG
	valBits uint64
}

// Return new Counter with the given initial value.
func NewGauge(initVal float64) *Gauge {
	g := &Gauge{}
	g.Set(initVal)
	return g
}

// Add the given value to the gauge.
func (g *Gauge) Add(v float64) {
	for {
		oldBits := atomic.LoadUint64(&g.valBits)
		newBits := math.Float64bits(math.Float64frombits(oldBits) + v)
		if atomic.CompareAndSwapUint64(&g.valBits, oldBits, newBits) {
			return
		}
	}
}

// Subtract the given value to the gauge.
func (g *Gauge) Sub(v float64) {
	g.Add(v * -1)
}

// Increment the gauge by 1.
func (g *Gauge) Inc() {
	g.Add(1)
}

// Decrement the gauge by 1.
func (g *Gauge) Dec() {
	g.Add(-1)
}

// Set an arbitary value to the gauge.
func (g *Gauge) Set(v float64) {
	atomic.StoreUint64(&g.valBits, math.Float64bits(v))
}

// Return the current value stored in the gauge.
func (g *Gauge) Value() float64 {
	return math.Float64frombits(atomic.LoadUint64(&g.valBits))
}
