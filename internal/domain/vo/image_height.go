package vo

import "errors"

var (
	ErrImageHeightNotPositive = errors.New("image height must be positive")
)

type ImageHeight struct {
	value int32
}

func NewImageHeight(value int32) (ImageHeight, error) {
	if value <= 0 {
		return ImageHeight{}, ErrImageHeightNotPositive
	}

	return ImageHeight{
		value: value,
	}, nil
}

func (h ImageHeight) Int() int {
	return int(h.value)
}

func (h ImageHeight) Int32() int32 {
	return h.value
}

func (h ImageHeight) Equal(other ImageHeight) bool {
	return h.value == other.value
}
