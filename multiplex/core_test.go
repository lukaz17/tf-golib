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

// Bounds each wait for workers to take messages or exit, so a missing worker fails the test instead of hanging until go test's timeout.
const workerTakeTimeout = 5 * time.Second

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
	isSet := svc.SetRouter(controller)
	assert.True(t, isSet)
	assert.NotNil(t, svc.Router(), "unexpected value for router, must not be nil")
}

func TestServiceCore_SetRouter_RefusedWithWorkers(t *testing.T) {
	tests := []struct {
		name  string
		start func(t *testing.T, svc *HoldService, release chan struct{})
	}{
		{
			name:  "workers started",
			start: func(t *testing.T, svc *HoldService, release chan struct{}) { svc.SetWorker(1) },
		},
		{
			name: "workers still exiting",
			start: func(t *testing.T, svc *HoldService, release chan struct{}) {
				svc.SetWorker(1)
				requireTaken(t, svc, 1, release)
				svc.SetWorker(0)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logger := diag.NewDebugLogger(10)
			svc := NewHoldService(logger)
			t.Cleanup(func() { svc.SetWorker(0) })
			release := make(chan struct{})
			defer close(release)
			controller := NewServiceController(diag.NewDebugLogger(10))
			svc.SetRouter(controller)
			tt.start(t, svc, release)

			isSet := svc.SetRouter(NewServiceController(diag.NewDebugLogger(10)))

			assert.False(t, isSet)
			assert.Same(t, controller.Router(), svc.Router(), "router must be unchanged")
			assert.Contains(t, logger.AllMessages(), "WARN Hold: Ignored router change: service has workers.")
		})
	}
}

func TestServiceCore_Worker(t *testing.T) {
	logger := diag.NewDebugLogger(10)
	svc := NewRandomService(logger)
	assert.Equal(t, uint64(0), svc.WorkerCount(), "mismatch worker count")
	svc.SetWorker(17)
	assert.Equal(t, uint64(17), svc.WorkerCount(), "mismatch worker count")
	svc.SetWorker(10)
	assert.Equal(t, uint64(10), svc.WorkerCount(), "mismatch worker count")
}

func TestServiceCore_SetWorker_GrowWhileStarting(t *testing.T) {
	svc := NewHoldService(diag.NewDebugLogger(10))
	t.Cleanup(func() { svc.SetWorker(0) })
	release := make(chan struct{})
	defer close(release)

	svc.SetWorker(1)
	svc.SetWorker(2)

	assert.Equal(t, uint64(2), svc.WorkerCount(), "mismatch worker count")
	requireTaken(t, svc, 2, release)
}

func TestServiceCore_SetWorker_GrowWhileStopping(t *testing.T) {
	svc := NewHoldService(diag.NewDebugLogger(10))
	t.Cleanup(func() { svc.SetWorker(0) })
	release := make(chan struct{})
	svc.SetWorker(2)
	requireTaken(t, svc, 2, release)

	svc.SetWorker(0)
	svc.SetWorker(1)
	close(release)

	assert.Equal(t, uint64(1), svc.WorkerCount(), "mismatch worker count")
	requireTaken(t, svc, 1, release)
}

func TestServiceCore_SetWorker_StopsAheadOfQueuedRequests(t *testing.T) {
	svc := NewHoldService(diag.NewDebugLogger(10))
	svc.SetBackground(true)
	t.Cleanup(func() { svc.SetWorker(0) })
	release := make(chan struct{})
	svc.SetWorker(1)
	requireTaken(t, svc, 1, release)
	queuedTaken := make(chan struct{}, 1)
	svc.Exec("hold", ExecParams{"taken": queuedTaken, "release": release})

	svc.SetWorker(0)
	close(release)
	requireExited(t, svc.i)

	assert.Len(t, queuedTaken, 0, "queued request taken before stopping")
	assert.Len(t, svc.i.MainChan, 1, "queued request must stay queued")
}

func TestServiceCore_SetWorker_StopsWorkerWhoseHookHandlesEverything(t *testing.T) {
	logger := diag.NewDebugLogger(10)
	svc := NewRandomService(logger)
	svc.SetBackground(true)
	svc.SetWorker(1)

	svc.SetWorker(0)
	requireExited(t, svc.i)

	assert.Equal(t, uint64(0), svc.WorkerCount(), "mismatch worker count")
	assert.Equal(t, []string{
		"INFO Random#1: Process started.",
		"INFO Random#1: Process exited.",
	}, logger.AllMessages())
}

