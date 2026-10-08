package entity

import (
	"errors"
	"testing"
	"time"

	"github.com/velonyapp/asset/internal/domain/event"
	"github.com/velonyapp/asset/internal/domain/vo"
)

func TestImage_New_WhenObjectKeysEqual(t *testing.T) {
	image, err := NewImage(
		vo.Tags{},
		vo.ObjectKey{},
		vo.ObjectKey{},
		time.Now().UTC(),
	)

	if !errors.Is(err, ErrImageObjectKeysEqual) {
		t.Fatalf("expected ErrImageObjectKeysEqual, got %v", err)
	}

	if image != nil {
		t.Error("expected no image")
	}
}

func TestImage_CanStartUploading(t *testing.T) {
	image := ReconstituteImage(
		vo.ImageID{},
		vo.Tags{},
		vo.ObjectKey{},
		vo.ObjectKey{},
		vo.ImageStateCreated,
		time.Time{},
		time.Time{},
	)

	err := image.CanStartUploading()

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if !image.IsCreated() {
		t.Error("expected image to remain created")
	}

	if !image.UpdateTime().IsZero() {
		t.Errorf("expected update time to remain zero, got %v", image.UpdateTime())
	}

	if events := image.PullEvents(); len(events) != 0 {
		t.Errorf("expected no events, got %d", len(events))
	}
}

func TestImage_CanStartUploading_WhenAlreadyUploading(t *testing.T) {
	image := ReconstituteImage(
		vo.ImageID{},
		vo.Tags{},
		vo.ObjectKey{},
		vo.ObjectKey{},
		vo.ImageStateUploading,
		time.Time{},
		time.Time{},
	)

	err := image.CanStartUploading()

	if !errors.Is(err, ErrImageAlreadyUploading) {
		t.Fatalf("expected ErrImageAlreadyUploading, got %v", err)
	}

	if !image.IsUploading() {
		t.Error("expected image to remain uploading")
	}

	if !image.UpdateTime().IsZero() {
		t.Errorf("expected update time to remain zero, got %v", image.UpdateTime())
	}

	if events := image.PullEvents(); len(events) != 0 {
		t.Errorf("expected no events, got %d", len(events))
	}
}

func TestImage_CanStartUploading_WhenUploaded(t *testing.T) {
	image := ReconstituteImage(
		vo.ImageID{},
		vo.Tags{},
		vo.ObjectKey{},
		vo.ObjectKey{},
		vo.ImageStateUploaded,
		time.Time{},
		time.Time{},
	)

	err := image.CanStartUploading()

	if !errors.Is(err, ErrImageUploaded) {
		t.Fatalf("expected ErrImageUploaded, got %v", err)
	}

	if !image.IsUploaded() {
		t.Error("expected image to remain uploaded")
	}

	if !image.UpdateTime().IsZero() {
		t.Errorf("expected update time to remain zero, got %v", image.UpdateTime())
	}

	if events := image.PullEvents(); len(events) != 0 {
		t.Errorf("expected no events, got %d", len(events))
	}
}

func TestImage_CanStartUploading_WhenProcessing(t *testing.T) {
	image := ReconstituteImage(
		vo.ImageID{},
		vo.Tags{},
		vo.ObjectKey{},
		vo.ObjectKey{},
		vo.ImageStateProcessing,
		time.Time{},
		time.Time{},
	)

	err := image.CanStartUploading()

	if !errors.Is(err, ErrImageProcessing) {
		t.Fatalf("expected ErrImageProcessing, got %v", err)
	}

	if !image.IsProcessing() {
		t.Error("expected image to remain processing")
	}

	if !image.UpdateTime().IsZero() {
		t.Errorf("expected update time to remain zero, got %v", image.UpdateTime())
	}

	if events := image.PullEvents(); len(events) != 0 {
		t.Errorf("expected no events, got %d", len(events))
	}
}

func TestImage_CanStartUploading_WhenProcessed(t *testing.T) {
	image := ReconstituteImage(
		vo.ImageID{},
		vo.Tags{},
		vo.ObjectKey{},
		vo.ObjectKey{},
		vo.ImageStateProcessed,
		time.Time{},
		time.Time{},
	)

	err := image.CanStartUploading()

	if !errors.Is(err, ErrImageProcessed) {
		t.Fatalf("expected ErrImageProcessed, got %v", err)
	}

	if !image.IsProcessed() {
		t.Error("expected image to remain processed")
	}

	if !image.UpdateTime().IsZero() {
		t.Errorf("expected update time to remain zero, got %v", image.UpdateTime())
	}

	if events := image.PullEvents(); len(events) != 0 {
		t.Errorf("expected no events, got %d", len(events))
	}
}

func TestImage_CanStartUploading_WhenDeleted(t *testing.T) {
	image := ReconstituteImage(
		vo.ImageID{},
		vo.Tags{},
		vo.ObjectKey{},
		vo.ObjectKey{},
		vo.ImageStateDeleted,
		time.Time{},
		time.Time{},
	)

	err := image.CanStartUploading()

	if !errors.Is(err, ErrImageDeleted) {
		t.Fatalf("expected ErrImageDeleted, got %v", err)
	}

	if !image.IsDeleted() {
		t.Error("expected image to remain deleted")
	}

	if !image.UpdateTime().IsZero() {
		t.Errorf("expected update time to remain zero, got %v", image.UpdateTime())
	}

	if events := image.PullEvents(); len(events) != 0 {
		t.Errorf("expected no events, got %d", len(events))
	}
}

func TestImage_StartUploading(t *testing.T) {
	image := ReconstituteImage(
		vo.ImageID{},
		vo.Tags{},
		vo.ObjectKey{},
		vo.ObjectKey{},
		vo.ImageStateCreated,
		time.Time{},
		time.Time{},
	)

	now := time.Now().UTC()

	err := image.StartUploading(now)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if !image.IsUploading() {
		t.Error("expected image to be uploading")
	}

	if !image.UpdateTime().Equal(now) {
		t.Errorf("expected update time %v, got %v", now, image.UpdateTime())
	}

	if events := image.PullEvents(); len(events) != 0 {
		t.Errorf("expected no events, got %d", len(events))
	}
}

func TestImage_StartUploading_WhenAlreadyUploading(t *testing.T) {
	image := ReconstituteImage(
		vo.ImageID{},
		vo.Tags{},
		vo.ObjectKey{},
		vo.ObjectKey{},
		vo.ImageStateUploading,
		time.Time{},
		time.Time{},
	)

	now := time.Now().UTC()

	err := image.StartUploading(now)

	if !errors.Is(err, ErrImageAlreadyUploading) {
		t.Fatalf("expected ErrImageAlreadyUploading, got %v", err)
	}

	if !image.IsUploading() {
		t.Error("expected image to remain uploading")
	}

	if !image.UpdateTime().IsZero() {
		t.Errorf("expected update time to remain zero, got %v", image.UpdateTime())
	}

	if events := image.PullEvents(); len(events) != 0 {
		t.Errorf("expected no events, got %d", len(events))
	}
}

func TestImage_StartUploading_WhenUploaded(t *testing.T) {
	image := ReconstituteImage(
		vo.ImageID{},
		vo.Tags{},
		vo.ObjectKey{},
		vo.ObjectKey{},
		vo.ImageStateUploaded,
		time.Time{},
		time.Time{},
	)

	now := time.Now().UTC()

	err := image.StartUploading(now)

	if !errors.Is(err, ErrImageUploaded) {
		t.Fatalf("expected ErrImageUploaded, got %v", err)
	}

	if !image.IsUploaded() {
		t.Error("expected image to remain uploaded")
	}

	if !image.UpdateTime().IsZero() {
		t.Errorf("expected update time to remain zero, got %v", image.UpdateTime())
	}

	if events := image.PullEvents(); len(events) != 0 {
		t.Errorf("expected no events, got %d", len(events))
	}
}

