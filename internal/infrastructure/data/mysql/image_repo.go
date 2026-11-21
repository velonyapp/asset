package mysql

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"log/slog"
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
	log        *slog.Logger
}

func NewImageRepo(
	db *sql.DB,
	dispatcher *domainevent.Dispatcher,
	log *slog.Logger,
) repo.Image {
	return &imageRepo{
		db:         db,
		dispatcher: dispatcher,
		log:        log,
	}
}

type repoImageScanner interface {
	Scan(dest ...any) error
}

func (r *imageRepo) FindByID(ctx context.Context, imageID vo.ImageID) (*entity.Image, error) {
	const query = `
		SELECT
			id,
			tags,
			metadata,
			source_object_key,
			object_key,
			state,
			width,
			height,
			format,
			create_time,
			update_time,
			version
		FROM images
		WHERE id = ?
		LIMIT 1
	`

	row := executor(ctx, r.db).QueryRowContext(ctx, query,
		imageID.String(),
	)

	image, err := scanRepoImage(row)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}

		r.log.ErrorContext(ctx,
			"failed to find image",
			slog.String("image_id", imageID.String()),
			slog.Any("error", err),
		)
		return nil, err
	}

	return image, nil
}

func (r *imageRepo) FindByObjectKey(ctx context.Context, objectKey vo.ObjectKey) (*entity.Image, error) {
	const query = `
		SELECT
			id,
			tags,
			metadata,
			source_object_key,
			object_key,
			state,
			width,
			height,
			format,
			create_time,
			update_time,
			version
		FROM images
		WHERE object_key = CAST(? AS BINARY)
		LIMIT 1
	`

	row := executor(ctx, r.db).QueryRowContext(ctx, query,
		objectKey.String(),
	)

	image, err := scanRepoImage(row)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}

		r.log.ErrorContext(ctx,
			"failed to find image by object key",
			slog.String("object_key", objectKey.String()),
			slog.Any("error", err),
		)
		return nil, err
	}

	return image, nil
}

func (r *imageRepo) Save(ctx context.Context, image *entity.Image) error {
	if !image.IsDirty() {
		return nil
	}

	if image.IsDeleted() {
		if !image.IsNew() {
			const query = `
				DELETE FROM images
				WHERE id = ?
			`

			if _, err := executor(ctx, r.db).ExecContext(ctx, query,
				image.ID().String(),
			); err != nil {
				r.log.ErrorContext(ctx,
					"failed to delete image",
					slog.String("image_id", image.ID().String()),
					slog.Any("error", err),
				)
				return err
			}
		}
	} else {
		if image.IsNew() {
			const query = `
				INSERT INTO images (
					id,
					tags,
					metadata,
					source_object_key,
					object_key,
					state,
					width,
					height,
					format,
					create_time,
					update_time,
					version
				)
				VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
			`

			metadata, err := json.Marshal(image.Metadata().Map())
			if err != nil {
				r.log.ErrorContext(ctx,
					"failed to marshal image metadata",
					slog.String("image_id", image.ID().String()),
					slog.Any("error", err),
				)
				return err
			}

			if _, err := executor(ctx, r.db).ExecContext(ctx, query,
				image.ID().String(),
				strings.Join(image.Tags().Strings(), ";"),
				string(metadata),
				image.SourceObjectKey().String(),
				image.ObjectKey().String(),
				image.State().String(),
				image.Width().Int32(),
				image.Height().Int32(),
				image.Format().String(),
				image.CreateTime(),
				image.UpdateTime(),
				image.Version(),
			); err != nil {
				r.log.ErrorContext(ctx,
					"failed to insert image",
					slog.String("image_id", image.ID().String()),
					slog.Any("error", err),
				)

				var mysqlErr *mysql.MySQLError
				if errors.As(err, &mysqlErr) && mysqlErr.Number == 1062 {
					switch {
					case strings.Contains(mysqlErr.Message, "uq_images_object_keys"):
						return repo.ErrObjectKeyAlreadyExists
					}
				}

				return err
			}
		} else {
			const query = `
			    UPDATE images
			    SET
			        state = ?,
			        update_time = ?,
			        version = ?
			    WHERE id = ?
			        AND version = ?
			`

			result, err := executor(ctx, r.db).ExecContext(ctx, query,
				image.State().String(),
				image.UpdateTime(),
				image.Version(),
				image.ID().String(),
				image.Version()-1,
			)
			if err != nil {
				r.log.ErrorContext(ctx,
					"failed to update image",
					slog.String("image_id", image.ID().String()),
					slog.Any("error", err),
				)
				return err
			}
			rowsAffected, err := result.RowsAffected()
			if err != nil {
				r.log.ErrorContext(ctx,
					"failed to get image rows affected",
					slog.String("image_id", image.ID().String()),
					slog.Any("error", err),
				)
				return err
			}
			if rowsAffected == 0 {
				return repo.ErrImageVersionConflict
			}
		}
	}

	for _, domainEvent := range image.PullEvents() {
		if err := r.dispatcher.Dispatch(ctx, domainEvent); err != nil {
			r.log.ErrorContext(ctx,
				"failed to dispatch image domain event",
				slog.String("image_id", image.ID().String()),
				slog.Any("error", err),
			)
			return err
		}
	}

	image.MarkPersisted()

	return nil
}

func scanRepoImage(scanner repoImageScanner) (*entity.Image, error) {
	var (
		idRaw              string
		tagsRaw            string
		metadataRaw        []byte
		sourceObjectKeyRaw string
		objectKeyRaw       string
		stateRaw           string
		widthRaw           int32
		heightRaw          int32
		formatRaw          string
		createTimeRaw      time.Time
		updateTimeRaw      time.Time
		versionRaw         uint32
	)

	if err := scanner.Scan(
		&idRaw,
		&tagsRaw,
		&metadataRaw,
		&sourceObjectKeyRaw,
		&objectKeyRaw,
		&stateRaw,
		&widthRaw,
		&heightRaw,
		&formatRaw,
		&createTimeRaw,
		&updateTimeRaw,
		&versionRaw,
	); err != nil {
		return nil, err
	}

	id, _ := vo.NewImageID(idRaw)
	tags, _ := vo.NewTags(nil)
	if tagsRaw != "" {
		tags, _ = vo.NewTags(strings.Split(tagsRaw, ";"))
	}
	var metadataValues map[string]string
	if err := json.Unmarshal(metadataRaw, &metadataValues); err != nil {
		return nil, err
	}
	metadata, err := vo.NewMetadata(metadataValues)
	if err != nil {
		return nil, err
	}
	sourceObjectKey, _ := vo.NewObjectKey(sourceObjectKeyRaw)
	objectKey, _ := vo.NewObjectKey(objectKeyRaw)
	state, _ := vo.NewImageState(stateRaw)
	width, _ := vo.NewImageWidth(widthRaw)
	height, _ := vo.NewImageHeight(heightRaw)
	format, _ := vo.NewImageFormat(formatRaw)
	createTime := createTimeRaw
	updateTime := updateTimeRaw
	version := versionRaw

	return entity.ReconstituteImage(
		id,
		tags,
		metadata,
		sourceObjectKey,
		objectKey,
		state,
		width,
		height,
		format,
		createTime,
		updateTime,
		version,
	), nil
}
