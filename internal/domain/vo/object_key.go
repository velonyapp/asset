package vo

import (
	"errors"
	"unicode/utf8"
)

var (
	ErrObjectKeyEmpty       = errors.New("object key must not be empty")
	ErrObjectKeyTooLong     = errors.New("object key must not exceed 128 bytes")
	ErrObjectKeyInvalidUTF8 = errors.New("object key must contain valid UTF-8")
)

type ObjectKey struct {
	value string
}

func NewObjectKey(value string) (ObjectKey, error) {
	if value == "" {
		return ObjectKey{}, ErrObjectKeyEmpty
	}

	if len(value) > 128 {
		return ObjectKey{}, ErrObjectKeyTooLong
	}

	if !utf8.ValidString(value) {
		return ObjectKey{}, ErrObjectKeyInvalidUTF8
	}

	return ObjectKey{
		value: value,
	}, nil
}

func (k ObjectKey) Value() string {
	return k.value
}
func (k ObjectKey) String() string {
	return k.value
}
