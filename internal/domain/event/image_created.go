package event

import (
	"time"

	"github.com/velonyapp/asset/internal/domain/vo"
)

var _ DomainEvent = (*ImageCreated)(nil)

type ImageCreated struct {
	BaseDomainEvent

	tags       []vo.Tag
	storageKey vo.StorageKey
}

func NewImageCreated(
	imageID vo.ImageID,
	tags []vo.Tag,
	storageKey vo.StorageKey,
	occurTime time.Time,
) *ImageCreated {
	return &ImageCreated{
		BaseDomainEvent: NewBaseDomainEvent(imageID.Value(), occurTime),

		tags:       tags,
		storageKey: storageKey,
	}
}

func (i *ImageCreated) Tags() []vo.Tag {
	return append([]vo.Tag(nil), i.tags...)
}

func (e *ImageCreated) StorageKey() vo.StorageKey {
	return e.storageKey
}
