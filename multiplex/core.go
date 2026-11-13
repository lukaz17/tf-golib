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
	// Number of unread stop signals ExitChan holds; signals beyond it are dropped.
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
	Router    *ServiceRouter

	MainChan   chan *ServiceMessage
	ExitChan   chan bool
	background bool

	workerID           uint64
	workerCount        uint64
	servingWorkerCount *Uint64ThreadSafe
	runningWorkerCount uint64
	workerExitChan     chan struct{}

	Logger diag.Logger

	CoreProcessHook func(workerID uint64, msg *ServiceMessage) *HookState
}

// Init ServiceCore internal and return the reference for later access.
// processHook receives every request first; a nil HookState counts as not handled, like Handled false.
// A request left unhandled, other than exit, is dropped with a warning: its sender receives ErrDropped unless the hook already returned a result.
// An exit command left unhandled stops the Process routine that takes it: its sender receives nil once the routine has exited, unless the hook already returned a result.
func (s *ServiceCore) InitServiceCore(serviceID string, logger diag.Logger, processHook func(workerID uint64, msg *ServiceMessage) *HookState) *ServiceCoreInternal {
	if s.i != nil {
		return s.i
	}
	s.i = &ServiceCoreInternal{
		ServiceID: serviceID,

		MainChan: make(chan *ServiceMessage, MainChainCapacity),
		ExitChan: make(chan bool, ExtraChanCapacity),

		servingWorkerCount: &Uint64ThreadSafe{},
		workerExitChan:     make(chan struct{}),

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

// Set router to use for inter-service operability, and report whether it was set.
// Call it before SetWorker: while the service has Process routines, including ones still starting or exiting, it keeps the current router and returns false.
// Never call it concurrently with Router or Dispatch.
func (s ServiceCore) SetRouter(controller *ServiceController) bool {
	s.i.servingWorkerCount.Lock()
	defer s.i.servingWorkerCount.Unlock()
	if s.i.workerCount > 0 || s.i.servingWorkerCount.ValueNoLock() > 0 {
		s.i.Logger.Warnf("%s: Ignored router change: service has workers.", s.i.ServiceID)
		return false
	}
	s.i.Router = controller.i.Router
	return true
}

// Set number of Process routines the service should use to handle requests.
// Lowering it stops the extra routines after their current request, ahead of queued requests; with 0, queued requests wait until routines are added again.
// A call made while routines from an earlier call are still starting or stopping also takes effect.
func (s ServiceCore) SetWorker(workerCount uint64) {
	s.i.servingWorkerCount.Lock()
	defer s.i.servingWorkerCount.Unlock()
	s.i.workerCount = workerCount
	routineCount := s.i.servingWorkerCount.ValueNoLock()
	if routineCount > workerCount {
		close(s.i.workerExitChan)
		s.i.workerExitChan = make(chan struct{})
	}
	for ; routineCount < workerCount; routineCount++ {
		s.i.workerID++
		s.i.servingWorkerCount.AddNoLock(1)
		s.i.runningWorkerCount++
		go s.process(s.i.workerID, s.i.workerExitChan)
	}
}

// Return number of Process routines the service is managing.
func (s ServiceCore) WorkerCount() uint64 {
	s.i.servingWorkerCount.Lock()
	defer s.i.servingWorkerCount.Unlock()
	return s.i.workerCount
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
// Without a router, it logs a warning and drops the request, returning ErrDropped to its sender; a later SetRouter does not deliver it.
func (s ServiceCore) Dispatch(serviceID string, command string, params ExecParams) {
	if s.i.Router == nil {
		s.i.Logger.Warnf("%s: Ignored command %q for service %q: no router.", s.i.ServiceID, command, serviceID)
		params.returnDropped()
		return
	}
	s.i.Router.Forward(serviceID, command, params)
}

// Process routine to handle the request.
func (s ServiceCore) process(workerID uint64, workerExitChan <-chan struct{}) {
	s.i.Logger.Infof("%s#%d: Process started.", s.i.ServiceID, workerID)
	var exitCommand *ServiceMessage
	status := InitState
	for status != ExitState {
		select {
		case <-workerExitChan:
			var shouldExit bool
			shouldExit, workerExitChan = s.processExitSignal()
			if shouldExit {
				status = ExitState
			}
			continue
		default:
		}
		var msg *ServiceMessage
		select {
		case <-workerExitChan:
			continue
		case msg = <-s.i.MainChan:
		}
		if s.i.CoreProcessHook != nil {
			hookState := s.i.CoreProcessHook(workerID, msg)
			if hookState != nil && hookState.Handled {
				continue
			}
		}
		if msg.Command == "exit" {
			s.processExitCommand()
			exitCommand = msg
			status = ExitState
			continue
		}
		s.i.Logger.Warnf("%s#%d: Ignored command %q: not handled.", s.i.ServiceID, workerID, msg.Command)
		msg.Params.returnDropped()
	}

	s.i.Logger.Infof("%s#%d: Process exited.", s.i.ServiceID, workerID)
	s.i.servingWorkerCount.Lock()
	s.i.runningWorkerCount--
	isLastWorker := s.i.runningWorkerCount == 0
	background := s.i.background
	s.i.servingWorkerCount.Unlock()

	if exitCommand != nil {
		exitCommand.Params.returnIfUnanswered(nil)
	}
	if background && isLastWorker {
		select {
		case s.i.ExitChan <- true:
		default: // ExitChan full: its unread signals already end the next wait.
		}
	}
}

// Remove the calling Process routine for an exit command, and lower the worker count by one unless it is already 0.
func (s ServiceCore) processExitCommand() {
	s.i.servingWorkerCount.Lock()
	defer s.i.servingWorkerCount.Unlock()
	s.i.servingWorkerCount.SubNoLock(1)
	if s.i.workerCount > 0 {
		s.i.workerCount--
	}
}

// Handle a closed workerExitChan: report whether the calling Process routine should exit, removing it if the routines outnumber the worker count, and otherwise return the current workerExitChan to wait on.
func (s ServiceCore) processExitSignal() (shouldExit bool, workerExit <-chan struct{}) {
	s.i.servingWorkerCount.Lock()
	defer s.i.servingWorkerCount.Unlock()
	if s.i.servingWorkerCount.ValueNoLock() <= s.i.workerCount {
		return false, s.i.workerExitChan
	}
	s.i.servingWorkerCount.SubNoLock(1)
	return true, nil
}

// Set whether the service signals ExitChan each time its last Process routine exits.
// A signal left unread ends the next wait on ExitChan at once; signals beyond ExitChan's capacity are dropped.
func (s ServiceCore) SetBackground(background bool) {
	s.i.servingWorkerCount.Lock()
	defer s.i.servingWorkerCount.Unlock()
	s.i.background = background
}