func TestImage_StartUploading_WhenProcessing(t *testing.T) {
	image := ReconstituteImage(
		vo.ImageID{},
		vo.Tags{},
		vo.ObjectKey{},
		vo.ObjectKey{},
		vo.ImageStateProcessing,
		time.Time{},
		time.Time{},
	)

	now := time.Now().UTC()

	err := image.StartUploading(now)

	if !errors.Is(err, ErrImageProcessing) {
		t.Fatalf("expected ErrImageProcessing, got %v", err)
	}

	if !image.IsProcessing() {
		t.Error("expected image to remain processing")
	}

	if !image.UpdateTime().IsZero() {
		t.Errorf("expected update time to remain zero, got %v", image.UpdateTime())
	}

	if events := image.PullEvents(); len(events) != 0 {
		t.Errorf("expected no events, got %d", len(events))
	}
}

func TestImage_StartUploading_WhenProcessed(t *testing.T) {
	image := ReconstituteImage(
		vo.ImageID{},
		vo.Tags{},
		vo.ObjectKey{},
		vo.ObjectKey{},
		vo.ImageStateProcessed,
		time.Time{},
		time.Time{},
	)

	now := time.Now().UTC()

	err := image.StartUploading(now)

	if !errors.Is(err, ErrImageProcessed) {
		t.Fatalf("expected ErrImageProcessed, got %v", err)
	}

	if !image.IsProcessed() {
		t.Error("expected image to remain processed")
	}

	if !image.UpdateTime().IsZero() {
		t.Errorf("expected update time to remain zero, got %v", image.UpdateTime())
	}

	if events := image.PullEvents(); len(events) != 0 {
		t.Errorf("expected no events, got %d", len(events))
	}
}

func TestImage_StartUploading_WhenDeleted(t *testing.T) {
	image := ReconstituteImage(
		vo.ImageID{},
		vo.Tags{},
		vo.ObjectKey{},
		vo.ObjectKey{},
		vo.ImageStateDeleted,
		time.Time{},
		time.Time{},
	)

	now := time.Now().UTC()

	err := image.StartUploading(now)

	if !errors.Is(err, ErrImageDeleted) {
		t.Fatalf("expected ErrImageDeleted, got %v", err)
	}

	if !image.IsDeleted() {
		t.Error("expected image to remain deleted")
	}

	if !image.UpdateTime().IsZero() {
		t.Errorf("expected update time to remain zero, got %v", image.UpdateTime())
	}

	if events := image.PullEvents(); len(events) != 0 {
		t.Errorf("expected no events, got %d", len(events))
	}
}

func TestImage_CanConfirmUpload_WhenCreated(t *testing.T) {
	image := ReconstituteImage(
		vo.ImageID{},
		vo.Tags{},
		vo.ObjectKey{},
		vo.ObjectKey{},
		vo.ImageStateCreated,
		time.Time{},
		time.Time{},
	)

	err := image.CanConfirmUpload()

	if !errors.Is(err, ErrImageNotUploading) {
		t.Fatalf("expected ErrImageNotUploading, got %v", err)
	}

	if !image.IsCreated() {
		t.Error("expected image to remain created")
	}

	if !image.UpdateTime().IsZero() {
		t.Errorf("expected update time to remain zero, got %v", image.UpdateTime())
	}

	if events := image.PullEvents(); len(events) != 0 {
		t.Errorf("expected no events, got %d", len(events))
	}
}

func TestImage_CanConfirmUpload(t *testing.T) {
	image := ReconstituteImage(
		vo.ImageID{},
		vo.Tags{},
		vo.ObjectKey{},
		vo.ObjectKey{},
		vo.ImageStateUploading,
		time.Time{},
		time.Time{},
	)

	err := image.CanConfirmUpload()

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if !image.IsUploading() {
		t.Error("expected image to remain uploading")
	}

	if !image.UpdateTime().IsZero() {
		t.Errorf("expected update time to remain zero, got %v", image.UpdateTime())
	}

	if events := image.PullEvents(); len(events) != 0 {
		t.Errorf("expected no events, got %d", len(events))
	}
}

func TestImage_CanConfirmUpload_WhenAlreadyUploaded(t *testing.T) {
	image := ReconstituteImage(
		vo.ImageID{},
		vo.Tags{},
		vo.ObjectKey{},
		vo.ObjectKey{},
		vo.ImageStateUploaded,
		time.Time{},
		time.Time{},
	)

	err := image.CanConfirmUpload()

	if !errors.Is(err, ErrImageAlreadyUploaded) {
		t.Fatalf("expected ErrImageAlreadyUploaded, got %v", err)
	}

	if !image.IsUploaded() {
		t.Error("expected image to remain uploaded")
	}

	if !image.UpdateTime().IsZero() {
		t.Errorf("expected update time to remain zero, got %v", image.UpdateTime())
	}

	if events := image.PullEvents(); len(events) != 0 {
		t.Errorf("expected no events, got %d", len(events))
	}
}

func TestImage_CanConfirmUpload_WhenProcessing(t *testing.T) {
	image := ReconstituteImage(
		vo.ImageID{},
		vo.Tags{},
		vo.ObjectKey{},
		vo.ObjectKey{},
		vo.ImageStateProcessing,
		time.Time{},
		time.Time{},
	)

	err := image.CanConfirmUpload()

	if !errors.Is(err, ErrImageProcessing) {
		t.Fatalf("expected ErrImageProcessing, got %v", err)
	}

	if !image.IsProcessing() {
		t.Error("expected image to remain processing")
	}

	if !image.UpdateTime().IsZero() {
		t.Errorf("expected update time to remain zero, got %v", image.UpdateTime())
	}

	if events := image.PullEvents(); len(events) != 0 {
		t.Errorf("expected no events, got %d", len(events))
	}
}

func TestImage_CanConfirmUpload_WhenProcessed(t *testing.T) {
	image := ReconstituteImage(
		vo.ImageID{},
		vo.Tags{},
		vo.ObjectKey{},
		vo.ObjectKey{},
		vo.ImageStateProcessed,
		time.Time{},
		time.Time{},
	)

	err := image.CanConfirmUpload()

	if !errors.Is(err, ErrImageProcessed) {
		t.Fatalf("expected ErrImageProcessed, got %v", err)
	}

	if !image.IsProcessed() {
		t.Error("expected image to remain processed")
	}

	if !image.UpdateTime().IsZero() {
		t.Errorf("expected update time to remain zero, got %v", image.UpdateTime())
	}

	if events := image.PullEvents(); len(events) != 0 {
		t.Errorf("expected no events, got %d", len(events))
	}
}

func TestImage_CanConfirmUpload_WhenDeleted(t *testing.T) {
	image := ReconstituteImage(
		vo.ImageID{},
		vo.Tags{},
		vo.ObjectKey{},
		vo.ObjectKey{},
		vo.ImageStateDeleted,
		time.Time{},
		time.Time{},
	)

	err := image.CanConfirmUpload()

	if !errors.Is(err, ErrImageDeleted) {
		t.Fatalf("expected ErrImageDeleted, got %v", err)
	}

	if !image.IsDeleted() {
		t.Error("expected image to remain deleted")
	}

	if !image.UpdateTime().IsZero() {
		t.Errorf("expected update time to remain zero, got %v", image.UpdateTime())
	}

	if events := image.PullEvents(); len(events) != 0 {
		t.Errorf("expected no events, got %d", len(events))
	}
}

