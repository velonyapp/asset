package port

import (
	"errors"
	"io"
)

var (
	ErrNoImageProcessingOptions       = errors.New("no image processing options provided")
	ErrInvalidResizeDimensions        = errors.New("invalid image resize dimensions")
	ErrUnsupportedResizeFit           = errors.New("unsupported image resize fit")
	ErrUnsupportedImageResizeGravity  = errors.New("unsupported image gravity")
	ErrUnsupportedImageEncodingFormat = errors.New("unsupported image format")
	ErrInvalidImageQuality            = errors.New("invalid image quality")
	ErrInvalidImageBackgroundColor    = errors.New("invalid image background color")
)

type ImageResizeFit string

const (
	ImageResizeFitContain ImageResizeFit = "contain"
	ImageResizeFitCover   ImageResizeFit = "cover"
	ImageResizeFitPad     ImageResizeFit = "pad"
	ImageResizeFitStretch ImageResizeFit = "stretch"
)

type ImageResizeGravity string

const (
	ImageResizeGravityCenter      ImageResizeGravity = "center"
	ImageResizeGravityTop         ImageResizeGravity = "top"
	ImageResizeGravityTopRight    ImageResizeGravity = "top_right"
	ImageResizeGravityRight       ImageResizeGravity = "right"
	ImageResizeGravityBottomRight ImageResizeGravity = "bottom_right"
	ImageResizeGravityBottom      ImageResizeGravity = "bottom"
	ImageResizeGravityBottomLeft  ImageResizeGravity = "bottom_left"
	ImageResizeGravityLeft        ImageResizeGravity = "left"
	ImageResizeGravityTopLeft     ImageResizeGravity = "top_left"
)

type ImageEncodingFormat string

const (
	ImageEncodingFormatJPEG ImageEncodingFormat = "jpeg"
	ImageEncodingFormatPNG  ImageEncodingFormat = "png"
	ImageEncodingFormatWebP ImageEncodingFormat = "webp"
	ImageEncodingFormatAVIF ImageEncodingFormat = "avif"
)

type ImageResize struct {
	Width  *uint32
	Height *uint32

	Fit     ImageResizeFit
	Gravity ImageResizeGravity

	BackgroundColor *string
	AllowUpscale    bool
}

type ImageEncoding struct {
	Format  ImageEncodingFormat
	Quality *uint32
}

type ImageProcessOptions struct {
	Resize         *ImageResize
	Encoding       *ImageEncoding
	AutoRotate     bool
	RemoveMetadata bool
}

type ImageProcessor interface {
	Process(src io.Reader, dst io.Writer, opts ImageProcessOptions) error
}
