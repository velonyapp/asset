package vo

import (
	"errors"

	"github.com/google/uuid"
)

type ImageID struct {
	value string
}

var (
	ErrImageIDInvalid = errors.New("invalid image ID")
)

func NewImageID(value string) (ImageID, error) {
	id, err := uuid.Parse(value)
	if err != nil {
		return ImageID{}, ErrImageIDInvalid
	}

	return ImageID{value: id.String()}, nil
}

func GenerateImageID() ImageID {
	return ImageID{
		value: uuid.Must(uuid.NewV7()).String(),
	}
}

func (id ImageID) String() string {
	return id.value
}

func (id ImageID) Equal(other ImageID) bool {
	return id.value == other.value
}
