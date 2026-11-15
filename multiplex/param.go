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
	"errors"
	"sync"
	"sync/atomic"
)

// ErrDropped is the result returned to a sender waiting for a request that was dropped instead of handled: dispatched without a router, sent to no service or to an unregistered one, or left unhandled by the receiving service's hook.
var ErrDropped = errors.New("multiplex: request dropped")

// ServiceMessage defines a request for processing by Service.
type ServiceMessage struct {
	Command string
	Params  ExecParams
	Extra   interface{}
}

// Return parameter value if any, or fallback to def.
func (m *ServiceMessage) GetParam(key string, def interface{}) interface{} {
	if m.Params == nil {
		return def
	}
	return m.Params.Get(key, def)
}

// Set parameter with specified key.
func (m *ServiceMessage) SetParam(key string, val interface{}) {
	if m.Params == nil {
		m.Params = make(ExecParams)
	}
	m.Params.Set(key, val)
}

// Delete parameter with specified key.
func (m *ServiceMessage) DeleteParam(key string) {
	if m.Params != nil {
		m.Params.Delete(key)
	}
}

// Return Singal of the param.
func (m *ServiceMessage) ReturnSignal() *sync.WaitGroup {
	if m.Params != nil {
		return m.Params.ReturnSignal()
	}
	return nil
}

// Return Result of the param.
func (m *ServiceMessage) ReturnResult() interface{} {
	if m.Params != nil {
		return m.Params.ReturnResult()
	}
	return nil
}

// Indicate that the request expect returning result.
// This is for sender side.
func (m *ServiceMessage) ExpectReturn() {
	if m.Params == nil {
		m.Params = make(ExecParams)
	}
	m.Params.ExpectReturn()
}

// Indicate that the request expect returning result using a custom signal.
// This is for sender side.
func (m *ServiceMessage) ExpectReturnCustomSignal(signal *sync.WaitGroup) {
	if m.Params == nil {
		m.Params = make(ExecParams)
	}
	m.Params.ExpectReturnCustomSignal(signal)
}

// Set the returning result then signal listener that the request has been completed, and report whether the sender receives result.
// Like a return statement, the first answer ends the request: a later Return, including one after the service answered on its own, changes nothing and reports false.
// Nothing will be done if the sender doesn't expect returns, and it reports false.
// This is for recipient side.
func (m *ServiceMessage) Return(result interface{}) bool {
	if m.Params != nil {
		return m.Params.Return(result)
	}
	return false
}

// Listen to the signal.
// The routine won't be blocked and receive nil if it doesn't expect returns.
// This is for sender side.
func (m *ServiceMessage) Wait() {
	if m.Params != nil {
		m.Params.Wait()
	}
}

// Listen to the signal and return received result, or ErrDropped if the request was dropped.
// The routine won't be blocked and receive nil if it doesn't expect returns.
// This is for sender side.
func (m *ServiceMessage) WaitForReturn() interface{} {
	if m.Params != nil {
		return m.Params.WaitForReturn()
	}
	return nil
}

// Collection of parameters as key-value mapping.
type ExecParams map[string]interface{}

// Return parameter value if any, or fallback to def.
func (p ExecParams) Get(key string, def interface{}) interface{} {
	if val, ok := p[key]; ok {
		return val
	}
	return def
}

// Set parameter with specified key.
func (p ExecParams) Set(key string, val interface{}) {
	p[key] = val
}

// Delete parameter with specified key.
func (p ExecParams) Delete(key string) {
	delete(p, key)
}

// Return Singal of the param.
func (p ExecParams) ReturnSignal() *sync.WaitGroup {
	if p["return"] != nil {
		ret := p.Get("return", nil).(*ReturnParams)
		return ret.signal
	}
	return nil
}

// Return Result of the param.
func (p ExecParams) ReturnResult() interface{} {
	if p["return"] != nil {
		ret := p.Get("return", nil).(*ReturnParams)
		return ret.result
	}
	return nil
}

// Indicate that the request expect returning result.
// This is for sender side.
func (p ExecParams) ExpectReturn() {
	signal := new(sync.WaitGroup)
	signal.Add(1)
	p.ExpectReturnCustomSignal(signal)
}

// Indicate that the request expect returning result using a custom signal.
// This is for sender side.
func (p ExecParams) ExpectReturnCustomSignal(signal *sync.WaitGroup) {
	p["return"] = &ReturnParams{
		signal: signal,
	}
}

// Set the returning result then signal listener that the request has been completed, and report whether the sender receives result.
// Like a return statement, the first answer ends the request: a later Return, including one after the service answered on its own, changes nothing and reports false.
// Nothing will be done if the sender doesn't expect returns, and it reports false.
// This is for recipient side.
func (p ExecParams) Return(result interface{}) bool {
	if p["return"] == nil {
		return false
	}
	ret := p.Get("return", nil).(*ReturnParams)
	// Claim the request before writing result: concurrent answers must never both write it, or both signal, which panics.
	if !atomic.CompareAndSwapInt32(&ret.returned, 0, 1) {
		return false
	}
	ret.result = result
	ret.signal.Done()
	return true
}

// Listen to the signal.
// The routine won't be blocked and receive nil if it doesn't expect returns.
// This is for sender side.
func (p ExecParams) Wait() {
	if p["return"] != nil {
		ret := p.Get("return", nil).(*ReturnParams)
		ret.signal.Wait()
	}
}

// Listen to the signal and return received result, or ErrDropped if the request was dropped.
// The routine won't be blocked and receive nil if it doesn't expect returns.
// This is for sender side.
func (p ExecParams) WaitForReturn() interface{} {
	if p["return"] != nil {
		ret := p.Get("return", nil).(*ReturnParams)
		ret.signal.Wait()
		return ret.result
	}
	return nil
}

// Return ErrDropped to the sender, unless the request already returned a result.
func (p ExecParams) returnDropped() {
	p.returnIfUnanswered(ErrDropped)
}

// Return result to the sender, unless the request already returned a result.
func (p ExecParams) returnIfUnanswered(result interface{}) {
	p.Return(result)
}

// ReturnParams comprises of a dyanmic type result and a signal for synchronous support.
type ReturnParams struct {
	signal   *sync.WaitGroup
	result   interface{}
	returned int32 // 1 once Return was called; int32 for sync/atomic, as Go 1.18 has no atomic.Bool
}

// Return Singal of the param.
func (p *ReturnParams) Signal() *sync.WaitGroup {
	return p.signal
}

// Return Result of the param.
func (p *ReturnParams) Result() interface{} {
	return p.result
}
