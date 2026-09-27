package mysql

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/velonyapp/asset/internal/application/domainevent"
	"github.com/velonyapp/asset/internal/domain/entity"
	"github.com/velonyapp/asset/internal/domain/repo"
	"github.com/velonyapp/asset/internal/domain/vo"
)

var _ repo.Image = (*imageRepo)(nil)

type imageRepo struct {
	db         *sql.DB
	dispatcher *domainevent.Dispatcher
}

func NewImageRepo(
	db *sql.DB,
	dispatcher *domainevent.Dispatcher,
) repo.Image {
	return &imageRepo{
		db:         db,
		dispatcher: dispatcher,
	}
}

type imageScanner interface {
	Scan(dest ...any) error
}

func (repo *imageRepo) FindByID(ctx context.Context, imageID vo.ImageID) (*entity.Image, error) {
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

func (repo *imageRepo) FindByStorageKey(ctx context.Context, storageKey vo.StorageKey) (*entity.Image, error) {
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

func (repo *imageRepo) Save(ctx context.Context, image *entity.Image) error {
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
			delete_time = ?
	`

	var deleteTime any
	if image.DeleteTime() != nil {
		deleteTime = image.DeleteTime()
	}

	if _, err := executor(ctx, repo.db).ExecContext(ctx, query,
		image.ID().Value(),
		image.StorageKey().Value(),
		image.IsReady(),
		image.CreateTime(),
		deleteTime,

		image.StorageKey().Value(),
		image.IsReady(),
		deleteTime,
	); err != nil {
		return err
	}

	for _, domainEvent := range image.PullEvents() {
		if err := repo.dispatcher.Dispatch(ctx, domainEvent); err != nil {
			return err
		}
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

	var deleteTimeVO *time.Time
	if deleteTime.Valid {
		value := deleteTime.Time
		deleteTimeVO = &value
	}

	return entity.ReconstituteImage(
		vo.NewImageID(id),
		storageKeyVO,
		ready,
		createTime,
		deleteTimeVO,
	), nil
}
