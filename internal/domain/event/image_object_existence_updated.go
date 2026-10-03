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
}

func NewImageObjectExistenceUpdated(
	imageID vo.ImageID,
	tags vo.Tags,
	objectExists bool,
	occurTime time.Time,
) *ImageObjectExistenceUpdated {
	return &ImageObjectExistenceUpdated{
		BaseDomainEvent: NewBaseDomainEvent(imageID.String(), occurTime),

		tags:         tags,
		objectExists: objectExists,
	}
}

func (e *ImageObjectExistenceUpdated) Tags() vo.Tags {
	return e.tags
}

func (e *ImageObjectExistenceUpdated) ObjectExists() bool {
	return e.objectExists
}
