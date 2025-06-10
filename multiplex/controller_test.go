// Copyright (C) 2025 T-Force I/O
//
// TF GoLib is licensed under the MIT license.
// You should receive a copy of MIT along with this software.
// If not, see <https://opensource.org/license/mit>

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

func TestServiceController_Run_Background(t *testing.T) {
	logger := diag.NewDebugLogger(10)
	svc := NewServiceController(logger)
	go func() {
		time.Sleep(100 * time.Millisecond)
		svc.Exec("exit", ExecParams{})
	}()
	svc.Run(true)
	time.Sleep(100 * time.Millisecond)
	assert.Equal(t, "INFO Controller#1: Process exited.", logger.LastMessage(), "routine not exited properly")
}

func TestServiceController_Dispatch(t *testing.T) {
	logger := diag.NewDebugLogger(10)
	svc := NewServiceController(logger)
	logger2 := diag.NewDebugLogger(10)
	echo := NewEchoService(logger2)
	echo.SetWorker(1)
	echo.SetRouter(svc)
	svc.Register(echo)
	svc.Dispatch(echo.ServiceID(), "", ExecParams{"message": "Hello, World!"})
	go func() {
		time.Sleep(100 * time.Millisecond)
		svc.Exec("exit", ExecParams{})
	}()
	svc.Run(true)
	time.Sleep(100 * time.Millisecond)
	assert.Equal(t, "INFO Controller#1: Process exited.", logger.LastMessage(), "routine not exited properly")
	assert.Equal(t, "INFO Echo#1: Message received: Hello, World!", logger2.LastMessage(), "message is not dispatched properly")
}

func TestServiceController_Integration(t *testing.T) {
	logger := diag.NewDebugLogger(10)
	svc := NewServiceController(logger)
	logger2 := diag.NewDebugLogger(10)
	hash := NewHashService(logger2)
	hash.SetWorker(1)
	hash.SetRouter(svc)
	svc.Register(hash)
	logger3 := diag.NewDebugLogger(10)
	random := NewRandomService(logger3)
	random.SetWorker(1)
	random.SetRouter(svc)
	svc.Register(random)
	logger4 := diag.NewDebugLogger(10)
	shutdown := NewShutdownService(logger4)
	shutdown.SetWorker(1)
	shutdown.SetRouter(svc)
	svc.Register(shutdown)
	svc.Dispatch("Hash", "sha256_random", ExecParams{})
	svc.Dispatch("Shutdown", "", ExecParams{})
	svc.Run(true)
	assert.Equal(t, "INFO Controller#1: Process exited.", logger.LastMessage(), "routine not exited properly")
	assert.Contains(t, logger2.LastMessage(), "INFO Hash#1: Value hashed: ", "invalid message")
	assert.Contains(t, logger3.LastMessage(), "INFO Random#1: Value randomed:", "invalid message")
}
