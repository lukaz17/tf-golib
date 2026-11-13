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

package multiplex_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/tforce-io/tf-golib/diag"
	"github.com/tforce-io/tf-golib/multiplex"
)

// Bounds each wait on a service, so a missing worker fails the test instead of hanging until go test's timeout.
const serviceTimeout = 5 * time.Second

func TestDownstreamService_Exec_ReturnsResult(t *testing.T) {
	logger := diag.NewDebugLogger(10)
	svc := NewGreeterService(logger)
	svc.SetWorker(1)
	t.Cleanup(func() { svc.SetWorker(0) })
	params := multiplex.ExecParams{"name": "world"}
	params.ExpectReturn()

	svc.Exec("greet", params)
	result := requireReturn(t, params)

	assert.Equal(t, "Hello, world", result)
	assert.Contains(t, logger.AllMessages(), "INFO Greeter#1: Greeted world.")
}

func TestDownstreamService_Exec_DropsUnhandledCommand(t *testing.T) {
	logger := diag.NewDebugLogger(10)
	svc := NewGreeterService(logger)
	svc.SetWorker(1)
	t.Cleanup(func() { svc.SetWorker(0) })
	params := multiplex.ExecParams{"name": "world"}
	params.ExpectReturn()

	svc.Exec("wave", multiplex.ExecParams{})
	svc.Exec("greet", params)
	result := requireReturn(t, params)

	assert.Equal(t, "Hello, world", result, "worker must keep running after an unhandled command")
	assert.Equal(t, uint64(1), svc.WorkerCount(), "mismatch worker count")
	assert.Contains(t, logger.AllMessages(), `WARN Greeter#1: Ignored command "wave": not handled.`)
}

func TestDownstreamService_Dispatch_ReachesRegisteredService(t *testing.T) {
	controller := multiplex.NewServiceController(diag.NewDebugLogger(10))
	svc := NewGreeterService(diag.NewDebugLogger(10))
	isSet := svc.SetRouter(controller)
	svc.SetWorker(1)
	t.Cleanup(func() { svc.SetWorker(0) })
	isRegistered := controller.Register(svc)
	stopped := make(chan struct{})
	go func() {
		defer close(stopped)
		controller.Run(true)
	}()
	params := multiplex.ExecParams{"name": "world"}
	params.ExpectReturn()

	controller.Dispatch(svc.ServiceID(), "greet", params)
	result := requireReturn(t, params)
	controller.Exec("exit", multiplex.ExecParams{})

	assert.True(t, isSet, "router must be set before workers start")
	assert.True(t, isRegistered)
	assert.Equal(t, "Hello, world", result)
	select {
	case <-stopped:
	case <-time.After(serviceTimeout):
		t.Fatal("controller did not stop on exit command")
	}
}

func TestDownstreamService_SetBackground_SignalsStop(t *testing.T) {
	svc := NewGreeterService(diag.NewDebugLogger(10))
	svc.SetBackground(true)
	svc.SetWorker(2)

	svc.Exec("exit", multiplex.ExecParams{})
	svc.Exec("exit", multiplex.ExecParams{})
	select {
	case <-svc.i.ExitChan:
	case <-time.After(serviceTimeout):
		t.Fatal("service did not signal its stop")
	}

	assert.Equal(t, uint64(0), svc.WorkerCount(), "mismatch worker count")
	assert.Len(t, svc.i.ExitChan, 0, "one stop must send one signal")
}

func TestDownstreamService_Dispatch_WithoutRouter_ReturnsErrDropped(t *testing.T) {
	tests := []struct {
		name      string
		serviceID string
	}{
		{name: "to a service", serviceID: "Other"},
		{name: "to the controller", serviceID: ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := NewGreeterService(diag.NewDebugLogger(10))
			params := multiplex.ExecParams{"name": "world"}
			params.ExpectReturn()

			svc.Dispatch(tt.serviceID, "greet", params)
			result := requireReturn(t, params)

			assert.Equal(t, multiplex.ErrDropped, result)
		})
	}
}

func TestDownstreamService_Exec_UnhandledCommand_ReturnsErrDropped(t *testing.T) {
	tests := []struct {
		name string
		hook func(workerID uint64, msg *multiplex.ServiceMessage) *multiplex.HookState
	}{
		{
			name: "hook returns no state",
			hook: func(uint64, *multiplex.ServiceMessage) *multiplex.HookState { return nil },
		},
		{
			name: "hook reports not handled",
			hook: func(uint64, *multiplex.ServiceMessage) *multiplex.HookState {
				return &multiplex.HookState{Handled: false}
			},
		},
		{
			name: "service without hook",
			hook: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logger := diag.NewDebugLogger(10)
			svc := NewHookService(logger, tt.hook)
			svc.SetWorker(1)
			t.Cleanup(func() { svc.SetWorker(0) })
			params := multiplex.ExecParams{}
			params.ExpectReturn()

			svc.Exec("wave", params)
			result := requireReturn(t, params)

			assert.Equal(t, multiplex.ErrDropped, result)
			assert.Contains(t, logger.AllMessages(), `WARN Hook#1: Ignored command "wave": not handled.`)
		})
	}
}

func TestDownstreamService_Exec_AnsweredUnhandledCommand_KeepsResult(t *testing.T) {
	svc := NewHookService(diag.NewDebugLogger(10), func(workerID uint64, msg *multiplex.ServiceMessage) *multiplex.HookState {
		msg.Return("waved")
		return nil
	})
	svc.SetWorker(1)
	t.Cleanup(func() { svc.SetWorker(0) })
	first := multiplex.ExecParams{}
	first.ExpectReturn()
	second := multiplex.ExecParams{}
	second.ExpectReturn()

	svc.Exec("wave", first)
	svc.Exec("wave", second)
	firstResult := requireReturn(t, first)
	secondResult := requireReturn(t, second)

	assert.Equal(t, "waved", firstResult)
	assert.Equal(t, "waved", secondResult, "worker must keep running after its hook answered a request")
}

