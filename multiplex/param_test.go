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

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestServiceMessageInitialization(t *testing.T) {
	msg := ServiceMessage{
		Command: "testCommand",
		Params:  ExecParams{"key": "value"},
		Extra:   "extraData",
	}
	assert.Equal(t, "testCommand", msg.Command, "Command should match")
	assert.Equal(t, "value", msg.Params["key"], "Params['key'] should match")
	assert.Equal(t, "extraData", msg.Extra, "Extra should match")
}

func TestServiceMessage_GetParam(t *testing.T) {
	msg := ServiceMessage{
		Params: ExecParams{"key": "value"},
	}
	msg2 := ServiceMessage{}
	assert.Equal(t, "value", msg.GetParam("key", "default"), "GetParam should return the correct value")
	assert.Equal(t, "default", msg.GetParam("missing", "default"), "GetParam should return the default value for missing keys")
	assert.Equal(t, "default", msg2.GetParam("key", "default"), "GetParam should return the default value for missing keys")
}

func TestServiceMessage_SetParam(t *testing.T) {
	msg := ServiceMessage{}
	msg.SetParam("key", "value")
	assert.NotNil(t, msg.Params, "Params should be initialized")
	assert.Equal(t, "value", msg.Params["key"], "SetParam should correctly set the key-value pair")
}

func TestServiceMessage_DeleteParam(t *testing.T) {
	msg := ServiceMessage{
		Params: ExecParams{"key": "value"},
	}
	msg.DeleteParam("key")
	_, exists := msg.Params["key"]
	assert.False(t, exists, "DeleteParam should remove the key")
}

func TestServiceMessage_ExpectReturn(t *testing.T) {
	msg := ServiceMessage{}
	msg.ExpectReturn()
	assert.NotNil(t, msg.Params["return"], "ExpectReturn should set the 'return' key")
}

func TestServiceMessage_ExpectReturnCustomSignal(t *testing.T) {
	msg := ServiceMessage{}
	signal := &sync.WaitGroup{}
	msg.ExpectReturnCustomSignal(signal)
	retParams, ok := msg.Params["return"].(*ReturnParams)
	assert.True(t, ok, "ExpectReturnCustomSignal should set the 'return' key with a ReturnParams value")
	assert.Equal(t, signal, retParams.signal, "ExpectReturnCustomSignal should set the correct signal")
}

func TestServiceMessage_Return(t *testing.T) {
	msg := ServiceMessage{}
	msg.ExpectReturn()
	go func() {
		msg.Return("result")
	}()
	assert.Equal(t, "result", msg.WaitForReturn(), "Return should set the result and signal completion")
}

func TestServiceMessage_Return_ReportsDelivery(t *testing.T) {
	tests := []struct {
		name          string
		arrange       func(msg *ServiceMessage)
		wantDelivered bool
		wantResult    interface{}
	}{
		{
			name:          "first answer reaches the sender",
			arrange:       func(m *ServiceMessage) { m.ExpectReturn() },
			wantDelivered: true,
			wantResult:    "answer",
		},
		{
			name: "answer after another Return is discarded",
			arrange: func(m *ServiceMessage) {
				m.ExpectReturn()
				m.Return("earlier")
			},
			wantDelivered: false,
			wantResult:    "earlier",
		},
		{
			name:          "answer to a message without params is discarded",
			arrange:       func(*ServiceMessage) {},
			wantDelivered: false,
			wantResult:    nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			msg := &ServiceMessage{}
			tt.arrange(msg)

			delivered := msg.Return("answer")

			assert.Equal(t, tt.wantDelivered, delivered)
			assert.Equal(t, tt.wantResult, msg.WaitForReturn())
		})
	}
}

func TestServiceMessage_WaitForReturn_NoReturn(t *testing.T) {
	msg := ServiceMessage{}
	assert.Nil(t, msg.WaitForReturn(), "WaitForReturn should return nil if no return is expected")
}

func TestServiceMessage_Wait(t *testing.T) {
	msg := ServiceMessage{}
	msg.ExpectReturn()
	go func() {
		msg.Return("result")
	}()
	msg.Wait()
	assert.Equal(t, "result", msg.WaitForReturn(), "Wait should block until the signal is completed")
}

func TestServiceMessage_ReturnSignal(t *testing.T) {
	msg := ServiceMessage{}
	assert.Nil(t, msg.ReturnSignal(), "ReturnSignal should return nil if no return is expected")
	msg.ExpectReturn()
	assert.NotNil(t, msg.ReturnSignal(), "ReturnSignal should return the signal")
}

func TestServiceMessage_ReturnResult(t *testing.T) {
	msg := ServiceMessage{}
	assert.Nil(t, msg.ReturnResult(), "ReturnResult should return nil if no return is expected")
	msg.ExpectReturn()
	go func() {
		msg.Return("result")
	}()
	msg.Wait()
	assert.Equal(t, "result", msg.ReturnResult(), "ReturnResult should return the correct result")
}

func TestExecParams_Get(t *testing.T) {
	params := ExecParams{"key": "value"}
	assert.Equal(t, "value", params.Get("key", "default"), "Get should return the correct value")
	assert.Equal(t, "default", params.Get("missing", "default"), "Get should return the default value for missing keys")
}

func TestExecParams_Set(t *testing.T) {
	params := ExecParams{}
	params.Set("key", "value")
	assert.Equal(t, "value", params["key"], "Set should correctly set the key-value pair")
}