func TestImage_ConfirmUpload_WhenCreated(t *testing.T) {
	image := ReconstituteImage(
		vo.ImageID{},
		vo.Tags{},
		vo.ObjectKey{},
		vo.ObjectKey{},
		vo.ImageStateCreated,
		time.Time{},
		time.Time{},
	)

	now := time.Now().UTC()

	err := image.ConfirmUpload(now)

	if !errors.Is(err, ErrImageNotUploading) {
		t.Fatalf("expected ErrImageNotUploading, got %v", err)
	}

	if !image.IsCreated() {
		t.Error("expected image to remain created")
	}

	if !image.UpdateTime().IsZero() {
		t.Errorf("expected update time to remain zero, got %v", image.UpdateTime())
	}

	if events := image.PullEvents(); len(events) != 0 {
		t.Errorf("expected no events, got %d", len(events))
	}
}

func TestImage_ConfirmUpload(t *testing.T) {
	image := ReconstituteImage(
		vo.ImageID{},
		vo.Tags{},
		vo.ObjectKey{},
		vo.ObjectKey{},
		vo.ImageStateUploading,
		time.Time{},
		time.Time{},
	)

	now := time.Now().UTC()

	err := image.ConfirmUpload(now)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if !image.IsUploaded() {
		t.Error("expected image to be uploaded")
	}

	if !image.UpdateTime().Equal(now) {
		t.Errorf("expected update time %v, got %v", now, image.UpdateTime())
	}

	events := image.PullEvents()

	if len(events) != 1 {
		t.Fatalf("expected 1 event, got %d", len(events))
	}

	if _, ok := events[0].(event.ImageUploaded); !ok {
		t.Errorf(
			"expected ImageUploaded, got %T",
			events[0],
		)
	}
}

func TestImage_ConfirmUpload_WhenAlreadyUploaded(t *testing.T) {
	image := ReconstituteImage(
		vo.ImageID{},
		vo.Tags{},
		vo.ObjectKey{},
		vo.ObjectKey{},
		vo.ImageStateUploaded,
		time.Time{},
		time.Time{},
	)

	now := time.Now().UTC()

	err := image.ConfirmUpload(now)

	if !errors.Is(err, ErrImageAlreadyUploaded) {
		t.Fatalf("expected ErrImageAlreadyUploaded, got %v", err)
	}

	if !image.IsUploaded() {
		t.Error("expected image to remain uploaded")
	}

	if !image.UpdateTime().IsZero() {
		t.Errorf("expected update time to remain zero, got %v", image.UpdateTime())
	}

	if events := image.PullEvents(); len(events) != 0 {
		t.Errorf("expected no events, got %d", len(events))
	}
}

func TestImage_ConfirmUpload_WhenProcessing(t *testing.T) {
	image := ReconstituteImage(
		vo.ImageID{},
		vo.Tags{},
		vo.ObjectKey{},
		vo.ObjectKey{},
		vo.ImageStateProcessing,
		time.Time{},
		time.Time{},
	)

	now := time.Now().UTC()

	err := image.ConfirmUpload(now)

	if !errors.Is(err, ErrImageProcessing) {
		t.Fatalf("expected ErrImageProcessing, got %v", err)
	}

	if !image.IsProcessing() {
		t.Error("expected image to remain processing")
	}

	if !image.UpdateTime().IsZero() {
		t.Errorf("expected update time to remain zero, got %v", image.UpdateTime())
	}

	if events := image.PullEvents(); len(events) != 0 {
		t.Errorf("expected no events, got %d", len(events))
	}
}

func TestImage_ConfirmUpload_WhenProcessed(t *testing.T) {
	image := ReconstituteImage(
		vo.ImageID{},
		vo.Tags{},
		vo.ObjectKey{},
		vo.ObjectKey{},
		vo.ImageStateProcessed,
		time.Time{},
		time.Time{},
	)

	now := time.Now().UTC()

	err := image.ConfirmUpload(now)

	if !errors.Is(err, ErrImageProcessed) {
		t.Fatalf("expected ErrImageProcessed, got %v", err)
	}

	if !image.IsProcessed() {
		t.Error("expected image to remain processed")
	}

	if !image.UpdateTime().IsZero() {
		t.Errorf("expected update time to remain zero, got %v", image.UpdateTime())
	}

	if events := image.PullEvents(); len(events) != 0 {
		t.Errorf("expected no events, got %d", len(events))
	}
}

func TestImage_ConfirmUpload_WhenDeleted(t *testing.T) {
	image := ReconstituteImage(
		vo.ImageID{},
		vo.Tags{},
		vo.ObjectKey{},
		vo.ObjectKey{},
		vo.ImageStateDeleted,
		time.Time{},
		time.Time{},
	)

	now := time.Now().UTC()

	err := image.ConfirmUpload(now)

	if !errors.Is(err, ErrImageDeleted) {
		t.Fatalf("expected ErrImageDeleted, got %v", err)
	}

	if !image.IsDeleted() {
		t.Error("expected image to remain deleted")
	}

	if !image.UpdateTime().IsZero() {
		t.Errorf("expected update time to remain zero, got %v", image.UpdateTime())
	}

	if events := image.PullEvents(); len(events) != 0 {
		t.Errorf("expected no events, got %d", len(events))
	}
}

func TestImage_CanStartProcessing_WhenCreated(t *testing.T) {
	image := ReconstituteImage(
		vo.ImageID{},
		vo.Tags{},
		vo.ObjectKey{},
		vo.ObjectKey{},
		vo.ImageStateCreated,
		time.Time{},
		time.Time{},
	)

	err := image.CanStartProcessing()

	if !errors.Is(err, ErrImageNotUploaded) {
		t.Fatalf("expected ErrImageNotUploaded, got %v", err)
	}

	if !image.IsCreated() {
		t.Error("expected image to remain created")
	}

	if !image.UpdateTime().IsZero() {
		t.Errorf("expected update time to remain zero, got %v", image.UpdateTime())
	}

	if events := image.PullEvents(); len(events) != 0 {
		t.Errorf("expected no events, got %d", len(events))
	}
}

func TestImage_CanStartProcessing_WhenUploading(t *testing.T) {
	image := ReconstituteImage(
		vo.ImageID{},
		vo.Tags{},
		vo.ObjectKey{},
		vo.ObjectKey{},
		vo.ImageStateUploading,
		time.Time{},
		time.Time{},
	)

	err := image.CanStartProcessing()

	if !errors.Is(err, ErrImageNotUploaded) {
		t.Fatalf("expected ErrImageNotUploaded, got %v", err)
	}

	if !image.IsUploading() {
		t.Error("expected image to remain uploading")
	}

	if !image.UpdateTime().IsZero() {
		t.Errorf("expected update time to remain zero, got %v", image.UpdateTime())
	}

	if events := image.PullEvents(); len(events) != 0 {
		t.Errorf("expected no events, got %d", len(events))
	}
}

func TestImage_CanStartProcessing(t *testing.T) {
	image := ReconstituteImage(
		vo.ImageID{},
		vo.Tags{},
		vo.ObjectKey{},
		vo.ObjectKey{},
		vo.ImageStateUploaded,
		time.Time{},
		time.Time{},
	)

	err := image.CanStartProcessing()

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if !image.IsUploaded() {
		t.Error("expected image to remain uploaded")
	}

	if !image.UpdateTime().IsZero() {
		t.Errorf("expected update time to remain zero, got %v", image.UpdateTime())
	}

	if events := image.PullEvents(); len(events) != 0 {
		t.Errorf("expected no events, got %d", len(events))
	}
}

