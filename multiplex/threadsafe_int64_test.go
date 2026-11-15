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

package multiplex

import (
	"sync"
	"testing"
)

func TestInt64ThreadSafe_SetAndValue(t *testing.T) {
	counter := &Int64ThreadSafe{}

	counter.Set(10)
	if counter.Value() != 10 {
		t.Errorf("Expected value to be 10, got %d", counter.Value())
	}
}

func TestInt64ThreadSafe_Add(t *testing.T) {
	counter := &Int64ThreadSafe{}

	counter.Set(10)
	counter.Add(5)
	if counter.Value() != 15 {
		t.Errorf("Expected value to be 15, got %d", counter.Value())
	}
}

func TestInt64ThreadSafe_Sub(t *testing.T) {
	counter := &Int64ThreadSafe{}

	counter.Set(15)
	counter.Sub(3)
	if counter.Value() != 12 {
		t.Errorf("Expected value to be 12, got %d", counter.Value())
	}
}

func TestInt64ThreadSafe_Mul(t *testing.T) {
	counter := &Int64ThreadSafe{}

	counter.Set(12)
	counter.Mul(2)
	if counter.Value() != 24 {
		t.Errorf("Expected value to be 24, got %d", counter.Value())
	}
}

func TestInt64ThreadSafe_Div(t *testing.T) {
	counter := &Int64ThreadSafe{}

	counter.Set(24)
	counter.Div(4)
	if counter.Value() != 6 {
		t.Errorf("Expected value to be 6, got %d", counter.Value())
	}
}

func TestInt64ThreadSafe_ConcurrentAdd(t *testing.T) {
	var wg sync.WaitGroup
	counter := &Int64ThreadSafe{}

	counter.Set(6)
	wg.Add(2)
	go func() {
		defer wg.Done()
		counter.Add(10)
	}()
	go func() {
		defer wg.Done()
		counter.Add(20)
	}()
	wg.Wait()
	if counter.Value() != 36 {
		t.Errorf("Expected value to be 36 after concurrent Add, got %d", counter.Value())
	}
}

func TestInt64ThreadSafe_LockAndUnlock(t *testing.T) {
	counter := &Int64ThreadSafe{}

	counter.Lock()
	go func() {
		defer counter.Unlock()
		counter.SetNoLock(50)
	}()
	counter.Lock()
	if counter.ValueNoLock() != 50 {
		t.Errorf("Expected value to be 50 after manual lock/unlock, got %d", counter.Value())
	}
	counter.Unlock()
}