func TestExecParams_Delete(t *testing.T) {
	params := ExecParams{"key": "value"}
	params.Delete("key")
	_, exists := params["key"]
	assert.False(t, exists, "Delete should remove the key")
}

func TestExecParams_ReturnSignal(t *testing.T) {
	params := ExecParams{}
	assert.Nil(t, params.ReturnSignal(), "ReturnSignal should return nil if no return is expected")
	params.ExpectReturn()
	signal := params.ReturnSignal()
	assert.NotNil(t, signal, "ReturnSignal should return the signal")
}

func TestExecParams_ReturnResult(t *testing.T) {
	params := ExecParams{}
	assert.Nil(t, params.ReturnResult(), "ReturnResult should return nil if no return is expected")
	params.ExpectReturn()
	go func() {
		params.Return("result")
	}()
	params.Wait()
	assert.Equal(t, "result", params.ReturnResult(), "ReturnResult should return the correct result")
}

func TestExecParams_ExpectReturn(t *testing.T) {
	params := ExecParams{}
	params.ExpectReturn()
	_, exists := params["return"]
	assert.True(t, exists, "ExpectReturn should set the 'return' key")
}

func TestExecParams_Return(t *testing.T) {
	params := ExecParams{}
	params.ExpectReturn()
	go func() {
		params.Return("result")
	}()
	assert.Equal(t, "result", params.WaitForReturn(), "Return should set the result and signal completion")
}

func TestExecParams_Return_ReportsDelivery(t *testing.T) {
	tests := []struct {
		name          string
		arrange       func(params ExecParams)
		wantDelivered bool
		wantResult    interface{}
	}{
		{
			name:          "first answer reaches the sender",
			arrange:       func(p ExecParams) { p.ExpectReturn() },
			wantDelivered: true,
			wantResult:    "answer",
		},
		{
			name: "answer after another Return is discarded",
			arrange: func(p ExecParams) {
				p.ExpectReturn()
				p.Return("earlier")
			},
			wantDelivered: false,
			wantResult:    "earlier",
		},
		{
			name: "answer after the service dropped the request is discarded",
			arrange: func(p ExecParams) {
				p.ExpectReturn()
				p.returnDropped()
			},
			wantDelivered: false,
			wantResult:    ErrDropped,
		},
		{
			name:          "answer to a request expecting none is discarded",
			arrange:       func(ExecParams) {},
			wantDelivered: false,
			wantResult:    nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			params := ExecParams{}
			tt.arrange(params)

			delivered := params.Return("answer")

			assert.Equal(t, tt.wantDelivered, delivered)
			assert.Equal(t, tt.wantResult, params.WaitForReturn())
		})
	}
}

func TestExecParams_Return_ConcurrentAnswersReachSenderOnce(t *testing.T) {
	// Each round gives the answers a fresh chance to interleave; correct code passes every round.
	const rounds = 1000
	tests := []struct {
		name          string
		returnAnswers int
		dropAnswers   int
	}{
		{name: "answers through Return", returnAnswers: 4},
		{name: "answer through Return against the worker's fallback", returnAnswers: 1, dropAnswers: 1},
		{name: "answers through the worker's fallback", dropAnswers: 2},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for round := 0; round < rounds; round++ {
				params := ExecParams{}
				params.ExpectReturn()
				delivered := make(chan interface{}, tt.returnAnswers)
				answers := []func(){}
				for i := 0; i < tt.returnAnswers; i++ {
					answer := i
					answers = append(answers, func() {
						if params.Return(answer) {
							delivered <- answer
						}
					})
				}
				for i := 0; i < tt.dropAnswers; i++ {
					answers = append(answers, params.returnDropped)
				}

				panics := answerConcurrently(answers)
				// A second signal leaves the sender's WaitGroup negative, and WaitForReturn would then block forever.
				require.Empty(t, panics, "round %d: request was answered more than once", round)
				result := params.WaitForReturn()
				close(delivered)

				winners := []interface{}{}
				for answer := range delivered {
					winners = append(winners, answer)
				}
				require.LessOrEqual(t, len(winners), 1, "round %d: more than one Return reported its answer delivered", round)
				want := interface{}(ErrDropped)
				if len(winners) == 1 {
					want = winners[0]
				}
				require.Equal(t, want, result, "round %d: sender must receive the answer reported delivered", round)
			}
		})
	}
}

func TestExecParams_WaitForReturn_NoReturn(t *testing.T) {
	params := ExecParams{}
	assert.Nil(t, params.WaitForReturn(), "WaitForReturn should return nil if no return is expected")
}

func TestReturnParams_Signal(t *testing.T) {
	signal := &sync.WaitGroup{}
	retParams := ReturnParams{signal: signal}
	assert.Equal(t, signal, retParams.Signal(), "Signal should return the correct signal")
}

func TestReturnParams_Result(t *testing.T) {
	retParams := ReturnParams{result: "testResult"}
	assert.Equal(t, "testResult", retParams.Result(), "Result should return the correct result")
}

// Run each answer on its own goroutine, all released at once, and return the panics they raised.
func answerConcurrently(answers []func()) []interface{} {
	start := make(chan struct{})
	raised := make(chan interface{}, len(answers))
	var wg sync.WaitGroup
	for _, answer := range answers {
		wg.Add(1)
		go func(answer func()) {
			defer wg.Done()
			defer func() {
				if r := recover(); r != nil {
					raised <- r
				}
			}()
			<-start
			answer()
		}(answer)
	}
	close(start)
	wg.Wait()
	close(raised)
	panics := []interface{}{}
	for r := range raised {
		panics = append(panics, r)
	}
	return panics
}
