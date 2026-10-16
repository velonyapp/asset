package event

import (
	"time"

	"github.com/velonyapp/asset/internal/domain/vo"
)

var _ DomainEvent = (*ImageActive)(nil)

type ImageActive struct {
	BaseDomainEvent

	tags       vo.Tags
	activeTime time.Time
}

func NewImageActive(
	imageID vo.ImageID,
	tags vo.Tags,
	activeTime time.Time,
) ImageActive {
	return ImageActive{
		BaseDomainEvent: NewBaseDomainEvent(imageID.String()),

		tags:       tags,
		activeTime: activeTime,
	}
}

func (e ImageActive) Tags() vo.Tags {
	return e.tags
}

func (e ImageActive) ActiveTime() time.Time {
	return e.activeTime
}
