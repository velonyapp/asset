package entity

import (
	"errors"

	"github.com/velonyapp/asset/internal/domain/event"
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
	Ready      bool
	CreateTime vo.Time
	DeleteTime *vo.Time

	domainEvents []event.DomainEvent
}

func NewImage(
	storageKey vo.StorageKey,
) *Image {
	now := vo.NewTimeNow()
	imageID := vo.NewImageIDRandom()

	image := &Image{
		ID:         imageID,
		StorageKey: storageKey,
		Ready:      false,
		CreateTime: now,
	}

	image.recordEvent(
		event.NewImageCreated(
			imageID,
			storageKey,
			now,
		),
	)

	return image
}

func (img *Image) Finalize() error {
	if img.DeleteTime != nil {
		return ErrImageDeleted
	}

	img.Ready = true

	return nil
}

func (img *Image) Delete() error {
	if img.DeleteTime != nil {
		return ErrImageDeleted
	}

	now := vo.NewTimeNow()

	img.DeleteTime = &now

	img.recordEvent(
		event.NewImageDeleted(
			img.ID,
			now,
		),
	)

	return nil
}

func (img *Image) PullEvents() []event.DomainEvent {
	pulled := img.domainEvents
	img.domainEvents = nil
	return pulled
}

func (img *Image) recordEvent(domainEvent event.DomainEvent) {
	img.domainEvents = append(img.domainEvents, domainEvent)
}