func TestImage_CanStartProcessing_WhenAlreadyProcessing(t *testing.T) {
	image := ReconstituteImage(
		vo.ImageID{},
		vo.Tags{},
		vo.ObjectKey{},
		vo.ObjectKey{},
		vo.ImageStateProcessing,
		time.Time{},
		time.Time{},
	)

	err := image.CanStartProcessing()

	if !errors.Is(err, ErrImageAlreadyProcessing) {
		t.Fatalf("expected ErrImageAlreadyProcessing, got %v", err)
	}

	if !image.IsProcessing() {
		t.Error("expected image to remain processing")
	}

	if !image.UpdateTime().IsZero() {
		t.Errorf("expected update time to remain zero, got %v", image.UpdateTime())
	}

	if events := image.PullEvents(); len(events) != 0 {
		t.Errorf("expected no events, got %d", len(events))
	}
}

func TestImage_CanStartProcessing_WhenProcessed(t *testing.T) {
	image := ReconstituteImage(
		vo.ImageID{},
		vo.Tags{},
		vo.ObjectKey{},
		vo.ObjectKey{},
		vo.ImageStateProcessed,
		time.Time{},
		time.Time{},
	)

	err := image.CanStartProcessing()

	if !errors.Is(err, ErrImageProcessed) {
		t.Fatalf("expected ErrImageProcessed, got %v", err)
	}

	if !image.IsProcessed() {
		t.Error("expected image to remain processed")
	}

	if !image.UpdateTime().IsZero() {
		t.Errorf("expected update time to remain zero, got %v", image.UpdateTime())
	}

	if events := image.PullEvents(); len(events) != 0 {
		t.Errorf("expected no events, got %d", len(events))
	}
}

func TestImage_CanStartProcessing_WhenDeleted(t *testing.T) {
	image := ReconstituteImage(
		vo.ImageID{},
		vo.Tags{},
		vo.ObjectKey{},
		vo.ObjectKey{},
		vo.ImageStateDeleted,
		time.Time{},
		time.Time{},
	)

	err := image.CanStartProcessing()

	if !errors.Is(err, ErrImageDeleted) {
		t.Fatalf("expected ErrImageDeleted, got %v", err)
	}

	if !image.IsDeleted() {
		t.Error("expected image to remain deleted")
	}

	if !image.UpdateTime().IsZero() {
		t.Errorf("expected update time to remain zero, got %v", image.UpdateTime())
	}

	if events := image.PullEvents(); len(events) != 0 {
		t.Errorf("expected no events, got %d", len(events))
	}
}

func TestImage_StartProcessing_WhenCreated(t *testing.T) {
	image := ReconstituteImage(
		vo.ImageID{},
		vo.Tags{},
		vo.ObjectKey{},
		vo.ObjectKey{},
		vo.ImageStateCreated,
		time.Time{},
		time.Time{},
	)

	now := time.Now().UTC()

	err := image.StartProcessing(now)

	if !errors.Is(err, ErrImageNotUploaded) {
		t.Fatalf("expected ErrImageNotUploaded, got %v", err)
	}

	if !image.IsCreated() {
		t.Error("expected image to remain created")
	}

	if !image.UpdateTime().IsZero() {
		t.Errorf("expected update time to remain zero, got %v", image.UpdateTime())
	}

	if events := image.PullEvents(); len(events) != 0 {
		t.Errorf("expected no events, got %d", len(events))
	}
}

func TestImage_StartProcessing_WhenUploading(t *testing.T) {
	image := ReconstituteImage(
		vo.ImageID{},
		vo.Tags{},
		vo.ObjectKey{},
		vo.ObjectKey{},
		vo.ImageStateUploading,
		time.Time{},
		time.Time{},
	)

	now := time.Now().UTC()

	err := image.StartProcessing(now)

	if !errors.Is(err, ErrImageNotUploaded) {
		t.Fatalf("expected ErrImageNotUploaded, got %v", err)
	}

	if !image.IsUploading() {
		t.Error("expected image to remain uploading")
	}

	if !image.UpdateTime().IsZero() {
		t.Errorf("expected update time to remain zero, got %v", image.UpdateTime())
	}

	if events := image.PullEvents(); len(events) != 0 {
		t.Errorf("expected no events, got %d", len(events))
	}
}

func TestImage_StartProcessing(t *testing.T) {
	image := ReconstituteImage(
		vo.ImageID{},
		vo.Tags{},
		vo.ObjectKey{},
		vo.ObjectKey{},
		vo.ImageStateUploaded,
		time.Time{},
		time.Time{},
	)

	now := time.Now().UTC()

	err := image.StartProcessing(now)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if !image.IsProcessing() {
		t.Error("expected image to be processing")
	}

	if !image.UpdateTime().Equal(now) {
		t.Errorf("expected update time %v, got %v", now, image.UpdateTime())
	}

	if events := image.PullEvents(); len(events) != 0 {
		t.Errorf("expected no events, got %d", len(events))
	}
}

func TestImage_StartProcessing_WhenAlreadyProcessing(t *testing.T) {
	image := ReconstituteImage(
		vo.ImageID{},
		vo.Tags{},
		vo.ObjectKey{},
		vo.ObjectKey{},
		vo.ImageStateProcessing,
		time.Time{},
		time.Time{},
	)

	now := time.Now().UTC()

	err := image.StartProcessing(now)

	if !errors.Is(err, ErrImageAlreadyProcessing) {
		t.Fatalf("expected ErrImageAlreadyProcessing, got %v", err)
	}

	if !image.IsProcessing() {
		t.Error("expected image to remain processing")
	}

	if !image.UpdateTime().IsZero() {
		t.Errorf("expected update time to remain zero, got %v", image.UpdateTime())
	}

	if events := image.PullEvents(); len(events) != 0 {
		t.Errorf("expected no events, got %d", len(events))
	}
}

func TestImage_StartProcessing_WhenProcessed(t *testing.T) {
	image := ReconstituteImage(
		vo.ImageID{},
		vo.Tags{},
		vo.ObjectKey{},
		vo.ObjectKey{},
		vo.ImageStateProcessed,
		time.Time{},
		time.Time{},
	)

	now := time.Now().UTC()

	err := image.StartProcessing(now)

	if !errors.Is(err, ErrImageProcessed) {
		t.Fatalf("expected ErrImageProcessed, got %v", err)
	}

	if !image.IsProcessed() {
		t.Error("expected image to remain processed")
	}

	if !image.UpdateTime().IsZero() {
		t.Errorf("expected update time to remain zero, got %v", image.UpdateTime())
	}

	if events := image.PullEvents(); len(events) != 0 {
		t.Errorf("expected no events, got %d", len(events))
	}
}

func TestImage_StartProcessing_WhenDeleted(t *testing.T) {
	image := ReconstituteImage(
		vo.ImageID{},
		vo.Tags{},
		vo.ObjectKey{},
		vo.ObjectKey{},
		vo.ImageStateDeleted,
		time.Time{},
		time.Time{},
	)

	now := time.Now().UTC()

	err := image.StartProcessing(now)

	if !errors.Is(err, ErrImageDeleted) {
		t.Fatalf("expected ErrImageDeleted, got %v", err)
	}

	if !image.IsDeleted() {
		t.Error("expected image to remain deleted")
	}

	if !image.UpdateTime().IsZero() {
		t.Errorf("expected update time to remain zero, got %v", image.UpdateTime())
	}

	if events := image.PullEvents(); len(events) != 0 {
		t.Errorf("expected no events, got %d", len(events))
	}
}

func TestImage_CanCancelProcessing_WhenCreated(t *testing.T) {
	image := ReconstituteImage(
		vo.ImageID{},
		vo.Tags{},
		vo.ObjectKey{},
		vo.ObjectKey{},
		vo.ImageStateCreated,
		time.Time{},
		time.Time{},
	)

	err := image.CanCancelProcessing()

	if !errors.Is(err, ErrImageNotProcessing) {
		t.Fatalf("expected ErrImageNotProcessing, got %v", err)
	}

	if !image.IsCreated() {
		t.Error("expected image to remain created")
	}

	if !image.UpdateTime().IsZero() {
		t.Errorf("expected update time to remain zero, got %v", image.UpdateTime())
	}

	if events := image.PullEvents(); len(events) != 0 {
		t.Errorf("expected no events, got %d", len(events))
	}
}

