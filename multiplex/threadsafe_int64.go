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

import "sync"

// Struct Int64ThreadSafe is a thread-safe counter with underlying value of int64.
type Int64ThreadSafe struct {
	value int64
	lock  sync.Mutex
}

// Return current value exclusively.
func (c *Int64ThreadSafe) Value() int64 {
	c.lock.Lock()
	defer c.lock.Unlock()
	return c.value
}

// Return current value.
// Recommended to use with Lock() and Unlock() to prevent unexpected result.
// Maybe inconsistent if there are many routines updating its value at the same time.
func (c *Int64ThreadSafe) ValueNoLock() int64 {
	return c.value
}

// Set value eclusively.
func (c *Int64ThreadSafe) Set(value int64) {
	c.lock.Lock()
	defer c.lock.Unlock()
	c.value = value
}

// Set value.
// Recommended to use with Lock() and Unlock() to prevent unexpected result.
func (c *Int64ThreadSafe) SetNoLock(value int64) {
	c.value = value
}

// Add current value with n exclusively.
func (c *Int64ThreadSafe) Add(n int64) {
	c.lock.Lock()
	defer c.lock.Unlock()
	c.value += n
}

// AAdd current value with n.
// Recommended to use with Lock() and Unlock() to prevent unexpected result.
func (c *Int64ThreadSafe) AddNoLock(n int64) {
	c.value += n
}

// Subtract current value with n exclusively.
func (c *Int64ThreadSafe) Sub(n int64) {
	c.lock.Lock()
	defer c.lock.Unlock()
	c.value -= n
}

// Subtract current value with n.
// Recommended to use with Lock() and Unlock() to prevent unexpected result.
func (c *Int64ThreadSafe) SubNoLock(n int64) {
	c.value -= n
}

// Multiply current value with n exclusively.
func (c *Int64ThreadSafe) Mul(n int64) {
	c.lock.Lock()
	defer c.lock.Unlock()
	c.value *= n
}

// Multiply current value with n.
// Recommended to use with Lock() and Unlock() to prevent unexpected result.
func (c *Int64ThreadSafe) MulNoLock(value int64) {
	c.value *= value
}

// Divide current value with n exclusively.
func (c *Int64ThreadSafe) Div(n int64) {
	c.lock.Lock()
	defer c.lock.Unlock()
	c.value /= n
}

// Divide current value with n.
// Recommended to use with Lock() and Unlock() to prevent unexpected result.
func (c *Int64ThreadSafe) DivNoLock(value int64) {
	c.value /= value
}

// Acquire the lock manually.
func (c *Int64ThreadSafe) Lock() {
	c.lock.Lock()
}

// Release the lock manually.
func (c *Int64ThreadSafe) Unlock() {
	c.lock.Unlock()
}
