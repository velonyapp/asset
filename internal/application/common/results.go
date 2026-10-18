package common

import "time"

type ImageResult struct {
	ID         string
	Tags       []string
	ObjectKey  string
	State      string
	Width      uint32
	Height     uint32
	Format     string
	CreateTime time.Time
}