func TestImage_CanCancelProcessing_WhenUploading(t *testing.T) {
	image := ReconstituteImage(
		vo.ImageID{},
		vo.Tags{},
		vo.ObjectKey{},
		vo.ObjectKey{},
		vo.ImageStateUploading,
		time.Time{},
		time.Time{},
	)

	err := image.CanCancelProcessing()

	if !errors.Is(err, ErrImageNotProcessing) {
		t.Fatalf("expected ErrImageNotProcessing, got %v", err)
	}

	if !image.IsUploading() {
		t.Error("expected image to remain uploading")
	}

	if !image.UpdateTime().IsZero() {
		t.Errorf("expected update time to remain zero, got %v", image.UpdateTime())
	}

	if events := image.PullEvents(); len(events) != 0 {
		t.Errorf("expected no events, got %d", len(events))
	}
}

func TestImage_CanCancelProcessing_WhenUploaded(t *testing.T) {
	image := ReconstituteImage(
		vo.ImageID{},
		vo.Tags{},
		vo.ObjectKey{},
		vo.ObjectKey{},
		vo.ImageStateUploaded,
		time.Time{},
		time.Time{},
	)

	err := image.CanCancelProcessing()

	if !errors.Is(err, ErrImageNotProcessing) {
		t.Fatalf("expected ErrImageNotProcessing, got %v", err)
	}

	if !image.IsUploaded() {
		t.Error("expected image to remain uploaded")
	}

	if !image.UpdateTime().IsZero() {
		t.Errorf("expected update time to remain zero, got %v", image.UpdateTime())
	}

	if events := image.PullEvents(); len(events) != 0 {
		t.Errorf("expected no events, got %d", len(events))
	}
}

func TestImage_CanCancelProcessing(t *testing.T) {
	image := ReconstituteImage(
		vo.ImageID{},
		vo.Tags{},
		vo.ObjectKey{},
		vo.ObjectKey{},
		vo.ImageStateProcessing,
		time.Time{},
		time.Time{},
	)

	err := image.CanCancelProcessing()

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if !image.IsProcessing() {
		t.Error("expected image to remain processing")
	}

	if !image.UpdateTime().IsZero() {
		t.Errorf("expected update time to remain zero, got %v", image.UpdateTime())
	}

	if events := image.PullEvents(); len(events) != 0 {
		t.Errorf("expected no events, got %d", len(events))
	}
}

func TestImage_CanCancelProcessing_WhenProcessed(t *testing.T) {
	image := ReconstituteImage(
		vo.ImageID{},
		vo.Tags{},
		vo.ObjectKey{},
		vo.ObjectKey{},
		vo.ImageStateProcessed,
		time.Time{},
		time.Time{},
	)

	err := image.CanCancelProcessing()

	if !errors.Is(err, ErrImageNotProcessing) {
		t.Fatalf("expected ErrImageNotProcessing, got %v", err)
	}

	if !image.IsProcessed() {
		t.Error("expected image to remain processed")
	}

	if !image.UpdateTime().IsZero() {
		t.Errorf("expected update time to remain zero, got %v", image.UpdateTime())
	}

	if events := image.PullEvents(); len(events) != 0 {
		t.Errorf("expected no events, got %d", len(events))
	}
}

func TestImage_CanCancelProcessing_WhenDeleted(t *testing.T) {
	image := ReconstituteImage(
		vo.ImageID{},
		vo.Tags{},
		vo.ObjectKey{},
		vo.ObjectKey{},
		vo.ImageStateDeleted,
		time.Time{},
		time.Time{},
	)

	err := image.CanCancelProcessing()

	if !errors.Is(err, ErrImageDeleted) {
		t.Fatalf("expected ErrImageDeleted, got %v", err)
	}

	if !image.IsDeleted() {
		t.Error("expected image to remain deleted")
	}

	if !image.UpdateTime().IsZero() {
		t.Errorf("expected update time to remain zero, got %v", image.UpdateTime())
	}

	if events := image.PullEvents(); len(events) != 0 {
		t.Errorf("expected no events, got %d", len(events))
	}
}

func TestImage_CancelProcessing_WhenCreated(t *testing.T) {
	image := ReconstituteImage(
		vo.ImageID{},
		vo.Tags{},
		vo.ObjectKey{},
		vo.ObjectKey{},
		vo.ImageStateCreated,
		time.Time{},
		time.Time{},
	)

	now := time.Now().UTC()

	err := image.CancelProcessing(now)

	if !errors.Is(err, ErrImageNotProcessing) {
		t.Fatalf("expected ErrImageNotProcessing, got %v", err)
	}

	if !image.IsCreated() {
		t.Error("expected image to remain created")
	}

	if !image.UpdateTime().IsZero() {
		t.Errorf("expected update time to remain zero, got %v", image.UpdateTime())
	}

	if events := image.PullEvents(); len(events) != 0 {
		t.Errorf("expected no events, got %d", len(events))
	}
}

func TestImage_CancelProcessing_WhenUploading(t *testing.T) {
	image := ReconstituteImage(
		vo.ImageID{},
		vo.Tags{},
		vo.ObjectKey{},
		vo.ObjectKey{},
		vo.ImageStateUploading,
		time.Time{},
		time.Time{},
	)

	now := time.Now().UTC()

	err := image.CancelProcessing(now)

	if !errors.Is(err, ErrImageNotProcessing) {
		t.Fatalf("expected ErrImageNotProcessing, got %v", err)
	}

	if !image.IsUploading() {
		t.Error("expected image to remain uploading")
	}

	if !image.UpdateTime().IsZero() {
		t.Errorf("expected update time to remain zero, got %v", image.UpdateTime())
	}

	if events := image.PullEvents(); len(events) != 0 {
		t.Errorf("expected no events, got %d", len(events))
	}
}

func TestImage_CancelProcessing_WhenUploaded(t *testing.T) {
	image := ReconstituteImage(
		vo.ImageID{},
		vo.Tags{},
		vo.ObjectKey{},
		vo.ObjectKey{},
		vo.ImageStateUploaded,
		time.Time{},
		time.Time{},
	)

	now := time.Now().UTC()

	err := image.CancelProcessing(now)

	if !errors.Is(err, ErrImageNotProcessing) {
		t.Fatalf("expected ErrImageNotProcessing, got %v", err)
	}

	if !image.IsUploaded() {
		t.Error("expected image to remain uploaded")
	}

	if !image.UpdateTime().IsZero() {
		t.Errorf("expected update time to remain zero, got %v", image.UpdateTime())
	}

	if events := image.PullEvents(); len(events) != 0 {
		t.Errorf("expected no events, got %d", len(events))
	}
}

func TestImage_CancelProcessing(t *testing.T) {
	image := ReconstituteImage(
		vo.ImageID{},
		vo.Tags{},
		vo.ObjectKey{},
		vo.ObjectKey{},
		vo.ImageStateProcessing,
		time.Time{},
		time.Time{},
	)

	now := time.Now().UTC()

	err := image.CancelProcessing(now)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if !image.IsUploaded() {
		t.Error("expected image to be uploaded")
	}

	if !image.UpdateTime().Equal(now) {
		t.Errorf("expected update time %v, got %v", now, image.UpdateTime())
	}

	if events := image.PullEvents(); len(events) != 0 {
		t.Errorf("expected no events, got %d", len(events))
	}
}

