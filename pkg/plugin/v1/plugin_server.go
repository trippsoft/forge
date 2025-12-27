// Copyright (c) Forge
// SPDX-License-Identifier: MPL-2.0

package pluginv1

import (
	context "context"
	sync "sync"

	"github.com/trippsoft/forge/pkg/info"
	"github.com/trippsoft/forge/pkg/plugin"
	"github.com/zclconf/go-cty/cty"
	"github.com/zclconf/go-cty/cty/json"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// PluginV1Server implements the PluginServiceServer and PluginV1ServiceServer interfaces.
type PluginV1Server struct {
	UnimplementedPluginV1ServiceServer
	plugin.UnimplementedPluginServiceServer

	modules map[string]PluginModule

	shutdownChan chan struct{}
	waitGroup    sync.WaitGroup
	mutex        sync.RWMutex
	isShutdown   bool
}

// GetAPIVersion retrieves the API version supported by the plugin server.
func (s *PluginV1Server) GetAPIVersion(
	ctx context.Context,
	request *plugin.GetAPIVersionRequest,
) (*plugin.GetAPIVersionResponse, error) {
	return &plugin.GetAPIVersionResponse{
		ApiVersion: 1,
	}, nil
}

// GetModules retrieves the plugin API version, available plugin modules, and their specifications.
func (s *PluginV1Server) GetModules(
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

// RunModule executes the specified plugin module with the provided input.
func (s *PluginV1Server) RunModule(
	ctx context.Context,
	request *RunModuleRequest,
) (*RunModuleResponse, error) {

	s.mutex.RLock()
	if s.isShutdown {
		s.mutex.RUnlock()
		return nil, status.Error(codes.Unavailable, "server is shutting down")
	}
	s.waitGroup.Add(1)
	s.mutex.RUnlock()

	defer s.waitGroup.Done()

	module, exists := s.modules[request.ModuleName]
	if !exists {
		return nil, status.Errorf(codes.NotFound, "module %q not found", request.ModuleName)
	}

	input := make(map[string]cty.Value, len(request.Input))
	for key, value := range request.Input {
		ctyValue, err := json.Unmarshal(value, cty.DynamicPseudoType)
		if err != nil {
			return nil, status.Error(codes.InvalidArgument, err.Error())
		}
		input[key] = ctyValue
	}

	hostInfo := info.NewHostInfo()
	hostInfo.CopyFrom(request.HostInfo)
	result := module.RunModule(hostInfo, input, request.WhatIf)

	return &RunModuleResponse{
		Result: result,
	}, nil
}

// Shutdown initiates a graceful shutdown of the plugin server.
func (s *PluginV1Server) Shutdown(
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

// NewPluginServer creates a new instance of PluginServer with the provided modules.
func NewPluginServer(modules ...PluginModule) *PluginV1Server {
	moduleMap := make(map[string]PluginModule, len(modules))
	for _, module := range modules {
		moduleMap[module.Name()] = module
	}

	return &PluginV1Server{
		modules:      moduleMap,
		shutdownChan: make(chan struct{}),
	}
}
