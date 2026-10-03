package vo

import (
	"errors"
	"testing"
)

func TestNewTags(t *testing.T) {
	tests := []struct {
		name     string
		value    []string
		expected []string
		err      error
	}{
		{
			name:     "empty tags",
			value:    []string{},
			expected: []string{},
		},
		{
			name: "single tag",
			value: []string{
				"image",
			},
			expected: []string{
				"image",
			},
		},
		{
			name: "multiple tags",
			value: []string{
				"image",
				"video",
				"audio",
			},
			expected: []string{
				"audio",
				"image",
				"video",
			},
		},
		{
			name: "removes duplicate tags",
			value: []string{
				"image",
				"video",
				"image",
				"audio",
				"video",
			},
			expected: []string{
				"audio",
				"image",
				"video",
			},
		},
		{
			name: "returns tags in canonical order",
			value: []string{
				"zebra",
				"alpha",
				"middle",
			},
			expected: []string{
				"alpha",
				"middle",
				"zebra",
			},
		},
		{
			name: "invalid tag",
			value: []string{
				"image",
				"Invalid",
			},
			err: ErrTagInvalidCharacter,
		},
		{
			name: "empty tag",
			value: []string{
				"image",
				"",
			},
			err: ErrTagEmpty,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tags, err := NewTags(tt.value)

			if tt.err != nil {
				if !errors.Is(err, tt.err) {
					t.Fatalf("expected error %v, got %v", tt.err, err)
				}

				return
			}

			if err != nil {
				t.Fatalf("expected no error, got %v", err)
			}

			actual := tags.Strings()

			if len(actual) != len(tt.expected) {
				t.Fatalf(
					"expected %d tags, got %d",
					len(tt.expected),
					len(actual),
				)
			}

			for i := range tt.expected {
				if actual[i] != tt.expected[i] {
					t.Errorf(
						"expected tag at index %d to be %q, got %q",
						i,
						tt.expected[i],
						actual[i],
					)
				}
			}
		})
	}
}

func TestTags_Values(t *testing.T) {
	tags, err := NewTags([]string{
		"video",
		"audio",
		"image",
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	values := tags.Values()

	expected := []string{
		"audio",
		"image",
		"video",
	}

	if len(values) != len(expected) {
		t.Fatalf(
			"expected %d tags, got %d",
			len(expected),
			len(values),
		)
	}

	for i := range expected {
		if values[i].String() != expected[i] {
			t.Errorf(
				"expected tag at index %d to be %q, got %q",
				i,
				expected[i],
				values[i].String(),
			)
		}
	}
}

func TestTags_Strings(t *testing.T) {
	tags, err := NewTags([]string{
		"video",
		"audio",
		"image",
		"video",
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	actual := tags.Strings()

	expected := []string{
		"audio",
		"image",
		"video",
	}

	if len(actual) != len(expected) {
		t.Fatalf(
			"expected %d tags, got %d",
			len(expected),
			len(actual),
		)
	}

	for i := range expected {
		if actual[i] != expected[i] {
			t.Errorf(
				"expected tag at index %d to be %q, got %q",
				i,
				expected[i],
				actual[i],
			)
		}
	}
}

func TestTags_Equal(t *testing.T) {
	tests := []struct {
		name     string
		value1   []string
		value2   []string
		expected bool
	}{
		{
			name:     "empty tags are equal",
			value1:   []string{},
			value2:   []string{},
			expected: true,
		},
		{
			name: "same tags",
			value1: []string{
				"audio",
				"image",
				"video",
			},
			value2: []string{
				"audio",
				"image",
				"video",
			},
			expected: true,
		},
		{
			name: "order does not matter",
			value1: []string{
				"audio",
				"image",
				"video",
			},
			value2: []string{
				"video",
				"audio",
				"image",
			},
			expected: true,
		},
		{
			name: "duplicates do not matter",
			value1: []string{
				"audio",
				"image",
				"image",
				"video",
			},
			value2: []string{
				"video",
				"audio",
				"image",
			},
			expected: true,
		},
		{
			name: "different tags",
			value1: []string{
				"audio",
				"image",
			},
			value2: []string{
				"audio",
				"video",
			},
			expected: false,
		},
		{
			name: "different number of tags",
			value1: []string{
				"audio",
				"image",
			},
			value2: []string{
				"audio",
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tags1, err := NewTags(tt.value1)
			if err != nil {
				t.Fatalf("failed to create tags1: %v", err)
			}

			tags2, err := NewTags(tt.value2)
			if err != nil {
				t.Fatalf("failed to create tags2: %v", err)
			}

			result := tags1.Equal(tags2)

			if result != tt.expected {
				t.Errorf(
					"expected Equal() to return %v, got %v",
					tt.expected,
					result,
				)
			}
		})
	}
}
