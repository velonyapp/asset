package processing

import (
	"io"

	"github.com/velonyapp/asset/internal/application/port"
	"github.com/velonyapp/asset/internal/domain/vo"

	"github.com/davidbyttow/govips/v2/vips"
)

var _ port.ImageProcessor = (*imageProcessor)(nil)

type imageProcessor struct{}

func NewImageProcessor() port.ImageProcessor {
	return &imageProcessor{}
}

func (p *imageProcessor) Process(src io.Reader, dst io.Writer, width vo.ImageWidth, height vo.ImageHeight, format vo.ImageFormat) error {
	data, err := io.ReadAll(src)
	if err != nil {
		return err
	}

	// Auto Rotate, Resize, Crop Center
	imageRef, err := vips.NewThumbnailWithSizeFromBuffer(
		data,
		width.Int(),
		height.Int(),
		vips.InterestingCentre,
		vips.SizeBoth,
	)
	if err != nil {
		return err
	}
	defer imageRef.Close()

	// Remove Metadata
	if err := imageRef.RemoveMetadata(); err != nil {
		return err
	}

	// Encode
	switch format.String() {
	case vo.ImageFormatJPEG.String():
		return imageRef.SaveToWriterJpeg(dst, nil)
	case vo.ImageFormatPNG.String():
		return imageRef.SaveToWriterPng(dst, nil)
	case vo.ImageFormatWebP.String():
		return imageRef.SaveToWriterWebp(dst, nil)
	case vo.ImageFormatAVIF.String():
		result, _, err := imageRef.ExportAvif(nil)
		if err != nil {
			return err
		}

		for len(result) > 0 {
			n, err := dst.Write(result)
			if err != nil {
				return err
			}

			if n == 0 {
				return io.ErrShortWrite
			}

			result = result[n:]
		}
	}

	return nil
}
