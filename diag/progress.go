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
	"time"
)

// Progress is a helper type to track progress.
type Progress struct {
	curVal *Counter
	ttlVal *Counter

	started time.Time

	updated   time.Time
	updatedMu sync.Mutex
}

// Return new Progress tracker with total items to complete.
// It panics unless total is positive and finite.
func NewProgress(total float64) *Progress {
	if !isPositiveFinite(total) {
		panic("total must be positive and finite")
	}
	return &Progress{
		curVal:  NewCounter(0),
		ttlVal:  NewCounter(total),
		started: time.Now(),
		updated: time.Now(),
	}
}

// Add v number of items as to be completed.
// It panics unless v is positive and finite.
func (p *Progress) Add(v float64) {
	if !isPositiveFinite(v) {
		panic("v must be positive and finite")
	}
	p.ttlVal.Add(v)
}

// Mark v number of items as completed.
// It panics unless v is positive and finite.
func (p *Progress) Complete(v float64) {
	if !isPositiveFinite(v) {
		panic("v must be positive and finite")
	}
	p.curVal.Add(v)
	p.updatedMu.Lock()
	p.updated = time.Now()
	p.updatedMu.Unlock()
}

// Estimate when the Progress will be completed.
// A projected duration beyond the longest time.Duration (~292 years) is clamped to it;
// so is one left undefined by amounts summing past float64's range.
func (p *Progress) EstimatedTime() time.Time {
	cur := p.curVal.Value()
	if cur == 0 {
		return time.Now().AddDate(999, 0, 0)
	}
	passed := time.Since(p.started)
	timeToComplete := clampedDuration(float64(passed) / cur * p.ttlVal.Value())
	return p.started.Add(timeToComplete)
}

func (p *Progress) Percent() float64 {
	return (p.curVal.Value() / p.ttlVal.Value()) * 100
}

// Estimate how long the Progress will be completed.
// A projection beyond the longest time.Duration (~292 years) is clamped to it;
// so is one left undefined by amounts summing past float64's range.
func (p *Progress) RemainTime() time.Duration {
	cur := p.curVal.Value()
	if cur == 0 {
		return 0
	}
	passed := time.Since(p.started)
	remain := p.ttlVal.Value() - cur
	remainTime := clampedDuration(float64(passed) / cur * remain)
	return remainTime
}

// Return current value, total value and last updated time.
func (p *Progress) Value() (float64, float64, time.Time) {
	return p.curVal.Value(), p.ttlVal.Value(), p.lastUpdated()
}

// Return the time of the last Complete call, or the creation time if there is none.
func (p *Progress) lastUpdated() time.Time {
	p.updatedMu.Lock()
	defer p.updatedMu.Unlock()
	return p.updated
}

// Return the Duration of a projected number of nanoseconds, or the longest Duration for a projection beyond it or undefined (NaN).
func clampedDuration(ns float64) time.Duration {
	// float64(math.MaxInt64) rounds up to 2^63, the first value out of range.
	// Projections are never below minus the elapsed time, so only the upper end can overflow.
	// NaN, from amounts summing past float64's range, fails every comparison, so it is checked apart.
	if ns >= float64(math.MaxInt64) || math.IsNaN(ns) {
		return math.MaxInt64
	}
	return time.Duration(ns)
}

// Report whether v is an amount Progress accepts: above 0, and neither infinite nor NaN.
func isPositiveFinite(v float64) bool {
	// v > 0 is false for NaN, which v <= 0 would let through.
	return v > 0 && !math.IsInf(v, 1)
}
