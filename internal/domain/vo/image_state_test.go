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
			name:     "uploaded",
			value:    "UPLOADED",
			expected: ImageStateUploaded,
		},
		{
			name:     "processed",
			value:    "PROCESSED",
			expected: ImageStateProcessed,
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
	state := ImageStateUploaded

	if state.String() != "UPLOADED" {
		t.Errorf("expected UPLOADED, got %s", state.String())
	}
}

func TestImageState_Equal(t *testing.T) {
	if !ImageStateProcessed.Equal(ImageStateProcessed) {
		t.Error("expected states to be equal")
	}

	if ImageStateProcessed.Equal(ImageStatePending) {
		t.Error("expected states to be different")
	}
}