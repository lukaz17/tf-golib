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
	"github.com/tforce-io/tf-golib/diag"
)

const (
	// Number of pending requests a service supports.
	MainChainCapacity = 256
	// Number of pending requests a service supports.
	ExtraChanCapacity = 16
)

// Service interface defines minimum functions for a service.
type Service interface {
	// Return Service Identifier string
	ServiceID() string

	// Return Router for inter-service operability.
	Router() *ServiceRouter

	// Set number of Process routines the service should use to handle requests.
	SetWorker(workerCount uint64)

	// Return number of Process routines the service is managing.
	WorkerCount() uint64

	// Enqueue the request.
	Exec(command string, params ExecParams)

	// Request other service to handle the request via configurated Router.
	Dispatch(serviceID string, command string, params ExecParams)
}

// ServiceCore is the base struct for deriving new service.
// New service need to embed ServiceCore to access to pre-defined pattern.
type ServiceCore struct {
	i *ServiceCoreInternal
}

// ServiceCoreInternal stores internal data of a ServiceCore.
type ServiceCoreInternal struct {
	ServiceID string
	WorkerID  uint64
	Router    *ServiceRouter

	MainChan   chan *ServiceMessage
	ExitChan   chan bool
	Background bool

	WorkerCounter *Uint64ThreadSafe
	WorkerCount   uint64

	Logger diag.Logger

	CoreProcessHook func(workerID uint64, msg *ServiceMessage) *HookState
}

// Init ServiceCore internal and return the reference for later access.
func (s *ServiceCore) InitServiceCore(serviceID string, logger diag.Logger, processHook func(workerID uint64, msg *ServiceMessage) *HookState) *ServiceCoreInternal {
	if s.i != nil {
		return s.i
	}
	s.i = &ServiceCoreInternal{
		ServiceID:     serviceID,
		MainChan:      make(chan *ServiceMessage, MainChainCapacity),
		ExitChan:      make(chan bool, ExtraChanCapacity),
		WorkerCounter: &Uint64ThreadSafe{},

		Logger: logger,

		CoreProcessHook: processHook,
	}
	return s.i
}

// Return Service Identifier string.
func (s ServiceCore) ServiceID() string {
	return s.i.ServiceID
}

// Return Router for inter-service operability.
func (s ServiceCore) Router() *ServiceRouter {
	return s.i.Router
}

// Set router to use for inter-service operability.
func (s ServiceCore) SetRouter(controller *ServiceController) {
	s.i.Router = controller.i.Router
}

// Set number of Process routines the service should use to handle requests.
// A call made while routines from an earlier call are still starting or exiting also takes effect.
func (s ServiceCore) SetWorker(workerCount uint64) {
	s.i.WorkerCounter.Lock()
	current := s.i.WorkerCount
	for i := current; i < workerCount; i++ {
		s.i.WorkerID++
		go s.process(s.i.WorkerID)
	}
	s.i.WorkerCount = workerCount
	s.i.WorkerCounter.Unlock()

	// Send after unlocking: a send blocks while MainChan is full, and the routines that would drain it need this lock to start.
	for i := current; i > workerCount; i-- {
		cmd := &ServiceMessage{
			Command: "exit",
		}
		s.i.MainChan <- cmd
	}
}

// Return number of Process routines the service is managing.
func (s ServiceCore) WorkerCount() uint64 {
	s.i.WorkerCounter.Lock()
	defer s.i.WorkerCounter.Unlock()
	return s.i.WorkerCount
}

// Enqueue the request.
func (s ServiceCore) Exec(command string, params ExecParams) {
	msg := &ServiceMessage{
		Command: command,
		Params:  params,
	}
	s.i.MainChan <- msg
}

// Request other service to handle the request via configurated Router.
func (s ServiceCore) Dispatch(serviceID string, command string, params ExecParams) {
	s.i.Router.Forward(serviceID, command, params)
}

// Process routine to handle the request.
func (s ServiceCore) process(workerID uint64) {
	s.i.WorkerCounter.Add(1)
	s.i.Logger.Infof("%s#%d: Process started.", s.i.ServiceID, workerID)
	status := InitState
	for status != ExitState {
		msg := <-s.i.MainChan
		if s.i.CoreProcessHook != nil {
			hookState := s.i.CoreProcessHook(workerID, msg)
			if hookState.Handled {
				continue
			}
		}
		if msg.Command == "exit" {
			status = ExitState
			continue
		}
	}
	s.i.WorkerCounter.Lock()
	s.i.WorkerCounter.SubNoLock(1)
	background := s.i.Background
	s.i.WorkerCounter.Unlock()
	s.i.Logger.Infof("%s#%d: Process exited.", s.i.ServiceID, workerID)
	if background {
		s.i.ExitChan <- true
	}
}

// Set whether exiting Process routines signal ExitChan.
func (s ServiceCore) setBackground(background bool) {
	s.i.WorkerCounter.Lock()
	defer s.i.WorkerCounter.Unlock()
	s.i.Background = background
}
