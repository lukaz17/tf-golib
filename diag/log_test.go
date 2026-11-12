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
	"bytes"
	"errors"
	"log"
	"strings"
	"testing"
)

func TestDefaultLogger_Error(t *testing.T) {
	logger := &DefaultLogger{}
	output := captureStdout(func() {
		logger.Error(errors.New("invalid"), "Message")
	})
	if !strings.HasSuffix(output, "ERROR invalid Message\n") {
		t.Errorf("Unexpected log output: %s", output)
	}
}

func TestDefaultLogger_Errorf(t *testing.T) {
	logger := &DefaultLogger{}
	output := captureStdout(func() {
		logger.Errorf(errors.New("invalid"), "%s", "Messagef")
	})
	if !strings.HasSuffix(output, "ERROR invalid Messagef\n") {
		t.Errorf("Unexpected log output: %s", output)
	}
}

func TestDefaultLogger_Warn(t *testing.T) {
	logger := &DefaultLogger{}
	output := captureStdout(func() {
		logger.Warn("Message")
	})
	if !strings.HasSuffix(output, "WARN Message\n") {
		t.Errorf("Unexpected log output: %s", output)
	}
}

func TestDefaultLogger_Warnf(t *testing.T) {
	logger := &DefaultLogger{}
	output := captureStdout(func() {
		logger.Warnf("%s", "Messagef")
	})
	if !strings.HasSuffix(output, "WARN Messagef\n") {
		t.Errorf("Unexpected log output: %s", output)
	}
}

func TestDefaultLogger_Info(t *testing.T) {
	logger := &DefaultLogger{}
	output := captureStdout(func() {
		logger.Info("Message")
	})
	if !strings.HasSuffix(output, "INFO Message\n") {
		t.Errorf("Unexpected log output: %s", output)
	}
}

func TestDefaultLogger_Infof(t *testing.T) {
	logger := &DefaultLogger{}
	output := captureStdout(func() {
		logger.Infof("%s", "Messagef")
	})
	if !strings.HasSuffix(output, "INFO Messagef\n") {
		t.Errorf("Unexpected log output: %s", output)
	}
}

func TestDefaultLogger_Debug(t *testing.T) {
	logger := &DefaultLogger{}
	output := captureStdout(func() {
		logger.Debug("Message")
	})
	if !strings.HasSuffix(output, "DEBUG Message\n") {
		t.Errorf("Unexpected log output: %s", output)
	}
}

func TestDefaultLogger_Debugf(t *testing.T) {
	logger := &DefaultLogger{}
	output := captureStdout(func() {
		logger.Debugf("%s", "Messagef")
	})
	if !strings.HasSuffix(output, "DEBUG Messagef\n") {
		t.Errorf("Unexpected log output: %s", output)
	}
}

func TestDefaultLogger_Trace(t *testing.T) {
	logger := &DefaultLogger{}
	output := captureStdout(func() {
		logger.Trace("Message")
	})
	if !strings.HasSuffix(output, "TRACE Message\n") {
		t.Errorf("Unexpected log output: %s", output)
	}
}

func TestDefaultLogger_Tracef(t *testing.T) {
	logger := &DefaultLogger{}
	output := captureStdout(func() {
		logger.Tracef("%s", "Messagef")
	})
	if !strings.HasSuffix(output, "TRACE Messagef\n") {
		t.Errorf("Unexpected log output: %s", output)
	}
}

func captureStdout(f func()) string {
	var buf bytes.Buffer
	log.SetOutput(&buf)
	defer log.SetOutput(nil)
	f()
	return buf.String()
}
