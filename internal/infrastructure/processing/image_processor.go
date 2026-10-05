package processing

import (
	"io"
	"strconv"
	"strings"

	"github.com/velonyapp/asset/internal/application/port"

	"github.com/davidbyttow/govips/v2/vips"
)

var _ port.ImageProcessor = (*imageProcessor)(nil)

type imageProcessor struct{}

func NewImageProcessor() port.ImageProcessor {
	return &imageProcessor{}
}

func (p *imageProcessor) Process(src io.Reader, dst io.Writer, opts port.ImageProcessOptions) error {
	if opts.Resize == nil &&
		opts.Encoding == nil &&
		!opts.AutoRotate &&
		!opts.RemoveMetadata {
		return port.ErrNoImageProcessingOptions
	}

	imageRef, err := vips.NewImageFromReader(src)
	if err != nil {
		return err
	}
	defer imageRef.Close()

	if opts.AutoRotate {
		if err := imageRef.AutoRotate(); err != nil {
			return err
		}
	}
	if opts.Resize != nil {
		if err := resize(imageRef, *opts.Resize); err != nil {
			return err
		}
	}
	if opts.RemoveMetadata {
		if err := imageRef.RemoveMetadata(); err != nil {
			return err
		}
	}
	if opts.Encoding != nil {
		return encode(imageRef, dst, *opts.Encoding)
	}

	return encodeNative(imageRef, dst)
}

func resize(imageRef *vips.ImageRef, opts port.ImageResize) error {
	if opts.Fit == "" {
		opts.Fit = port.ImageResizeFitContain
	}
	if opts.Gravity == "" {
		opts.Gravity = port.ImageGravityCenter
	}

	switch opts.Fit {
	case port.ImageResizeFitContain:
		return resizeContain(imageRef, opts)
	case port.ImageResizeFitCover:
		return resizeCover(imageRef, opts)
	case port.ImageResizeFitPad:
		return resizePad(imageRef, opts)
	case port.ImageResizeFitStretch:
		return resizeStretch(imageRef, opts)
	default:
		return port.ErrUnsupportedResizeFit
	}
}

func resizeContain(imageRef *vips.ImageRef, opts port.ImageResize) error {
	if opts.Width == nil && opts.Height == nil {
		return port.ErrInvalidResizeDimensions
	}

	scale := 0.0

	if opts.Width != nil {
		if *opts.Width == 0 {
			return port.ErrInvalidResizeDimensions
		}

		scale = float64(*opts.Width) / float64(imageRef.Width())
	}

	if opts.Height != nil {
		if *opts.Height == 0 {
			return port.ErrInvalidResizeDimensions
		}

		heightScale := float64(*opts.Height) / float64(imageRef.Height())
		if scale == 0 || heightScale < scale {
			scale = heightScale
		}
	}

	if !opts.AllowUpscale && scale > 1 {
		scale = 1
	}

	if scale == 1 {
		return nil
	}

	return imageRef.Resize(scale, vips.KernelLanczos3)
}

func resizeCover(imageRef *vips.ImageRef, opts port.ImageResize) error {
	if opts.Width == nil || opts.Height == nil || *opts.Width == 0 || *opts.Height == 0 {
		return port.ErrInvalidResizeDimensions
	}

	width := int(*opts.Width)
	height := int(*opts.Height)

	widthScale := float64(width) / float64(imageRef.Width())
	heightScale := float64(height) / float64(imageRef.Height())

	scale := widthScale
	if heightScale > scale {
		scale = heightScale
	}

	if !opts.AllowUpscale && scale > 1 {
		scale = 1

		imageRatio := float64(imageRef.Width()) / float64(imageRef.Height())
		targetRatio := float64(width) / float64(height)

		if imageRatio > targetRatio {
			width = int(float64(imageRef.Height()) * targetRatio)
			height = imageRef.Height()
		} else {
			width = imageRef.Width()
			height = int(float64(imageRef.Width()) / targetRatio)
		}
	}

	if scale != 1 {
		if err := imageRef.Resize(scale, vips.KernelLanczos3); err != nil {
			return err
		}
	}

	maxX := imageRef.Width() - width
	maxY := imageRef.Height() - height

	x, y, err := gravityOffset(maxX, maxY, opts.Gravity)
	if err != nil {
		return err
	}

	return imageRef.Crop(x, y, width, height)
}

func resizePad(imageRef *vips.ImageRef, opts port.ImageResize) error {
	if opts.Width == nil || opts.Height == nil || *opts.Width == 0 || *opts.Height == 0 {
		return port.ErrInvalidResizeDimensions
	}

	width := int(*opts.Width)
	height := int(*opts.Height)

	widthScale := float64(width) / float64(imageRef.Width())
	heightScale := float64(height) / float64(imageRef.Height())

	scale := widthScale
	if heightScale < scale {
		scale = heightScale
	}

	if !opts.AllowUpscale && scale > 1 {
		scale = 1
	}

	if scale != 1 {
		if err := imageRef.Resize(scale, vips.KernelLanczos3); err != nil {
			return err
		}
	}

	maxX := width - imageRef.Width()
	maxY := height - imageRef.Height()

	x, y, err := gravityOffset(maxX, maxY, opts.Gravity)
	if err != nil {
		return err
	}

	background, err := parseBackgroundColor(opts.BackgroundColor)
	if err != nil {
		return err
	}

	return imageRef.EmbedBackgroundRGBA(x, y, width, height, background)
}

