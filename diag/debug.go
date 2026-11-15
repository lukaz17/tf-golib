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
	"container/ring"
	"fmt"
	"sync"
	"time"
)

// DebugLogger implement Logger interface that store log contents in a circular list.
type DebugLogger struct {
	i *DebugLoggerInternal
}

// DebugLoggerInternal stores internal data of a DebugLogger.
type DebugLoggerInternal struct {
	Cache *ring.Ring
	LogMu *sync.Mutex
}

// Return new DebugLogger with maximum capacity for cache.
func NewDebugLogger(capacity int) *DebugLogger {
	return &DebugLogger{
		&DebugLoggerInternal{
			Cache: ring.New(capacity),
			LogMu: &sync.Mutex{},
		},
	}
}

// Print a message with Error level.
func (l DebugLogger) Error(err error, v ...interface{}) {
	v2 := append([]interface{}{"ERROR ", err, " "}, v...)
	l.write(ErrorLevel, v2...)

}

// Print a message with Error level with format.
func (l DebugLogger) Errorf(err error, format string, v ...interface{}) {
	v2 := append([]interface{}{"ERROR ", err}, v...)
	l.writef(ErrorLevel, "%s%v "+format, v2...)
}

// Print a message with Warn level.
func (l DebugLogger) Warn(v ...interface{}) {
	v2 := append([]interface{}{"WARN "}, v...)
	l.write(WarnLevel, v2...)
}

// Print a message with Warn level with format.
func (l DebugLogger) Warnf(format string, v ...interface{}) {
	v2 := append([]interface{}{"WARN "}, v...)
	l.writef(WarnLevel, "%s"+format, v2...)
}

// Print a message with Info level.
func (l DebugLogger) Info(v ...interface{}) {
	v2 := append([]interface{}{"INFO "}, v...)
	l.write(InfoLevel, v2...)
}

// Print a message with Info level with format.
func (l DebugLogger) Infof(format string, v ...interface{}) {
	v2 := append([]interface{}{"INFO "}, v...)
	l.writef(InfoLevel, "%s"+format, v2...)
}

// Print a message with Debug level.
func (l DebugLogger) Debug(v ...interface{}) {
	v2 := append([]interface{}{"DEBUG "}, v...)
	l.write(DebugLevel, v2...)
}

// Print a message with Debug level with format.
func (l DebugLogger) Debugf(format string, v ...interface{}) {
	v2 := append([]interface{}{"DEBUG "}, v...)
	l.writef(DebugLevel, "%s"+format, v2...)
}

// Print a message with Trace level.
func (l DebugLogger) Trace(v ...interface{}) {
	v2 := append([]interface{}{"TRACE "}, v...)
	l.write(TraceLevel, v2...)
}

// Print a message with Trace level with format.
func (l DebugLogger) Tracef(format string, v ...interface{}) {
	v2 := append([]interface{}{"TRACE "}, v...)
	l.writef(TraceLevel, "%s"+format, v2...)
}

// Return the context for last log message.
func (l *DebugLogger) Last() *LogContext {
	l.i.LogMu.Lock()
	defer l.i.LogMu.Unlock()
	v := l.i.Cache.Value
	if v == nil {
		return nil
	}
	return v.(*LogContext)
}

// Return the message for last log message.
func (l *DebugLogger) LastMessage() string {
	last := l.Last()
	if last == nil {
		return ""
	} else {
		return last.Message
	}
}

// Return the context for all log messages in chronological order.
func (l *DebugLogger) All() []*LogContext {
	l.i.LogMu.Lock()
	defer l.i.LogMu.Unlock()
	cap := l.i.Cache.Len()
	contexts := []*LogContext{}
	for i := 0; i < cap; i++ {
		l.i.Cache = l.i.Cache.Next()
		if l.i.Cache.Value == nil {
			continue
		}
		contexts = append(contexts, l.i.Cache.Value.(*LogContext))
	}
	return contexts
}

// Return the message for all log messages in chronological order.
func (l *DebugLogger) AllMessages() []string {
	l.i.LogMu.Lock()
	defer l.i.LogMu.Unlock()
	cap := l.i.Cache.Len()
	messages := []string{}
	for i := 0; i < cap; i++ {
		l.i.Cache = l.i.Cache.Next()
		if l.i.Cache.Value == nil {
			continue
		}
		messages = append(messages, l.i.Cache.Value.(*LogContext).Message)
	}
	return messages
}

// Return the context for all log messages in reverse chronological order.
func (l *DebugLogger) AllReverse() []*LogContext {
	l.i.LogMu.Lock()
	defer l.i.LogMu.Unlock()
	cap := l.i.Cache.Len()
	contexts := []*LogContext{}
	for i := 0; i < cap; i++ {
		if l.i.Cache.Value == nil {
			break
		}
		contexts = append(contexts, l.i.Cache.Value.(*LogContext))
		l.i.Cache = l.i.Cache.Prev()
	}
	return contexts
}

// Return the message for all log messages in reverse chronological order.
func (l *DebugLogger) AllMessagesReverse() []string {
	l.i.LogMu.Lock()
	defer l.i.LogMu.Unlock()
	cap := l.i.Cache.Len()
	messages := []string{}
	for i := 0; i < cap; i++ {
		if l.i.Cache.Value == nil {
			break
		}
		messages = append(messages, l.i.Cache.Value.(*LogContext).Message)
		l.i.Cache = l.i.Cache.Prev()
	}
	return messages
}

// Write the message context into cache.
func (l *DebugLogger) write(level LogLevel, v ...interface{}) {
	l.i.LogMu.Lock()
	defer l.i.LogMu.Unlock()
	l.i.Cache = l.i.Cache.Next()
	l.i.Cache.Value = &LogContext{
		Time:    time.Now(),
		Level:   level,
		Message: fmt.Sprint(v...),
	}
}

// Write the message context into cache with format.
func (l *DebugLogger) writef(level LogLevel, format string, v ...interface{}) {
	l.i.LogMu.Lock()
	defer l.i.LogMu.Unlock()
	l.i.Cache = l.i.Cache.Next()
	l.i.Cache.Value = &LogContext{
		Time:    time.Now(),
		Level:   level,
		Message: fmt.Sprintf(format, v...),
	}
}
