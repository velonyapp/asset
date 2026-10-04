package middleware

import (
	"context"
	"errors"

	apiv1 "github.com/velonyapp/asset/gen/api/v1"
	applicationport "github.com/velonyapp/asset/internal/application/port"
	domainvo "github.com/velonyapp/asset/internal/domain/vo"
	presentationapi "github.com/velonyapp/asset/internal/presentation/api"

	kerrors "github.com/go-kratos/kratos/v3/errors"
	"github.com/go-kratos/kratos/v3/middleware"
)

type ErrorMapper middleware.Middleware

func NewErrorMapperMiddleware() ErrorMapper {
	return ErrorMapper(
		func(next middleware.Handler) middleware.Handler {
			return func(ctx context.Context, req any) (any, error) {
				reply, err := next(ctx, req)
				if err == nil {
					return reply, nil
				}

				var kratosErr *kerrors.Error
				if errors.As(err, &kratosErr) {
					return reply, err
				}

				switch {
				case errors.Is(err, applicationport.ErrInvalidResizeDimensions):
					return reply, kerrors.BadRequest(
						apiv1.ErrorReason_INVALID_IMAGE_RESIZE.String(),
						applicationport.ErrInvalidResizeDimensions.Error(),
					)

				case errors.Is(err, applicationport.ErrUnsupportedResizeFit):
					return reply, kerrors.BadRequest(
						apiv1.ErrorReason_INVALID_IMAGE_RESIZE.String(),
						applicationport.ErrUnsupportedResizeFit.Error(),
					)

				case errors.Is(err, applicationport.ErrUnsupportedImageGravity):
					return reply, kerrors.BadRequest(
						apiv1.ErrorReason_INVALID_IMAGE_RESIZE.String(),
						applicationport.ErrUnsupportedImageGravity.Error(),
					)

				case errors.Is(err, applicationport.ErrInvalidImageBackgroundColor):
					return reply, kerrors.BadRequest(
						apiv1.ErrorReason_INVALID_IMAGE_RESIZE.String(),
						applicationport.ErrInvalidImageBackgroundColor.Error(),
					)

				case errors.Is(err, applicationport.ErrUnsupportedImageFormat):
					return reply, kerrors.BadRequest(
						apiv1.ErrorReason_INVALID_IMAGE_ENCODING.String(),
						applicationport.ErrUnsupportedImageFormat.Error(),
					)

				case errors.Is(err, applicationport.ErrInvalidImageQuality):
					return reply, kerrors.BadRequest(
						apiv1.ErrorReason_INVALID_IMAGE_ENCODING.String(),
						applicationport.ErrInvalidImageQuality.Error(),
					)

				case errors.Is(err, domainvo.ErrObjectKeyEmpty):
					return reply, kerrors.BadRequest(
						apiv1.ErrorReason_INVALID_STORAGE_KEY.String(),
						domainvo.ErrObjectKeyEmpty.Error(),
					)

				case errors.Is(err, domainvo.ErrObjectKeyTooLong):
					return reply, kerrors.BadRequest(
						apiv1.ErrorReason_INVALID_STORAGE_KEY.String(),
						domainvo.ErrObjectKeyTooLong.Error(),
					)

				case errors.Is(err, domainvo.ErrObjectKeyElementEmpty):
					return reply, kerrors.BadRequest(
						apiv1.ErrorReason_INVALID_STORAGE_KEY.String(),
						domainvo.ErrObjectKeyElementEmpty.Error(),
					)

				case errors.Is(err, domainvo.ErrObjectKeyInvalidCharacter):
					return reply, kerrors.BadRequest(
						apiv1.ErrorReason_INVALID_STORAGE_KEY.String(),
						domainvo.ErrObjectKeyInvalidCharacter.Error(),
					)

				case errors.Is(err, presentationapi.ErrInvalidImageResourceName):
					return reply, kerrors.BadRequest(
						"",
						presentationapi.ErrInvalidImageResourceName.Error(),
					)

				default:
					return reply, err // debug

					// Production:
					// return reply, kerrors.InternalServer(
					// 	"",
					// 	"internal server error",
					// )
				}
			}
		},
	)
}

func (m ErrorMapper) Middleware() middleware.Middleware {
	return middleware.Middleware(m)
}