func TestImage_CancelProcessing_WhenProcessed(t *testing.T) {
	image := ReconstituteImage(
		vo.ImageID{},
		vo.Tags{},
		vo.ObjectKey{},
		vo.ObjectKey{},
		vo.ImageStateProcessed,
		time.Time{},
		time.Time{},
	)

	now := time.Now().UTC()

	err := image.CancelProcessing(now)

	if !errors.Is(err, ErrImageNotProcessing) {
		t.Fatalf("expected ErrImageNotProcessing, got %v", err)
	}

	if !image.IsProcessed() {
		t.Error("expected image to remain processed")
	}

	if !image.UpdateTime().IsZero() {
		t.Errorf("expected update time to remain zero, got %v", image.UpdateTime())
	}

	if events := image.PullEvents(); len(events) != 0 {
		t.Errorf("expected no events, got %d", len(events))
	}
}

func TestImage_CancelProcessing_WhenDeleted(t *testing.T) {
	image := ReconstituteImage(
		vo.ImageID{},
		vo.Tags{},
		vo.ObjectKey{},
		vo.ObjectKey{},
		vo.ImageStateDeleted,
		time.Time{},
		time.Time{},
	)

	now := time.Now().UTC()

	err := image.CancelProcessing(now)

	if !errors.Is(err, ErrImageDeleted) {
		t.Fatalf("expected ErrImageDeleted, got %v", err)
	}

	if !image.IsDeleted() {
		t.Error("expected image to remain deleted")
	}

	if !image.UpdateTime().IsZero() {
		t.Errorf("expected update time to remain zero, got %v", image.UpdateTime())
	}

	if events := image.PullEvents(); len(events) != 0 {
		t.Errorf("expected no events, got %d", len(events))
	}
}

func TestImage_CanCompleteProcessing_WhenCreated(t *testing.T) {
	image := ReconstituteImage(
		vo.ImageID{},
		vo.Tags{},
		vo.ObjectKey{},
		vo.ObjectKey{},
		vo.ImageStateCreated,
		time.Time{},
		time.Time{},
	)

	err := image.CanCompleteProcessing()

	if !errors.Is(err, ErrImageNotProcessing) {
		t.Fatalf("expected ErrImageNotProcessing, got %v", err)
	}

	if !image.IsCreated() {
		t.Error("expected image to remain created")
	}

	if !image.UpdateTime().IsZero() {
		t.Errorf("expected update time to remain zero, got %v", image.UpdateTime())
	}

	if events := image.PullEvents(); len(events) != 0 {
		t.Errorf("expected no events, got %d", len(events))
	}
}

func TestImage_CanCompleteProcessing_WhenUploading(t *testing.T) {
	image := ReconstituteImage(
		vo.ImageID{},
		vo.Tags{},
		vo.ObjectKey{},
		vo.ObjectKey{},
		vo.ImageStateUploading,
		time.Time{},
		time.Time{},
	)

	err := image.CanCompleteProcessing()

	if !errors.Is(err, ErrImageNotProcessing) {
		t.Fatalf("expected ErrImageNotProcessing, got %v", err)
	}

	if !image.IsUploading() {
		t.Error("expected image to remain uploading")
	}

	if !image.UpdateTime().IsZero() {
		t.Errorf("expected update time to remain zero, got %v", image.UpdateTime())
	}

	if events := image.PullEvents(); len(events) != 0 {
		t.Errorf("expected no events, got %d", len(events))
	}
}

func TestImage_CanCompleteProcessing_WhenUploaded(t *testing.T) {
	image := ReconstituteImage(
		vo.ImageID{},
		vo.Tags{},
		vo.ObjectKey{},
		vo.ObjectKey{},
		vo.ImageStateUploaded,
		time.Time{},
		time.Time{},
	)

	err := image.CanCompleteProcessing()

	if !errors.Is(err, ErrImageNotProcessing) {
		t.Fatalf("expected ErrImageNotProcessing, got %v", err)
	}

	if !image.IsUploaded() {
		t.Error("expected image to remain uploaded")
	}

	if !image.UpdateTime().IsZero() {
		t.Errorf("expected update time to remain zero, got %v", image.UpdateTime())
	}

	if events := image.PullEvents(); len(events) != 0 {
		t.Errorf("expected no events, got %d", len(events))
	}
}

func TestImage_CanCompleteProcessing(t *testing.T) {
	image := ReconstituteImage(
		vo.ImageID{},
		vo.Tags{},
		vo.ObjectKey{},
		vo.ObjectKey{},
		vo.ImageStateProcessing,
		time.Time{},
		time.Time{},
	)

	err := image.CanCompleteProcessing()

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if !image.IsProcessing() {
		t.Error("expected image to remain processing")
	}

	if !image.UpdateTime().IsZero() {
		t.Errorf("expected update time to remain zero, got %v", image.UpdateTime())
	}

	if events := image.PullEvents(); len(events) != 0 {
		t.Errorf("expected no events, got %d", len(events))
	}
}

func TestImage_CanCompleteProcessing_WhenAlreadyProcessed(t *testing.T) {
	image := ReconstituteImage(
		vo.ImageID{},
		vo.Tags{},
		vo.ObjectKey{},
		vo.ObjectKey{},
		vo.ImageStateProcessed,
		time.Time{},
		time.Time{},
	)

	err := image.CanCompleteProcessing()

	if !errors.Is(err, ErrImageAlreadyProcessed) {
		t.Fatalf("expected ErrImageAlreadyProcessed, got %v", err)
	}

	if !image.IsProcessed() {
		t.Error("expected image to remain processed")
	}

	if !image.UpdateTime().IsZero() {
		t.Errorf("expected update time to remain zero, got %v", image.UpdateTime())
	}

	if events := image.PullEvents(); len(events) != 0 {
		t.Errorf("expected no events, got %d", len(events))
	}
}

func TestImage_CanCompleteProcessing_WhenDeleted(t *testing.T) {
	image := ReconstituteImage(
		vo.ImageID{},
		vo.Tags{},
		vo.ObjectKey{},
		vo.ObjectKey{},
		vo.ImageStateDeleted,
		time.Time{},
		time.Time{},
	)

	err := image.CanCompleteProcessing()

	if !errors.Is(err, ErrImageDeleted) {
		t.Fatalf("expected ErrImageDeleted, got %v", err)
	}

	if !image.IsDeleted() {
		t.Error("expected image to remain deleted")
	}

	if !image.UpdateTime().IsZero() {
		t.Errorf("expected update time to remain zero, got %v", image.UpdateTime())
	}

	if events := image.PullEvents(); len(events) != 0 {
		t.Errorf("expected no events, got %d", len(events))
	}
}

func TestImage_CompleteProcessing_WhenCreated(t *testing.T) {
	image := ReconstituteImage(
		vo.ImageID{},
		vo.Tags{},
		vo.ObjectKey{},
		vo.ObjectKey{},
		vo.ImageStateCreated,
		time.Time{},
		time.Time{},
	)

	now := time.Now().UTC()

	err := image.CompleteProcessing(now)

	if !errors.Is(err, ErrImageNotProcessing) {
		t.Fatalf("expected ErrImageNotProcessing, got %v", err)
	}

	if !image.IsCreated() {
		t.Error("expected image to remain created")
	}

	if !image.UpdateTime().IsZero() {
		t.Errorf("expected update time to remain zero, got %v", image.UpdateTime())
	}

	if events := image.PullEvents(); len(events) != 0 {
		t.Errorf("expected no events, got %d", len(events))
	}
}

