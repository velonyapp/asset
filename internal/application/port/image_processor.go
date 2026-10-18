package port

import (
	"io"

	"github.com/velonyapp/asset/internal/domain/vo"
)

type ImageProcessor interface {
	Process(src io.Reader, dst io.Writer, width vo.ImageWidth, height vo.ImageHeight, format vo.ImageFormat) error
}
