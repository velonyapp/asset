package event

import (
	"time"

	"github.com/velonyapp/asset/internal/domain/vo"
)

var _ DomainEvent = (*ImageUpdated)(nil)

type ImageUpdated struct {
	BaseDomainEvent

	tags       vo.Tags
	state      vo.ImageState
	updateTime time.Time
}

func NewImageUpdated(
	imageID vo.ImageID,
	tags vo.Tags,
	state vo.ImageState,
	updateTime time.Time,
) ImageUpdated {
	return ImageUpdated{
		BaseDomainEvent: NewBaseDomainEvent(imageID.String()),

		tags:       tags,
		state:      state,
		updateTime: updateTime,
	}
}

func (e ImageUpdated) Tags() vo.Tags {
	return e.tags
}

func (e ImageUpdated) State() vo.ImageState {
	return e.state
}

func (e ImageUpdated) UpdateTime() time.Time {
	return e.updateTime
}
