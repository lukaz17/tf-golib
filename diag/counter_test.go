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
	"sync/atomic"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewCounter(t *testing.T) {
	// Test with zero value
	c := NewCounter(0)
	assert.Equal(t, 0.0, c.Value())

	// Test with positive integer
	c = NewCounter(10)
	assert.Equal(t, 10.0, c.Value())

	// Test with negative integer
	c = NewCounter(-5)
	assert.Equal(t, -5.0, c.Value())

	// Test with float
	c = NewCounter(3.14)
	assert.InDelta(t, 3.14, c.Value(), 0.001)

	// Test with negative float
	c = NewCounter(-2.71)
	assert.InDelta(t, -2.71, c.Value(), 0.001)
}

func TestCounter_Add(t *testing.T) {
	c := NewCounter(0)

	// Test adding integers
	c.Add(5)
	assert.Equal(t, 5.0, c.Value())

	c.Add(10)
	assert.Equal(t, 15.0, c.Value())

	// Test adding floats
	c.Add(2.5)
	assert.InDelta(t, 17.5, c.Value(), 0.001)

	// Test adding negative values
	c.Add(-3)
	assert.InDelta(t, 14.5, c.Value(), 0.001)

	c.Add(-1.5)
	assert.InDelta(t, 13.0, c.Value(), 0.001)
}

func TestCounter_Sub(t *testing.T) {
	c := NewCounter(20)

	// Test subtracting integers
	c.Sub(5)
	assert.Equal(t, 15.0, c.Value())

	// Test subtracting floats
	c.Sub(2.5)
	assert.InDelta(t, 12.5, c.Value(), 0.001)

	// Test subtracting negative values (should add)
	c.Sub(-3)
	assert.InDelta(t, 15.5, c.Value(), 0.001)
}

func TestCounter_Inc(t *testing.T) {
	c := NewCounter(0)

	c.Inc()
	assert.Equal(t, 1.0, c.Value())

	c.Inc()
	assert.Equal(t, 2.0, c.Value())

	// Test with initial float value
	c = NewCounter(2.5)
	c.Inc()
	assert.InDelta(t, 3.5, c.Value(), 0.001)
}

func TestCounter_Dec(t *testing.T) {
	c := NewCounter(5)

	c.Dec()
	assert.Equal(t, 4.0, c.Value())

	c.Dec()
	assert.Equal(t, 3.0, c.Value())

	// Test with initial float value
	c = NewCounter(2.5)
	c.Dec()
	assert.InDelta(t, 1.5, c.Value(), 0.001)
}

func TestCounter_Set(t *testing.T) {
	c := NewCounter(7.5)

	// Test setting integer
	c.Set(25)
	assert.Equal(t, 25.0, c.Value())

	// Test setting float
	c.Set(15.75)
	assert.InDelta(t, 15.75, c.Value(), 0.001)

	// Test setting zero
	c.Set(0)
	assert.Equal(t, 0.0, c.Value())

	// Test setting negative
	c.Set(-10.5)
	assert.InDelta(t, -10.5, c.Value(), 0.001)
}

func TestCounter_Value(t *testing.T) {
	c := NewCounter(0)

	// Test initial value
	assert.Equal(t, 0.0, c.Value())

	// Test after operations
	c.Add(10)
	c.Sub(3)
	c.Inc()
	c.Dec()
	assert.Equal(t, 7.0, c.Value())

	// Test with mixed integer and float operations
	c.Set(0)
	c.Add(5)    // integer
	c.Add(2.5)  // float
	c.Sub(1.25) // float
	assert.InDelta(t, 6.25, c.Value(), 0.001)
}

