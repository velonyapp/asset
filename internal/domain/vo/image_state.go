package vo

import "errors"

var (
	ErrImageStateInvalid = errors.New("invalid image state")
)

type ImageState string

const (
	ImageStatePending   ImageState = "PENDING"
	ImageStateUploaded  ImageState = "UPLOADED"
	ImageStateProcessed ImageState = "PROCESSED"
)

func NewImageState(value string) (ImageState, error) {
	state := ImageState(value)

	switch state {
	case ImageStatePending,
		ImageStateUploaded,
		ImageStateProcessed:
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
