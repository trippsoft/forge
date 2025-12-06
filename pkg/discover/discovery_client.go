// Copyright (c) Forge
// SPDX-License-Identifier: MPL-2.0

package discover

import (
	context "context"
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// DiscoveryClient represents a client for discovering managed system information.
type DiscoveryClient struct {
	address string
	cleanup func()
}

// NewDiscoveryClient creates a new DiscoveryClient with the specified address and cleanup function.
func NewDiscoveryClient(port uint16, cleanup func()) *DiscoveryClient {
	return &DiscoveryClient{
		address: fmt.Sprintf("127.0.0.1:%d", port),
		cleanup: cleanup,
	}
}

func (c *DiscoveryClient) Discover() (*DiscoverInfoResponse, error) {
	connection, err := grpc.NewClient(c.address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		c.cleanup()
		return nil, err
	}
	defer connection.Close()
	defer c.cleanup()

	client := NewDiscoveryPluginClient(connection)
	response, err := client.DiscoverInfo(context.Background(), &DiscoverInfoRequest{})
	if err != nil {
		return nil, err
	}

	_, _ = client.Shutdown(context.Background(), &ShutdownRequest{})

	return response, nil
}
