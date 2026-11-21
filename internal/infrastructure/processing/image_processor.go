package processing

import (
	"context"
	"io"
	"log/slog"

	"github.com/velonyapp/asset/internal/application/port"
	"github.com/velonyapp/asset/internal/domain/vo"

	"github.com/davidbyttow/govips/v2/vips"
)

var _ port.ImageProcessor = (*imageProcessor)(nil)

type imageProcessor struct {
	log *slog.Logger
}

func NewImageProcessor(
	log *slog.Logger,
) port.ImageProcessor {
	return &imageProcessor{
		log: log,
	}
}

func (p *imageProcessor) Process(ctx context.Context, src io.Reader, dst io.Writer, width vo.ImageWidth, height vo.ImageHeight, format vo.ImageFormat) error {
	data, err := io.ReadAll(src)
	if err != nil {
		p.log.ErrorContext(ctx,
			"failed to read source image",
			slog.Any("error", err),
		)
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
		p.log.ErrorContext(ctx,
			"failed to create image thumbnail",
			slog.Any("error", err),
		)
		return err
	}
	defer imageRef.Close()

	// Remove Metadata
	if err := imageRef.RemoveMetadata(); err != nil {
		p.log.ErrorContext(ctx,
			"failed to remove image metadata",
			slog.Any("error", err),
		)
		return err
	}

	// Encode
	switch format.String() {
	case vo.ImageFormatJPEG.String():
		if err := imageRef.SaveToWriterJpeg(dst, nil); err != nil {
			p.log.ErrorContext(ctx,
				"failed to encode JPEG image",
				slog.Any("error", err),
			)
			return err
		}
	case vo.ImageFormatPNG.String():
		if err := imageRef.SaveToWriterPng(dst, nil); err != nil {
			p.log.ErrorContext(ctx,
				"failed to encode PNG image",
				slog.Any("error", err),
			)
			return err
		}
	case vo.ImageFormatWebP.String():
		if err := imageRef.SaveToWriterWebp(dst, nil); err != nil {
			p.log.ErrorContext(ctx,
				"failed to encode WebP image",
				slog.Any("error", err),
			)
			return err
		}
	case vo.ImageFormatAVIF.String():
		result, _, err := imageRef.ExportAvif(nil)
		if err != nil {
			p.log.ErrorContext(ctx,
				"failed to encode AVIF image",
				slog.Any("error", err),
			)
			return err
		}

		for len(result) > 0 {
			n, err := dst.Write(result)
			if err != nil {
				p.log.ErrorContext(ctx,
					"failed to write AVIF image",
					slog.Any("error", err),
				)
				return err
			}

			if n == 0 {
				p.log.ErrorContext(ctx,
					"failed to write AVIF image",
					slog.Any("error", io.ErrShortWrite),
				)
				return io.ErrShortWrite
			}

			result = result[n:]
		}
	}

	return nil
}
