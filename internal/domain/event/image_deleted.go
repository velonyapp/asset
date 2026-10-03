package event

import (
	"time"

	"github.com/velonyapp/asset/internal/domain/vo"
)

var _ DomainEvent = (*ImageDeleted)(nil)

type ImageDeleted struct {
	BaseDomainEvent

	tags       vo.Tags
	deleteTime time.Time
}

func NewImageDeleted(
	imageID vo.ImageID,
	tags vo.Tags,
	deleteTime time.Time,
) *ImageDeleted {
	return &ImageDeleted{
		BaseDomainEvent: NewBaseDomainEvent(imageID.String()),

		tags:       tags,
		deleteTime: deleteTime,
	}
}

func (e *ImageDeleted) Tags() vo.Tags {
	return e.tags
}

func (e *ImageDeleted) DeleteTime() time.Time {
	return e.deleteTime
}