func TestServiceCore_SetWorker_SignalsExitChanAfterLastWorkerExits(t *testing.T) {
	logger := diag.NewDebugLogger(10)
	svc := NewHoldService(logger)
	svc.SetBackground(true)
	release := make(chan struct{})
	svc.SetWorker(2)
	requireTaken(t, svc, 1, release)

	svc.SetWorker(0)
	close(release)
	requireExited(t, svc.i)

	assert.ElementsMatch(t, []string{
		"INFO Hold#1: Process started.",
		"INFO Hold#2: Process started.",
		"INFO Hold#1: Process exited.",
		"INFO Hold#2: Process exited.",
	}, logger.AllMessages())
	assert.Len(t, svc.i.ExitChan, 0, "one stop must send one signal")
}

func TestServiceCore_SetWorker_SignalsExitChanOncePerStop(t *testing.T) {
	logger := diag.NewDebugLogger(10)
	svc := NewRandomService(logger)
	svc.SetBackground(true)
	svc.SetWorker(4)

	svc.SetWorker(0)
	requireExited(t, svc.i)

	assert.ElementsMatch(t, []string{
		"INFO Random#1: Process started.",
		"INFO Random#2: Process started.",
		"INFO Random#3: Process started.",
		"INFO Random#4: Process started.",
		"INFO Random#1: Process exited.",
		"INFO Random#2: Process exited.",
		"INFO Random#3: Process exited.",
		"INFO Random#4: Process exited.",
	}, logger.AllMessages())
	assert.Len(t, svc.i.ExitChan, 0, "one stop must send one signal")
}

func TestServiceCore_WorkerCount_DuringSetWorker(t *testing.T) {
	svc := NewHoldService(diag.NewDebugLogger(10))
	t.Cleanup(func() { svc.SetWorker(0) })
	set := make(chan struct{})

	go func() {
		defer close(set)
		svc.SetWorker(3)
	}()
	during := svc.WorkerCount()
	<-set

	assert.Contains(t, []uint64{0, 3}, during, "worker count read during SetWorker")
	assert.Equal(t, uint64(3), svc.WorkerCount(), "mismatch worker count")
}

func TestServiceCore_Exec(t *testing.T) {
	logger := diag.NewDebugLogger(10)
	svc := NewEchoService(logger)
	svc.SetWorker(1)
	logged := make(chan struct{}, 1)
	svc.Exec("", ExecParams{
		"message": "Hello, World!",
		"logged":  logged,
	})
	requireSignal(t, logged, "message not logged")
	assert.Equal(t, "INFO Echo#1: Message received: Hello, World!", logger.LastMessage(), "invalid message")
}

func TestServiceCore_Exec_Exit_LowersWorkerCount(t *testing.T) {
	svc := NewHoldService(diag.NewDebugLogger(10))
	svc.SetBackground(true)
	t.Cleanup(func() { svc.SetWorker(0) })
	release := make(chan struct{})
	defer close(release)
	svc.SetWorker(2)

	svc.Exec("exit", ExecParams{})
	svc.Exec("exit", ExecParams{})
	requireExited(t, svc.i)
	workerCount := svc.WorkerCount()
	svc.SetWorker(2)

	assert.Equal(t, uint64(0), workerCount, "mismatch worker count")
	requireTaken(t, svc, 2, release)
}

func TestServiceCore_Exec_NilHookState_CountsAsNotHandled(t *testing.T) {
	logger := diag.NewDebugLogger(10)
	svc := NewNilStateService(logger)
	svc.SetBackground(true)
	svc.SetWorker(1)

	svc.Exec("hello", ExecParams{})
	svc.Exec("exit", ExecParams{})
	requireExited(t, svc.i)

	assert.Equal(t, uint64(0), svc.WorkerCount(), "mismatch worker count")
	assert.Equal(t, []string{
		"INFO NilState#1: Process started.",
		"INFO NilState#1: Message received: hello",
		`WARN NilState#1: Ignored command "hello": not handled.`,
		"INFO NilState#1: Message received: exit",
		"INFO NilState#1: Process exited.",
	}, logger.AllMessages())
}

func TestServiceCore_Dispatch_WithoutRouter_DropsCommand(t *testing.T) {
	tests := []struct {
		name      string
		serviceID string
		wantLog   string
	}{
		{
			name:      "to a service",
			serviceID: "Random",
			wantLog:   `WARN Echo: Ignored command "hello" for service "Random": no router.`,
		},
		{
			name:      "to the controller",
			serviceID: "",
			wantLog:   `WARN Echo: Ignored command "hello" for service "": no router.`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logger := diag.NewDebugLogger(10)
			svc := NewEchoService(logger)
			controller := NewServiceController(diag.NewDebugLogger(10))

			assert.NotPanics(t, func() { svc.Dispatch(tt.serviceID, "hello", ExecParams{}) })
			isSet := svc.SetRouter(controller)

			assert.True(t, isSet)
			assert.Equal(t, []string{tt.wantLog}, logger.AllMessages())
			assert.Len(t, controller.i.MainChan, 0, "dropped command must not be delivered after SetRouter")
		})
	}
}

