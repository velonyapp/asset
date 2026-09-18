package event

import "github.com/velonyapp/asset/internal/domain/vo"

type ImageDeleted struct {
	BaseDomainEvent

	DeleteTime vo.Time
}

func NewImageDeleted(
	imageID vo.ImageID,
	createTime vo.Time,
) ImageDeleted {
	return ImageDeleted{
		BaseDomainEvent: NewBaseDomainEvent(imageID.String()),

		DeleteTime: createTime,
	}
}
