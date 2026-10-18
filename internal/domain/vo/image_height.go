package vo

import "errors"

var (
	ErrImageHeightNotPositive = errors.New("image height must be positive")
)

type ImageHeight struct {
	value int
}

func NewImageHeight(value int) (ImageHeight, error) {
	if value <= 0 {
		return ImageHeight{}, ErrImageHeightNotPositive
	}

	return ImageHeight{
		value: value,
	}, nil
}

func (h ImageHeight) Int() int {
	return h.value
}

func (h ImageHeight) UInt32() uint32 {
	return uint32(h.value)
}

func (h ImageHeight) Equal(other ImageHeight) bool {
	return h.value == other.value
}
