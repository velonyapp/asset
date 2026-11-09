package common

import "time"

type ImageResult struct {
	ID         string
	Tags       []string
	Metadata   map[string]string
	ObjectKey  string
	State      string
	Width      int32
	Height     int32
	Format     string
	CreateTime time.Time
}
