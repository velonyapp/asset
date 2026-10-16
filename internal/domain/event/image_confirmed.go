package event

import (
	"time"

	"github.com/velonyapp/asset/internal/domain/vo"
)

var _ DomainEvent = (*ImageConfirmed)(nil)

type ImageConfirmed struct {
	BaseDomainEvent

	tags        vo.Tags
	confirmTime time.Time
}

func NewImageConfirmed(
	imageID vo.ImageID,
	tags vo.Tags,
	confirmTime time.Time,
) ImageConfirmed {
	return ImageConfirmed{
		BaseDomainEvent: NewBaseDomainEvent(imageID.String()),

		tags:        tags,
		confirmTime: confirmTime,
	}
}

func (e ImageConfirmed) Tags() vo.Tags {
	return e.tags
}

func (e ImageConfirmed) ConfirmTime() time.Time {
	return e.confirmTime
}
