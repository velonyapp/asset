package s3

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/velonyapp/asset/internal/application/port"
	"github.com/velonyapp/asset/internal/conf"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/aws/smithy-go"
	smithyhttp "github.com/aws/smithy-go/transport/http"
)

type Storage struct {
	client        *s3.Client
	presignClient *s3.PresignClient
	c             *conf.Data
}

func NewStorage(client *s3.Client, c *conf.Data) port.Storage {
	return &Storage{
		client:        client,
		presignClient: s3.NewPresignClient(client),
		c:             c,
	}
}

func (storage *Storage) Exists(ctx context.Context, key string) (bool, error) {
	_, err := storage.client.HeadObject(ctx, &s3.HeadObjectInput{
		Bucket: aws.String(storage.c.S3.Bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		var responseError *smithyhttp.ResponseError
		if errors.As(err, &responseError) && responseError.HTTPStatusCode() == http.StatusNotFound {
			return false, nil
		}

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

func (storage *Storage) Delete(ctx context.Context, key string) error {
	if _, err := storage.client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(storage.c.S3.Bucket),
		Key:    aws.String(key),
	}); err != nil {
		return err
	}

	return nil
}

func (storage *Storage) DeleteMany(ctx context.Context, keys []string) error {
	const batchSize = 1000

	for start := 0; start < len(keys); start += batchSize {
		end := start + batchSize
		if end > len(keys) {
			end = len(keys)
		}

		objects := make([]types.ObjectIdentifier, 0, end-start)

		for _, key := range keys[start:end] {
			objects = append(objects, types.ObjectIdentifier{
				Key: aws.String(key),
			})
		}

		result, err := storage.client.DeleteObjects(ctx, &s3.DeleteObjectsInput{
			Bucket: aws.String(storage.c.S3.Bucket),
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

func (storage *Storage) PresignGet(ctx context.Context, key string, expiresIn time.Duration) (string, error) {
	result, err := storage.presignClient.PresignGetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(storage.c.S3.Bucket),
		Key:    aws.String(key),
	}, func(options *s3.PresignOptions) {
		options.Expires = expiresIn
	})
	if err != nil {
		return "", err
	}

	return result.URL, nil
}

func (storage *Storage) PresignPut(ctx context.Context, key string, expiresIn time.Duration) (string, error) {
	result, err := storage.presignClient.PresignPutObject(ctx, &s3.PutObjectInput{
		Bucket: aws.String(storage.c.S3.Bucket),
		Key:    aws.String(key),
	}, func(options *s3.PresignOptions) {
		options.Expires = expiresIn
	})
	if err != nil {
		return "", err
	}

	return result.URL, nil
}
