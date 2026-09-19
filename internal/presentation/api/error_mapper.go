package api

import (
	"errors"

	apiv1 "github.com/velonyapp/asset/gen/api/v1"
	applicationport "github.com/velonyapp/asset/internal/application/port"
	domainvo "github.com/velonyapp/asset/internal/domain/vo"

	kerrors "github.com/go-kratos/kratos/v3/errors"
)

func mapError(err error) error {
	if err == nil {
		return nil
	}

	switch {
	case errors.Is(err, applicationport.ErrInvalidResizeDimensions):
		return kerrors.BadRequest(
			apiv1.ErrorReason_INVALID_IMAGE_RESIZE.String(),
			applicationport.ErrInvalidResizeDimensions.Error(),
		)

	case errors.Is(err, applicationport.ErrUnsupportedResizeFit):
		return kerrors.BadRequest(
			apiv1.ErrorReason_INVALID_IMAGE_RESIZE.String(),
			applicationport.ErrUnsupportedResizeFit.Error(),
		)

	case errors.Is(err, applicationport.ErrUnsupportedImageGravity):
		return kerrors.BadRequest(
			apiv1.ErrorReason_INVALID_IMAGE_RESIZE.String(),
			applicationport.ErrUnsupportedImageGravity.Error(),
		)

	case errors.Is(err, applicationport.ErrInvalidImageBackgroundColor):
		return kerrors.BadRequest(
			apiv1.ErrorReason_INVALID_IMAGE_RESIZE.String(),
			applicationport.ErrInvalidImageBackgroundColor.Error(),
		)

	case errors.Is(err, applicationport.ErrImageUpscaleNotAllowed):
		return kerrors.BadRequest(
			apiv1.ErrorReason_INVALID_IMAGE_RESIZE.String(),
			applicationport.ErrImageUpscaleNotAllowed.Error(),
		)

	case errors.Is(err, applicationport.ErrUnsupportedImageFormat):
		return kerrors.BadRequest(
			apiv1.ErrorReason_INVALID_IMAGE_ENCODING.String(),
			applicationport.ErrUnsupportedImageFormat.Error(),
		)

	case errors.Is(err, applicationport.ErrInvalidImageQuality):
		return kerrors.BadRequest(
			apiv1.ErrorReason_INVALID_IMAGE_ENCODING.String(),
			applicationport.ErrInvalidImageQuality.Error(),
		)

	case errors.Is(err, applicationport.ErrInvalidUploadToken):
		return kerrors.Unauthorized(
			apiv1.ErrorReason_INVALID_UPLOAD_TOKEN.String(),
			applicationport.ErrInvalidUploadToken.Error(),
		)

	case errors.Is(err, applicationport.ErrExpiredUploadToken):
		return kerrors.Unauthorized(
			apiv1.ErrorReason_INVALID_UPLOAD_TOKEN.String(),
			applicationport.ErrInvalidUploadToken.Error(),
		)

	case errors.Is(err, domainvo.ErrStorageKeyEmpty):
		return kerrors.BadRequest(
			apiv1.ErrorReason_INVALID_STORAGE_KEY.String(),
			domainvo.ErrStorageKeyEmpty.Error(),
		)

	case errors.Is(err, domainvo.ErrStorageKeyTooLong):
		return kerrors.BadRequest(
			apiv1.ErrorReason_INVALID_STORAGE_KEY.String(),
			domainvo.ErrStorageKeyTooLong.Error(),
		)

	case errors.Is(err, domainvo.ErrStorageKeyInvalidUTF8):
		return kerrors.BadRequest(
			apiv1.ErrorReason_INVALID_STORAGE_KEY.String(),
			domainvo.ErrStorageKeyInvalidUTF8.Error(),
		)

	default:
		return err // debug

		// Production:
		// return kerrors.InternalServer(
		// 	"",
		// 	"internal server error",
		// )
	}
}
