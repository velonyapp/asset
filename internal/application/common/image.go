package common

import "time"

type ImageResult struct {
	ID           string
	Tags         []string
	ObjectKey    string
	ObjectExists bool
	CreateTime   time.Time
}
