package s3

import (
	"context"
	"log/slog"

	"github.com/velonyapp/asset/internal/conf"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"go.opentelemetry.io/contrib/instrumentation/github.com/aws/aws-sdk-go-v2/otelaws"
)

func NewClient(c *conf.Data, log *slog.Logger) (*s3.Client, error) {
	cfg, err := config.LoadDefaultConfig(context.Background(),
		config.WithRegion("auto"),
		config.WithCredentialsProvider(
			credentials.NewStaticCredentialsProvider(
				c.S3.AccessKeyId,
				c.S3.SecretAccessKey,
				"",
			),
		),
	)
	if err != nil {
		log.Error(
			"failed to load S3 configuration",
			slog.Any("error", err),
		)
		return nil, err
	}

	otelaws.AppendMiddlewares(&cfg.APIOptions)

	client := s3.NewFromConfig(cfg, func(options *s3.Options) {
		options.BaseEndpoint = aws.String(c.S3.Endpoint)
	})

	log.Info(
		"S3 client initialized",
	)

	return client, nil
}