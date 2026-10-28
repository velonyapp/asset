package vo

import (
	"errors"
	"maps"
	"testing"
)

func TestNewMetadata(t *testing.T) {
	tests := []struct {
		name  string
		value map[string]string
		err   error
	}{
		{
			name: "nil metadata",
		},
		{
			name:  "empty metadata",
			value: map[string]string{},
		},
		{
			name:  "empty key",
			value: map[string]string{"": "value"},
			err:   ErrMetadataKeyEmpty,
		},
		{
			name:  "empty value",
			value: map[string]string{"key": ""},
		},
		{
			name:  "unrestricted keys and values",
			value: map[string]string{" ": "", "任意 Key!": "任意 Value!"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			metadata, err := NewMetadata(tt.value)
			if !errors.Is(err, tt.err) {
				t.Fatalf("expected %v, got %v", tt.err, err)
			}
			if tt.err != nil {
				return
			}

			if !maps.Equal(metadata.Map(), tt.value) {
				t.Errorf("expected %v, got %v", tt.value, metadata.Map())
			}
		})
	}
}

func TestMetadata_Map(t *testing.T) {
	value := map[string]string{"key": "value"}
	metadata, err := NewMetadata(value)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	value["key"] = "changed input"
	values := metadata.Map()
	values["key"] = "changed output"
	delete(values, "key")

	if metadata.Map()["key"] != "value" {
		t.Errorf("expected metadata to remain unchanged, got %v", metadata.Map())
	}
}
