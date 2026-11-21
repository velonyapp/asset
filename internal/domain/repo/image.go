package repo

import (
	"context"
	"errors"

	"github.com/velonyapp/asset/internal/domain/entity"
	"github.com/velonyapp/asset/internal/domain/vo"
)

var (
	ErrObjectKeyAlreadyExists = errors.New("object key already exists")
	ErrImageVersionConflict   = errors.New("image version conflict")
)

type Image interface {
	FindByID(ctx context.Context, imageID vo.ImageID) (*entity.Image, error)

	Save(ctx context.Context, image *entity.Image) error
}
