package common

import (
	"time"

	"github.com/velonyapp/asset/internal/domain/entity"
)

type ImageResult struct {
	ID         string
	StorageKey string
	Ready      bool
	CreateTime time.Time
}

func NewImageResult(image *entity.Image) *ImageResult {
	if image == nil {
		return nil
	}

	result := &ImageResult{
		ID:         image.ID().Value(),
		StorageKey: image.StorageKey().Value(),
		Ready:      image.IsReady(),
		CreateTime: image.CreateTime(),
	}

	return result
}