// EchoService logs the "message" param of every message, then sends on the message's "logged" channel.
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
	msg.Params["logged"].(chan struct{}) <- struct{}{}
	return &HookState{Handled: true}
}

// HashService logs the SHA-256 of the "input" param for "sha256"; for "sha256_random", it hashes a value
// requested from the Random service, then sends on the message's "hashed" channel, or logs a warning and sends nothing without a value.
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
		randomResult := randomParams.WaitForReturn()
		message, ok := randomResult.(string)
		if !ok {
			// Send nothing on "hashed": ShutdownService's timeout ends the run, and the test fails on its assertions instead of crashing.
			s.i.Logger.Warnf("%s#%d: No value to hash: %v.", s.i.ServiceID, workerID, randomResult)
			return &HookState{Handled: true}
		}
		messageBuf := []byte(message)
		hashBuf := sha256.Sum256(messageBuf)
		hash := hex.EncodeToString(hashBuf[:])
		s.i.Logger.Infof("%s#%d: Value hashed: %s.", s.i.ServiceID, workerID, hash)
		msg.Params["hashed"].(chan struct{}) <- struct{}{}
	}
	return &HookState{Handled: true}
}

// HoldService keeps the worker that takes a message until the message's release channel
// is closed, and passes "exit" on, so an exit command stops the worker that takes it.
type HoldService struct {
	ServiceCore
	i *ServiceCoreInternal
}

func NewHoldService(logger diag.Logger) *HoldService {
	svc := &HoldService{}
	svc.i = svc.InitServiceCore("Hold", logger, svc.coreProcessHook)
	return svc
}

func (s *HoldService) coreProcessHook(workerID uint64, msg *ServiceMessage) *HookState {
	if msg.Command == "exit" {
		return &HookState{Handled: false}
	}
	msg.Params["taken"].(chan struct{}) <- struct{}{}
	<-msg.Params["release"].(chan struct{})
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
	// Log before returning, so a sender that got the value also sees the log.
	s.i.Logger.Infof("%s#%d: Value randomed: %s.", s.i.ServiceID, workerID, hex)
	msg.Return(hex)
	return &HookState{Handled: true}
}

// ShutdownService sends an exit command to the controller once the message's "after" channel signals,
// or after workerTakeTimeout if it never does.
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
	select {
	case <-msg.Params["after"].(chan struct{}):
	case <-time.After(workerTakeTimeout): // Stop anyway, so a missing signal fails the test instead of hanging Run.
	}
	s.Dispatch("", "exit", ExecParams{})
	return &HookState{Handled: true}
}

// NilStateService logs the command of every message and returns a nil HookState.
type NilStateService struct {
	ServiceCore
	i *ServiceCoreInternal
}

func NewNilStateService(logger diag.Logger) *NilStateService {
	svc := &NilStateService{}
	svc.i = svc.InitServiceCore("NilState", logger, svc.coreProcessHook)
	return svc
}

func (s *NilStateService) coreProcessHook(workerID uint64, msg *ServiceMessage) *HookState {
	s.i.Logger.Infof("%s#%d: Message received: %s", s.i.ServiceID, workerID, msg.Command)
	return nil
}

// Queue count messages, each keeping the worker that takes it until release is closed, and fails the test unless workers take all of them;
// before release is closed, that needs count workers running at once.
func requireTaken(t *testing.T, svc Service, count int, release chan struct{}) {
	t.Helper()
	taken := make(chan struct{}, count)
	for i := 0; i < count; i++ {
		svc.Exec("hold", ExecParams{"taken": taken, "release": release})
	}
	for i := 0; i < count; i++ {
		select {
		case <-taken:
		case <-time.After(workerTakeTimeout):
			t.Fatalf("workers took %d of %d messages", i, count)
		}
	}
}

// Wait until the last Process routine of svc, which must be in background mode, exits, and fail the test unless it does.
func requireExited(t *testing.T, svc *ServiceCoreInternal) {
	t.Helper()
	select {
	case <-svc.ExitChan:
	case <-time.After(workerTakeTimeout):
		t.Fatal("service did not stop")
	}
}

// Wait for a signal on signal, and fail the test with failure unless it arrives.
func requireSignal(t *testing.T, signal <-chan struct{}, failure string) {
	t.Helper()
	select {
	case <-signal:
	case <-time.After(workerTakeTimeout):
		t.Fatal(failure)
	}
}
