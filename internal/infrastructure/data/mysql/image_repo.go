package mysql

import (
	"context"
	"database/sql"
	"errors"
	"strings"
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
			tags,
			object_key,
			object_exists,
			create_time
		FROM images
		WHERE id = ?
		LIMIT 1
	`

	row := executor(ctx, repo.db).QueryRowContext(ctx, query, imageID.String())

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
	if image.IsDeleted() {
		const query = `
			DELETE FROM images
			WHERE id = ?
		`

		_, err := executor(ctx, repo.db).ExecContext(ctx, query, image.ID().String())
		return err
	}

	const query = `
		INSERT INTO images (
			id,
			tags,
			object_key,
			object_exists,
			create_time
		)
		VALUES (?, ?, ?, ?, ?)
		ON DUPLICATE KEY UPDATE
			object_exists = ?
	`

	if _, err := executor(ctx, repo.db).ExecContext(ctx, query,
		image.ID().String(),
		strings.Join(image.Tags().Strings(), ";"),
		image.ObjectKey().String(),
		image.ObjectExists(),
		image.CreateTime(),

		image.ObjectExists(),
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
		idRaw           string
		tagsRaw         string
		objectKeyRaw    string
		objectExistsRaw bool
		createTimeRaw   time.Time
	)

	if err := scanner.Scan(
		&idRaw,
		&tagsRaw,
		&objectKeyRaw,
		&objectExistsRaw,
		&createTimeRaw,
	); err != nil {
		return nil, err
	}

	id, err := vo.NewImageID(idRaw)
	if err != nil {
		return nil, err
	}
	tags, err := vo.NewTags(strings.Split(tagsRaw, ";"))
	if err != nil {
		return nil, err
	}
	objectKey, err := vo.NewObjectKey(objectKeyRaw)
	if err != nil {
		return nil, err
	}
	objectExists := objectExistsRaw
	createTime := createTimeRaw

	return entity.ReconstituteImage(
		id,
		tags,
		objectKey,
		objectExists,
		createTime,
		nil,
	), nil
}
