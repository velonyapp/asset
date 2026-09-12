package common

import (
	"github.com/velonyapp/asset/internal/domain/entity"
)

type ImageResult struct {
	ID  string
	Key string
}

func NewImageResult(image *entity.Image) *ImageResult {
	result := &ImageResult{
		ID:  image.ID.String(),
		Key: image.ID.String(),
	}

	return result
}
