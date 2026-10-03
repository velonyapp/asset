package vo

import (
	"errors"
)

var (
	ErrTagEmpty            = errors.New("tag must not be empty")
	ErrTagTooLong          = errors.New("tag must not exceed 64 characters")
	ErrTagInvalidCharacter = errors.New("tag contains an invalid character")
)

type Tag struct {
	value string
}

func NewTag(value string) (Tag, error) {
	if value == "" {
		return Tag{}, ErrTagEmpty
	}

	for i := 0; i < len(value); i++ {
		c := value[i]

		if !((c >= 'a' && c <= 'z') ||
			(c >= '0' && c <= '9') ||
			c == '-' ||
			c == '_' ||
			c == '.') {
			return Tag{}, ErrTagInvalidCharacter
		}
	}

	if len(value) > 64 {
		return Tag{}, ErrTagTooLong
	}

	return Tag{
		value: value,
	}, nil
}

func (t Tag) String() string {
	return t.value
}

func (t Tag) Equal(other Tag) bool {
	return t.value == other.value
}