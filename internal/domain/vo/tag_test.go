package vo

import (
	"errors"
	"strings"
	"testing"
)

func TestNewTag(t *testing.T) {
	tests := []struct {
		name     string
		value    string
		expected string
		err      error
	}{
		{
			name:     "valid lowercase tag",
			value:    "billing",
			expected: "billing",
		},
		{
			name:     "allows numbers",
			value:    "service123",
			expected: "service123",
		},
		{
			name:     "allows dash underscore and dot",
			value:    "billing.invoice_created-v2",
			expected: "billing.invoice_created-v2",
		},
		{
			name:     "exactly 64 characters",
			value:    strings.Repeat("a", 64),
			expected: strings.Repeat("a", 64),
		},
		{
			name:  "empty value",
			value: "",
			err:   ErrTagEmpty,
		},
		{
			name:  "uppercase character",
			value: "Billing",
			err:   ErrTagInvalidCharacter,
		},
		{
			name:  "contains space",
			value: "billing service",
			err:   ErrTagInvalidCharacter,
		},
		{
			name:  "leading space",
			value: " billing",
			err:   ErrTagInvalidCharacter,
		},
		{
			name:  "trailing space",
			value: "billing ",
			err:   ErrTagInvalidCharacter,
		},
		{
			name:  "contains slash",
			value: "billing/invoice",
			err:   ErrTagInvalidCharacter,
		},
		{
			name:  "contains at sign",
			value: "billing@invoice",
			err:   ErrTagInvalidCharacter,
		},
		{
			name:  "contains hash",
			value: "billing#invoice",
			err:   ErrTagInvalidCharacter,
		},
		{
			name:  "contains unicode character",
			value: "café",
			err:   ErrTagInvalidCharacter,
		},
		{
			name:  "exceeds 64 characters",
			value: strings.Repeat("a", 65),
			err:   ErrTagTooLong,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tag, err := NewTag(tt.value)

			if tt.err != nil {
				if !errors.Is(err, tt.err) {
					t.Fatalf("expected error %v, got %v", tt.err, err)
				}

				return
			}

			if err != nil {
				t.Fatalf("expected no error, got %v", err)
			}

			if tag.String() != tt.expected {
				t.Errorf(
					"expected %q, got %q",
					tt.expected,
					tag.String(),
				)
			}
		})
	}
}

func TestTag_String(t *testing.T) {
	tag, err := NewTag("billing.invoice")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	expected := "billing.invoice"

	if tag.String() != expected {
		t.Errorf("expected %q, got %q", expected, tag.String())
	}
}

func TestTag_Equal(t *testing.T) {
	tests := []struct {
		name     string
		value1   string
		value2   string
		expected bool
	}{
		{
			name:     "equal tags",
			value1:   "billing.invoice",
			value2:   "billing.invoice",
			expected: true,
		},
		{
			name:     "different tags",
			value1:   "billing.invoice",
			value2:   "billing.payment",
			expected: false,
		},
		{
			name:     "different separators",
			value1:   "billing.invoice",
			value2:   "billing-invoice",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tag1, err := NewTag(tt.value1)
			if err != nil {
				t.Fatalf("failed to create tag1: %v", err)
			}

			tag2, err := NewTag(tt.value2)
			if err != nil {
				t.Fatalf("failed to create tag2: %v", err)
			}

			result := tag1.Equal(tag2)

			if result != tt.expected {
				t.Errorf(
					"expected Equal() to return %v for %q and %q, got %v",
					tt.expected,
					tag1.String(),
					tag2.String(),
					result,
				)
			}
		})
	}
}
