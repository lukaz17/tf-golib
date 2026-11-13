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
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/tforce-io/tf-golib/diag"
)

func TestServiceController_Register(t *testing.T) {
	logger := diag.NewDebugLogger(10)
	svc := NewServiceController(logger)
	echo := NewEchoService(logger)
	isRegistered := svc.Register(echo)
	assert.False(t, isRegistered, "service can't be registered without correct router")
	echo.SetRouter(svc)
	isRegistered = svc.Register(echo)
	assert.True(t, isRegistered)
}

func TestServiceController_Unregister(t *testing.T) {
	logger := diag.NewDebugLogger(10)
	svc := NewServiceController(logger)
	echo := NewEchoService(logger)
	echo.SetRouter(svc)
	svc.Register(echo)
	_, found := svc.services[echo.ServiceID()]
	assert.True(t, found, "service should be registered")
	svc.Unregister(echo.ServiceID())
	_, found = svc.services[echo.ServiceID()]
	assert.False(t, found, "service must be unregistered")
	svc.Unregister(echo.ServiceID())
	_, found = svc.services[echo.ServiceID()]
	assert.False(t, found, "service must be unregistered")
}

func TestServiceController_Register_WhileRunning(t *testing.T) {
	logger := diag.NewDebugLogger(10)
	svc := NewServiceController(logger)
	hold := NewHoldService(diag.NewDebugLogger(10))
	hold.SetRouter(svc)
	hold.SetWorker(1)
	t.Cleanup(func() { hold.SetWorker(0) })
	taken := make(chan struct{}, 1)
	release := make(chan struct{})
	close(release)
	svc.Dispatch("NoSuchService", "hello", ExecParams{})
	exited := runInBackground(svc)

	isRegistered := svc.Register(hold)
	svc.Dispatch(hold.ServiceID(), "hold", ExecParams{"taken": taken, "release": release})
	svc.Exec("exit", ExecParams{})
	<-exited

	assert.True(t, isRegistered)
	assert.Equal(t, []string{
		"INFO Controller#1: Process started.",
		`WARN Controller#1: Ignored command "hello": service "NoSuchService" is not registered.`,
		"INFO Controller#1: Process exited.",
	}, logger.AllMessages())
	select {
	case <-taken:
	case <-time.After(workerTakeTimeout):
		t.Fatal("service registered while running did not take the dispatched message")
	}
}

func TestServiceController_Unregister_WhileRunning(t *testing.T) {
	logger := diag.NewDebugLogger(10)
	svc := NewServiceController(logger)
	hold := NewHoldService(diag.NewDebugLogger(10))
	hold.SetRouter(svc)
	svc.Register(hold)
	svc.Dispatch("NoSuchService", "hello", ExecParams{})
	exited := runInBackground(svc)

	svc.Unregister(hold.ServiceID())
	svc.Dispatch(hold.ServiceID(), "hold", ExecParams{})
	svc.Exec("exit", ExecParams{})
	<-exited

	assert.Equal(t, []string{
		"INFO Controller#1: Process started.",
		`WARN Controller#1: Ignored command "hello": service "NoSuchService" is not registered.`,
		`WARN Controller#1: Ignored command "hold": service "Hold" is not registered.`,
		"INFO Controller#1: Process exited.",
	}, logger.AllMessages())
}

func TestServiceController_Run_Background(t *testing.T) {
	logger := diag.NewDebugLogger(10)
	svc := NewServiceController(logger)
	echo := NewEchoService(diag.NewDebugLogger(10))
	echo.SetRouter(svc)
	echo.SetWorker(1)
	t.Cleanup(func() { echo.SetWorker(0) })
	svc.Register(echo)
	logged := make(chan struct{}, 1)
	exited := runInBackground(svc)
	svc.Dispatch(echo.ServiceID(), "", ExecParams{"message": "Hello, World!", "logged": logged})
	requireSignal(t, logged, "controller did not route the message")

	svc.Exec("exit", ExecParams{})
	requireSignal(t, exited, "Run did not return after exit")

	assert.Equal(t, "INFO Controller#1: Process exited.", logger.LastMessage(), "routine not exited properly")
}

func TestServiceController_Run_Background_ExitAlreadyQueued(t *testing.T) {
	logger := diag.NewDebugLogger(10)
	svc := NewServiceController(logger)
	svc.Exec("exit", ExecParams{})

	svc.Run(true)

	assert.Equal(t, "INFO Controller#1: Process exited.", logger.LastMessage(), "routine not exited properly")
}

