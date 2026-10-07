package entity

import (
	"errors"
	"time"

	"github.com/velonyapp/asset/internal/domain/event"
	"github.com/velonyapp/asset/internal/domain/vo"
)

var (
	ErrImageDeleted          = errors.New("image is deleted")
	ErrImageAlreadyDeleted   = errors.New("image is already deleted")
	ErrImageNotPending       = errors.New("image is not pending")
	ErrImageAlreadyUploaded  = errors.New("image is already uploaded")
	ErrImageNotUploaded      = errors.New("image is not uploaded")
	ErrImageProcessed        = errors.New("image is processed")
	ErrImageAlreadyProcessed = errors.New("image is already processed")
	ErrImageObjectKeysEqual  = errors.New("image source object key and object key must differ")
)

type Image struct {
	id              vo.ImageID
	tags            vo.Tags
	sourceObjectKey vo.ObjectKey
	objectKey       vo.ObjectKey
	state           vo.ImageState
	createTime      time.Time
	deleteTime      *time.Time

	domainEvents []event.DomainEvent
}

func NewImage(
	tags vo.Tags,
	sourceObjectKey vo.ObjectKey,
	objectKey vo.ObjectKey,
	now time.Time,
) (*Image, error) {
	if sourceObjectKey.Equal(objectKey) {
		return nil, ErrImageObjectKeysEqual
	}

	id := vo.GenerateImageID()
	state := vo.ImageStatePending

	a := &Image{
		id:              id,
		tags:            tags,
		sourceObjectKey: sourceObjectKey,
		objectKey:       objectKey,
		state:           state,
		createTime:      now,
	}

	a.recordEvent(
		event.NewImageCreated(
			id,
			tags,
			objectKey,
			state,
			now,
		),
	)

	return a, nil
}

func ReconstituteImage(
	id vo.ImageID,
	tags vo.Tags,
	sourceObjectKey vo.ObjectKey,
	objectKey vo.ObjectKey,
	state vo.ImageState,
	createTime time.Time,
	deleteTime *time.Time,
) *Image {
	a := &Image{
		id:              id,
		tags:            tags,
		sourceObjectKey: sourceObjectKey,
		objectKey:       objectKey,
		state:           state,
		createTime:      createTime,
	}

	if deleteTime != nil {
		value := *deleteTime
		a.deleteTime = &value
	}

	return a
}

func (a *Image) ID() vo.ImageID {
	return a.id
}

func (a *Image) Tags() vo.Tags {
	return a.tags
}

func (a *Image) SourceObjectKey() vo.ObjectKey {
	return a.sourceObjectKey
}

func (a *Image) ObjectKey() vo.ObjectKey {
	return a.objectKey
}

func (a *Image) State() vo.ImageState {
	return a.state
}

func (a *Image) CreateTime() time.Time {
	return a.createTime
}

func (a *Image) DeleteTime() *time.Time {
	if a.deleteTime == nil {
		return nil
	}

	value := *a.deleteTime
	return &value
}

func (a *Image) IsPending() bool {
	return a.state.Equal(vo.ImageStatePending)
}

func (a *Image) IsUploaded() bool {
	return a.state.Equal(vo.ImageStateUploaded)
}

func (a *Image) IsProcessed() bool {
	return a.state.Equal(vo.ImageStateProcessed)
}

func (a *Image) IsDeleted() bool {
	return a.deleteTime != nil
}

func (a *Image) CanUpload() error {
	if a.IsDeleted() {
		return ErrImageDeleted
	}
	if a.IsUploaded() {
		return ErrImageAlreadyUploaded
	}
	if a.IsProcessed() {
		return ErrImageProcessed
	}
	if !a.IsPending() {
		return ErrImageNotPending
	}

	return nil
}

func (a *Image) CanProcess() error {
	if a.IsDeleted() {
		return ErrImageDeleted
	}
	if a.IsProcessed() {
		return ErrImageAlreadyProcessed
	}
	if !a.IsUploaded() {
		return ErrImageNotUploaded
	}

	return nil
}

func (a *Image) CanDelete() error {
	if a.IsDeleted() {
		return ErrImageAlreadyDeleted
	}

	return nil
}

func (a *Image) Upload(now time.Time) error {
	if err := a.CanUpload(); err != nil {
		return err
	}

	newState := vo.ImageStateUploaded

	a.state = newState

	a.recordEvent(
		event.NewImageUpdated(
			a.id,
			a.tags,
			newState,
			now,
		),
	)

	return nil
}

func (a *Image) Process(now time.Time) error {
	if err := a.CanProcess(); err != nil {
		return err
	}

	newState := vo.ImageStateProcessed

	a.state = newState

	a.recordEvent(
		event.NewImageUpdated(
			a.id,
			a.tags,
			newState,
			now,
		),
	)

	return nil
}

func (a *Image) Delete(now time.Time) error {
	if err := a.CanDelete(); err != nil {
		return err
	}

	a.deleteTime = &now

	a.recordEvent(
		event.NewImageDeleted(
			a.id,
			a.tags,
			now,
		),
	)

	return nil
}

func (a *Image) PullEvents() []event.DomainEvent {
	pulled := a.domainEvents
	a.domainEvents = nil
	return pulled
}

func (a *Image) recordEvent(domainEvent event.DomainEvent) {
	a.domainEvents = append(a.domainEvents, domainEvent)
}
