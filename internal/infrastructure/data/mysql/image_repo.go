package mysql

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/velonyapp/asset/internal/domain/entity"
	"github.com/velonyapp/asset/internal/domain/repo"
	"github.com/velonyapp/asset/internal/domain/vo"
)

type ImageRepo struct {
	db *sql.DB
}

func NewImageRepo(
	db *sql.DB,
) repo.Image {
	return &ImageRepo{
		db: db,
	}
}

type imageScanner interface {
	Scan(dest ...any) error
}

func (repo *ImageRepo) FindByID(ctx context.Context, imageID vo.ImageID) (*entity.Image, error) {
	const query = `
		SELECT
			id,
			storage_key,
			ready,
			create_time,
			delete_time
		FROM images
		WHERE id = ?
		LIMIT 1
	`

	row := executor(ctx, repo.db).QueryRowContext(ctx, query, imageID.Value())

	image, err := scanImage(row)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}

		return nil, err
	}

	return image, nil
}

func (repo *ImageRepo) FindByStorageKey(ctx context.Context, storageKey vo.StorageKey) (*entity.Image, error) {
	const query = `
		SELECT
			id,
			storage_key,
			ready,
			create_time,
			delete_time
		FROM images
		WHERE storage_key = ?
		LIMIT 1
	`

	row := executor(ctx, repo.db).QueryRowContext(ctx, query, storageKey.Value())

	image, err := scanImage(row)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}

		return nil, err
	}

	return image, nil
}

func (repo *ImageRepo) Save(ctx context.Context, image *entity.Image) error {
	const query = `
		INSERT INTO images (
			id,
			storage_key,
			ready,
			create_time,
			delete_time
		)
		VALUES (?, ?, ?, ?, ?)
		ON DUPLICATE KEY UPDATE
			storage_key = ?,
			ready = ?,
			create_time = ?,
			delete_time = ?
	`

	var deleteTime any
	if image.DeleteTime != nil {
		deleteTime = image.DeleteTime.Value()
	}

	if _, err := executor(ctx, repo.db).ExecContext(ctx, query,
		image.ID.Value(),
		image.StorageKey.Value(),
		image.Ready,
		image.CreateTime.Value(),
		deleteTime,

		image.StorageKey.Value(),
		image.Ready,
		image.CreateTime.Value(),
		deleteTime,
	); err != nil {
		return err
	}

	return nil
}

func scanImage(scanner imageScanner) (*entity.Image, error) {
	var (
		id         string
		storageKey string
		ready      bool
		createTime time.Time
		deleteTime sql.NullTime
	)

	if err := scanner.Scan(
		&id,
		&storageKey,
		&ready,
		&createTime,
		&deleteTime,
	); err != nil {
		return nil, err
	}

	storageKeyVO, err := vo.NewStorageKey(storageKey)
	if err != nil {
		return nil, err
	}

	var deleteTimeVO *vo.Time
	if deleteTime.Valid {
		value := vo.NewTime(deleteTime.Time)
		deleteTimeVO = &value
	}

	return &entity.Image{
		ID:         vo.NewImageID(id),
		StorageKey: storageKeyVO,
		Ready:      ready,
		CreateTime: vo.NewTime(createTime),
		DeleteTime: deleteTimeVO,
	}, nil
}
