package s3

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"log/slog"

	"github.com/velonyapp/asset/internal/application/port"
	"github.com/velonyapp/asset/internal/domain/vo"
)

var _ port.SourceObjectKeyGenerator = (*sourceObjectKeyGenerator)(nil)

type sourceObjectKeyGenerator struct {
	log *slog.Logger
}

func NewSourceObjectKeyGenerator(
	log *slog.Logger,
) port.SourceObjectKeyGenerator {
	return &sourceObjectKeyGenerator{
		log: log,
	}
}

func (g *sourceObjectKeyGenerator) Generate(ctx context.Context) (vo.ObjectKey, error) {
	var b [32]byte

	if _, err := rand.Read(b[:]); err != nil {
		g.log.ErrorContext(ctx,
			"failed to generate source object key bytes",
			slog.Any("error", err),
		)
		return vo.ObjectKey{}, err
	}

	objectKey, err := vo.NewObjectKey("tmp/sources/" + hex.EncodeToString(b[:]))
	if err != nil {
		g.log.ErrorContext(ctx,
			"failed to create source object key value object",
			slog.Any("error", err),
		)
		return vo.ObjectKey{}, err
	}

	return objectKey, nil
}
