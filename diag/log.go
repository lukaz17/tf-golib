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

import "log"

// Logger defines required function for all log levels used for application logging.
type Logger interface {
	// Print a message with Error level.
	Error(err error, v ...interface{})
	// Print a message with Error level with format.
	Errorf(err error, format string, v ...interface{})

	// Print a message with Warn level.
	Warn(v ...interface{})
	// Print a message with Warn level with format.
	Warnf(format string, v ...interface{})

	// Print a message with Info level.
	Info(v ...interface{})
	// Print a message with Info level with format.
	Infof(format string, v ...interface{})

	// Print a message with Debug level.
	Debug(v ...interface{})
	// Print a message with Debug level with format.
	Debugf(format string, v ...interface{})

	// Print a message with Trace level.
	Trace(v ...interface{})
	// Print a message with Trace level with format.
	Tracef(format string, v ...interface{})
}

// DefaultLogger implement Logger interface that prints log message to stdout
// using global logger instance of Go.
type DefaultLogger struct {
}

// Print a message with Error level.
func (l DefaultLogger) Error(err error, v ...interface{}) {
	v2 := append([]interface{}{"ERROR ", err, " "}, v...)
	log.Print(v2...)
}

// Print a message with Error level with format.
func (l DefaultLogger) Errorf(err error, format string, v ...interface{}) {
	v2 := append([]interface{}{"ERROR ", err}, v...)
	log.Printf("%s%v "+format, v2...)
}

// Print a message with Warn level.
func (l DefaultLogger) Warn(v ...interface{}) {
	v2 := append([]interface{}{"WARN "}, v...)
	log.Print(v2...)
}

// Print a message with Warn level with format.
func (l DefaultLogger) Warnf(format string, v ...interface{}) {
	v2 := append([]interface{}{"WARN "}, v...)
	log.Printf("%s"+format, v2...)
}

// Print a message with Info level.
func (l DefaultLogger) Info(v ...interface{}) {
	v2 := append([]interface{}{"INFO "}, v...)
	log.Print(v2...)
}

// Print a message with Info level with format.
func (l DefaultLogger) Infof(format string, v ...interface{}) {
	v2 := append([]interface{}{"INFO "}, v...)
	log.Printf("%s"+format, v2...)
}

// Print a message with Debug level.
func (l DefaultLogger) Debug(v ...interface{}) {
	v2 := append([]interface{}{"DEBUG "}, v...)
	log.Print(v2...)
}

// Print a message with Debug level with format.
func (l DefaultLogger) Debugf(format string, v ...interface{}) {
	v2 := append([]interface{}{"DEBUG "}, v...)
	log.Printf("%s"+format, v2...)
}

// Print a message with Trace level.
func (l DefaultLogger) Trace(v ...interface{}) {
	v2 := append([]interface{}{"TRACE "}, v...)
	log.Print(v2...)
}

// Print a message with Trace level with format.
func (l DefaultLogger) Tracef(format string, v ...interface{}) {
	v2 := append([]interface{}{"TRACE "}, v...)
	log.Printf("%s"+format, v2...)
}