func TestImage_CompleteProcessing_WhenUploading(t *testing.T) {
	image := ReconstituteImage(
		vo.ImageID{},
		vo.Tags{},
		vo.ObjectKey{},
		vo.ObjectKey{},
		vo.ImageStateUploading,
		time.Time{},
		time.Time{},
	)

	now := time.Now().UTC()

	err := image.CompleteProcessing(now)

	if !errors.Is(err, ErrImageNotProcessing) {
		t.Fatalf("expected ErrImageNotProcessing, got %v", err)
	}

	if !image.IsUploading() {
		t.Error("expected image to remain uploading")
	}

	if !image.UpdateTime().IsZero() {
		t.Errorf("expected update time to remain zero, got %v", image.UpdateTime())
	}

	if events := image.PullEvents(); len(events) != 0 {
		t.Errorf("expected no events, got %d", len(events))
	}
}

func TestImage_CompleteProcessing_WhenUploaded(t *testing.T) {
	image := ReconstituteImage(
		vo.ImageID{},
		vo.Tags{},
		vo.ObjectKey{},
		vo.ObjectKey{},
		vo.ImageStateUploaded,
		time.Time{},
		time.Time{},
	)

	now := time.Now().UTC()

	err := image.CompleteProcessing(now)

	if !errors.Is(err, ErrImageNotProcessing) {
		t.Fatalf("expected ErrImageNotProcessing, got %v", err)
	}

	if !image.IsUploaded() {
		t.Error("expected image to remain uploaded")
	}

	if !image.UpdateTime().IsZero() {
		t.Errorf("expected update time to remain zero, got %v", image.UpdateTime())
	}

	if events := image.PullEvents(); len(events) != 0 {
		t.Errorf("expected no events, got %d", len(events))
	}
}

func TestImage_CompleteProcessing(t *testing.T) {
	image := ReconstituteImage(
		vo.ImageID{},
		vo.Tags{},
		vo.ObjectKey{},
		vo.ObjectKey{},
		vo.ImageStateProcessing,
		time.Time{},
		time.Time{},
	)

	now := time.Now().UTC()

	err := image.CompleteProcessing(now)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if !image.IsProcessed() {
		t.Error("expected image to be processed")
	}

	if !image.UpdateTime().Equal(now) {
		t.Errorf("expected update time %v, got %v", now, image.UpdateTime())
	}

	events := image.PullEvents()

	if len(events) != 1 {
		t.Fatalf("expected 1 event, got %d", len(events))
	}

	if _, ok := events[0].(event.ImageProcessed); !ok {
		t.Errorf(
			"expected ImageProcessed, got %T",
			events[0],
		)
	}
}

func TestImage_CompleteProcessing_WhenAlreadyProcessed(t *testing.T) {
	image := ReconstituteImage(
		vo.ImageID{},
		vo.Tags{},
		vo.ObjectKey{},
		vo.ObjectKey{},
		vo.ImageStateProcessed,
		time.Time{},
		time.Time{},
	)

	now := time.Now().UTC()

	err := image.CompleteProcessing(now)

	if !errors.Is(err, ErrImageAlreadyProcessed) {
		t.Fatalf("expected ErrImageAlreadyProcessed, got %v", err)
	}

	if !image.IsProcessed() {
		t.Error("expected image to remain processed")
	}

	if !image.UpdateTime().IsZero() {
		t.Errorf("expected update time to remain zero, got %v", image.UpdateTime())
	}

	if events := image.PullEvents(); len(events) != 0 {
		t.Errorf("expected no events, got %d", len(events))
	}
}

func TestImage_CompleteProcessing_WhenDeleted(t *testing.T) {
	image := ReconstituteImage(
		vo.ImageID{},
		vo.Tags{},
		vo.ObjectKey{},
		vo.ObjectKey{},
		vo.ImageStateDeleted,
		time.Time{},
		time.Time{},
	)

	now := time.Now().UTC()

	err := image.CompleteProcessing(now)

	if !errors.Is(err, ErrImageDeleted) {
		t.Fatalf("expected ErrImageDeleted, got %v", err)
	}

	if !image.IsDeleted() {
		t.Error("expected image to remain deleted")
	}

	if !image.UpdateTime().IsZero() {
		t.Errorf("expected update time to remain zero, got %v", image.UpdateTime())
	}

	if events := image.PullEvents(); len(events) != 0 {
		t.Errorf("expected no events, got %d", len(events))
	}
}

func TestImage_CanDelete(t *testing.T) {
	image := ReconstituteImage(
		vo.ImageID{},
		vo.Tags{},
		vo.ObjectKey{},
		vo.ObjectKey{},
		vo.ImageStateCreated,
		time.Time{},
		time.Time{},
	)

	err := image.CanDelete()

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if !image.IsCreated() {
		t.Error("expected image to remain created")
	}

	if !image.UpdateTime().IsZero() {
		t.Errorf("expected update time to remain zero, got %v", image.UpdateTime())
	}

	if events := image.PullEvents(); len(events) != 0 {
		t.Errorf("expected no events, got %d", len(events))
	}
}

func TestImage_CanDelete_WhenUploading(t *testing.T) {
	image := ReconstituteImage(
		vo.ImageID{},
		vo.Tags{},
		vo.ObjectKey{},
		vo.ObjectKey{},
		vo.ImageStateUploading,
		time.Time{},
		time.Time{},
	)

	err := image.CanDelete()

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if !image.IsUploading() {
		t.Error("expected image to remain uploading")
	}

	if !image.UpdateTime().IsZero() {
		t.Errorf("expected update time to remain zero, got %v", image.UpdateTime())
	}

	if events := image.PullEvents(); len(events) != 0 {
		t.Errorf("expected no events, got %d", len(events))
	}
}

func TestImage_CanDelete_WhenUploaded(t *testing.T) {
	image := ReconstituteImage(
		vo.ImageID{},
		vo.Tags{},
		vo.ObjectKey{},
		vo.ObjectKey{},
		vo.ImageStateUploaded,
		time.Time{},
		time.Time{},
	)

	err := image.CanDelete()

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if !image.IsUploaded() {
		t.Error("expected image to remain uploaded")
	}

	if !image.UpdateTime().IsZero() {
		t.Errorf("expected update time to remain zero, got %v", image.UpdateTime())
	}

	if events := image.PullEvents(); len(events) != 0 {
		t.Errorf("expected no events, got %d", len(events))
	}
}

func TestImage_CanDelete_WhenProcessing(t *testing.T) {
	image := ReconstituteImage(
		vo.ImageID{},
		vo.Tags{},
		vo.ObjectKey{},
		vo.ObjectKey{},
		vo.ImageStateProcessing,
		time.Time{},
		time.Time{},
	)

	err := image.CanDelete()

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if !image.IsProcessing() {
		t.Error("expected image to remain processing")
	}

	if !image.UpdateTime().IsZero() {
		t.Errorf("expected update time to remain zero, got %v", image.UpdateTime())
	}

	if events := image.PullEvents(); len(events) != 0 {
		t.Errorf("expected no events, got %d", len(events))
	}
}

func TestImage_CanDelete_WhenProcessed(t *testing.T) {
	image := ReconstituteImage(
		vo.ImageID{},
		vo.Tags{},
		vo.ObjectKey{},
		vo.ObjectKey{},
		vo.ImageStateProcessed,
		time.Time{},
		time.Time{},
	)

	err := image.CanDelete()

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if !image.IsProcessed() {
		t.Error("expected image to remain processed")
	}

	if !image.UpdateTime().IsZero() {
		t.Errorf("expected update time to remain zero, got %v", image.UpdateTime())
	}

	if events := image.PullEvents(); len(events) != 0 {
		t.Errorf("expected no events, got %d", len(events))
	}
}

