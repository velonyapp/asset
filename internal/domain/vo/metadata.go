package vo

import (
	"errors"
	"maps"
)

var (
	ErrMetadataKeyEmpty = errors.New("metadata key must not be empty")
)

type Metadata struct {
	value map[string]string
}

func NewMetadata(value map[string]string) (Metadata, error) {
	for key := range value {
		if key == "" {
			return Metadata{}, ErrMetadataKeyEmpty
		}
	}

	return Metadata{
		value: maps.Clone(value),
	}, nil
}

func (m Metadata) Map() map[string]string {
	return maps.Clone(m.value)
}

func (m Metadata) Equal(other Metadata) bool {
	return maps.Equal(m.value, other.value)
}
