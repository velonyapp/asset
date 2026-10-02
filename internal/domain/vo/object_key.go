package vo

import (
	"errors"
	"strings"
)

var (
	ErrObjectKeyEmpty            = errors.New("object key must not be empty")
	ErrObjectKeyElementEmpty     = errors.New("object key element must not be empty")
	ErrObjectKeyTooLong          = errors.New("object key must not exceed 128 characters")
	ErrObjectKeyInvalidCharacter = errors.New("object key contains an invalid character")
)

type ObjectKey struct {
	value string
}

func NewObjectKey(value string) (ObjectKey, error) {
	if value == "" {
		return ObjectKey{}, ErrObjectKeyEmpty
	}

	elements := strings.Split(value, "/")

	for _, element := range elements {
		if element == "" {
			return ObjectKey{}, ErrObjectKeyElementEmpty
		}

		for i := 0; i < len(element); i++ {
			c := element[i]

			if !((c >= 'a' && c <= 'z') ||
				(c >= 'A' && c <= 'Z') ||
				(c >= '0' && c <= '9') ||
				c == '-' ||
				c == '_' ||
				c == '.') {
				return ObjectKey{}, ErrObjectKeyInvalidCharacter
			}
		}
	}

	if len(value) > 128 {
		return ObjectKey{}, ErrObjectKeyTooLong
	}

	return ObjectKey{
		value: value,
	}, nil
}

func (k ObjectKey) String() string {
	return k.value
}

func (k ObjectKey) Equal(other ObjectKey) bool {
	return k.value == other.value
}
