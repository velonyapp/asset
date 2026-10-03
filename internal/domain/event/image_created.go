package event

import (
	"time"

	"github.com/velonyapp/asset/internal/domain/vo"
)

var _ DomainEvent = (*ImageCreated)(nil)

type ImageCreated struct {
	BaseDomainEvent

	tags      vo.Tags
	objectKey vo.ObjectKey
}

func NewImageCreated(
	imageID vo.ImageID,
	tags vo.Tags,
	objectKey vo.ObjectKey,
	occurTime time.Time,
) *ImageCreated {
	return &ImageCreated{
		BaseDomainEvent: NewBaseDomainEvent(imageID.String(), occurTime),

		tags:      tags,
		objectKey: objectKey,
	}
}

func (e *ImageCreated) Tags() vo.Tags {
	return e.tags
}

func (e *ImageCreated) ObjectKey() vo.ObjectKey {
	return e.objectKey
}
