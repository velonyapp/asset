package event

import (
	"time"

	"github.com/velonyapp/asset/internal/domain/vo"
)

var _ DomainEvent = (*ImageUploaded)(nil)

type ImageUploaded struct {
	BaseDomainEvent

	tags       vo.Tags
	uploadTime time.Time
}

func NewImageUploaded(
	imageID vo.ImageID,
	tags vo.Tags,
	uploadTime time.Time,
) ImageUploaded {
	return ImageUploaded{
		BaseDomainEvent: NewBaseDomainEvent(imageID.String()),

		tags:       tags,
		uploadTime: uploadTime,
	}
}

func (e ImageUploaded) Tags() vo.Tags {
	return e.tags
}

func (e ImageUploaded) UploadTime() time.Time {
	return e.uploadTime
}