func TestDownstreamService_MisroutedRequest_ReturnsErrDropped(t *testing.T) {
	tests := []struct {
		name string
		send func(controller *multiplex.ServiceController, params multiplex.ExecParams)
	}{
		{
			name: "unregistered service",
			send: func(c *multiplex.ServiceController, p multiplex.ExecParams) { c.Dispatch("NoSuchService", "greet", p) },
		},
		{
			name: "no target service",
			send: func(c *multiplex.ServiceController, p multiplex.ExecParams) { c.Exec("greet", p) },
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			controller := multiplex.NewServiceController(diag.NewDebugLogger(10))
			params := multiplex.ExecParams{"name": "world"}
			params.ExpectReturn()
			tt.send(controller, params)
			controller.Exec("exit", multiplex.ExecParams{})

			controller.Run(true)
			result := requireReturn(t, params)

			assert.Equal(t, multiplex.ErrDropped, result)
		})
	}
}

func TestDownstreamService_ExitCommand_ReturnsNilOnceStopped(t *testing.T) {
	tests := []struct {
		name string
		stop func(params multiplex.ExecParams) (stopped multiplex.Service, logger *diag.DebugLogger)
	}{
		{
			name: "sent to a service",
			stop: func(params multiplex.ExecParams) (multiplex.Service, *diag.DebugLogger) {
				logger := diag.NewDebugLogger(10)
				svc := NewGreeterService(logger)
				svc.SetWorker(1)
				svc.Exec("exit", params)
				return svc, logger
			},
		},
		{
			name: "dispatched to a service",
			stop: func(params multiplex.ExecParams) (multiplex.Service, *diag.DebugLogger) {
				logger := diag.NewDebugLogger(10)
				svc := NewGreeterService(logger)
				controller := multiplex.NewServiceController(diag.NewDebugLogger(10))
				svc.SetRouter(controller)
				svc.SetWorker(1)
				controller.Register(svc)
				controller.Dispatch(svc.ServiceID(), "exit", params)
				controller.Exec("exit", multiplex.ExecParams{})
				controller.Run(true)
				return svc, logger
			},
		},
		{
			name: "sent to the controller",
			stop: func(params multiplex.ExecParams) (multiplex.Service, *diag.DebugLogger) {
				logger := diag.NewDebugLogger(10)
				controller := multiplex.NewServiceController(logger)
				controller.Exec("exit", params)
				controller.Run(true)
				return controller, logger
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			params := multiplex.ExecParams{}
			params.ExpectReturn()

			stopped, logger := tt.stop(params)
			result := requireReturn(t, params)

			assert.Nil(t, result)
			assert.Contains(t, logger.AllMessages(), "INFO "+stopped.ServiceID()+"#1: Process exited.", "result must arrive once the routine has exited")
			assert.Equal(t, uint64(0), stopped.WorkerCount(), "mismatch worker count")
		})
	}
}

func TestDownstreamService_ExitCommand_AnsweredByHook_KeepsResult(t *testing.T) {
	svc := NewHookService(diag.NewDebugLogger(10), func(workerID uint64, msg *multiplex.ServiceMessage) *multiplex.HookState {
		msg.Return("bye")
		return nil
	})
	svc.SetBackground(true)
	svc.SetWorker(1)
	params := multiplex.ExecParams{}
	params.ExpectReturn()

	svc.Exec("exit", params)
	select {
	case <-svc.i.ExitChan:
	case <-time.After(serviceTimeout):
		t.Fatal("service did not signal its stop")
	}
	result := requireReturn(t, params)

	assert.Equal(t, "bye", result)
	assert.Equal(t, uint64(0), svc.WorkerCount(), "exit command must stop the worker")
}

// GreeterService is written the way a downstream project writes a service: it embeds ServiceCore and implements only the hook.
type GreeterService struct {
	multiplex.ServiceCore
	i *multiplex.ServiceCoreInternal
}

func NewGreeterService(logger diag.Logger) *GreeterService {
	svc := &GreeterService{}
	svc.i = svc.InitServiceCore("Greeter", logger, svc.coreProcessHook)
	return svc
}

// Answer "greet" with a greeting for the "name" param, and leave every other command, "exit" included, to the worker.
func (s *GreeterService) coreProcessHook(workerID uint64, msg *multiplex.ServiceMessage) *multiplex.HookState {
	if msg.Command != "greet" {
		return nil
	}
	name := msg.GetParam("name", "stranger").(string)
	s.i.Logger.Infof("%s#%d: Greeted %s.", s.i.ServiceID, workerID, name)
	msg.Return("Hello, " + name)
	return &multiplex.HookState{Handled: true}
}

// HookService runs the hook it is created with, so each test states how its service treats requests.
type HookService struct {
	multiplex.ServiceCore
	i *multiplex.ServiceCoreInternal
}

func NewHookService(logger diag.Logger, hook func(workerID uint64, msg *multiplex.ServiceMessage) *multiplex.HookState) *HookService {
	svc := &HookService{}
	svc.i = svc.InitServiceCore("Hook", logger, hook)
	return svc
}

// Wait for the result of the request sent with params, and fail the test unless it arrives.
func requireReturn(t *testing.T, params multiplex.ExecParams) interface{} {
	t.Helper()
	result := make(chan interface{}, 1)
	go func() { result <- params.WaitForReturn() }()
	select {
	case r := <-result:
		return r
	case <-time.After(serviceTimeout):
		t.Fatal("no result returned")
		return nil
	}
}
