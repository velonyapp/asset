package s3

import (
	"context"
	"errors"
	"io"
	"time"

	"github.com/velonyapp/asset/internal/application/port"
	"github.com/velonyapp/asset/internal/conf"
	"github.com/velonyapp/asset/internal/domain/vo"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/aws/smithy-go"
)

var _ port.Storage = (*storage)(nil)

type storage struct {
	client        *s3.Client
	presignClient *s3.PresignClient
	c             *conf.Data
}

func NewStorage(client *s3.Client, c *conf.Data) port.Storage {
	return &storage{
		client:        client,
		presignClient: s3.NewPresignClient(client),
		c:             c,
	}
}

func (s *storage) Exists(ctx context.Context, key vo.ObjectKey) (bool, error) {
	_, err := s.client.HeadObject(ctx, &s3.HeadObjectInput{
		Bucket: aws.String(s.c.S3.Bucket),
		Key:    aws.String(key.Value()),
	})
	if err != nil {
		var apiError smithy.APIError
		if errors.As(err, &apiError) {
			switch apiError.ErrorCode() {
			case "NotFound", "NoSuchKey", "NoSuchObject":
				return false, nil
			}
		}

		return false, err
	}

	return true, nil
}

func (s *storage) Get(ctx context.Context, key vo.ObjectKey) (io.ReadCloser, error) {
	result, err := s.client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(s.c.S3.Bucket),
		Key:    aws.String(key.Value()),
	})
	if err != nil {
		return nil, err
	}

	return result.Body, nil
}

func (s *storage) Put(ctx context.Context, key vo.ObjectKey, body io.Reader) error {
	if _, err := s.client.PutObject(ctx, &s3.PutObjectInput{
		Bucket: aws.String(s.c.S3.Bucket),
		Key:    aws.String(key.Value()),
		Body:   body,
	}); err != nil {
		return err
	}

	return nil
}

func (s *storage) Delete(ctx context.Context, key vo.ObjectKey) error {
	if _, err := s.client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(s.c.S3.Bucket),
		Key:    aws.String(key.Value()),
	}); err != nil {
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
				Key: aws.String(key.Value()),
			})
		}

		result, err := s.client.DeleteObjects(ctx, &s3.DeleteObjectsInput{
			Bucket: aws.String(s.c.S3.Bucket),
			Delete: &types.Delete{
				Objects: objects,
			},
		})
		if err != nil {
			return err
		}

		if len(result.Errors) > 0 {
			deleteError := result.Errors[0]

			if deleteError.Message != nil {
				return errors.New(aws.ToString(deleteError.Message))
			}

			return errors.New(aws.ToString(deleteError.Code))
		}
	}

	return nil
}

func (s *storage) PresignGet(ctx context.Context, key vo.ObjectKey, ttl time.Duration) (string, error) {
	result, err := s.presignClient.PresignGetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(s.c.S3.Bucket),
		Key:    aws.String(key.Value()),
	}, func(options *s3.PresignOptions) {
		options.Expires = ttl
	})
	if err != nil {
		return "", err
	}

	return result.URL, nil
}

func (s *storage) PresignPut(ctx context.Context, key vo.ObjectKey, ttl time.Duration) (string, error) {
	result, err := s.presignClient.PresignPutObject(ctx, &s3.PutObjectInput{
		Bucket: aws.String(s.c.S3.Bucket),
		Key:    aws.String(key.Value()),
	}, func(options *s3.PresignOptions) {
		options.Expires = ttl
	})
	if err != nil {
		return "", err
	}

	return result.URL, nil
}
