// Copyright (c) Forge
// SPDX-License-Identifier: MPL-2.0

package network

import (
	"errors"
	"fmt"
	"math/rand"
	"net"
)

func GetListenerAndPortInRange(minPort, maxPort uint16) (net.Listener, uint16, error) {
	if minPort > maxPort {
		return nil, 0, errors.New("minimum port cannot be greater than maximum port")
	}

	if minPort < 1024 {
		return nil, 0, errors.New("minimum port cannot be less than 1024")
	}

	if minPort == maxPort {
		address := fmt.Sprintf("127.0.0.1:%d", minPort)
		listener, err := net.Listen("tcp", address)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to listen on port %d: %w", minPort, err)
		}

		return listener, minPort, nil
	}

	triedPorts := make(map[uint16]struct{})
	for {
		if len(triedPorts) >= int(maxPort-minPort+1) {
			return nil, 0, fmt.Errorf("no available ports in range %d-%d", minPort, maxPort)
		}

		randomPort := minPort + uint16(rand.Intn(int(maxPort-minPort+1)))
		if _, tried := triedPorts[randomPort]; tried {
			continue
		}

		triedPorts[randomPort] = struct{}{}
		address := fmt.Sprintf("127.0.0.1:%d", randomPort)

		listener, err := net.Listen("tcp", address)
		if err == nil {
			return listener, randomPort, nil
		}
	}
}
