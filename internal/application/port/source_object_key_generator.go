package port

import (
	"context"

	"github.com/velonyapp/asset/internal/domain/vo"
)

type SourceObjectKeyGenerator interface {
	Generate(ctx context.Context) (vo.ObjectKey, error)
}
