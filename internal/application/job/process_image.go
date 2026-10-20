package job

var _ Job = (*ProcessImage)(nil)

type ProcessImage struct {
	BaseJob

	imageID string
}

func NewProcessImage(
	imageID string,
) ProcessImage {
	return ProcessImage{
		BaseJob: NewBaseJob(),

		imageID: imageID,
	}
}

func (e ProcessImage) Type() string {
	return "asset.image.process"
}

func (e ProcessImage) ImageID() string {
	return e.imageID
}
