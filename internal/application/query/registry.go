package query

type HandlerRegistry struct{}

func NewHandlerRegistry(
	getImage GetImageHandler,
	presignImage PresignImageHandler,
) *HandlerRegistry {
	Register(getImage)
	Register(presignImage)

	return &HandlerRegistry{}
}
