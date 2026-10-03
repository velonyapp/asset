package event

import (
	"time"

	"github.com/velonyapp/asset/internal/domain/vo"
)

var _ DomainEvent = (*ImageObjectExistenceUpdated)(nil)

type ImageObjectExistenceUpdated struct {
	BaseDomainEvent

	tags         vo.Tags
	objectExists bool
	updateTime   time.Time
}

func NewImageObjectExistenceUpdated(
	imageID vo.ImageID,
	tags vo.Tags,
	objectExists bool,
	updateTime time.Time,
) *ImageObjectExistenceUpdated {
	return &ImageObjectExistenceUpdated{
		BaseDomainEvent: NewBaseDomainEvent(imageID.String()),

		tags:         tags,
		objectExists: objectExists,
		updateTime:   updateTime,
	}
}

func (e *ImageObjectExistenceUpdated) Tags() vo.Tags {
	return e.tags
}

func (e *ImageObjectExistenceUpdated) ObjectExists() bool {
	return e.objectExists
}

func (e *ImageObjectExistenceUpdated) UpdateTime() time.Time {
	return e.updateTime
}
