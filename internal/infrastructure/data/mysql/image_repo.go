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

func (r *imageRepo) FindByID(ctx context.Context, imageID vo.ImageID) (*entity.Image, error) {
	const query = `
		SELECT
			id,
			tags,
			source_object_key,
			object_key,
			state,
			width,
			height,
			format,
			create_time,
			update_time
		FROM images
		WHERE id = ?
			AND state <> 'DELETED'
		LIMIT 1
	`

	row := executor(ctx, r.db).QueryRowContext(ctx, query,
		imageID.String(),
	)

	image, err := scanImage(row)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}

		return nil, err
	}

	return image, nil
}

func (r *imageRepo) FindByIDForUpdate(ctx context.Context, imageID vo.ImageID) (*entity.Image, error) {
	const query = `
		SELECT
			id,
			tags,
			source_object_key,
			object_key,
			state,
			width,
			height,
			format,
			create_time,
			update_time
		FROM images
		WHERE id = ?
			AND state <> 'DELETED'
		LIMIT 1
		FOR UPDATE
	`

	row := executor(ctx, r.db).QueryRowContext(ctx, query,
		imageID.String(),
	)

	image, err := scanImage(row)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}

		return nil, err
	}

	return image, nil
}

func (r *imageRepo) FindByAnyObjectKeyForUpdate(ctx context.Context, objectKey vo.ObjectKey) (*entity.Image, error) {
	const query = `
		SELECT
			id,
			tags,
			source_object_key,
			object_key,
			state,
			width,
			height,
			format,
			create_time,
			update_time
		FROM images
		WHERE (source_object_key = ? OR object_key = ?)
  			AND state <> 'DELETED'
		LIMIT 1
		FOR UPDATE
	`

	row := executor(ctx, r.db).QueryRowContext(ctx, query,
		objectKey.String(),
		objectKey.String(),
	)

	image, err := scanImage(row)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}

		return nil, err
	}

	return image, nil
}

func (r *imageRepo) Save(ctx context.Context, image *entity.Image) error {
	const query = `
	    INSERT INTO images (
	        id,
	        tags,
	        source_object_key,
	        object_key,
	        state,
			width,
			height,
			format,
	        create_time,
	        update_time
	    )
	    VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?) AS new
	    ON DUPLICATE KEY UPDATE
	        state = new.state,
	        update_time = new.update_time
	`

	if _, err := executor(ctx, r.db).ExecContext(ctx, query,
		image.ID().String(),
		strings.Join(image.Tags().Strings(), ";"),
		image.SourceObjectKey().String(),
		image.ObjectKey().String(),
		image.State().String(),
		image.Width().Int(),
		image.Height().Int(),
		image.Format().String(),
		image.CreateTime(),
		image.UpdateTime(),
	); err != nil {
		return err
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
		idRaw              string
		tagsRaw            string
		sourceObjectKeyRaw string
		objectKeyRaw       string
		stateRaw           string
		widthRaw           int
		heightRaw          int
		formatRaw          string
		createTimeRaw      time.Time
		updateTimeRaw      time.Time
	)

	if err := scanner.Scan(
		&idRaw,
		&tagsRaw,
		&sourceObjectKeyRaw,
		&objectKeyRaw,
		&stateRaw,
		&widthRaw,
		&heightRaw,
		&formatRaw,
		&createTimeRaw,
		&updateTimeRaw,
	); err != nil {
		return nil, err
	}

	id, _ := vo.NewImageID(idRaw)
	tags, _ := vo.NewTags(nil)
	if tagsRaw != "" {
		tags, _ = vo.NewTags(strings.Split(tagsRaw, ";"))
	}
	sourceObjectKey, _ := vo.NewObjectKey(sourceObjectKeyRaw)
	objectKey, _ := vo.NewObjectKey(objectKeyRaw)
	state, _ := vo.NewImageState(stateRaw)
	width, _ := vo.NewImageWidth(widthRaw)
	height, _ := vo.NewImageHeight(heightRaw)
	format, _ := vo.NewImageFormat(formatRaw)
	createTime := createTimeRaw
	updateTime := updateTimeRaw

	return entity.ReconstituteImage(
		id,
		tags,
		sourceObjectKey,
		objectKey,
		state,
		width,
		height,
		format,
		createTime,
		updateTime,
	), nil
}
