package middleware

import (
	"context"
	"errors"
	"net/http"

	apiv1 "github.com/velonyapp/asset/gen/api/v1"
	applicationcommand "github.com/velonyapp/asset/internal/application/command"
	applicationport "github.com/velonyapp/asset/internal/application/port"
	domainrepo "github.com/velonyapp/asset/internal/domain/repo"
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

				// Image Repo
				case errors.Is(err, domainrepo.ErrImageNotFound):
					return reply, kerrors.NotFound(
						apiv1.ErrorReason_IMAGE_NOT_FOUND.String(),
						domainrepo.ErrImageNotFound.Error(),
					)
				case errors.Is(err, domainrepo.ErrObjectKeyAlreadyExists):
					return reply, kerrors.Conflict(
						apiv1.ErrorReason_OBJECT_KEY_ALREADY_EXISTS.String(),
						domainrepo.ErrObjectKeyAlreadyExists.Error(),
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

				// Process Image Command
				case errors.Is(err, applicationcommand.ErrImageObjectDoesntExist):
					return reply, kerrors.New(http.StatusPreconditionFailed,
						apiv1.ErrorReason_OBJECT_NOT_FOUND.String(),
						applicationcommand.ErrImageObjectDoesntExist.Error(),
					)

				// API Service
				case errors.Is(err, presentationapi.ErrInvalidImageResourceName):
					return reply, kerrors.BadRequest(
						"",
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
