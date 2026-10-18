package middleware

import (
	"context"
	"errors"

	apiv1 "github.com/velonyapp/asset/gen/api/v1"
	applicationcommand "github.com/velonyapp/asset/internal/application/command"
	applicationcommon "github.com/velonyapp/asset/internal/application/common"
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

				// Image Width
				case errors.Is(err, domainvo.ErrImageWidthNotPositive):
					return reply, kerrors.BadRequest("",
						domainvo.ErrImageWidthNotPositive.Error(),
					)

				// Image Height
				case errors.Is(err, domainvo.ErrImageHeightNotPositive):
					return reply, kerrors.BadRequest("",
						domainvo.ErrImageHeightNotPositive.Error(),
					)

				// Image Format
				case errors.Is(err, domainvo.ErrImageFormatInvalid):
					return reply, kerrors.BadRequest("",
						domainvo.ErrImageFormatInvalid.Error(),
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
				case errors.Is(err, domainentity.ErrImageNotAwaitingUpload):
					return reply, kerrors.BadRequest("",
						domainentity.ErrImageNotAwaitingUpload.Error(),
					)
				case errors.Is(err, domainentity.ErrImageAlreadyPending):
					return reply, kerrors.BadRequest("",
						domainentity.ErrImageAlreadyPending.Error(),
					)
				case errors.Is(err, domainentity.ErrImageNotPending):
					return reply, kerrors.BadRequest("",
						domainentity.ErrImageNotPending.Error(),
					)
				case errors.Is(err, domainentity.ErrImageProcessing):
					return reply, kerrors.BadRequest("",
						domainentity.ErrImageProcessing.Error(),
					)
				case errors.Is(err, domainentity.ErrImageAlreadyProcessing):
					return reply, kerrors.BadRequest("",
						domainentity.ErrImageAlreadyProcessing.Error(),
					)
				case errors.Is(err, domainentity.ErrImageNotProcessing):
					return reply, kerrors.BadRequest("",
						domainentity.ErrImageNotProcessing.Error(),
					)
				case errors.Is(err, domainentity.ErrImageActive):
					return reply, kerrors.BadRequest("",
						domainentity.ErrImageActive.Error(),
					)
				case errors.Is(err, domainentity.ErrImageAlreadyActive):
					return reply, kerrors.BadRequest("",
						domainentity.ErrImageAlreadyActive.Error(),
					)
				case errors.Is(err, domainentity.ErrImageObjectKeysEqual):
					return reply, kerrors.BadRequest("",
						domainentity.ErrImageObjectKeysEqual.Error(),
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
