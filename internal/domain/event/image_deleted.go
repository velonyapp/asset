package event

import (
	"time"

	"github.com/velonyapp/asset/internal/domain/vo"
)

var _ DomainEvent = (*ImageDeleted)(nil)

type ImageDeleted struct {
	BaseDomainEvent

	tags vo.Tags
}

func NewImageDeleted(
	imageID vo.ImageID,
	tags vo.Tags,
	occurTime time.Time,
) *ImageDeleted {
	return &ImageDeleted{
		BaseDomainEvent: NewBaseDomainEvent(imageID.String(), occurTime),

		tags: tags,
	}
}

func (e *ImageDeleted) Tags() vo.Tags {
	return e.tags
}
