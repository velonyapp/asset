package service

import (
	"context"
	"errors"

	"github.com/velonyapp/asset/internal/domain/repo"
	"github.com/velonyapp/asset/internal/domain/vo"
)

var (
	ErrObjectKeyAlreadyExists = errors.New("object key already exists")
)

type ObjectKeyPolicy struct {
	imageRepo repo.Image
}

func NewObjectKeyPolicy(
	imageRepo repo.Image,
) *ObjectKeyPolicy {
	return &ObjectKeyPolicy{
		imageRepo: imageRepo,
	}
}

func (p *ObjectKeyPolicy) CanUse(ctx context.Context, objectKey vo.ObjectKey) error {
	image, err := p.imageRepo.FindByAnyObjectKey(ctx, objectKey)
	if err != nil {
		return err
	}
	if image != nil {
		return ErrObjectKeyAlreadyExists
	}

	return nil
}