func TestCounter_ConcurrentOperations(t *testing.T) {
	c := NewCounter(0)
	numGoroutines := 100
	numOperations := 100 // At 1000, go1.18's race detector spends ~1s on the contended atomics.

	var wg sync.WaitGroup
	wg.Add(numGoroutines * 4) // 4 types of operations

	// Concurrent Inc operations
	for i := 0; i < numGoroutines; i++ {
		go func() {
			defer wg.Done()
			for j := 0; j < numOperations; j++ {
				c.Inc()
			}
		}()
	}

	// Concurrent Dec operations
	for i := 0; i < numGoroutines; i++ {
		go func() {
			defer wg.Done()
			for j := 0; j < numOperations; j++ {
				c.Dec()
			}
		}()
	}

	// Concurrent Add operations
	for i := 0; i < numGoroutines; i++ {
		go func() {
			defer wg.Done()
			for j := 0; j < numOperations; j++ {
				c.Add(1)
			}
		}()
	}

	// Concurrent Sub operations
	for i := 0; i < numGoroutines; i++ {
		go func() {
			defer wg.Done()
			for j := 0; j < numOperations; j++ {
				c.Sub(1)
			}
		}()
	}

	wg.Wait()

	// Since we have equal Inc/Dec and Add/Sub operations, final value should be 0
	assert.Equal(t, 0.0, c.Value())
}

func TestCounter_ConcurrentReadWrite(t *testing.T) {
	c := NewCounter(1000)
	done := make(chan bool)

	// Reader goroutines
	for i := 0; i < 10; i++ {
		go func() {
			for {
				select {
				case <-done:
					return
				default:
					c.Value() // Just read, don't assert specific value due to concurrent writes
				}
			}
		}()
	}

	// Writer goroutines
	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < 100; j++ {
				c.Inc()
				c.Dec()
				c.Add(0.5)
				c.Sub(0.5)
			}
		}()
	}

	// Let it run for a bit
	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := 0; i < 1000; i++ {
			c.Value()
		}
	}()

	wg.Wait()
	close(done)

	// Final value should be close to initial (1000) since Inc/Dec and Add/Sub cancel out
	assert.InDelta(t, 1000.0, c.Value(), 1.0)
}

func TestCounter_Value_ReturnsOnlyHeldValuesDuringWrites(t *testing.T) {
	tests := []struct {
		name  string
		cycle func(c *Counter)
		held  map[float64]bool
	}{
		{
			name:  "set between increments",
			cycle: func(c *Counter) { c.Set(0); c.Inc() },
			held:  map[float64]bool{0: true, 1: true},
		},
		{
			name:  "fractional adds between whole adds",
			cycle: func(c *Counter) { c.Add(0.5); c.Add(-0.5); c.Inc(); c.Dec() },
			held:  map[float64]bool{0: true, 0.5: true, 1: true},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			const writerCycles = 100000
			const readerCount = 4
			c := NewCounter(0)
			var unheldReads int64
			done := make(chan struct{})
			var readers sync.WaitGroup
			for i := 0; i < readerCount; i++ {
				readers.Add(1)
				go func() {
					defer readers.Done()
					for {
						select {
						case <-done:
							return
						default:
							if !tt.held[c.Value()] {
								atomic.AddInt64(&unheldReads, 1)
							}
						}
					}
				}()
			}

			for i := 0; i < writerCycles; i++ {
				tt.cycle(c)
			}
			close(done)
			readers.Wait()

			assert.Zero(t, unheldReads, "reads of a value the counter never held")
		})
	}
}

func TestCounter_EdgeCases(t *testing.T) {
	// Test very large numbers
	c := NewCounter(0)
	largeNum := 1e15
	c.Add(largeNum)
	assert.InDelta(t, largeNum, c.Value(), 1e10)

	// Test very small numbers
	c.Set(0)
	smallNum := 1e-10
	c.Add(smallNum)
	assert.InDelta(t, smallNum, c.Value(), 1e-10)

	// Test mixed large and small operations
	c.Set(1e15)
	c.Add(1e-10)
	c.Sub(1e-10)
	assert.InDelta(t, 1e15, c.Value(), 1e10)
}

func TestCounter_ZeroOperations(t *testing.T) {
	c := NewCounter(5)

	// Adding zero should not change value
	c.Add(0)
	assert.Equal(t, 5.0, c.Value())

	// Subtracting zero should not change value
	c.Sub(0)
	assert.Equal(t, 5.0, c.Value())

	// Setting to current value
	c.Set(5)
	assert.Equal(t, 5.0, c.Value())
}
