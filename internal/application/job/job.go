package job

import "github.com/google/uuid"

type Job interface {
	ID() string
	Type() string
}

type BaseJob struct {
	id string
}

func NewBaseJob() BaseJob {
	return BaseJob{
		id: uuid.Must(uuid.NewV7()).String(),
	}
}

func (e BaseJob) ID() string { return e.id }
