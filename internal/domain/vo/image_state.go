package vo

import "errors"

var (
	ErrImageStateInvalid = errors.New("invalid image state")
)

type ImageState string

const (
	ImageStateCreated    ImageState = "CREATED"
	ImageStateUploading  ImageState = "UPLOADING"
	ImageStateUploaded   ImageState = "UPLOADED"
	ImageStateProcessing ImageState = "PROCESSING"
	ImageStateProcessed  ImageState = "PROCESSED"
	ImageStateDeleted    ImageState = "DELETED"
)

func NewImageState(value string) (ImageState, error) {
	state := ImageState(value)

	switch state {
	case ImageStateCreated,
		ImageStateUploading,
		ImageStateUploaded,
		ImageStateProcessing,
		ImageStateProcessed,
		ImageStateDeleted:
		return state, nil
	default:
		return "", ErrImageStateInvalid
	}
}

func (s ImageState) String() string {
	return string(s)
}

func (s ImageState) Equal(other ImageState) bool {
	return s == other
}
