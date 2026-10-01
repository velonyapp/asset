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
		BaseDomainEvent: NewBaseDomainEvent(imageID.Value(), occurTime),

		tags: tags,
	}
}

func (i *ImageDeleted) Tags() []vo.Tag {
	return append([]vo.Tag(nil), i.tags...)
}