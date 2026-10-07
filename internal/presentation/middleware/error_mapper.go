package middleware

import (
	"context"
	"errors"

	apiv1 "github.com/velonyapp/asset/gen/api/v1"
	applicationcommand "github.com/velonyapp/asset/internal/application/command"
	applicationcommon "github.com/velonyapp/asset/internal/application/common"
	applicationport "github.com/velonyapp/asset/internal/application/port"
	domainentity "github.com/velonyapp/asset/internal/domain/entity"
	domainservice "github.com/velonyapp/asset/internal/domain/service"
	domainvo "github.com/velonyapp/asset/internal/domain/vo"
	presentationapi "github.com/velonyapp/asset/internal/presentation/api"

	kerrors "github.com/go-kratos/kratos/v3/errors"
	"github.com/go-kratos/kratos/v3/middleware"
)

type ErrorMapper middleware.Middleware

func NewErrorMapper() ErrorMapper {
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
				// Image ID
				case errors.Is(err, domainvo.ErrImageIDInvalid):
					return reply, kerrors.BadRequest("",
						domainvo.ErrImageIDInvalid.Error(),
					)

				// Image State
				case errors.Is(err, domainvo.ErrImageStateInvalid):
					return reply, kerrors.BadRequest("",
						domainvo.ErrImageStateInvalid.Error(),
					)

				// Object Key
				case errors.Is(err, domainvo.ErrObjectKeyEmpty):
					return reply, kerrors.BadRequest("",
						domainvo.ErrObjectKeyEmpty.Error(),
					)
				case errors.Is(err, domainvo.ErrObjectKeyElementEmpty):
					return reply, kerrors.BadRequest("",
						domainvo.ErrObjectKeyElementEmpty.Error(),
					)
				case errors.Is(err, domainvo.ErrObjectKeyTooLong):
					return reply, kerrors.BadRequest("",
						domainvo.ErrObjectKeyTooLong.Error(),
					)
				case errors.Is(err, domainvo.ErrObjectKeyInvalidCharacter):
					return reply, kerrors.BadRequest("",
						domainvo.ErrObjectKeyInvalidCharacter.Error(),
					)

				// Tag
				case errors.Is(err, domainvo.ErrTagEmpty):
					return reply, kerrors.BadRequest("",
						domainvo.ErrTagEmpty.Error(),
					)
				case errors.Is(err, domainvo.ErrTagTooLong):
					return reply, kerrors.BadRequest("",
						domainvo.ErrTagTooLong.Error(),
					)
				case errors.Is(err, domainvo.ErrTagInvalidCharacter):
					return reply, kerrors.BadRequest("",
						domainvo.ErrTagInvalidCharacter.Error(),
					)

				// Object Key Policy
				case errors.Is(err, domainservice.ErrObjectKeyAlreadyExists):
					return reply, kerrors.Conflict("",
						domainservice.ErrObjectKeyAlreadyExists.Error(),
					)

				// Image
				case errors.Is(err, domainentity.ErrImageDeleted):
					return reply, kerrors.BadRequest("",
						domainentity.ErrImageDeleted.Error(),
					)
				case errors.Is(err, domainentity.ErrImageAlreadyDeleted):
					return reply, kerrors.BadRequest("",
						domainentity.ErrImageAlreadyDeleted.Error(),
					)
				case errors.Is(err, domainentity.ErrImageNotPending):
					return reply, kerrors.BadRequest("",
						domainentity.ErrImageNotPending.Error(),
					)
				case errors.Is(err, domainentity.ErrImageAlreadyUploaded):
					return reply, kerrors.BadRequest("",
						domainentity.ErrImageAlreadyUploaded.Error(),
					)
				case errors.Is(err, domainentity.ErrImageNotUploaded):
					return reply, kerrors.BadRequest("",
						domainentity.ErrImageNotUploaded.Error(),
					)
				case errors.Is(err, domainentity.ErrImageProcessed):
					return reply, kerrors.BadRequest("",
						domainentity.ErrImageProcessed.Error(),
					)
				case errors.Is(err, domainentity.ErrImageAlreadyProcessed):
					return reply, kerrors.BadRequest("",
						domainentity.ErrImageAlreadyProcessed.Error(),
					)
				case errors.Is(err, domainentity.ErrImageObjectKeysEqual):
					return reply, kerrors.BadRequest("",
						domainentity.ErrImageObjectKeysEqual.Error(),
					)

				// Image Processor
				case errors.Is(err, applicationport.ErrNoImageProcessingOptions):
					return reply, kerrors.BadRequest("",
						applicationport.ErrNoImageProcessingOptions.Error(),
					)
				case errors.Is(err, applicationport.ErrInvalidResizeDimensions):
					return reply, kerrors.BadRequest("",
						applicationport.ErrInvalidResizeDimensions.Error(),
					)
				case errors.Is(err, applicationport.ErrUnsupportedResizeFit):
					return reply, kerrors.BadRequest("",
						applicationport.ErrUnsupportedResizeFit.Error(),
					)
				case errors.Is(err, applicationport.ErrUnsupportedImageGravity):
					return reply, kerrors.BadRequest("",
						applicationport.ErrUnsupportedImageGravity.Error(),
					)
				case errors.Is(err, applicationport.ErrUnsupportedImageFormat):
					return reply, kerrors.BadRequest("",
						applicationport.ErrUnsupportedImageFormat.Error(),
					)
				case errors.Is(err, applicationport.ErrInvalidImageQuality):
					return reply, kerrors.BadRequest("",
						applicationport.ErrInvalidImageQuality.Error(),
					)
				case errors.Is(err, applicationport.ErrInvalidImageBackgroundColor):
					return reply, kerrors.BadRequest("",
						applicationport.ErrInvalidImageBackgroundColor.Error(),
					)

				// Application Common
				case errors.Is(err, applicationcommon.ErrImageNotFound):
					return reply, kerrors.NotFound(
						apiv1.ErrorReason_IMAGE_NOT_FOUND.String(),
						applicationcommon.ErrImageNotFound.Error(),
					)

				// Confirm Image Upload
				case errors.Is(err, applicationcommand.ErrSourceObjectNotFound):
					return reply, kerrors.BadRequest(
						apiv1.ErrorReason_SOURCE_OBJECT_NOT_FOUND.String(),
						applicationcommand.ErrSourceObjectNotFound.Error(),
					)

				// Create Image Command
				case errors.Is(err, applicationcommand.ErrObjectAlreadyExists):
					return reply, kerrors.Conflict(
						apiv1.ErrorReason_OBJECT_KEY_ALREADY_EXISTS.String(),
						applicationcommand.ErrObjectAlreadyExists.Error(),
					)

				// API Service
				case errors.Is(err, presentationapi.ErrInvalidImageResourceName):
					return reply, kerrors.BadRequest("",
						presentationapi.ErrInvalidImageResourceName.Error(),
					)

				default:
					return reply, kerrors.InternalServer("",
						"internal server error",
					).WithCause(err)
				}
			}
		},
	)
}

func (m ErrorMapper) Middleware() middleware.Middleware {
	return middleware.Middleware(m)
}
