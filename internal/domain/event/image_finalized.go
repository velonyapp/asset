package event

import (
	"time"

	"github.com/velonyapp/asset/internal/domain/vo"
)

var _ DomainEvent = (*ImageFinalized)(nil)

type ImageFinalized struct {
	BaseDomainEvent

	tags []vo.Tag
}

func NewImageFinalized(
	imageID vo.ImageID,
	tags []vo.Tag,
	occurTime time.Time,
) *ImageFinalized {
	return &ImageFinalized{
		BaseDomainEvent: NewBaseDomainEvent(imageID.Value(), occurTime),

		tags: tags,
	}
}

func (i *ImageFinalized) Tags() []vo.Tag {
	return append([]vo.Tag(nil), i.tags...)
}
