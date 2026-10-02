package common

import (
	"errors"
	"time"
)

var ErrImageNotFound = errors.New("image not found")

type ImageResult struct {
	ID           string
	Tags         []string
	ObjectKey    string
	ObjectExists bool
	CreateTime   time.Time
}
