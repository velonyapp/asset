package vo

import "errors"

var (
	ErrImageWidthNotPositive = errors.New("image width must be positive")
)

type ImageWidth struct {
	value int32
}

func NewImageWidth(value int32) (ImageWidth, error) {
	if value <= 0 {
		return ImageWidth{}, ErrImageWidthNotPositive
	}

	return ImageWidth{
		value: value,
	}, nil
}

func (w ImageWidth) Int() int {
	return int(w.value)
}

func (w ImageWidth) Int32() int32 {
	return w.value
}

func (w ImageWidth) Equal(other ImageWidth) bool {
	return w.value == other.value
}
