package vo

import (
	"errors"
	"strings"
	"testing"
)

func TestNewObjectKey(t *testing.T) {
	tests := []struct {
		name     string
		value    string
		expected string
		err      error
	}{
		{
			name:     "valid single element",
			value:    "image.jpg",
			expected: "image.jpg",
		},
		{
			name:     "valid multiple elements",
			value:    "users/123/profile.jpg",
			expected: "users/123/profile.jpg",
		},
		{
			name:     "allows dash underscore and dot",
			value:    "users/my-folder/image_file.test.jpg",
			expected: "users/my-folder/image_file.test.jpg",
		},
		{
			name:     "allows uppercase letters",
			value:    "Users/ABC/Image.JPG",
			expected: "Users/ABC/Image.JPG",
		},
		{
			name:     "exactly 128 characters",
			value:    strings.Repeat("a", 128),
			expected: strings.Repeat("a", 128),
		},
		{
			name:  "empty value",
			value: "",
			err:   ErrObjectKeyEmpty,
		},
		{
			name:  "leading slash creates empty element",
			value: "/users/image.jpg",
			err:   ErrObjectKeyElementEmpty,
		},
		{
			name:  "trailing slash creates empty element",
			value: "users/image.jpg/",
			err:   ErrObjectKeyElementEmpty,
		},
		{
			name:  "double slash creates empty element",
			value: "users//image.jpg",
			err:   ErrObjectKeyElementEmpty,
		},
		{
			name:  "only slash",
			value: "/",
			err:   ErrObjectKeyElementEmpty,
		},
		{
			name:  "contains space",
			value: "users/my image.jpg",
			err:   ErrObjectKeyInvalidCharacter,
		},
		{
			name:  "contains at sign",
			value: "users/image@123.jpg",
			err:   ErrObjectKeyInvalidCharacter,
		},
		{
			name:  "contains hash",
			value: "users/image#123.jpg",
			err:   ErrObjectKeyInvalidCharacter,
		},
		{
			name:  "contains backslash",
			value: `users\image.jpg`,
			err:   ErrObjectKeyInvalidCharacter,
		},
		{
			name:  "contains unicode character",
			value: "users/café.jpg",
			err:   ErrObjectKeyInvalidCharacter,
		},
		{
			name:  "exceeds 128 characters",
			value: strings.Repeat("a", 129),
			err:   ErrObjectKeyTooLong,
		},
		{
			name:  "multiple elements exceed 128 characters",
			value: strings.Repeat("a", 64) + "/" + strings.Repeat("b", 64),
			err:   ErrObjectKeyTooLong,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			key, err := NewObjectKey(tt.value)

			if tt.err != nil {
				if !errors.Is(err, tt.err) {
					t.Fatalf("expected error %v, got %v", tt.err, err)
				}

				return
			}

			if err != nil {
				t.Fatalf("expected no error, got %v", err)
			}

			if key.String() != tt.expected {
				t.Errorf(
					"expected %q, got %q",
					tt.expected,
					key.String(),
				)
			}
		})
	}
}

func TestObjectKey_String(t *testing.T) {
	key, err := NewObjectKey("users/123/profile.jpg")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	expected := "users/123/profile.jpg"

	if key.String() != expected {
		t.Errorf("expected %q, got %q", expected, key.String())
	}
}

func TestObjectKey_Equal(t *testing.T) {
	tests := []struct {
		name     string
		value1   string
		value2   string
		expected bool
	}{
		{
			name:     "equal keys",
			value1:   "users/123/profile.jpg",
			value2:   "users/123/profile.jpg",
			expected: true,
		},
		{
			name:     "different filename",
			value1:   "users/123/profile.jpg",
			value2:   "users/123/avatar.jpg",
			expected: false,
		},
		{
			name:     "different path",
			value1:   "users/123/profile.jpg",
			value2:   "users/456/profile.jpg",
			expected: false,
		},
		{
			name:     "case sensitive",
			value1:   "Users/123/profile.jpg",
			value2:   "users/123/profile.jpg",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			key1, err := NewObjectKey(tt.value1)
			if err != nil {
				t.Fatalf("failed to create key1: %v", err)
			}

			key2, err := NewObjectKey(tt.value2)
			if err != nil {
				t.Fatalf("failed to create key2: %v", err)
			}

			result := key1.Equal(key2)

			if result != tt.expected {
				t.Errorf(
					"expected Equal() to return %v for %q and %q, got %v",
					tt.expected,
					key1.String(),
					key2.String(),
					result,
				)
			}
		})
	}
}
