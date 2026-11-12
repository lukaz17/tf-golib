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
	"crypto/sha256"
	"encoding/hex"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/tforce-io/tf-golib/diag"
	"github.com/tforce-io/tf-golib/random/securerng"
)

func TestServiceCore_ServiceID(t *testing.T) {
	logger := diag.NewDebugLogger(10)
	svc := NewRandomService(logger)
	assert.Equal(t, "Random", svc.ServiceID(), "invalid service ID")
}

func TestServiceCore_Router(t *testing.T) {
	logger := diag.NewDebugLogger(10)
	svc := NewRandomService(logger)
	assert.Nil(t, svc.Router(), "unexpected value for router, must be nil")
	controller := NewServiceController(logger)
	svc.SetRouter(controller)
	assert.NotNil(t, svc.Router(), "unexpected value for router, must not be nil")
}

func TestServiceCore_Worker(t *testing.T) {
	logger := diag.NewDebugLogger(10)
	svc := NewRandomService(logger)
	assert.Equal(t, uint64(0), svc.WorkerCount(), "mismatch worker count")
	svc.SetWorker(17)
	assert.Equal(t, uint64(17), svc.WorkerCount(), "mismatch worker count")
	time.Sleep(100 * time.Millisecond)
	svc.SetWorker(10)
	assert.Equal(t, uint64(10), svc.WorkerCount(), "mismatch worker count")
}

func TestServiceCore_Exec(t *testing.T) {
	logger := diag.NewDebugLogger(10)
	svc := NewEchoService(logger)
	svc.SetWorker(1)
	svc.Exec("", ExecParams{
		"message": "Hello, World!",
	})
	time.Sleep(10 * time.Millisecond)
	assert.Equal(t, "INFO Echo#1: Message received: Hello, World!", logger.LastMessage(), "invalid message")
}

type EchoService struct {
	ServiceCore
	i *ServiceCoreInternal
}

func NewEchoService(logger diag.Logger) *EchoService {
	svc := &EchoService{}
	svc.i = svc.InitServiceCore("Echo", logger, svc.coreProcessHook)
	return svc
}

func (s *EchoService) coreProcessHook(workerID uint64, msg *ServiceMessage) *HookState {
	message := msg.Params["message"].(string)
	s.i.Logger.Infof("%s#%d: Message received: %s", s.i.ServiceID, workerID, message)
	return &HookState{Handled: true}
}

type HashService struct {
	ServiceCore
	i *ServiceCoreInternal
}

func NewHashService(logger diag.Logger) *HashService {
	svc := &HashService{}
	svc.i = svc.InitServiceCore("Hash", logger, svc.coreProcessHook)
	return svc
}

func (s *HashService) coreProcessHook(workerID uint64, msg *ServiceMessage) *HookState {
	switch msg.Command {
	case "sha256":
		message := msg.Params["input"].(string)
		messageBuf := []byte(message)
		hashBuf := sha256.Sum256(messageBuf)
		hash := hex.EncodeToString(hashBuf[:])
		s.i.Logger.Infof("%s#%d: Value hashed: %s.", s.i.ServiceID, workerID, hash)
	case "sha256_random":
		randomParams := ExecParams{}
		randomParams.ExpectReturn()
		s.Dispatch("Random", "", randomParams)
		message := randomParams.WaitForReturn().(string)
		messageBuf := []byte(message)
		hashBuf := sha256.Sum256(messageBuf)
		hash := hex.EncodeToString(hashBuf[:])
		s.i.Logger.Infof("%s#%d: Value hashed: %s.", s.i.ServiceID, workerID, hash)
	}
	return &HookState{Handled: true}
}

type RandomService struct {
	ServiceCore
	i *ServiceCoreInternal
}

func NewRandomService(logger diag.Logger) *RandomService {
	svc := &RandomService{}
	svc.i = svc.InitServiceCore("Random", logger, svc.coreProcessHook)
	return svc
}

func (s *RandomService) coreProcessHook(workerID uint64, msg *ServiceMessage) *HookState {
	hex := securerng.Hex(16)
	msg.Return(hex)
	s.i.Logger.Infof("%s#%d: Value randomed: %s.", s.i.ServiceID, workerID, hex)
	return &HookState{Handled: true}
}

type ShutdownService struct {
	ServiceCore
	i *ServiceCoreInternal
}

func NewShutdownService(logger diag.Logger) *ShutdownService {
	svc := &ShutdownService{}
	svc.i = svc.InitServiceCore("Shutdown", logger, svc.coreProcessHook)
	return svc
}

func (s *ShutdownService) coreProcessHook(workerID uint64, msg *ServiceMessage) *HookState {
	timeout := msg.GetParam("timeout", int64(100*time.Millisecond)).(int64)
	time.Sleep(time.Duration(timeout))
	s.Dispatch("", "exit", ExecParams{})
	return &HookState{Handled: true}
}
