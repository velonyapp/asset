package event

import (
	"time"

	"github.com/velonyapp/asset/internal/domain/vo"
)

var _ DomainEvent = (*ImageCreated)(nil)

type ImageCreated struct {
	BaseDomainEvent

	storageKey vo.StorageKey
}

func NewImageCreated(
	imageID vo.ImageID,
	storageKey vo.StorageKey,
	occurTime time.Time,
) *ImageCreated {
	return &ImageCreated{
		BaseDomainEvent: NewBaseDomainEvent(imageID.Value(), occurTime),

		storageKey: storageKey,
	}
}

func (e *ImageCreated) StorageKey() vo.StorageKey {
	return e.storageKey
}
