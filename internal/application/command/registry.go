package command

type HandlerRegistry struct{}

func NewHandlerRegistry(
	createImage CreateImageHandler,
	processImage ProcessImageHandler,
	reconcileImage ReconcileImageHandler,
	deleteImage DeleteImageHandler,
) *HandlerRegistry {
	Register(createImage)
	Register(processImage)
	Register(reconcileImage)
	Register(deleteImage)

	return &HandlerRegistry{}
}
