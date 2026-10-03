package event

import (
	"time"

	"github.com/velonyapp/asset/internal/domain/vo"
)

var _ DomainEvent = (*ImageCreated)(nil)

type ImageCreated struct {
	BaseDomainEvent

	tags       vo.Tags
	objectKey  vo.ObjectKey
	createTime time.Time
}

func NewImageCreated(
	imageID vo.ImageID,
	tags vo.Tags,
	objectKey vo.ObjectKey,
	createTime time.Time,
) *ImageCreated {
	return &ImageCreated{
		BaseDomainEvent: NewBaseDomainEvent(imageID.String()),

		tags:       tags,
		objectKey:  objectKey,
		createTime: createTime,
	}
}

func (e *ImageCreated) Tags() vo.Tags {
	return e.tags
}

func (e *ImageCreated) ObjectKey() vo.ObjectKey {
	return e.objectKey
}

func (e *ImageCreated) CreateTime() time.Time {
	return e.createTime
}
