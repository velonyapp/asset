package image

import (
	"bytes"
	"io"

	"github.com/velonyapp/asset/internal/application/port"

	"github.com/davidbyttow/govips/v2/vips"
)

const webpQuality = 85

type Processor struct{}

func NewProcessor() port.ImageProcessor {
	return &Processor{}
}

func (processor *Processor) Process(image io.Reader, resizeOptions *port.ImageResizeOptions,
) (io.Reader, error) {
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

	if resizeOptions != nil {
		if resizeOptions.Width == 0 || resizeOptions.Height == 0 {
			return nil, port.ErrInvalidResizeDimensions
		}

		switch resizeOptions.Fit {
		case port.ImageResizeFitContain:
			err = thing.ThumbnailWithSize(
				int(resizeOptions.Width),
				int(resizeOptions.Height),
				vips.InterestingNone,
				vips.SizeBoth,
			)

		case port.ImageResizeFitCover:
			err = thing.ThumbnailWithSize(
				int(resizeOptions.Width),
				int(resizeOptions.Height),
				vips.InterestingCentre,
				vips.SizeBoth,
			)

		default:
			return nil, port.ErrUnsupportedResizeFit
		}
		if err != nil {
			return nil, err
		}
	}

	result, _, err := thing.ExportWebp(&vips.WebpExportParams{
		Quality:       webpQuality,
		StripMetadata: true,
	})
	if err != nil {
		return nil, err
	}

	return bytes.NewReader(result), nil
}
