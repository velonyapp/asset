
package vo

import (
	"errors"
	"testing"
)

func TestNewImageFormat(t *testing.T) {
	tests := []struct {
		name     string
		value    string
		expected ImageFormat
	}{
		{
			name:     "jpeg",
			value:    "JPEG",
			expected: ImageFormatJPEG,
		},
		{
			name:     "png",
			value:    "PNG",
			expected: ImageFormatPNG,
		},
		{
			name:     "webp",
			value:    "WEBP",
			expected: ImageFormatWebP,
		},
		{
			name:     "avif",
			value:    "AVIF",
			expected: ImageFormatAVIF,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			format, err := NewImageFormat(tt.value)
			if err != nil {
				t.Fatalf("expected no error, got %v", err)
			}

			if !format.Equal(tt.expected) {
				t.Errorf("expected %s, got %s", tt.expected, format)
			}
		})
	}
}

func TestNewImageFormat_WhenInvalid(t *testing.T) {
	format, err := NewImageFormat("INVALID")

	if !errors.Is(err, ErrImageFormatInvalid) {
		t.Fatalf("expected ErrImageFormatInvalid, got %v", err)
	}

	if format != "" {
		t.Errorf("expected empty format, got %s", format)
	}
}

func TestImageFormat_String(t *testing.T) {
	format := ImageFormatJPEG

	if format.String() != "JPEG" {
		t.Errorf("expected JPEG, got %s", format.String())
	}
}

func TestImageFormat_Equal(t *testing.T) {
	if !ImageFormatJPEG.Equal(ImageFormatJPEG) {
		t.Error("expected formats to be equal")
	}

	if ImageFormatJPEG.Equal(ImageFormatPNG) {
		t.Error("expected formats to be different")
	}
}
