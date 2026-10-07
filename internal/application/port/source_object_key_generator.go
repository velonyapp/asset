package port

import "github.com/velonyapp/asset/internal/domain/vo"

type SourceObjectKeyGenerator interface {
	Generate() (vo.ObjectKey, error)
}
