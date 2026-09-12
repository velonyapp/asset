package vo

import (
	"errors"
	"unicode/utf8"
)

var (
	ErrAssetKeyEmpty = errors.New(
		"asset key must not be empty",
	)
	ErrAssetKeyTooLong = errors.New(
		"asset key must not exceed 1024 bytes",
	)
	ErrAssetKeyInvalidUTF8 = errors.New(
		"asset key must contain valid UTF-8",
	)
)

type AssetKey struct {
	value string
}

func NewAssetKey(value string) (AssetKey, error) {
	if value == "" {
		return AssetKey{}, ErrAssetKeyEmpty
	}

	if len(value) > 1024 {
		return AssetKey{}, ErrAssetKeyTooLong
	}

	if !utf8.ValidString(value) {
		return AssetKey{}, ErrAssetKeyInvalidUTF8
	}

	return AssetKey{
		value: value,
	}, nil
}

func (r AssetKey) Value() string {
	return r.value
}

func (r AssetKey) String() string {
	return r.value
}
