package vo

import "errors"

var ErrInvalidImageStatus = errors.New("invalid image status")

type ImageStatus string

const (
	ImageStatusPending    ImageStatus = "pending"
	ImageStatusProcessing ImageStatus = "processing"
	ImageStatusReady      ImageStatus = "ready"
	ImageStatusFailed     ImageStatus = "failed"
)

func NewImageStatus(value string) (ImageStatus, error) {
	status := ImageStatus(value)

	switch status {
	case ImageStatusPending,
		ImageStatusProcessing,
		ImageStatusReady,
		ImageStatusFailed:
		return status, nil
	default:
		return "", ErrInvalidImageStatus
	}
}

func (s ImageStatus) String() string {
	return string(s)
}
