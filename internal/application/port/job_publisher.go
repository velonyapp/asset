package port

import (
	"context"

	"github.com/velonyapp/asset/internal/application/job"
)

type JobPublisher interface {
	Publish(ctx context.Context, job job.Job) error
}
