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
	ErrImageAlreadyUploading  = errors.New("image is already uploading")
	ErrImageNotUploading      = errors.New("image is not uploading")
	ErrImageUploaded          = errors.New("image is already uploaded")
	ErrImageAlreadyUploaded   = errors.New("image is already uploaded")
	ErrImageNotUploaded       = errors.New("image is not uploaded")
	ErrImageProcessing        = errors.New("image is already processing")
	ErrImageAlreadyProcessing = errors.New("image is already processing")
	ErrImageNotProcessing     = errors.New("image is not processing")
	ErrImageProcessed         = errors.New("image is processed")
	ErrImageAlreadyProcessed  = errors.New("image is already processed")
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
	state := vo.ImageStateCreated

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

func (a *Image) IsCreated() bool {
	return a.state.Equal(vo.ImageStateCreated)
}

func (a *Image) IsUploading() bool {
	return a.state.Equal(vo.ImageStateUploading)
}

func (a *Image) IsUploaded() bool {
	return a.state.Equal(vo.ImageStateUploaded)
}

func (a *Image) IsProcessing() bool {
	return a.state.Equal(vo.ImageStateProcessing)
}

func (a *Image) IsProcessed() bool {
	return a.state.Equal(vo.ImageStateProcessed)
}

func (a *Image) IsDeleted() bool {
	return a.state.Equal(vo.ImageStateDeleted)
}

func (a *Image) CanStartUploading() error {
	if a.IsDeleted() {
		return ErrImageDeleted
	}
	if a.IsUploaded() {
		return ErrImageUploaded
	}
	if a.IsProcessing() {
		return ErrImageProcessing
	}
	if a.IsProcessed() {
		return ErrImageProcessed
	}
	if a.IsUploading() {
		return ErrImageAlreadyUploading
	}

	return nil
}

func (a *Image) CanConfirmUpload() error {
	if a.IsDeleted() {
		return ErrImageDeleted
	}
	if a.IsProcessing() {
		return ErrImageProcessing
	}
	if a.IsProcessed() {
		return ErrImageProcessed
	}
	if a.IsUploaded() {
		return ErrImageAlreadyUploaded
	}
	if !a.IsUploading() {
		return ErrImageNotUploading
	}

	return nil
}

func (a *Image) CanStartProcessing() error {
	if a.IsDeleted() {
		return ErrImageDeleted
	}
	if a.IsProcessed() {
		return ErrImageProcessed
	}
	if a.IsProcessing() {
		return ErrImageAlreadyProcessing
	}
	if !a.IsUploaded() {
		return ErrImageNotUploaded
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
	if a.IsProcessed() {
		return ErrImageAlreadyProcessed
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

func (a *Image) StartUploading(now time.Time) error {
	if err := a.CanStartUploading(); err != nil {
		return err
	}

	newState := vo.ImageStateUploading

	a.state = newState
	a.updateTime = now

	return nil
}

func (a *Image) ConfirmUpload(now time.Time) error {
	if err := a.CanConfirmUpload(); err != nil {
		return err
	}

	newState := vo.ImageStateUploaded

	a.state = newState
	a.updateTime = now

	a.recordEvent(event.NewImageUploaded(a.id, a.tags, now))

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

	newState := vo.ImageStateUploaded

	a.state = newState
	a.updateTime = now

	return nil
}

func (a *Image) CompleteProcessing(now time.Time) error {
	if err := a.CanCompleteProcessing(); err != nil {
		return err
	}

	newState := vo.ImageStateProcessed

	a.state = newState
	a.updateTime = now

	a.recordEvent(event.NewImageProcessed(a.id, a.tags, now))

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
