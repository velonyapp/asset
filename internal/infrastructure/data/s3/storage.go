package s3

import (
	"context"
	"errors"
	"io"
	"log/slog"
	"time"

	"github.com/velonyapp/asset/internal/application/port"
	"github.com/velonyapp/asset/internal/conf"
	"github.com/velonyapp/asset/internal/domain/vo"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/s3/transfermanager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/aws/smithy-go"
)

var _ port.Storage = (*storage)(nil)

type storage struct {
	client                *s3.Client
	presignClient         *s3.PresignClient
	transferManagerClient *transfermanager.Client
	c                     *conf.Data
	log                   *slog.Logger
}

func NewStorage(
	client *s3.Client,
	c *conf.Data,
	log *slog.Logger,
) port.Storage {
	return &storage{
		client:                client,
		presignClient:         s3.NewPresignClient(client),
		transferManagerClient: transfermanager.New(client),
		c:                     c,
		log:                   log,
	}
}

func (s *storage) Exists(ctx context.Context, key vo.ObjectKey) (bool, error) {
	_, err := s.client.HeadObject(ctx, &s3.HeadObjectInput{
		Bucket: aws.String(s.c.S3.Bucket),
		Key:    aws.String(key.String()),
	})
	if err != nil {
		var apiError smithy.APIError
		if errors.As(err, &apiError) {
			switch apiError.ErrorCode() {
			case "NotFound", "NoSuchKey", "NoSuchObject":
				return false, nil
			}
		}

		s.log.ErrorContext(ctx,
			"failed to check S3 object existence",
			slog.String("object_key", key.String()),
			slog.Any("error", err),
		)
		return false, err
	}

	return true, nil
}

func (s *storage) Get(ctx context.Context, key vo.ObjectKey) (io.ReadCloser, error) {
	result, err := s.client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(s.c.S3.Bucket),
		Key:    aws.String(key.String()),
	})
	if err != nil {
		s.log.ErrorContext(ctx,
			"failed to get S3 object",
			slog.String("object_key", key.String()),
			slog.Any("error", err),
		)
		return nil, err
	}

	return result.Body, nil
}

func (s *storage) Put(ctx context.Context, key vo.ObjectKey, body io.Reader) error {
	if _, err := s.transferManagerClient.UploadObject(ctx, &transfermanager.UploadObjectInput{
		Bucket: aws.String(s.c.S3.Bucket),
		Key:    aws.String(key.String()),
		Body:   body,
	}); err != nil {
		s.log.ErrorContext(ctx,
			"failed to put S3 object",
			slog.String("object_key", key.String()),
			slog.Any("error", err),
		)
		return err
	}

	return nil
}

func (s *storage) Delete(ctx context.Context, key vo.ObjectKey) error {
	if _, err := s.client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(s.c.S3.Bucket),
		Key:    aws.String(key.String()),
	}); err != nil {
		s.log.ErrorContext(ctx,
			"failed to delete S3 object",
			slog.String("object_key", key.String()),
			slog.Any("error", err),
		)
		return err
	}

	return nil
}

func (s *storage) DeleteMany(ctx context.Context, keys []vo.ObjectKey) error {
	const batchSize = 1000

	for start := 0; start < len(keys); start += batchSize {
		end := start + batchSize
		if end > len(keys) {
			end = len(keys)
		}

		objects := make([]types.ObjectIdentifier, 0, end-start)

		for _, key := range keys[start:end] {
			objects = append(objects, types.ObjectIdentifier{
				Key: aws.String(key.String()),
			})
		}

		result, err := s.client.DeleteObjects(ctx, &s3.DeleteObjectsInput{
			Bucket: aws.String(s.c.S3.Bucket),
			Delete: &types.Delete{
				Objects: objects,
			},
		})
		if err != nil {
			s.log.ErrorContext(ctx,
				"failed to delete S3 objects",
				slog.Int("object_count", len(objects)),
				slog.Any("error", err),
			)
			return err
		}

		if len(result.Errors) > 0 {
			deleteError := result.Errors[0]

			if deleteError.Message != nil {
				err := errors.New(aws.ToString(deleteError.Message))
				s.log.ErrorContext(ctx,
					"failed to delete S3 objects",
					slog.Int("object_count", len(objects)),
					slog.Any("error", err),
				)
				return err
			}

			err := errors.New(aws.ToString(deleteError.Code))
			s.log.ErrorContext(ctx,
				"failed to delete S3 objects",
				slog.Int("object_count", len(objects)),
				slog.Any("error", err),
			)
			return err
		}
	}

	return nil
}

func (s *storage) PresignGet(ctx context.Context, key vo.ObjectKey, ttl time.Duration) (string, error) {
	result, err := s.presignClient.PresignGetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(s.c.S3.Bucket),
		Key:    aws.String(key.String()),
	}, func(options *s3.PresignOptions) {
		options.Expires = ttl
	})
	if err != nil {
		s.log.ErrorContext(ctx,
			"failed to presign S3 get object",
			slog.String("object_key", key.String()),
			slog.Any("error", err),
		)
		return "", err
	}

	return result.URL, nil
}

func (s *storage) PresignPut(ctx context.Context, key vo.ObjectKey, ttl time.Duration) (string, error) {
	result, err := s.presignClient.PresignPutObject(ctx, &s3.PutObjectInput{
		Bucket: aws.String(s.c.S3.Bucket),
		Key:    aws.String(key.String()),
	}, func(options *s3.PresignOptions) {
		options.Expires = ttl
	})
	if err != nil {
		s.log.ErrorContext(ctx,
			"failed to presign S3 put object",
			slog.String("object_key", key.String()),
			slog.Any("error", err),
		)
		return "", err
	}

	return result.URL, nil
}
