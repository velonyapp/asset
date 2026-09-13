package repo

import (
	"context"

	"github.com/velonyapp/asset/internal/domain/entity"
	"github.com/velonyapp/asset/internal/domain/vo"
)

type Image interface {
	FindByID(ctx context.Context, id vo.ImageID) (*entity.Image, error)
	FindByStorageKey(ctx context.Context, storageKey vo.StorageKey) (*entity.Image, error)

	Save(ctx context.Context, image *entity.Image) error
}
