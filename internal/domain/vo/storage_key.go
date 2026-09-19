package vo

import (
	"errors"
	"unicode/utf8"
)

var (
	ErrStorageKeyEmpty = errors.New(
		"storage key must not be empty",
	)
	ErrStorageKeyTooLong = errors.New(
		"storage key must not exceed 128 bytes",
	)
	ErrStorageKeyInvalidUTF8 = errors.New(
		"storage key must contain valid UTF-8",
	)
)

type StorageKey struct {
	value string
}

func NewStorageKey(value string) (StorageKey, error) {
	if value == "" {
		return StorageKey{}, ErrStorageKeyEmpty
	}

	if len(value) > 128 {
		return StorageKey{}, ErrStorageKeyTooLong
	}

	if !utf8.ValidString(value) {
		return StorageKey{}, ErrStorageKeyInvalidUTF8
	}

	return StorageKey{
		value: value,
	}, nil
}

func (k StorageKey) Value() string {
	return k.value
}
func (k StorageKey) String() string {
	return k.value
}
