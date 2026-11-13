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
	"sync"

	"github.com/tforce-io/tf-golib/diag"
)

// ServiceController is a controller for managing services and routing messages between them.
type ServiceController struct {
	ServiceCore

	services   map[string]Service
	servicesMu sync.RWMutex
}

// Return new ServiceRouter.
func NewServiceController(logger diag.Logger) *ServiceController {
	svc := &ServiceController{}
	svc.InitServiceCore("Controller", logger, svc.coreProcessHook)
	svc.i.Router = &ServiceRouter{
		c: svc,
	}
	svc.services = make(map[string]Service)
	return svc
}

// Register service for routing.
// Safe to call while the controller is running.
func (s *ServiceController) Register(service Service) bool {
	if service.Router() != s.Router() {
		s.i.Logger.Warnf("Service %s's router is invalid", service.ServiceID())
		return false
	}
	serviceID := service.ServiceID()
	s.servicesMu.Lock()
	defer s.servicesMu.Unlock()
	s.services[serviceID] = service
	return true
}

// Unregister a service by serviceID.
// Safe to call while the controller is running.
func (s *ServiceController) Unregister(serviceID string) {
	s.servicesMu.Lock()
	defer s.servicesMu.Unlock()
	delete(s.services, serviceID)
}

// Run the controller.
func (s *ServiceController) Run(background bool) {
	if background {
		s.SetBackground(true)
	}
	s.SetWorker(1)

	if background {
		<-s.i.ExitChan
		s.SetBackground(false)
	}
}

// coreProcessHook is responsible for processing messages in the controller.
func (s *ServiceController) coreProcessHook(workerID uint64, msg *ServiceMessage) *HookState {
	if msg.Extra == nil && msg.Command == "exit" {
		return &HookState{Handled: false}
	}
	extra, ok := msg.Extra.(*ControllerExtra)
	if !ok || extra == nil {
		s.i.Logger.Warnf("%s#%d: Ignored command %q: no target service.", s.i.ServiceID, workerID, msg.Command)
		msg.Params.returnDropped()
		return &HookState{Handled: true}
	}
	// Unlock before Exec, which blocks while target's queue is full.
	s.servicesMu.RLock()
	service, found := s.services[extra.ServiceID]
	s.servicesMu.RUnlock()
	if !found {
		s.i.Logger.Warnf("%s#%d: Ignored command %q: service %q is not registered.", s.i.ServiceID, workerID, msg.Command, extra.ServiceID)
		msg.Params.returnDropped()
		return &HookState{Handled: true}
	}
	service.Exec(msg.Command, msg.Params)
	return &HookState{Handled: true}
}

// ServiceRouter is responsible for routing messages between services.
type ServiceRouter struct {
	c *ServiceController
}

// Forward the message to the specified serviceID.
func (s *ServiceRouter) Forward(serviceID, command string, params ExecParams) {
	msg := &ServiceMessage{
		Command: command,
		Params:  params,
	}
	if serviceID != "" {
		msg.Extra = &ControllerExtra{
			ServiceID: serviceID,
		}
	}
	s.c.i.MainChan <- msg
}

// ControllerExtra contains additional information for request to the controller.
type ControllerExtra struct {
	ServiceID string
}
