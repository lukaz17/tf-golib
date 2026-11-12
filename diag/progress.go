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
func NewProgress(total float64) *Progress {
	if total <= 0 {
		panic("total must be positive")
	}
	return &Progress{
		curVal:  NewCounter(0),
		ttlVal:  NewCounter(total),
		started: time.Now(),
		updated: time.Now(),
	}
}

// Add v number of items as to be completed.
func (p *Progress) Add(v float64) {
	if v <= 0 {
		panic("v must be positive")
	}
	p.ttlVal.Add(v)
}

// Mark v number of items as completed.
func (p *Progress) Complete(v float64) {
	if v <= 0 {
		panic("v must be positive")
	}
	p.curVal.Add(v)
	p.updatedMu.Lock()
	p.updated = time.Now()
	p.updatedMu.Unlock()
}

// Estimate when the Progress will be completed.
func (p *Progress) EstimatedTime() time.Time {
	if p.curVal.Value() == 0 {
		return time.Now().AddDate(999, 0, 0)
	}
	passed := time.Since(p.started)
	timeToComplete := passed / time.Duration(p.curVal.Value()) * time.Duration(p.ttlVal.Value())
	return p.lastUpdated().Add(timeToComplete)
}

func (p *Progress) Percent() float64 {
	return (p.curVal.Value() / p.ttlVal.Value()) * 100
}

// Estimate how long the Progress will be completed.
func (p *Progress) RemainTime() time.Duration {
	if p.curVal.Value() == 0 {
		return 0
	}
	passed := time.Since(p.started)
	remain := p.ttlVal.Value() - p.curVal.Value()
	remainTime := passed / time.Duration(p.curVal.Value()) * time.Duration(remain)
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
