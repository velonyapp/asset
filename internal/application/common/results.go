package common

import "time"

type ImageResult struct {
	ID         string
	Tags       []string
	ObjectKey  string
	State      string
	CreateTime time.Time
	UpdateTime time.Time
}
