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
	tags       []vo.Tag
	storageKey vo.StorageKey
	ready      bool
	createTime time.Time
	deleteTime *time.Time

	domainEvents []event.DomainEvent
}

func NewImage(
	tags []vo.Tag,
	storageKey vo.StorageKey,
	now time.Time,
) *Image {
	imageID := vo.NewImageIDRandom()

	i := &Image{
		id:         imageID,
		tags:       tags,
		storageKey: storageKey,
		ready:      false,
		createTime: now,
	}

	i.recordEvent(
		event.NewImageCreated(
			imageID,
			tags,
			storageKey,
			now,
		),
	)

	return i
}

func ReconstituteImage(
	id vo.ImageID,
	tags []vo.Tag,
	storageKey vo.StorageKey,
	ready bool,
	createTime time.Time,
	deleteTime *time.Time,
) *Image {
	i := &Image{
		id:         id,
		tags:       tags,
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

func (i *Image) Tags() []vo.Tag {
	return append([]vo.Tag(nil), i.tags...)
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

	if i.IsReady() {
		return nil
	}

	i.ready = true

	i.recordEvent(
		event.NewImageFinalized(
			i.id,
			i.tags,
			now,
		),
	)

	return nil
}

func (i *Image) Delete(now time.Time) error {
	if i.IsDeleted() {
		return nil
	}

	i.deleteTime = &now

	i.recordEvent(
		event.NewImageDeleted(
			i.id,
			i.tags,
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
