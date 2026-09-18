package event

import "github.com/velonyapp/asset/internal/domain/vo"

type ImageCreated struct {
	BaseDomainEvent

	StorageKey vo.StorageKey
	CreateTime vo.Time
}

func NewImageCreated(
	imageID vo.ImageID,
	storageKey vo.StorageKey,
	createTime vo.Time,
) ImageCreated {
	return ImageCreated{
		BaseDomainEvent: NewBaseDomainEvent(imageID.String()),

		StorageKey: storageKey,
		CreateTime: createTime,
	}
}
