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

// Struct Uint64ThreadSafe is a thread-safe counter with underlying value of uint64.
type Uint64ThreadSafe struct {
	value uint64
	lock  sync.Mutex
}

// Return current value exclusively.
func (c *Uint64ThreadSafe) Value() uint64 {
	c.lock.Lock()
	defer c.lock.Unlock()
	return c.value
}

// Return current value.
// Recommended to use with Lock() and Unlock() to prevent unexpected result.
// Maybe inconsistent if there are many routines updating its value at the same time.
func (c *Uint64ThreadSafe) ValueNoLock() uint64 {
	return c.value
}

// Set value eclusively.
func (c *Uint64ThreadSafe) Set(value uint64) {
	c.lock.Lock()
	defer c.lock.Unlock()
	c.value = value
}

// Set value.
// Recommended to use with Lock() and Unlock() to prevent unexpected result.
func (c *Uint64ThreadSafe) SetNoLock(value uint64) {
	c.value = value
}

// Add current value with n exclusively.
func (c *Uint64ThreadSafe) Add(n uint64) {
	c.lock.Lock()
	defer c.lock.Unlock()
	c.value += n
}

// AAdd current value with n.
// Recommended to use with Lock() and Unlock() to prevent unexpected result.
func (c *Uint64ThreadSafe) AddNoLock(n uint64) {
	c.value += n
}

// Subtract current value with n exclusively.
func (c *Uint64ThreadSafe) Sub(n uint64) {
	c.lock.Lock()
	defer c.lock.Unlock()
	c.value -= n
}

// Subtract current value with n.
// Recommended to use with Lock() and Unlock() to prevent unexpected result.
func (c *Uint64ThreadSafe) SubNoLock(n uint64) {
	c.value -= n
}

// Multiply current value with n exclusively.
func (c *Uint64ThreadSafe) Mul(n uint64) {
	c.lock.Lock()
	defer c.lock.Unlock()
	c.value *= n
}

// Multiply current value with n.
// Recommended to use with Lock() and Unlock() to prevent unexpected result.
func (c *Uint64ThreadSafe) MulNoLock(value uint64) {
	c.value *= value
}

// Divide current value with n exclusively.
func (c *Uint64ThreadSafe) Div(n uint64) {
	c.lock.Lock()
	defer c.lock.Unlock()
	c.value /= n
}

// Divide current value with n.
// Recommended to use with Lock() and Unlock() to prevent unexpected result.
func (c *Uint64ThreadSafe) DivNoLock(value uint64) {
	c.value /= value
}

// Acquire the lock manually.
func (c *Uint64ThreadSafe) Lock() {
	c.lock.Lock()
}

// Release the lock manually.
func (c *Uint64ThreadSafe) Unlock() {
	c.lock.Unlock()
}
