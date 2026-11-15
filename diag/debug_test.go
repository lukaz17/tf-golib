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
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDebugLogger_Error(t *testing.T) {
	logger := NewDebugLogger(10)
	logger.Error(errors.New("invalid"), "Message")
	output := logger.LastMessage()
	assert.Equal(t, "ERROR invalid Message", output)
}

func TestDebugLogger_Errorf(t *testing.T) {
	logger := NewDebugLogger(10)
	logger.Errorf(errors.New("invalid"), "%s", "Messagef")
	output := logger.LastMessage()
	assert.Equal(t, "ERROR invalid Messagef", output)
}

func TestDebugLogger_Warn(t *testing.T) {
	logger := NewDebugLogger(10)
	logger.Warn("Message")
	output := logger.LastMessage()
	assert.Equal(t, "WARN Message", output)
}

func TestDebugLogger_Warnf(t *testing.T) {
	logger := NewDebugLogger(10)
	logger.Warnf("%s", "Messagef")
	output := logger.LastMessage()
	assert.Equal(t, "WARN Messagef", output)
}

func TestDebugLogger_Info(t *testing.T) {
	logger := NewDebugLogger(10)
	logger.Info("Message")
	output := logger.LastMessage()
	assert.Equal(t, "INFO Message", output)
}

func TestDebugLogger_Infof(t *testing.T) {
	logger := NewDebugLogger(10)
	logger.Infof("%s", "Messagef")
	output := logger.LastMessage()
	assert.Equal(t, "INFO Messagef", output)
}

func TestDebugLogger_Debug(t *testing.T) {
	logger := NewDebugLogger(10)
	logger.Debug("Message")
	output := logger.LastMessage()
	assert.Equal(t, "DEBUG Message", output)
}

func TestDebugLogger_Debugf(t *testing.T) {
	logger := NewDebugLogger(10)
	logger.Debugf("%s", "Messagef")
	output := logger.LastMessage()
	assert.Equal(t, "DEBUG Messagef", output)
}

func TestDebugLogger_Trace(t *testing.T) {
	logger := NewDebugLogger(10)
	logger.Trace("Message")
	output := logger.LastMessage()
	assert.Equal(t, "TRACE Message", output)
}

func TestDebugLogger_Tracef(t *testing.T) {
	logger := NewDebugLogger(10)
	logger.Tracef("%s", "Messagef")
	output := logger.LastMessage()
	assert.Equal(t, "TRACE Messagef", output)
}

func TestDebugLogger_Last_Empty(t *testing.T) {
	logger := NewDebugLogger(10)
	output := logger.Last()
	assert.Nil(t, output)
}

func TestDebugLogger_LastMessage_Empty(t *testing.T) {
	logger := NewDebugLogger(10)
	output := logger.LastMessage()
	assert.Empty(t, output)
}

func TestDebugLogger_All(t *testing.T) {
	logger := NewDebugLogger(5)
	logger.Info("Message 1")
	logger.Warn("Message 2")
	logger.Debug("Message 3")
	logger.Info("Message 4")

	contexts := logger.All()
	assert.Len(t, contexts, 4)

	expectedMessages := []string{"INFO Message 1", "WARN Message 2", "DEBUG Message 3", "INFO Message 4"}
	for i, ctx := range contexts {
		assert.Equal(t, expectedMessages[i], ctx.Message)
	}
}

func TestDebugLogger_All_Overflown(t *testing.T) {
	logger := NewDebugLogger(3)
	logger.Info("Message 1")
	logger.Warn("Message 2")
	logger.Debug("Message 3")
	logger.Info("Message 4")

	contexts := logger.All()
	assert.Len(t, contexts, 3)

	expectedMessages := []string{"WARN Message 2", "DEBUG Message 3", "INFO Message 4"}
	for i, ctx := range contexts {
		assert.Equal(t, expectedMessages[i], ctx.Message)
	}
}

func TestDebugLogger_AllMessages(t *testing.T) {
	logger := NewDebugLogger(5)
	logger.Info("Message 1")
	logger.Warn("Message 2")
	logger.Debug("Message 3")
	logger.Info("Message 4")

	messages := logger.AllMessages()
	assert.Len(t, messages, 4)

	expectedMessages := []string{"INFO Message 1", "WARN Message 2", "DEBUG Message 3", "INFO Message 4"}
	for i, msg := range messages {
		assert.Equal(t, expectedMessages[i], msg)
	}
}

