package event

import (
	"time"

	"github.com/velonyapp/asset/internal/domain/vo"
)

var _ DomainEvent = (*ImageFinalized)(nil)

type ImageFinalized struct {
	BaseDomainEvent
}

func NewImageFinalized(
	imageID vo.ImageID,
	occurTime time.Time,
) *ImageFinalized {
	return &ImageFinalized{
		BaseDomainEvent: NewBaseDomainEvent(imageID.Value(), occurTime),
	}
}
