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

	case errors.Is(err, domainvo.ErrObjectKeyEmpty):
		return kerrors.BadRequest(
			apiv1.ErrorReason_INVALID_STORAGE_KEY.String(),
			domainvo.ErrObjectKeyEmpty.Error(),
		)

	case errors.Is(err, domainvo.ErrObjectKeyTooLong):
		return kerrors.BadRequest(
			apiv1.ErrorReason_INVALID_STORAGE_KEY.String(),
			domainvo.ErrObjectKeyTooLong.Error(),
		)

	case errors.Is(err, domainvo.ErrObjectKeyElementEmpty):
		return kerrors.BadRequest(
			apiv1.ErrorReason_INVALID_STORAGE_KEY.String(),
			domainvo.ErrObjectKeyElementEmpty.Error(),
		)

	case errors.Is(err, domainvo.ErrObjectKeyInvalidCharacter):
		return kerrors.BadRequest(
			apiv1.ErrorReason_INVALID_STORAGE_KEY.String(),
			domainvo.ErrObjectKeyInvalidCharacter.Error(),
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
