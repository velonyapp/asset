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
	id           vo.ImageID
	tags         []vo.Tag
	objectKey    vo.ObjectKey
	objectExists bool
	createTime   time.Time
	deleteTime   *time.Time

	domainEvents []event.DomainEvent
}

func NewImage(
	tags []vo.Tag,
	objectKey vo.ObjectKey,
	now time.Time,
) *Image {
	imageID := vo.NewImageIDRandom()

	i := &Image{
		id:           imageID,
		tags:         tags,
		objectKey:    objectKey,
		objectExists: false,
		createTime:   now,
	}

	i.recordEvent(
		event.NewImageCreated(
			imageID,
			tags,
			objectKey,
			now,
		),
	)

	return i
}

func ReconstituteImage(
	id vo.ImageID,
	tags []vo.Tag,
	objectKey vo.ObjectKey,
	objectExist bool,
	createTime time.Time,
	deleteTime *time.Time,
) *Image {
	i := &Image{
		id:           id,
		tags:         tags,
		objectKey:    objectKey,
		objectExists: objectExist,
		createTime:   createTime,
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

func (i *Image) ObjectKey() vo.ObjectKey {
	return i.objectKey
}

func (i *Image) ObjectExists() bool {
	return i.objectExists
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

func (i *Image) IsDeleted() bool {
	return i.deleteTime != nil
}

func (i *Image) UpdateObjectExistence(value bool, now time.Time) error {
	if i.IsDeleted() {
		return ErrImageDeleted
	}

	if value == i.objectExists {
		return nil
	}

	i.recordEvent(
		event.NewImageObjectExistenceUpdated(
			i.id,
			i.tags,
			i.objectExists,
			now,
		),
	)

	return nil
}

func (i *Image) Delete(now time.Time) {
	if i.IsDeleted() {
		return
	}

	i.deleteTime = &now

	i.recordEvent(
		event.NewImageDeleted(
			i.id,
			i.tags,
			now,
		),
	)
}

func (i *Image) PullEvents() []event.DomainEvent {
	pulled := i.domainEvents
	i.domainEvents = nil
	return pulled
}

func (i *Image) recordEvent(domainEvent event.DomainEvent) {
	i.domainEvents = append(i.domainEvents, domainEvent)
}
