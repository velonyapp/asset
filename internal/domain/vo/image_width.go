package vo

import "errors"

var (
	ErrImageWidthNotPositive = errors.New("image width must be positive")
)

type ImageWidth struct {
	value int
}

func NewImageWidth(value int) (ImageWidth, error) {
	if value <= 0 {
		return ImageWidth{}, ErrImageWidthNotPositive
	}

	return ImageWidth{
		value: value,
	}, nil
}

func (w ImageWidth) Int() int {
	return w.value
}

func (w ImageWidth) UInt32() uint32 {
	return uint32(w.value)
}

func (w ImageWidth) Equal(other ImageWidth) bool {
	return w.value == other.value
}
