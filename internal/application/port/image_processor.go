package port

import (
	"errors"
	"io"
)

var (
	ErrInvalidResizeDimensions = errors.New("invalid image resize dimensions")
	ErrUnsupportedResizeFit    = errors.New("unsupported image resize fit")
)

type ImageResizeFit string

const (
	ImageResizeFitContain ImageResizeFit = "contain"
	ImageResizeFitCover   ImageResizeFit = "cover"
)

type ImageResizeOptions struct {
	Width  uint32
	Height uint32
	Fit    ImageResizeFit
}

type ImageProcessor interface {
	Process(image io.Reader, options *ImageResizeOptions) (io.Reader, error)
}
