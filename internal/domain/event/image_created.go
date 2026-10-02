package event

import (
	"time"

	"github.com/velonyapp/asset/internal/domain/vo"
)

var _ DomainEvent = (*ImageCreated)(nil)

type ImageCreated struct {
	BaseDomainEvent

	tags      []vo.Tag
	objectKey vo.ObjectKey
}

func NewImageCreated(
	imageID vo.ImageID,
	tags []vo.Tag,
	objectKey vo.ObjectKey,
	occurTime time.Time,
) *ImageCreated {
	return &ImageCreated{
		BaseDomainEvent: NewBaseDomainEvent(imageID.String(), occurTime),

		tags:      tags,
		objectKey: objectKey,
	}
}

func (e *ImageCreated) Tags() []vo.Tag {
	return append([]vo.Tag(nil), e.tags...)
}

func (e *ImageCreated) ObjectKey() vo.ObjectKey {
	return e.objectKey
}
