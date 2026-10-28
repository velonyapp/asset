package event

import (
	"time"

	"github.com/velonyapp/asset/internal/domain/vo"
)

var _ DomainEvent = (*ImageDeleted)(nil)

type ImageDeleted struct {
	BaseDomainEvent

	tags       vo.Tags
	metadata   vo.Metadata
	deleteTime time.Time
}

func NewImageDeleted(
	imageID vo.ImageID,
	tags vo.Tags,
	metadata vo.Metadata,
	deleteTime time.Time,
) ImageDeleted {
	return ImageDeleted{
		BaseDomainEvent: NewBaseDomainEvent(imageID.String()),

		tags:       tags,
		metadata:   metadata,
		deleteTime: deleteTime,
	}
}

func (e ImageDeleted) Tags() vo.Tags {
	return e.tags
}

func (e ImageDeleted) Metadata() vo.Metadata {
	return e.metadata
}

func (e ImageDeleted) DeleteTime() time.Time {
	return e.deleteTime
}
