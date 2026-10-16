package entity

import (
	"errors"
	"time"

	"github.com/velonyapp/asset/internal/domain/event"
	"github.com/velonyapp/asset/internal/domain/vo"
)

var (
	ErrImageDeleted           = errors.New("image is deleted")
	ErrImageAlreadyDeleted    = errors.New("image is already deleted")
	ErrImageNotPending        = errors.New("image is not pending")
	ErrImageAlreadyConfirmed  = errors.New("image is already confirmed")
	ErrImageNotConfirmed      = errors.New("image is not confirmed")
	ErrImageProcessing        = errors.New("image is processing")
	ErrImageAlreadyProcessing = errors.New("image is already processing")
	ErrImageNotProcessing     = errors.New("image is not processing")
	ErrImageActive            = errors.New("image is active")
	ErrImageAlreadyActive     = errors.New("image is already active")
	ErrImageObjectKeysEqual   = errors.New("image source object key and object key must differ")
)

type Image struct {
	id              vo.ImageID
	tags            vo.Tags
	sourceObjectKey vo.ObjectKey
	objectKey       vo.ObjectKey
	state           vo.ImageState
	createTime      time.Time
	updateTime      time.Time

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
		updateTime:      now,
	}

	a.recordEvent(event.NewImageCreated(id, tags, objectKey, now))

	return a, nil
}

func ReconstituteImage(
	id vo.ImageID,
	tags vo.Tags,
	sourceObjectKey vo.ObjectKey,
	objectKey vo.ObjectKey,
	state vo.ImageState,
	createTime time.Time,
	updateTime time.Time,
) *Image {
	return &Image{
		id:              id,
		tags:            tags,
		sourceObjectKey: sourceObjectKey,
		objectKey:       objectKey,
		state:           state,
		createTime:      createTime,
		updateTime:      updateTime,
	}
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

func (a *Image) UpdateTime() time.Time {
	return a.updateTime
}

func (a *Image) IsPending() bool {
	return a.state.Equal(vo.ImageStatePending)
}

func (a *Image) IsConfirmed() bool {
	return a.state.Equal(vo.ImageStateConfirmed)
}

func (a *Image) IsProcessing() bool {
	return a.state.Equal(vo.ImageStateProcessing)
}

func (a *Image) IsActive() bool {
	return a.state.Equal(vo.ImageStateActive)
}

func (a *Image) IsDeleted() bool {
	return a.state.Equal(vo.ImageStateDeleted)
}

func (a *Image) CanConfirm() error {
	if a.IsDeleted() {
		return ErrImageDeleted
	}
	if a.IsProcessing() {
		return ErrImageProcessing
	}
	if a.IsActive() {
		return ErrImageActive
	}
	if a.IsConfirmed() {
		return ErrImageAlreadyConfirmed
	}
	if !a.IsPending() {
		return ErrImageNotPending
	}

	return nil
}

func (a *Image) CanStartProcessing() error {
	if a.IsDeleted() {
		return ErrImageDeleted
	}
	if a.IsActive() {
		return ErrImageActive
	}
	if a.IsProcessing() {
		return ErrImageAlreadyProcessing
	}
	if !a.IsConfirmed() {
		return ErrImageNotConfirmed
	}

	return nil
}

func (a *Image) CanCancelProcessing() error {
	if a.IsDeleted() {
		return ErrImageDeleted
	}
	if !a.IsProcessing() {
		return ErrImageNotProcessing
	}

	return nil
}

func (a *Image) CanCompleteProcessing() error {
	if a.IsDeleted() {
		return ErrImageDeleted
	}
	if a.IsActive() {
		return ErrImageAlreadyActive
	}
	if !a.IsProcessing() {
		return ErrImageNotProcessing
	}

	return nil
}

func (a *Image) CanDelete() error {
	if a.IsDeleted() {
		return ErrImageAlreadyDeleted
	}

	return nil
}

func (a *Image) Confirm(now time.Time) error {
	if err := a.CanConfirm(); err != nil {
		return err
	}

	newState := vo.ImageStateConfirmed

	a.state = newState
	a.updateTime = now

	a.recordEvent(event.NewImageConfirmed(a.id, a.tags, now))

	return nil
}

func (a *Image) StartProcessing(now time.Time) error {
	if err := a.CanStartProcessing(); err != nil {
		return err
	}

	newState := vo.ImageStateProcessing

	a.state = newState
	a.updateTime = now

	return nil
}

func (a *Image) CancelProcessing(now time.Time) error {
	if err := a.CanCancelProcessing(); err != nil {
		return err
	}

	newState := vo.ImageStateConfirmed

	a.state = newState
	a.updateTime = now

	return nil
}

func (a *Image) CompleteProcessing(now time.Time) error {
	if err := a.CanCompleteProcessing(); err != nil {
		return err
	}

	newState := vo.ImageStateActive

	a.state = newState
	a.updateTime = now

	a.recordEvent(event.NewImageActive(a.id, a.tags, now))

	return nil
}

func (a *Image) Delete(now time.Time) error {
	if err := a.CanDelete(); err != nil {
		return err
	}

	newState := vo.ImageStateDeleted

	a.state = newState
	a.updateTime = now

	a.recordEvent(event.NewImageDeleted(a.id, a.tags, now))

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
