package query

import (
	"context"

	"github.com/velonyapp/asset/internal/application/common"
	"github.com/velonyapp/asset/internal/domain/vo"
)

type Image interface {
	Find(ctx context.Context, imageID vo.ImageID) (*common.ImageResult, error)
}
