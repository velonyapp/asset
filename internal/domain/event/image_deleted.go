package event

import (
	"time"

	"github.com/velonyapp/asset/internal/domain/vo"
)

var _ DomainEvent = (*ImageDeleted)(nil)

type ImageDeleted struct {
	BaseDomainEvent

	tags []vo.Tag
}

func NewImageDeleted(
	imageID vo.ImageID,
	tags []vo.Tag,
	occurTime time.Time,
) *ImageDeleted {
	return &ImageDeleted{
		BaseDomainEvent: NewBaseDomainEvent(imageID.String(), occurTime),

		tags: tags,
	}
}

func (e *ImageDeleted) Tags() []vo.Tag {
	return append([]vo.Tag(nil), e.tags...)
}
