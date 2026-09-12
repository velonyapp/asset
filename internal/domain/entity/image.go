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
	Key        vo.AssetKey
	CreateTime vo.Time
	DeleteTime *vo.Time
}

func NewImage(
	ID vo.ImageID,
	Key vo.AssetKey,
	CreateTime vo.Time,
	DeleteTime *vo.Time,
) *Image {
	now := vo.NewTimeNow()
	imageID := vo.NewImageIDRandom()

	image := &Image{
		ID:         imageID,
		Key:        Key,
		CreateTime: now,
	}

	return image
}

func (u *Image) Delete() error {
	if u.DeleteTime != nil {
		return ErrImageDeleted
	}

	now := vo.NewTimeNow()

	u.DeleteTime = &now

	return nil
}
