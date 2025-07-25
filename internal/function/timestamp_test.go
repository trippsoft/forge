package function

import (
	"testing"
	"time"

	"github.com/zclconf/go-cty/cty"
)

func TestTimestamp(t *testing.T) {
	tests := []struct {
		name    string
		wantErr bool
	}{
		{
			name:    "timestamp generation",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := Timestamp()
			if (err != nil) != tt.wantErr {
				t.Errorf("Timestamp() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				// Verify it's a string
				if result.Type() != cty.String {
					t.Errorf("Timestamp() returned type %v, want %v", result.Type(), cty.String)
				}

				// Verify the timestamp is in RFC3339 format by parsing it
				timestampStr := result.AsString()
				_, parseErr := time.Parse(time.RFC3339, timestampStr)
				if parseErr != nil {
					t.Errorf("Timestamp() returned invalid RFC3339 format: %v", parseErr)
				}

				// Verify it's recent (within last minute)
				parsedTime, _ := time.Parse(time.RFC3339, timestampStr)
				now := time.Now().UTC()
				diff := now.Sub(parsedTime)
				if diff > time.Minute || diff < -time.Minute {
					t.Errorf("Timestamp() returned time %v, which is not recent (diff: %v)", parsedTime, diff)
				}
			}
		})
	}
}

func TestTimestampFunc(t *testing.T) {
	tests := []struct {
		name    string
		input   []cty.Value
		wantErr bool
	}{
		{
			name:    "function call with no arguments",
			input:   []cty.Value{},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := TimestampFunc.Call(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("TimestampFunc.Call() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				// Verify it's a string
				if result.Type() != cty.String {
					t.Errorf("TimestampFunc.Call() returned type %v, want %v", result.Type(), cty.String)
				}

				// Verify the timestamp is in RFC3339 format by parsing it
				timestampStr := result.AsString()
				_, parseErr := time.Parse(time.RFC3339, timestampStr)
				if parseErr != nil {
					t.Errorf("TimestampFunc.Call() returned invalid RFC3339 format: %v", parseErr)
				}
			}
		})
	}
}

func TestTimestampConsistency(t *testing.T) {
	// Test that two calls to timestamp are reasonably close in time
	result1, err1 := Timestamp()
	result2, err2 := Timestamp()

	if err1 != nil || err2 != nil {
		t.Errorf("Timestamp() errors: %v, %v", err1, err2)
		return
	}

	time1, parseErr1 := time.Parse(time.RFC3339, result1.AsString())
	time2, parseErr2 := time.Parse(time.RFC3339, result2.AsString())

	if parseErr1 != nil || parseErr2 != nil {
		t.Errorf("Parse errors: %v, %v", parseErr1, parseErr2)
		return
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