func resizeStretch(imageRef *vips.ImageRef, opts port.ImageResize) error {
	if opts.Width == nil || opts.Height == nil || *opts.Width == 0 || *opts.Height == 0 {
		return port.ErrInvalidResizeDimensions
	}

	width := int(*opts.Width)
	height := int(*opts.Height)

	if !opts.AllowUpscale {
		if width > imageRef.Width() {
			width = imageRef.Width()
		}
		if height > imageRef.Height() {
			height = imageRef.Height()
		}
	}

	return imageRef.ThumbnailWithSize(
		width,
		height,
		vips.InterestingNone,
		vips.SizeForce,
	)
}

func gravityOffset(
	maxX int,
	maxY int,
	gravity port.ImageGravity,
) (int, int, error) {
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
		return 0, 0, port.ErrUnsupportedImageGravity
	}

	return x, y, nil
}

func parseBackgroundColor(value *string) (*vips.ColorRGBA, error) {
	if value == nil {
		return &vips.ColorRGBA{R: 0, G: 0, B: 0, A: 0}, nil
	}

	hex := strings.TrimPrefix(*value, "#")
	if len(hex) != 6 && len(hex) != 8 {
		return nil, port.ErrInvalidImageBackgroundColor
	}

	r, err := strconv.ParseUint(hex[0:2], 16, 8)
	if err != nil {
		return nil, port.ErrInvalidImageBackgroundColor
	}
	g, err := strconv.ParseUint(hex[2:4], 16, 8)
	if err != nil {
		return nil, port.ErrInvalidImageBackgroundColor
	}
	b, err := strconv.ParseUint(hex[4:6], 16, 8)
	if err != nil {
		return nil, port.ErrInvalidImageBackgroundColor
	}
	a := uint64(255)

	if len(hex) == 8 {
		a, err = strconv.ParseUint(hex[6:8], 16, 8)
		if err != nil {
			return nil, port.ErrInvalidImageBackgroundColor
		}
	}

	return &vips.ColorRGBA{R: uint8(r), G: uint8(g), B: uint8(b), A: uint8(a)}, nil
}

func encodeNative(imageRef *vips.ImageRef, dst io.Writer) error {
	result, _, err := imageRef.ExportNative()
	if err != nil {
		return err
	}

	return writeAll(dst, result)
}

func encode(imageRef *vips.ImageRef, dst io.Writer, opts port.ImageEncoding) error {
	if opts.Quality != nil && (*opts.Quality == 0 || *opts.Quality > 100) {
		return port.ErrInvalidImageQuality
	}

	switch opts.Format {
	case port.ImageFormatJPEG:
		return encodeJPEG(imageRef, dst, opts)
	case port.ImageFormatPNG:
		return encodePNG(imageRef, dst, opts)
	case port.ImageFormatWebP:
		return encodeWebP(imageRef, dst, opts)
	case port.ImageFormatAVIF:
		return encodeAVIF(imageRef, dst, opts)
	default:
		return port.ErrUnsupportedImageFormat
	}
}

func encodeJPEG(
	imageRef *vips.ImageRef,
	dst io.Writer,
	opts port.ImageEncoding,
) error {
	params := vips.NewJpegExportParams()

	if opts.Quality != nil {
		params.Quality = int(*opts.Quality)
	}

	result, _, err := imageRef.ExportJpeg(params)
	if err != nil {
		return err
	}

	return writeAll(dst, result)
}

func encodePNG(imageRef *vips.ImageRef, dst io.Writer, opts port.ImageEncoding) error {
	params := vips.NewPngExportParams()

	if opts.Quality != nil {
		params.Quality = int(*opts.Quality)
	}

	result, _, err := imageRef.ExportPng(params)
	if err != nil {
		return err
	}

	return writeAll(dst, result)
}

func encodeWebP(imageRef *vips.ImageRef, dst io.Writer, opts port.ImageEncoding) error {
	params := vips.NewWebpExportParams()

	if opts.Quality != nil {
		params.Quality = int(*opts.Quality)
	}

	result, _, err := imageRef.ExportWebp(params)
	if err != nil {
		return err
	}

	return writeAll(dst, result)
}

func encodeAVIF(imageRef *vips.ImageRef, dst io.Writer, opts port.ImageEncoding) error {
	params := vips.NewAvifExportParams()

	if opts.Quality != nil {
		params.Quality = int(*opts.Quality)
	}

	result, _, err := imageRef.ExportAvif(params)
	if err != nil {
		return err
	}

	return writeAll(dst, result)
}

func writeAll(dst io.Writer, data []byte) error {
	for len(data) > 0 {
		n, err := dst.Write(data)
		if err != nil {
			return err
		}

		if n == 0 {
			return io.ErrShortWrite
		}

		data = data[n:]
	}

	return nil
}
