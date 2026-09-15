package service

import (
	"bytes"
	"io"
	"strconv"
	"strings"

	"github.com/velonyapp/asset/internal/application/port"

	"github.com/davidbyttow/govips/v2/vips"
)

type ImageProcessor struct{}

func NewImageProcessor() port.ImageProcessor {
	return &ImageProcessor{}
}

func (p *ImageProcessor) Process(image io.Reader, transform *port.ImageTransform) (io.Reader, error) {
	data, err := io.ReadAll(image)
	if err != nil {
		return nil, err
	}

	thing, err := vips.NewImageFromBuffer(data)
	if err != nil {
		return nil, err
	}
	defer thing.Close()

	if err := thing.AutoRotate(); err != nil {
		return nil, err
	}

	if transform != nil && transform.Resize != nil {
		resize := transform.Resize

		fit := resize.Fit
		if fit == "" {
			fit = port.ImageResizeFitContain
		}

		gravity := resize.Gravity
		if gravity == "" {
			gravity = port.ImageGravityCenter
		}

		switch fit {
		case port.ImageResizeFitContain:
			if resize.Width == nil && resize.Height == nil {
				return nil, port.ErrInvalidResizeDimensions
			}

			scale := 0.0

			if resize.Width != nil {
				if *resize.Width == 0 {
					return nil, port.ErrInvalidResizeDimensions
				}

				scale = float64(*resize.Width) / float64(thing.Width())
			}

			if resize.Height != nil {
				if *resize.Height == 0 {
					return nil, port.ErrInvalidResizeDimensions
				}

				heightScale := float64(*resize.Height) / float64(thing.Height())

				if scale == 0 || heightScale < scale {
					scale = heightScale
				}
			}

			if !resize.AllowUpscale && scale > 1 {
				scale = 1
			}

			if scale != 1 {
				if err := thing.Resize(scale, vips.KernelLanczos3); err != nil {
					return nil, err
				}
			}

		case port.ImageResizeFitCover:
			if resize.Width == nil || resize.Height == nil ||
				*resize.Width == 0 || *resize.Height == 0 {
				return nil, port.ErrInvalidResizeDimensions
			}

			width := int(*resize.Width)
			height := int(*resize.Height)

			widthScale := float64(width) / float64(thing.Width())
			heightScale := float64(height) / float64(thing.Height())

			scale := widthScale
			if heightScale > scale {
				scale = heightScale
			}

			if !resize.AllowUpscale && scale > 1 {
				return nil, port.ErrImageUpscaleNotAllowed
			}

			if scale != 1 {
				if err := thing.Resize(scale, vips.KernelLanczos3); err != nil {
					return nil, err
				}
			}

			maxX := thing.Width() - width
			maxY := thing.Height() - height

			x := maxX / 2
			y := maxY / 2

			switch gravity {
			case port.ImageGravityCenter:
			case port.ImageGravityTop:
				y = 0
			case port.ImageGravityTopRight:
				x = maxX
				y = 0
			case port.ImageGravityRight:
				x = maxX
			case port.ImageGravityBottomRight:
				x = maxX
				y = maxY
			case port.ImageGravityBottom:
				y = maxY
			case port.ImageGravityBottomLeft:
				x = 0
				y = maxY
			case port.ImageGravityLeft:
				x = 0
			case port.ImageGravityTopLeft:
				x = 0
				y = 0
			default:
				return nil, port.ErrUnsupportedImageGravity
			}

			if err := thing.Crop(x, y, width, height); err != nil {
				return nil, err
			}

		case port.ImageResizeFitPad:
			if resize.Width == nil || resize.Height == nil ||
				*resize.Width == 0 || *resize.Height == 0 {
				return nil, port.ErrInvalidResizeDimensions
			}

			width := int(*resize.Width)
			height := int(*resize.Height)

			widthScale := float64(width) / float64(thing.Width())
			heightScale := float64(height) / float64(thing.Height())

			scale := widthScale
			if heightScale < scale {
				scale = heightScale
			}

			if !resize.AllowUpscale && scale > 1 {
				scale = 1
			}

			if scale != 1 {
				if err := thing.Resize(scale, vips.KernelLanczos3); err != nil {
					return nil, err
				}
			}

			maxX := width - thing.Width()
			maxY := height - thing.Height()

			x := maxX / 2
			y := maxY / 2

			switch gravity {
			case port.ImageGravityCenter:
			case port.ImageGravityTop:
				y = 0
			case port.ImageGravityTopRight:
				x = maxX
				y = 0
			case port.ImageGravityRight:
				x = maxX
			case port.ImageGravityBottomRight:
				x = maxX
				y = maxY
			case port.ImageGravityBottom:
				y = maxY
			case port.ImageGravityBottomLeft:
				x = 0
				y = maxY
			case port.ImageGravityLeft:
				x = 0
			case port.ImageGravityTopLeft:
				x = 0
				y = 0
			default:
				return nil, port.ErrUnsupportedImageGravity
			}

			background := &vips.ColorRGBA{
				R: 0,
				G: 0,
				B: 0,
				A: 0,
			}

			if resize.BackgroundColor != nil {
				value := strings.TrimPrefix(*resize.BackgroundColor, "#")

				if len(value) != 6 && len(value) != 8 {
					return nil, port.ErrInvalidImageBackgroundColor
				}

				r, err := strconv.ParseUint(value[0:2], 16, 8)
				if err != nil {
					return nil, port.ErrInvalidImageBackgroundColor
				}

				g, err := strconv.ParseUint(value[2:4], 16, 8)
				if err != nil {
					return nil, port.ErrInvalidImageBackgroundColor
				}

				b, err := strconv.ParseUint(value[4:6], 16, 8)
				if err != nil {
					return nil, port.ErrInvalidImageBackgroundColor
				}

				a := uint64(255)

				if len(value) == 8 {
					a, err = strconv.ParseUint(value[6:8], 16, 8)
					if err != nil {
						return nil, port.ErrInvalidImageBackgroundColor
					}
				}

				background = &vips.ColorRGBA{
					R: uint8(r),
					G: uint8(g),
					B: uint8(b),
					A: uint8(a),
				}
			}

			if err := thing.EmbedBackgroundRGBA(
				x,
				y,
				width,
				height,
				background,
			); err != nil {
				return nil, err
			}

		case port.ImageResizeFitStretch:
			if resize.Width == nil || resize.Height == nil ||
				*resize.Width == 0 || *resize.Height == 0 {
				return nil, port.ErrInvalidResizeDimensions
			}

			width := int(*resize.Width)
			height := int(*resize.Height)

			if !resize.AllowUpscale &&
				(width > thing.Width() || height > thing.Height()) {
				return nil, port.ErrImageUpscaleNotAllowed
			}

			if err := thing.ThumbnailWithSize(
				width,
				height,
				vips.InterestingNone,
				vips.SizeForce,
			); err != nil {
				return nil, err
			}

		default:
			return nil, port.ErrUnsupportedResizeFit
		}
	}

	if err := thing.RemoveMetadata(); err != nil {
		return nil, err
	}

	if transform == nil ||
		transform.Encoding == nil ||
		transform.Encoding.Format == "" {
		result, _, err := thing.ExportNative()
		if err != nil {
			return nil, err
		}

		return bytes.NewReader(result), nil
	}

	encoding := transform.Encoding

	if encoding.Quality != nil {
		if *encoding.Quality == 0 || *encoding.Quality > 100 {
			return nil, port.ErrInvalidImageQuality
		}
	}

	var result []byte

	switch encoding.Format {
	case port.ImageFormatJPEG:
		params := vips.NewJpegExportParams()

		if encoding.Quality != nil {
			params.Quality = int(*encoding.Quality)
		}

		result, _, err = thing.ExportJpeg(params)

	case port.ImageFormatPNG:
		params := vips.NewPngExportParams()

		if encoding.Quality != nil {
			params.Quality = int(*encoding.Quality)
		}

		result, _, err = thing.ExportPng(params)

	case port.ImageFormatWebP:
		params := vips.NewWebpExportParams()

		if encoding.Quality != nil {
			params.Quality = int(*encoding.Quality)
		}

		result, _, err = thing.ExportWebp(params)

	case port.ImageFormatAVIF:
		params := vips.NewAvifExportParams()

		if encoding.Quality != nil {
			params.Quality = int(*encoding.Quality)
		}

		result, _, err = thing.ExportAvif(params)

	default:
		return nil, port.ErrUnsupportedImageFormat
	}

	if err != nil {
		return nil, err
	}

	return bytes.NewReader(result), nil
}