func TestServiceController_Run_Background_RunsAgainAfterExit(t *testing.T) {
	logger := diag.NewDebugLogger(10)
	svc := NewServiceController(logger)
	svc.Exec("exit", ExecParams{})
	svc.Run(true)
	svc.Exec("exit", ExecParams{})

	exited := runInBackground(svc)

	select {
	case <-exited:
	case <-time.After(workerTakeTimeout):
		t.Fatal("second Run did not start a worker")
	}
	assert.Equal(t, uint64(0), svc.WorkerCount(), "mismatch worker count")
	assert.Equal(t, []string{
		"INFO Controller#1: Process started.",
		"INFO Controller#1: Process exited.",
		"INFO Controller#2: Process started.",
		"INFO Controller#2: Process exited.",
	}, logger.AllMessages())
}

func TestServiceController_Run_DropsMisroutedMessage(t *testing.T) {
	noTarget := `WARN Controller#1: Ignored command "hello": no target service.`
	tests := []struct {
		name    string
		send    func(svc *ServiceController)
		warning string
	}{
		{
			name:    "unregistered service",
			send:    func(svc *ServiceController) { svc.Dispatch("NoSuchService", "hello", ExecParams{}) },
			warning: `WARN Controller#1: Ignored command "hello": service "NoSuchService" is not registered.`,
		},
		{
			name:    "dispatch without service",
			send:    func(svc *ServiceController) { svc.Dispatch("", "hello", ExecParams{}) },
			warning: noTarget,
		},
		{
			name:    "exec on controller",
			send:    func(svc *ServiceController) { svc.Exec("hello", ExecParams{}) },
			warning: noTarget,
		},
		{
			name: "nil target",
			send: func(svc *ServiceController) {
				svc.i.MainChan <- &ServiceMessage{Command: "hello", Extra: (*ControllerExtra)(nil)}
			},
			warning: noTarget,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logger := diag.NewDebugLogger(10)
			svc := NewServiceController(logger)
			tt.send(svc)
			svc.Exec("exit", ExecParams{})

			svc.Run(true)

			assert.Equal(t, []string{
				"INFO Controller#1: Process started.",
				tt.warning,
				"INFO Controller#1: Process exited.",
			}, logger.AllMessages())
		})
	}
}

func TestServiceController_Dispatch(t *testing.T) {
	logger := diag.NewDebugLogger(10)
	svc := NewServiceController(logger)
	logger2 := diag.NewDebugLogger(10)
	echo := NewEchoService(logger2)
	echo.SetRouter(svc)
	echo.SetWorker(1)
	svc.Register(echo)
	logged := make(chan struct{}, 1)
	svc.Dispatch(echo.ServiceID(), "", ExecParams{"message": "Hello, World!", "logged": logged})
	svc.Exec("exit", ExecParams{})
	svc.Run(true)
	requireSignal(t, logged, "message is not dispatched")
	assert.Equal(t, "INFO Controller#1: Process exited.", logger.LastMessage(), "routine not exited properly")
	assert.Equal(t, "INFO Echo#1: Message received: Hello, World!", logger2.LastMessage(), "message is not dispatched properly")
}

func TestServiceController_Integration(t *testing.T) {
	logger := diag.NewDebugLogger(10)
	svc := NewServiceController(logger)
	logger2 := diag.NewDebugLogger(10)
	hash := NewHashService(logger2)
	hash.SetRouter(svc)
	hash.SetWorker(1)
	svc.Register(hash)
	logger3 := diag.NewDebugLogger(10)
	random := NewRandomService(logger3)
	random.SetRouter(svc)
	random.SetWorker(1)
	svc.Register(random)
	logger4 := diag.NewDebugLogger(10)
	shutdown := NewShutdownService(logger4)
	shutdown.SetRouter(svc)
	shutdown.SetWorker(1)
	svc.Register(shutdown)
	hashed := make(chan struct{}, 1)
	svc.Dispatch("Hash", "sha256_random", ExecParams{"hashed": hashed})
	svc.Dispatch("Shutdown", "", ExecParams{"after": hashed})
	svc.Run(true)
	assert.Equal(t, "INFO Controller#1: Process exited.", logger.LastMessage(), "routine not exited properly")
	assert.Contains(t, logger2.LastMessage(), "INFO Hash#1: Value hashed: ", "invalid message")
	assert.Contains(t, logger3.LastMessage(), "INFO Random#1: Value randomed:", "invalid message")
}

// Run the controller in background on its own goroutine; the returned channel is closed once its worker has processed an exit command.
func runInBackground(svc *ServiceController) <-chan struct{} {
	exited := make(chan struct{})
	go func() {
		defer close(exited)
		svc.Run(true)
	}()
	return exited
}
