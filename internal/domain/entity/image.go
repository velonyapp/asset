package entity

import (
	"errors"
	"time"

	"github.com/velonyapp/asset/internal/domain/event"
	"github.com/velonyapp/asset/internal/domain/vo"
)

var (
	ErrImageDeleted = errors.New("image is deleted")
)

type Image struct {
	id         vo.ImageID
	storageKey vo.StorageKey
	ready      bool
	createTime time.Time
	deleteTime *time.Time

	domainEvents []event.DomainEvent
}

func NewImage(
	storageKey vo.StorageKey,
	now time.Time,
) *Image {
	imageID := vo.NewImageIDRandom()

	i := &Image{
		id:         imageID,
		storageKey: storageKey,
		ready:      false,
		createTime: now,
	}

	i.recordEvent(
		event.NewImageCreated(
			i.id,
			i.storageKey,
			now,
		),
	)

	return i
}

func ReconstituteImage(
	id vo.ImageID,
	storageKey vo.StorageKey,
	ready bool,
	createTime time.Time,
	deleteTime *time.Time,
) *Image {
	i := &Image{
		id:         id,
		storageKey: storageKey,
		ready:      ready,
		createTime: createTime,
	}

	if deleteTime != nil {
		value := *deleteTime
		i.deleteTime = &value
	}

	return i
}

func (i *Image) ID() vo.ImageID {
	return i.id
}

func (i *Image) StorageKey() vo.StorageKey {
	return i.storageKey
}

func (i *Image) CreateTime() time.Time {
	return i.createTime
}

func (i *Image) DeleteTime() *time.Time {
	if i.deleteTime == nil {
		return nil
	}

	value := *i.deleteTime
	return &value
}

func (i *Image) IsReady() bool {
	return i.ready
}

func (i *Image) IsDeleted() bool {
	return i.deleteTime != nil
}

func (i *Image) Finalize(now time.Time) error {
	if i.IsDeleted() {
		return ErrImageDeleted
	}

	i.ready = true

	i.recordEvent(
		event.NewImageFinalized(
			i.id,
			now,
		),
	)

	return nil
}

func (i *Image) Delete(now time.Time) error {
	if i.IsDeleted() {
		return ErrImageDeleted
	}

	i.deleteTime = &now

	i.recordEvent(
		event.NewImageDeleted(
			i.id,
			now,
		),
	)

	return nil
}

func (i *Image) PullEvents() []event.DomainEvent {
	pulled := i.domainEvents
	i.domainEvents = nil
	return pulled
}

func (i *Image) recordEvent(domainEvent event.DomainEvent) {
	i.domainEvents = append(i.domainEvents, domainEvent)
}
