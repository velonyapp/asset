package vo

import "github.com/google/uuid"

type ImageID struct {
	value string
}

func NewImageID(value string) ImageID {
	return ImageID{value: value}
}

func NewImageIDRandom() ImageID {
	return ImageID{value: uuid.Must(uuid.NewV7()).String()}
}

func (id ImageID) Value() string {
	return id.value
}

func (id ImageID) String() string {
	return id.value
}
