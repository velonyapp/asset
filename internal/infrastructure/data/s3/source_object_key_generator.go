package s3

import (
	"crypto/rand"
	"encoding/hex"

	"github.com/velonyapp/asset/internal/application/port"
	"github.com/velonyapp/asset/internal/domain/vo"
)

var _ port.SourceObjectKeyGenerator = (*sourceObjectKeyGenerator)(nil)

type sourceObjectKeyGenerator struct{}

func NewSourceObjectKeyGenerator() port.SourceObjectKeyGenerator {
	return &sourceObjectKeyGenerator{}
}

func (g *sourceObjectKeyGenerator) Generate() (vo.ObjectKey, error) {
	var b [32]byte

	if _, err := rand.Read(b[:]); err != nil {
		return vo.ObjectKey{}, err
	}

	return vo.NewObjectKey("tmp/sources/" + hex.EncodeToString(b[:]))
}
