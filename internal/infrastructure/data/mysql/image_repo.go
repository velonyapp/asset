package mysql

import (
	"context"
	"database/sql"
	"encoding/json"
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
			tags,
			object_key,
			object_exists,
			create_time,
			delete_time
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
	const query = `
		INSERT INTO images (
			id,
			tags,
			object_key,
			object_exists,
			create_time,
			delete_time
		)
		VALUES (?, ?, ?, ?, ?, ?)
		ON DUPLICATE KEY UPDATE
			tags = ?,
			object_key = ?,
			object_exists = ?,
			delete_time = ?
	`

	imageTags := image.Tags()

	tags := make([]string, len(imageTags))
	for i, tag := range imageTags {
		tags[i] = tag.Value()
	}

	tagsJSON, err := json.Marshal(tags)
	if err != nil {
		return err
	}

	var deleteTime any
	if image.DeleteTime() != nil {
		deleteTime = image.DeleteTime()
	}

	if _, err := executor(ctx, repo.db).ExecContext(ctx, query,
		image.ID().String(),
		tagsJSON,
		image.ObjectKey().String(),
		image.ObjectExists(),
		image.CreateTime(),
		deleteTime,

		tagsJSON,
		image.ObjectKey().String(),
		image.ObjectExists(),
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
		id           string
		tagsJSON     []byte
		objectKey    string
		objectExists bool
		createTime   time.Time
		deleteTime   sql.NullTime
	)

	if err := scanner.Scan(
		&id,
		&tagsJSON,
		&objectKey,
		&objectExists,
		&createTime,
		&deleteTime,
	); err != nil {
		return nil, err
	}

	imageID, err := vo.NewImageID(id)
	if err != nil {
		return nil, err
	}

	var tagValues []string
	if err := json.Unmarshal(tagsJSON, &tagValues); err != nil {
		return nil, err
	}

	tags := make([]vo.Tag, 0, len(tagValues))
	for _, value := range tagValues {
		tag, err := vo.NewTag(value)
		if err != nil {
			return nil, err
		}

		tags = append(tags, tag)
	}

	objectKeyVO, err := vo.NewObjectKey(objectKey)
	if err != nil {
		return nil, err
	}

	var deleteTimeVO *time.Time
	if deleteTime.Valid {
		value := deleteTime.Time
		deleteTimeVO = &value
	}

	return entity.ReconstituteImage(
		imageID,
		tags,
		objectKeyVO,
		objectExists,
		createTime,
		deleteTimeVO,
	), nil
}
