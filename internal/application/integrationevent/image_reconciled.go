package integrationevent

import "time"

var _ IntegrationEvent = (*ImageReconciled)(nil)

type ImageReconciled struct {
	BaseIntegrationEvent

	objectExists bool
}

func NewImageReconciled(
	imageID string,
	tags []string,
	occurTime time.Time,
	objectExists bool,
) ImageReconciled {
	return ImageReconciled{
		BaseIntegrationEvent: NewBaseIntegrationEvent(imageID, tags, occurTime),

		objectExists: objectExists,
	}
}

func (e ImageReconciled) Type() string {
	return "asset.image.reconciled"
}

func (e ImageReconciled) AggregateType() string {
	return "image"
}

func (e ImageReconciled) ObjectExists() bool {
	return e.objectExists
}
