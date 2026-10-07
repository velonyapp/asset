package vo

import (
	"errors"
	"testing"

	"github.com/google/uuid"
)

func TestNewImageID_Valid(t *testing.T) {
	input := "550e8400-e29b-41d4-a716-446655440000"

	id, err := NewImageID(input)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if id.String() != input {
		t.Errorf("expected %q, got %q", input, id.String())
	}
}

func TestNewImageID_Invalid(t *testing.T) {
	input := "not-a-valid-uuid"

	id, err := NewImageID(input)

	if !errors.Is(err, ErrImageIDInvalid) {
		t.Fatalf("expected ErrImageIDInvalid, got %v", err)
	}

	if id.String() != "" {
		t.Errorf("expected empty ImageID, got %q", id.String())
	}
}

func TestGenerateImageID(t *testing.T) {
	id := GenerateImageID()

	if id.String() == "" {
		t.Fatal("expected generated ID to not be empty")
	}

	parsed, err := uuid.Parse(id.String())
	if err != nil {
		t.Fatalf("expected valid UUID, got %q", id.String())
	}

	if parsed.Version() != uuid.Version(7) {
		t.Errorf("expected UUID v7, got version %d", parsed.Version())
	}
}

func TestImageID_Equal(t *testing.T) {
	value := "550e8400-e29b-41d4-a716-446655440000"

	id1, err := NewImageID(value)
	if err != nil {
		t.Fatalf("failed to create id1: %v", err)
	}

	id2, err := NewImageID(value)
	if err != nil {
		t.Fatalf("failed to create id2: %v", err)
	}

	if !id1.Equal(id2) {
		t.Errorf("expected %q to equal %q", id1.String(), id2.String())
	}
}

func TestImageID_NotEqual(t *testing.T) {
	id1, err := NewImageID("550e8400-e29b-41d4-a716-446655440000")
	if err != nil {
		t.Fatalf("failed to create id1: %v", err)
	}

	id2, err := NewImageID("550e8400-e29b-41d4-a716-446655440001")
	if err != nil {
		t.Fatalf("failed to create id2: %v", err)
	}

	if id1.Equal(id2) {
		t.Errorf("expected %q to not equal %q", id1.String(), id2.String())
	}
}
