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

	"github.com/go-sql-driver/mysql"
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

func (r *imageRepo) FindByID(ctx context.Context, imageID vo.ImageID) (*entity.Image, error) {
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

	row := executor(ctx, r.db).QueryRowContext(ctx, query, imageID.String())

	image, err := scanImage(row)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, repo.ErrImageNotFound
		}

		return nil, err
	}

	return image, nil
}

func (r *imageRepo) Save(ctx context.Context, image *entity.Image) error {
	if image.IsDeleted() {
		const query = `
			DELETE FROM images
			WHERE id = ?
		`

		if _, err := executor(ctx, r.db).ExecContext(ctx, query,
			image.ID().String(),
		); err != nil {
			return err
		}
	} else {
		const query = `
			UPDATE images
			SET object_exists = ?
			WHERE id = ?
		`

		result, err := executor(ctx, r.db).ExecContext(ctx, query,
			image.ObjectExists(),

			image.ID().String(),
		)
		if err != nil {
			return err
		}

		affected, err := result.RowsAffected()
		if err != nil {
			return err
		}

		if affected == 0 {
			const query = `
				INSERT INTO images (
					id,
					tags,
					object_key,
					object_exists,
					create_time
				)
				VALUES (?, ?, ?, ?, ?)
			`

			if _, err := executor(ctx, r.db).ExecContext(ctx, query,
				image.ID().String(),
				strings.Join(image.Tags().Strings(), ";"),
				image.ObjectKey().String(),
				image.ObjectExists(),
				image.CreateTime(),
			); err != nil {
				var mysqlErr *mysql.MySQLError
				if errors.As(err, &mysqlErr) &&
					mysqlErr.Number == 1062 &&
					strings.Contains(mysqlErr.Message, "uq_images_object_key") {
					return repo.ErrObjectKeyConflict
				}

				return err
			}
		}
	}

	for _, domainEvent := range image.PullEvents() {
		if err := r.dispatcher.Dispatch(ctx, domainEvent); err != nil {
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

	id, _ := vo.NewImageID(idRaw)
	tags, _ := vo.NewTags(nil)
	if tagsRaw != "" {
		tags, _ = vo.NewTags(strings.Split(tagsRaw, ";"))
	}
	objectKey, _ := vo.NewObjectKey(objectKeyRaw)
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
