package vo

import (
	"errors"
	"testing"
)

func TestNewImageState(t *testing.T) {
	tests := []struct {
		name     string
		value    string
		expected ImageState
	}{
		{
			name:     "pending",
			value:    "PENDING",
			expected: ImageStatePending,
		},
		{
			name:     "confirmed",
			value:    "CONFIRMED",
			expected: ImageStateConfirmed,
		},
		{
			name:     "processing",
			value:    "PROCESSING",
			expected: ImageStateProcessing,
		},
		{
			name:     "active",
			value:    "ACTIVE",
			expected: ImageStateActive,
		},
		{
			name:     "deleted",
			value:    "DELETED",
			expected: ImageStateDeleted,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			state, err := NewImageState(tt.value)
			if err != nil {
				t.Fatalf("expected no error, got %v", err)
			}

			if !state.Equal(tt.expected) {
				t.Errorf("expected %s, got %s", tt.expected, state)
			}
		})
	}
}

func TestNewImageState_WhenInvalid(t *testing.T) {
	state, err := NewImageState("INVALID")

	if !errors.Is(err, ErrImageStateInvalid) {
		t.Fatalf("expected ErrImageStateInvalid, got %v", err)
	}

	if state != "" {
		t.Errorf("expected empty state, got %s", state)
	}
}

func TestImageState_String(t *testing.T) {
	state := ImageStateConfirmed

	if state.String() != "CONFIRMED" {
		t.Errorf("expected CONFIRMED, got %s", state.String())
	}
}

func TestImageState_Equal(t *testing.T) {
	if !ImageStateActive.Equal(ImageStateActive) {
		t.Error("expected states to be equal")
	}

	if ImageStateActive.Equal(ImageStateConfirmed) {
		t.Error("expected states to be different")
	}
}
