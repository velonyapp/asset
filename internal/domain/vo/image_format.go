package vo

import "errors"

var (
	ErrImageFormatInvalid = errors.New("invalid image format")
)

type ImageFormat string

const (
	ImageFormatJPEG ImageFormat = "JPEG"
	ImageFormatPNG  ImageFormat = "PNG"
	ImageFormatWebP ImageFormat = "WEBP"
	ImageFormatAVIF ImageFormat = "AVIF"
)

func NewImageFormat(value string) (ImageFormat, error) {
	format := ImageFormat(value)

	switch format {
	case ImageFormatJPEG,
		ImageFormatPNG,
		ImageFormatWebP,
		ImageFormatAVIF:
		return format, nil
	default:
		return "", ErrImageFormatInvalid
	}
}

func (f ImageFormat) String() string {
	return string(f)
}

func (f ImageFormat) Equal(other ImageFormat) bool {
	return f == other
}
