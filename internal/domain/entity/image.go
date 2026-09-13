package entity

import (
	"errors"

	"github.com/velonyapp/asset/internal/domain/vo"
)

var (
	ErrImageDeleted = errors.New(
		"image is deleted",
	)
)

type Image struct {
	ID         vo.ImageID
	StorageKey vo.StorageKey
	Status     vo.ImageStatus
	CreateTime vo.Time
	DeleteTime *vo.Time
}

func NewImage(
	StorageKey vo.StorageKey,
	Status vo.ImageStatus,
) *Image {
	now := vo.NewTimeNow()
	imageID := vo.NewImageIDRandom()

	image := &Image{
		ID:         imageID,
		StorageKey: StorageKey,
		Status:     Status,
		CreateTime: now,
	}

	return image
}

func (i *Image) ChangeStatus(status vo.ImageStatus) error {
	if i.DeleteTime != nil {
		return ErrImageDeleted
	}

	i.Status = status

	return nil
}

func (i *Image) Delete() error {
	if i.DeleteTime != nil {
		return ErrImageDeleted
	}

	now := vo.NewTimeNow()

	i.DeleteTime = &now

	return nil
}
