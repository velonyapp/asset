package port

import (
	"errors"
	"io"
)

var (
	ErrInvalidResizeDimensions     = errors.New("invalid image resize dimensions")
	ErrUnsupportedResizeFit        = errors.New("unsupported image resize fit")
	ErrUnsupportedImageGravity     = errors.New("unsupported image gravity")
	ErrUnsupportedImageFormat      = errors.New("unsupported image format")
	ErrInvalidImageQuality         = errors.New("invalid image quality")
	ErrInvalidImageBackgroundColor = errors.New("invalid image background color")
	ErrImageUpscaleNotAllowed      = errors.New("image upscale not allowed")
)

type ImageResizeFit string

const (
	ImageResizeFitContain ImageResizeFit = "contain"
	ImageResizeFitCover   ImageResizeFit = "cover"
	ImageResizeFitPad     ImageResizeFit = "pad"
	ImageResizeFitStretch ImageResizeFit = "stretch"
)

type ImageGravity string

const (
	ImageGravityCenter      ImageGravity = "center"
	ImageGravityTop         ImageGravity = "top"
	ImageGravityTopRight    ImageGravity = "top_right"
	ImageGravityRight       ImageGravity = "right"
	ImageGravityBottomRight ImageGravity = "bottom_right"
	ImageGravityBottom      ImageGravity = "bottom"
	ImageGravityBottomLeft  ImageGravity = "bottom_left"
	ImageGravityLeft        ImageGravity = "left"
	ImageGravityTopLeft     ImageGravity = "top_left"
)

type ImageFormat string

const (
	ImageFormatJPEG ImageFormat = "jpeg"
	ImageFormatPNG  ImageFormat = "png"
	ImageFormatWebP ImageFormat = "webp"
	ImageFormatAVIF ImageFormat = "avif"
)

type ImageResize struct {
	Width  *uint32
	Height *uint32

	Fit     ImageResizeFit
	Gravity ImageGravity

	BackgroundColor *string
	AllowUpscale    bool
}

type ImageEncoding struct {
	Format  ImageFormat
	Quality *uint32
}

type ImageTransform struct {
	Resize   *ImageResize
	Encoding *ImageEncoding
}

type ImageProcessor interface {
	Process(image io.Reader, transform *ImageTransform) (io.Reader, error)
}
