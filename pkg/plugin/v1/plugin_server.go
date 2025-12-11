// Copyright (c) Forge
// SPDX-License-Identifier: MPL-2.0

package pluginv1

import (
	context "context"
	sync "sync"

	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// PluginServer implements the PluginServiceServer interface.
//
// It serves as a placeholder for future plugin-related functionalities.
type PluginServer struct {
	UnimplementedPluginServiceServer

	modules map[string]PluginModule

	shutdownChan chan struct{}
	waitGroup    sync.WaitGroup
	mutex        sync.RWMutex
	isShutdown   bool
}

func (s *PluginServer) GetModules(
	ctx context.Context,
	request *GetModulesRequest,
) (*GetModulesResponse, error) {

	s.mutex.RLock()
	if s.isShutdown {
		s.mutex.RUnlock()
		return nil, status.Error(codes.Unavailable, "server is shutting down")
	}
	s.waitGroup.Add(1)
	s.mutex.RUnlock()

	defer s.waitGroup.Done()

	modules := make(map[string]*ModuleSpec, len(s.modules))
	for name, module := range s.modules {
		spec, err := module.InputSpec().ToProtobuf()
		if err != nil {
			return nil, status.Error(codes.Internal, err.Error())
		}
		modules[name] = &ModuleSpec{
			Spec: spec,
		}
	}

	return &GetModulesResponse{
		Modules: modules,
	}, nil
}

func (s *PluginServer) Shutdown(
	ctx context.Context,
	request *ShutdownRequest,
) (*ShutdownResponse, error) {

	s.mutex.Lock()
	if s.isShutdown {
		s.mutex.Unlock()
		return nil, status.Error(codes.Unavailable, "server is already shutting down")
	}
	s.isShutdown = true
	close(s.shutdownChan)
	s.mutex.Unlock()

	done := make(chan struct{})
	go func() {
		s.waitGroup.Wait()
		close(done)
	}()

	select {
	case <-done:
		return &ShutdownResponse{}, nil
	case <-ctx.Done():
		return nil, status.Error(codes.DeadlineExceeded, "shutdown timed out")
	}
}
