package event

import (
	"time"

	"github.com/velonyapp/asset/internal/domain/vo"
)

var _ DomainEvent = (*ImageReconciled)(nil)

type ImageReconciled struct {
	BaseDomainEvent

	tags         vo.Tags
	objectExists bool
	updateTime   time.Time
}

func NewImageReconciled(
	imageID vo.ImageID,
	tags vo.Tags,
	objectExists bool,
	updateTime time.Time,
) ImageReconciled {
	return ImageReconciled{
		BaseDomainEvent: NewBaseDomainEvent(imageID.String()),

		tags:         tags,
		objectExists: objectExists,
		updateTime:   updateTime,
	}
}

func (e ImageReconciled) Tags() vo.Tags {
	return e.tags
}

func (e ImageReconciled) ObjectExists() bool {
	return e.objectExists
}

func (e ImageReconciled) UpdateTime() time.Time {
	return e.updateTime
}
