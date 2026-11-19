package mysql

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"strings"
	"time"

	"github.com/velonyapp/asset/internal/application/common"
	"github.com/velonyapp/asset/internal/application/query"
	"github.com/velonyapp/asset/internal/domain/vo"
)

var _ query.Image = (*imageQuery)(nil)

type imageQuery struct {
	db *sql.DB
}

func NewImageQuery(
	db *sql.DB,
) query.Image {
	return &imageQuery{
		db: db,
	}
}

type queryImageScanner interface {
	Scan(dest ...any) error
}

func (q *imageQuery) Find(ctx context.Context, imageID vo.ImageID) (*common.ImageResult, error) {
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
			create_time
		FROM images
		WHERE id = ?
			AND state <> 'DELETED'
		LIMIT 1
	`

	row := executor(ctx, q.db).QueryRowContext(ctx, query,
		imageID.String(),
	)

	image, err := scanQueryImage(row)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}

		return nil, err
	}

	return &image, nil
}

func scanQueryImage(scanner queryImageScanner) (common.ImageResult, error) {
	var (
		id              string
		tagsRaw         string
		metadataRaw     []byte
		sourceObjectKey string
		objectKey       string
		state           string
		width           int32
		height          int32
		format          string
		createTime      time.Time
	)

	if err := scanner.Scan(
		&id,
		&tagsRaw,
		&metadataRaw,
		&sourceObjectKey,
		&objectKey,
		&state,
		&width,
		&height,
		&format,
		&createTime,
	); err != nil {
		return common.ImageResult{}, err
	}

	tags := make([]string, 0)
	if tagsRaw != "" {
		tags = strings.Split(tagsRaw, ";")
	}
	var metadata map[string]string
	if err := json.Unmarshal(metadataRaw, &metadata); err != nil {
		return common.ImageResult{}, err
	}

	return common.ImageResult{
		ID:              id,
		Tags:            tags,
		Metadata:        metadata,
		SourceObjectKey: sourceObjectKey,
		ObjectKey:       objectKey,
		State:           state,
		Width:           width,
		Height:          height,
		Format:          format,
		CreateTime:      createTime,
	}, nil
}
