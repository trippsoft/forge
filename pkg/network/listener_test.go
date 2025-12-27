// Copyright (c) Forge
// SPDX-License-Identifier: MPL-2.0

package network

import (
	"fmt"
	"net"
	"testing"
)

func TestGetListenerAndPortInRange_SinglePortSuccess(t *testing.T) {
	listener, port, err := GetListenerAndPortInRange(9000, 9000)
	defer func() {
		if listener != nil {
			listener.Close()
		}
	}()

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if port != 9000 {
		t.Fatalf("Expected port 9000, got %d", port)
	}

	if listener == nil {
		t.Fatal("Expected valid listener, got nil")
	}
}

func TestGetListenerAndPortInRange_PortRangeSuccess(t *testing.T) {
	listener, port, err := GetListenerAndPortInRange(9100, 9110)
	defer func() {
		if listener != nil {
			listener.Close()
		}
	}()

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if port < 9100 || port > 9110 {
		t.Fatalf("Expected port in range 9100-9110, got %d", port)
	}

	if listener == nil {
		t.Fatal("Expected valid listener, got nil")
	}
}

func TestGetListenerAndPortInRange_MinPortGreaterThanMax(t *testing.T) {
	_, _, err := GetListenerAndPortInRange(9200, 9100)

	if err == nil {
		t.Fatal("Expected error for minPort > maxPort, got nil")
	}

	if err.Error() != "minimum port cannot be greater than maximum port" {
		t.Fatalf("Expected specific error message, got %v", err)
	}
}

func TestGetListenerAndPortInRange_MinPortLessThan1024(t *testing.T) {
	_, _, err := GetListenerAndPortInRange(1023, 1025)

	if err == nil {
		t.Fatal("Expected error for minPort < 1024, got nil")
	}

	if err.Error() != "minimum port cannot be less than 1024" {
		t.Fatalf("Expected specific error message, got %v", err)
	}
}

func TestGetListenerAndPortInRange_SinglePortUnavailable(t *testing.T) {
	// Create a listener on a specific port
	existingListener, err := net.Listen("tcp", "127.0.0.1:9300")
	if err != nil {
		t.Fatalf("Failed to create blocking listener: %v", err)
	}
	defer existingListener.Close()

	// Try to listen on the same port
	_, _, err = GetListenerAndPortInRange(9300, 9300)

	if err == nil {
		t.Fatal("Expected error when port is already in use, got nil")
	}
}

func TestGetListenerAndPortInRange_RangeExhausted(t *testing.T) {
	// Block all ports in a small range
	listeners := make([]net.Listener, 0)
	defer func() {
		for _, l := range listeners {
			l.Close()
		}
	}()

	for port := 9400; port <= 9402; port++ {
		listener, err := net.Listen("tcp", net.JoinHostPort("127.0.0.1", fmt.Sprintf("%d", port)))
		if err != nil {
			t.Fatalf("Failed to block port %d: %v", port, err)
		}

		listeners = append(listeners, listener)
	}

	_, _, err := GetListenerAndPortInRange(9400, 9402)

	if err == nil {
		t.Fatal("Expected error when all ports in range are exhausted, got nil")
	}

	if err.Error() != "no available ports in range 9400-9402" {
		t.Fatalf("Expected exhaustion error, got %v", err)
	}
}

func TestGetListenerAndPortInRange_MinPortBoundary(t *testing.T) {
	listener, port, err := GetListenerAndPortInRange(1024, 1024)
	defer func() {
		if listener != nil {
			listener.Close()
		}
	}()

	if err != nil {
		t.Fatalf("Expected no error for port 1024, got %v", err)
	}

	if port != 1024 {
		t.Fatalf("Expected port 1024, got %d", port)
	}
}

func TestGetListenerAndPortInRange_RangeWithSkip(t *testing.T) {
	blockedListener, err := net.Listen("tcp", "127.0.0.1:9500")
	if err != nil {
		t.Fatalf("Failed to block port: %v", err)
	}
	defer blockedListener.Close()

	listener, port, err := GetListenerAndPortInRange(9500, 9501)
	defer func() {
		if listener != nil {
			listener.Close()
		}
	}()

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if port == 9500 {
		t.Fatalf("Expected port other than 9500, got %d", port)
	}

	if port < 9500 || port > 9501 {
		t.Fatalf("Expected port in range 9500-9501, got %d", port)
	}
}
