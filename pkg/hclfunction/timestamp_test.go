// Copyright (c) Forge
// SPDX-License-Identifier: MPL-2.0

package hclfunction

import (
	"testing"
	"time"

	"github.com/zclconf/go-cty/cty"
)

func TestTimestamp(t *testing.T) {
	actual, err := Timestamp()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if actual.Type() != cty.String {
		t.Fatalf("expected type %v, got %v", cty.String.FriendlyName(), actual.Type().FriendlyName())
	}

	timestamp := actual.AsString()

	// Verify the timestamp is in RFC3339 format by parsing it
	_, err = time.Parse(time.RFC3339, timestamp)
	if err != nil {
		t.Fatalf("expected valid RFC3339 format, got error while parsing: %v", err)
	}

	// Verify it's recent (within last minute)
	parsedTime, _ := time.Parse(time.RFC3339, timestamp)
	now := time.Now().UTC()
	diff := now.Sub(parsedTime)
	if diff > time.Minute || diff < -time.Minute {
		t.Errorf("Timestamp() returned time %v, which is not recent (diff: %v)", parsedTime, diff)
	}
}

func TestTimestamp_Consistency(t *testing.T) {
	result1, err1 := Timestamp()
	result2, err2 := Timestamp()
	if err1 != nil {
		t.Fatalf("expected no error on first call, got %v", err1)
	}

	if err2 != nil {
		t.Fatalf("expected no error on second call, got %v", err2)
	}

	time1, err := time.Parse(time.RFC3339, result1.AsString())
	if err != nil {
		t.Fatalf("error parsing first timestamp: %v", err)
	}

	time2, err := time.Parse(time.RFC3339, result2.AsString())
	if err != nil {
		t.Fatalf("error parsing second timestamp: %v", err)
	}

	diff := time2.Sub(time1)
	if diff < 0 {
		diff = -diff
	}

	// Two timestamps should be within a second of each other
	if diff > time.Second {
		t.Errorf("Timestamps too far apart: %v", diff)
	}
}

func TestTimestampFunc(t *testing.T) {
	actual, err := TimestampFunc.Call([]cty.Value{})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if actual.Type() != cty.String {
		t.Fatalf("expected type %v, got %v", cty.String.FriendlyName(), actual.Type().FriendlyName())
	}

	timestamp := actual.AsString()
	// Verify the timestamp is in RFC3339 format by parsing it
	_, err = time.Parse(time.RFC3339, timestamp)
	if err != nil {
		t.Fatalf("expected valid RFC3339 format, got error while parsing: %v", err)
	}

	// Verify it's recent (within last minute)
	parsedTime, _ := time.Parse(time.RFC3339, timestamp)
	now := time.Now().UTC()
	diff := now.Sub(parsedTime)
	if diff > time.Minute || diff < -time.Minute {
		t.Errorf("Timestamp() returned time %v, which is not recent (diff: %v)", parsedTime, diff)
	}
}

func TestTimestampFunc_Consistency(t *testing.T) {
	result1, err1 := TimestampFunc.Call([]cty.Value{})
	result2, err2 := TimestampFunc.Call([]cty.Value{})
	if err1 != nil {
		t.Fatalf("expected no error on first call, got %v", err1)
	}

	if err2 != nil {
		t.Fatalf("expected no error on second call, got %v", err2)
	}

	time1, err := time.Parse(time.RFC3339, result1.AsString())
	if err != nil {
		t.Fatalf("error parsing first timestamp: %v", err)
	}

	time2, err := time.Parse(time.RFC3339, result2.AsString())
	if err != nil {
		t.Fatalf("error parsing second timestamp: %v", err)
	}

	diff := time2.Sub(time1)
	if diff < 0 {
		diff = -diff
	}

	// Two timestamps should be within a second of each other
	if diff > time.Second {
		t.Errorf("Timestamps too far apart: %v", diff)
	}
}