func TestImage_CanDelete_WhenAlreadyDeleted(t *testing.T) {
	image := ReconstituteImage(
		vo.ImageID{},
		vo.Tags{},
		vo.ObjectKey{},
		vo.ObjectKey{},
		vo.ImageStateDeleted,
		time.Time{},
		time.Time{},
	)

	err := image.CanDelete()

	if !errors.Is(err, ErrImageAlreadyDeleted) {
		t.Fatalf("expected ErrImageAlreadyDeleted, got %v", err)
	}

	if !image.IsDeleted() {
		t.Error("expected image to remain deleted")
	}

	if !image.UpdateTime().IsZero() {
		t.Errorf("expected update time to remain zero, got %v", image.UpdateTime())
	}

	if events := image.PullEvents(); len(events) != 0 {
		t.Errorf("expected no events, got %d", len(events))
	}
}

func TestImage_Delete(t *testing.T) {
	image := ReconstituteImage(
		vo.ImageID{},
		vo.Tags{},
		vo.ObjectKey{},
		vo.ObjectKey{},
		vo.ImageStateCreated,
		time.Time{},
		time.Time{},
	)

	now := time.Now().UTC()

	err := image.Delete(now)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if !image.IsDeleted() {
		t.Error("expected image to be deleted")
	}

	if !image.UpdateTime().Equal(now) {
		t.Errorf("expected update time %v, got %v", now, image.UpdateTime())
	}

	events := image.PullEvents()

	if len(events) != 1 {
		t.Fatalf("expected 1 event, got %d", len(events))
	}

	if _, ok := events[0].(event.ImageDeleted); !ok {
		t.Errorf(
			"expected ImageDeleted, got %T",
			events[0],
		)
	}
}

func TestImage_Delete_WhenUploading(t *testing.T) {
	image := ReconstituteImage(
		vo.ImageID{},
		vo.Tags{},
		vo.ObjectKey{},
		vo.ObjectKey{},
		vo.ImageStateUploading,
		time.Time{},
		time.Time{},
	)

	now := time.Now().UTC()

	err := image.Delete(now)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if !image.IsDeleted() {
		t.Error("expected image to be deleted")
	}

	if !image.UpdateTime().Equal(now) {
		t.Errorf("expected update time %v, got %v", now, image.UpdateTime())
	}

	events := image.PullEvents()

	if len(events) != 1 {
		t.Fatalf("expected 1 event, got %d", len(events))
	}

	if _, ok := events[0].(event.ImageDeleted); !ok {
		t.Errorf(
			"expected ImageDeleted, got %T",
			events[0],
		)
	}
}

func TestImage_Delete_WhenUploaded(t *testing.T) {
	image := ReconstituteImage(
		vo.ImageID{},
		vo.Tags{},
		vo.ObjectKey{},
		vo.ObjectKey{},
		vo.ImageStateUploaded,
		time.Time{},
		time.Time{},
	)

	now := time.Now().UTC()

	err := image.Delete(now)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if !image.IsDeleted() {
		t.Error("expected image to be deleted")
	}

	if !image.UpdateTime().Equal(now) {
		t.Errorf("expected update time %v, got %v", now, image.UpdateTime())
	}

	events := image.PullEvents()

	if len(events) != 1 {
		t.Fatalf("expected 1 event, got %d", len(events))
	}

	if _, ok := events[0].(event.ImageDeleted); !ok {
		t.Errorf(
			"expected ImageDeleted, got %T",
			events[0],
		)
	}
}

func TestImage_Delete_WhenProcessing(t *testing.T) {
	image := ReconstituteImage(
		vo.ImageID{},
		vo.Tags{},
		vo.ObjectKey{},
		vo.ObjectKey{},
		vo.ImageStateProcessing,
		time.Time{},
		time.Time{},
	)

	now := time.Now().UTC()

	err := image.Delete(now)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if !image.IsDeleted() {
		t.Error("expected image to be deleted")
	}

	if !image.UpdateTime().Equal(now) {
		t.Errorf("expected update time %v, got %v", now, image.UpdateTime())
	}

	events := image.PullEvents()

	if len(events) != 1 {
		t.Fatalf("expected 1 event, got %d", len(events))
	}

	if _, ok := events[0].(event.ImageDeleted); !ok {
		t.Errorf(
			"expected ImageDeleted, got %T",
			events[0],
		)
	}
}

func TestImage_Delete_WhenProcessed(t *testing.T) {
	image := ReconstituteImage(
		vo.ImageID{},
		vo.Tags{},
		vo.ObjectKey{},
		vo.ObjectKey{},
		vo.ImageStateProcessed,
		time.Time{},
		time.Time{},
	)

	now := time.Now().UTC()

	err := image.Delete(now)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if !image.IsDeleted() {
		t.Error("expected image to be deleted")
	}

	if !image.UpdateTime().Equal(now) {
		t.Errorf("expected update time %v, got %v", now, image.UpdateTime())
	}

	events := image.PullEvents()

	if len(events) != 1 {
		t.Fatalf("expected 1 event, got %d", len(events))
	}

	if _, ok := events[0].(event.ImageDeleted); !ok {
		t.Errorf(
			"expected ImageDeleted, got %T",
			events[0],
		)
	}
}

func TestImage_Delete_WhenAlreadyDeleted(t *testing.T) {
	image := ReconstituteImage(
		vo.ImageID{},
		vo.Tags{},
		vo.ObjectKey{},
		vo.ObjectKey{},
		vo.ImageStateDeleted,
		time.Time{},
		time.Time{},
	)

	now := time.Now().UTC()

	err := image.Delete(now)

	if !errors.Is(err, ErrImageAlreadyDeleted) {
		t.Fatalf("expected ErrImageAlreadyDeleted, got %v", err)
	}

	if !image.IsDeleted() {
		t.Error("expected image to remain deleted")
	}

	if !image.UpdateTime().IsZero() {
		t.Errorf("expected update time to remain zero, got %v", image.UpdateTime())
	}

	if events := image.PullEvents(); len(events) != 0 {
		t.Errorf("expected no events, got %d", len(events))
	}
}

func TestImage_PullEvents(t *testing.T) {
	image := ReconstituteImage(
		vo.ImageID{},
		vo.Tags{},
		vo.ObjectKey{},
		vo.ObjectKey{},
		vo.ImageStateCreated,
		time.Time{},
		time.Time{},
	)

	err := image.Delete(time.Now().UTC())
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	first := image.PullEvents()

	if len(first) != 1 {
		t.Fatalf("expected 1 event, got %d", len(first))
	}

	second := image.PullEvents()

	if len(second) != 0 {
		t.Errorf("expected events to be consumed, got %d", len(second))
	}
}

func TestImage_PullEvents_WhenMultipleEvents(t *testing.T) {
	image := ReconstituteImage(
		vo.ImageID{},
		vo.Tags{},
		vo.ObjectKey{},
		vo.ObjectKey{},
		vo.ImageStateUploading,
		time.Time{},
		time.Time{},
	)

	now := time.Now().UTC()

	err := image.ConfirmUpload(now)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	err = image.StartProcessing(now)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	err = image.CompleteProcessing(now)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	err = image.Delete(now)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	events := image.PullEvents()

	if len(events) != 3 {
		t.Fatalf("expected 3 events, got %d", len(events))
	}

	if _, ok := events[0].(event.ImageUploaded); !ok {
		t.Errorf(
			"expected ImageUploaded, got %T",
			events[0],
		)
	}

	if _, ok := events[1].(event.ImageProcessed); !ok {
		t.Errorf(
			"expected ImageProcessed, got %T",
			events[1],
		)
	}

	if _, ok := events[2].(event.ImageDeleted); !ok {
		t.Errorf(
			"expected ImageDeleted, got %T",
			events[2],
		)
	}

	if events := image.PullEvents(); len(events) != 0 {
		t.Errorf("expected events to be consumed, got %d", len(events))
	}
}
