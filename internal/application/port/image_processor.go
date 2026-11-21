package port

import (
	"context"
	"io"

	"github.com/velonyapp/asset/internal/domain/vo"
)

type ImageProcessor interface {
	Process(ctx context.Context, src io.Reader, dst io.Writer, width vo.ImageWidth, height vo.ImageHeight, format vo.ImageFormat) error
}
