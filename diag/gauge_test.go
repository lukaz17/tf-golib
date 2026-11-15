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
	"runtime"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewGauge(t *testing.T) {
	// Test with zero value
	g := NewGauge(0)
	assert.Equal(t, 0.0, g.Value())

	// Test with positive integer
	g = NewGauge(10)
	assert.Equal(t, 10.0, g.Value())

	// Test with negative integer
	g = NewGauge(-5)
	assert.Equal(t, -5.0, g.Value())

	// Test with float
	g = NewGauge(3.14)
	assert.InDelta(t, 3.14, g.Value(), 0.001)

	// Test with negative float
	g = NewGauge(-2.71)
	assert.InDelta(t, -2.71, g.Value(), 0.001)
}

func TestGauge_Add(t *testing.T) {
	g := NewGauge(0)

	// Test adding integers
	g.Add(5)
	assert.Equal(t, 5.0, g.Value())

	g.Add(10)
	assert.Equal(t, 15.0, g.Value())

	// Test adding floats
	g.Add(2.5)
	assert.InDelta(t, 17.5, g.Value(), 0.001)

	// Test adding negative values
	g.Add(-3)
	assert.InDelta(t, 14.5, g.Value(), 0.001)

	g.Add(-1.5)
	assert.InDelta(t, 13.0, g.Value(), 0.001)
}

func TestGauge_Sub(t *testing.T) {
	g := NewGauge(20)

	// Test subtracting integers
	g.Sub(5)
	assert.Equal(t, 15.0, g.Value())

	// Test subtracting floats
	g.Sub(2.5)
	assert.InDelta(t, 12.5, g.Value(), 0.001)

	// Test subtracting negative values (should add)
	g.Sub(-3)
	assert.InDelta(t, 15.5, g.Value(), 0.001)
}

func TestGauge_Inc(t *testing.T) {
	g := NewGauge(0)

	g.Inc()
	assert.Equal(t, 1.0, g.Value())

	g.Inc()
	assert.Equal(t, 2.0, g.Value())

	// Test with initial float value
	g = NewGauge(2.5)
	g.Inc()
	assert.InDelta(t, 3.5, g.Value(), 0.001)
}

func TestGauge_Dec(t *testing.T) {
	g := NewGauge(5)

	g.Dec()
	assert.Equal(t, 4.0, g.Value())

	g.Dec()
	assert.Equal(t, 3.0, g.Value())

	// Test with initial float value
	g = NewGauge(2.5)
	g.Dec()
	assert.InDelta(t, 1.5, g.Value(), 0.001)
}

func TestGauge_Set(t *testing.T) {
	g := NewGauge(7.5)

	// Test setting integer
	g.Set(25)
	assert.Equal(t, 25.0, g.Value())

	// Test setting float
	g.Set(15.75)
	assert.InDelta(t, 15.75, g.Value(), 0.001)

	// Test setting zero
	g.Set(0)
	assert.Equal(t, 0.0, g.Value())

	// Test setting negative
	g.Set(-10.5)
	assert.InDelta(t, -10.5, g.Value(), 0.001)
}

func TestGauge_Value(t *testing.T) {
	g := NewGauge(0)

	// Test initial value
	assert.Equal(t, 0.0, g.Value())

	// Test after operations
	g.Add(10)
	g.Sub(3)
	g.Inc()
	g.Dec()
	assert.Equal(t, 7.0, g.Value())

	// Test with mixed integer and float operations
	g.Set(0)
	g.Add(5)    // integer
	g.Add(2.5)  // float
	g.Sub(1.25) // float
	assert.InDelta(t, 6.25, g.Value(), 0.001)
}

func TestGauge_ConcurrentOperations(t *testing.T) {
	g := NewGauge(0)
	numGoroutines := 100
	numOperations := 100 // At 1000, go1.18's race detector spends ~8s on the contended atomics.

	var wg sync.WaitGroup
	wg.Add(numGoroutines * 4) // 4 types of operations

	// Concurrent Inc operations
	for i := 0; i < numGoroutines; i++ {
		go func() {
			defer wg.Done()
			for j := 0; j < numOperations; j++ {
				g.Inc()
			}
		}()
	}

	// Concurrent Dec operations
	for i := 0; i < numGoroutines; i++ {
		go func() {
			defer wg.Done()
			for j := 0; j < numOperations; j++ {
				g.Dec()
			}
		}()
	}

	// Concurrent Add operations
	for i := 0; i < numGoroutines; i++ {
		go func() {
			defer wg.Done()
			for j := 0; j < numOperations; j++ {
				g.Add(1)
			}
		}()
	}

	// Concurrent Sub operations
	for i := 0; i < numGoroutines; i++ {
		go func() {
			defer wg.Done()
			for j := 0; j < numOperations; j++ {
				g.Sub(1)
			}
		}()
	}

	wg.Wait()

	// Since we have equal Inc/Dec and Add/Sub operations, final value should be 0
	assert.Equal(t, 0.0, g.Value())
}

func TestGauge_ConcurrentReadWrite(t *testing.T) {
	g := NewGauge(1000)
	done := make(chan bool)

	// Reader goroutines
	for i := 0; i < 10; i++ {
		go func() {
			for {
				select {
				case <-done:
					return
				default:
					// Just read, don't assert specific value due to concurrent writes
					g.Value()
					// Without yielding, spinning readers slow the writers' CAS loops for seconds under -race.
					runtime.Gosched()
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
				g.Inc()
				g.Dec()
				g.Add(0.5)
				g.Sub(0.5)
			}
		}()
	}

	// Let it run for a bit
	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := 0; i < 1000; i++ {
			g.Value()
		}
	}()

	wg.Wait()
	close(done)

	// Final value should be close to initial (1000) since Inc/Dec and Add/Sub cancel out
	assert.InDelta(t, 1000.0, g.Value(), 1.0)
}

func TestGauge_EdgeCases(t *testing.T) {
	// Test very large numbers
	g := NewGauge(0)
	largeNum := 1e15
	g.Add(largeNum)
	assert.InDelta(t, largeNum, g.Value(), 1e10)

	// Test very small numbers
	g.Set(0)
	smallNum := 1e-10
	g.Add(smallNum)
	assert.InDelta(t, smallNum, g.Value(), 1e-10)

	// Test mixed large and small operations
	g.Set(1e15)
	g.Add(1e-10)
	g.Sub(1e-10)
	assert.InDelta(t, 1e15, g.Value(), 1e10)
}

func TestGauge_ZeroOperations(t *testing.T) {
	g := NewGauge(5)

	// Adding zero should not change value
	g.Add(0)
	assert.Equal(t, 5.0, g.Value())

	// Subtracting zero should not change value
	g.Sub(0)
	assert.Equal(t, 5.0, g.Value())

	// Setting to current value
	g.Set(5)
	assert.Equal(t, 5.0, g.Value())
}