func TestDebugLogger_AllMessages_Overflown(t *testing.T) {
	logger := NewDebugLogger(3)
	logger.Info("Message 1")
	logger.Warn("Message 2")
	logger.Debug("Message 3")
	logger.Info("Message 4")

	messages := logger.AllMessages()
	assert.Len(t, messages, 3)

	expectedMessages := []string{"WARN Message 2", "DEBUG Message 3", "INFO Message 4"}
	for i, msg := range messages {
		assert.Equal(t, expectedMessages[i], msg)
	}
}

func TestDebugLogger_AllReverse(t *testing.T) {
	logger := NewDebugLogger(5)
	logger.Info("Message 1")
	logger.Warn("Message 2")
	logger.Debug("Message 3")
	logger.Info("Message 4")

	contexts := logger.AllReverse()
	assert.Len(t, contexts, 4)

	expectedMessages := []string{"INFO Message 4", "DEBUG Message 3", "WARN Message 2", "INFO Message 1"}
	for i, ctx := range contexts {
		assert.Equal(t, expectedMessages[i], ctx.Message)
	}
}

func TestDebugLogger_AllReverse_Overflown(t *testing.T) {
	logger := NewDebugLogger(3)
	logger.Info("Message 1")
	logger.Warn("Message 2")
	logger.Debug("Message 3")
	logger.Info("Message 4")

	contexts := logger.AllReverse()
	assert.Len(t, contexts, 3)

	expectedMessages := []string{"INFO Message 4", "DEBUG Message 3", "WARN Message 2"}
	for i, ctx := range contexts {
		assert.Equal(t, expectedMessages[i], ctx.Message)
	}
}

func TestDebugLogger_AllReverse_KeepsLastMessage(t *testing.T) {
	logger := NewDebugLogger(5)
	logger.Info("Message 1")
	logger.Info("Message 2")

	_ = logger.AllReverse()

	assert.Equal(t, "INFO Message 2", logger.LastMessage())
}

func TestDebugLogger_AllReverse_KeepsOldestMessageOnNextWrite(t *testing.T) {
	logger := NewDebugLogger(5)
	logger.Info("Message 1")
	logger.Info("Message 2")
	_ = logger.AllReverse()

	logger.Info("Message 3")

	assert.Equal(t, []string{"INFO Message 1", "INFO Message 2", "INFO Message 3"}, logger.AllMessages())
}

func TestDebugLogger_AllMessagesReverse(t *testing.T) {
	logger := NewDebugLogger(5)
	logger.Info("Message 1")
	logger.Warn("Message 2")
	logger.Debug("Message 3")
	logger.Info("Message 4")

	messages := logger.AllMessagesReverse()
	assert.Len(t, messages, 4)

	expectedMessages := []string{"INFO Message 4", "DEBUG Message 3", "WARN Message 2", "INFO Message 1"}
	for i, msg := range messages {
		assert.Equal(t, expectedMessages[i], msg)
	}
}

func TestDebugLogger_AllMessagesReverse_Overflown(t *testing.T) {
	logger := NewDebugLogger(3)
	logger.Info("Message 1")
	logger.Warn("Message 2")
	logger.Debug("Message 3")
	logger.Info("Message 4")

	messages := logger.AllMessagesReverse()
	assert.Len(t, messages, 3)

	expectedMessages := []string{"INFO Message 4", "DEBUG Message 3", "WARN Message 2"}
	for i, msg := range messages {
		assert.Equal(t, expectedMessages[i], msg)
	}
}

func TestDebugLogger_AllMessagesReverse_KeepsLastMessage(t *testing.T) {
	logger := NewDebugLogger(5)
	logger.Info("Message 1")
	logger.Info("Message 2")

	_ = logger.AllMessagesReverse()

	assert.Equal(t, "INFO Message 2", logger.LastMessage())
}

func TestDebugLogger_AllMessagesReverse_KeepsOldestMessageOnNextWrite(t *testing.T) {
	logger := NewDebugLogger(5)
	logger.Info("Message 1")
	logger.Info("Message 2")
	_ = logger.AllMessagesReverse()

	logger.Info("Message 3")

	assert.Equal(t, []string{"INFO Message 1", "INFO Message 2", "INFO Message 3"}, logger.AllMessages())
}
