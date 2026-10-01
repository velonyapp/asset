package vo

import (
	"errors"
	"strings"
	"unicode/utf8"
)

var (
	ErrTagEmpty       = errors.New("tag must not be empty")
	ErrTagTooLong     = errors.New("tag must not exceed 64 characters")
	ErrTagInvalidUTF8 = errors.New("tag must contain valid UTF-8")
)

type Tag struct {
	value string
}

func NewTag(value string) (Tag, error) {
	value = strings.TrimSpace(value)

	if value == "" {
		return Tag{}, ErrTagEmpty
	}

	if !utf8.ValidString(value) {
		return Tag{}, ErrTagInvalidUTF8
	}

	if utf8.RuneCountInString(value) > 64 {
		return Tag{}, ErrTagTooLong
	}

	return Tag{
		value: value,
	}, nil
}

func (t Tag) Value() string {
	return t.value
}

func (t Tag) String() string {
	return t.value
}
