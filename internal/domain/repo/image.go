package repo

import (
	"context"

	"github.com/velonyapp/asset/internal/domain/entity"
	"github.com/velonyapp/asset/internal/domain/vo"
)

type Image interface {
	FindByID(ctx context.Context, imageID vo.ImageID) (*entity.Image, error)
	FindByIDForUpdate(ctx context.Context, imageID vo.ImageID) (*entity.Image, error)
	FindByAnyObjectKeyForUpdate(ctx context.Context, objectKey vo.ObjectKey) (*entity.Image, error)

	Save(ctx context.Context, image *entity.Image) error
}
