// Copyright (c) Forge
// SPDX-License-Identifier: MPL-2.0

package info

import (
	context "context"
	"sync"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// DiscoveryServer implements the DiscoveryPluginServer interface.
//
// It provides information about the host environment and manages graceful shutdown.
type DiscoveryServer struct {
	UnimplementedDiscoveryPluginServer

	shutdownChan chan struct{}
	waitGroup    sync.WaitGroup
	mutex        sync.RWMutex
	isShutdown   bool
}

// DiscoverInfo provides information about the host environment.
func (s *DiscoveryServer) DiscoverInfo(
	ctx context.Context,
	request *DiscoverInfoRequest,
) (*DiscoverInfoResponse, error) {

	s.mutex.RLock()
	if s.isShutdown {
		s.mutex.RUnlock()
		return nil, status.Error(codes.Unavailable, "server is shutting down")
	}
	s.waitGroup.Add(1)
	s.mutex.RUnlock()

	defer s.waitGroup.Done()

	hostInfo, err := discoverHostInfo()
	if err != nil {
		return nil, err
	}

	return &DiscoverInfoResponse{
		HostInfo: hostInfo,
	}, nil
}

// Shutdown initiates a graceful shutdown of the server.
func (s *DiscoveryServer) Shutdown(
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

// WaitForShutdown blocks until the server has completed its shutdown process.
func (s *DiscoveryServer) WaitForShutdown() {
	<-s.shutdownChan
	s.waitGroup.Wait()
}

// NewDiscoveryServer creates a new instance of DiscoveryServer.
func NewDiscoveryServer() *DiscoveryServer {
	return &DiscoveryServer{
		shutdownChan: make(chan struct{}),
	}
}
