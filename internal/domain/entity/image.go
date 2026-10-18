package entity

import (
	"errors"
	"time"

	"github.com/velonyapp/asset/internal/domain/event"
	"github.com/velonyapp/asset/internal/domain/vo"
)

var (
	ErrImageDeleted           = errors.New("image is deleted")
	ErrImageNotAwaitingUpload = errors.New("image is not awaiting upload")
	ErrImageAlreadyDeleted    = errors.New("image is already deleted")
	ErrImageAlreadyPending    = errors.New("image is already pending")
	ErrImageNotPending        = errors.New("image is not pending")
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
	width           vo.ImageWidth
	height          vo.ImageHeight
	format          vo.ImageFormat
	createTime      time.Time
	updateTime      time.Time

	domainEvents []event.DomainEvent
}

func NewImage(
	tags vo.Tags,
	sourceObjectKey vo.ObjectKey,
	objectKey vo.ObjectKey,
	width vo.ImageWidth,
	height vo.ImageHeight,
	format vo.ImageFormat,
	now time.Time,
) (*Image, error) {
	if sourceObjectKey.Equal(objectKey) {
		return nil, ErrImageObjectKeysEqual
	}

	id := vo.GenerateImageID()
	state := vo.ImageStateAwaitingUpload

	a := &Image{
		id:              id,
		tags:            tags,
		sourceObjectKey: sourceObjectKey,
		objectKey:       objectKey,
		state:           state,
		width:           width,
		height:          height,
		format:          format,
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
	width vo.ImageWidth,
	height vo.ImageHeight,
	format vo.ImageFormat,
	createTime time.Time,
	updateTime time.Time,
) *Image {
	return &Image{
		id:              id,
		tags:            tags,
		sourceObjectKey: sourceObjectKey,
		objectKey:       objectKey,
		state:           state,
		width:           width,
		height:          height,
		format:          format,
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

func (a *Image) Width() vo.ImageWidth {
	return a.width
}

func (a *Image) Height() vo.ImageHeight {
	return a.height
}

func (a *Image) Format() vo.ImageFormat {
	return a.format
}

func (a *Image) CreateTime() time.Time {
	return a.createTime
}

func (a *Image) UpdateTime() time.Time {
	return a.updateTime
}

func (a *Image) IsAwaitingUpload() bool {
	return a.state.Equal(vo.ImageStateAwaitingUpload)
}

func (a *Image) IsPending() bool {
	return a.state.Equal(vo.ImageStatePending)
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

func (a *Image) CanConfirmUpload() error {
	if a.IsDeleted() {
		return ErrImageDeleted
	}
	if a.IsProcessing() {
		return ErrImageProcessing
	}
	if a.IsActive() {
		return ErrImageActive
	}
	if a.IsPending() {
		return ErrImageAlreadyPending
	}
	if !a.IsAwaitingUpload() {
		return ErrImageNotAwaitingUpload
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
	if !a.IsPending() {
		return ErrImageNotPending
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

func (a *Image) ConfirmUpload(now time.Time) error {
	if err := a.CanConfirmUpload(); err != nil {
		return err
	}

	newState := vo.ImageStatePending

	a.state = newState
	a.updateTime = now

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

	newState := vo.ImageStatePending

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
