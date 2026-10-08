package event

import (
	"time"

	"github.com/velonyapp/asset/internal/domain/vo"
)

var _ DomainEvent = (*ImageProcessed)(nil)

type ImageProcessed struct {
	BaseDomainEvent

	tags        vo.Tags
	processTime time.Time
}

func NewImageProcessed(
	imageID vo.ImageID,
	tags vo.Tags,
	processTime time.Time,
) ImageProcessed {
	return ImageProcessed{
		BaseDomainEvent: NewBaseDomainEvent(imageID.String()),

		tags:        tags,
		processTime: processTime,
	}
}

func (e ImageProcessed) Tags() vo.Tags {
	return e.tags
}

func (e ImageProcessed) ProcessTime() time.Time {
	return e.processTime
}
