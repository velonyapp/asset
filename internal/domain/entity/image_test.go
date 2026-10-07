package entity

import (
	"errors"
	"testing"
	"time"

	"github.com/velonyapp/asset/internal/domain/event"
	"github.com/velonyapp/asset/internal/domain/vo"
)

func TestImage_CanUpload(t *testing.T) {
	image := ReconstituteImage(
		vo.ImageID{},
		vo.Tags{},
		vo.ObjectKey{},
		vo.ObjectKey{},
		vo.ImageStatePending,
		time.Time{},
		nil,
	)

	err := image.CanUpload()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if !image.IsPending() {
		t.Error("expected image to remain pending")
	}

	if events := image.PullEvents(); len(events) != 0 {
		t.Errorf("expected no events, got %d", len(events))
	}
}

func TestImage_CanUpload_WhenAlreadyUploaded(t *testing.T) {
	image := ReconstituteImage(
		vo.ImageID{},
		vo.Tags{},
		vo.ObjectKey{},
		vo.ObjectKey{},
		vo.ImageStateUploaded,
		time.Time{},
		nil,
	)

	err := image.CanUpload()

	if !errors.Is(err, ErrImageAlreadyUploaded) {
		t.Fatalf("expected ErrImageAlreadyUploaded, got %v", err)
	}

	if !image.IsUploaded() {
		t.Error("expected image to remain uploaded")
	}

	if events := image.PullEvents(); len(events) != 0 {
		t.Errorf("expected no events, got %d", len(events))
	}
}

func TestImage_CanUpload_WhenProcessed(t *testing.T) {
	image := ReconstituteImage(
		vo.ImageID{},
		vo.Tags{},
		vo.ObjectKey{},
		vo.ObjectKey{},
		vo.ImageStateProcessed,
		time.Time{},
		nil,
	)

	err := image.CanUpload()

	if !errors.Is(err, ErrImageProcessed) {
		t.Fatalf("expected ErrImageProcessed, got %v", err)
	}

	if !image.IsProcessed() {
		t.Error("expected image to remain processed")
	}

	if events := image.PullEvents(); len(events) != 0 {
		t.Errorf("expected no events, got %d", len(events))
	}
}

func TestImage_CanUpload_WhenDeleted(t *testing.T) {
	deleteTime := time.Now().UTC()

	image := ReconstituteImage(
		vo.ImageID{},
		vo.Tags{},
		vo.ObjectKey{},
		vo.ObjectKey{},
		vo.ImageStatePending,
		time.Time{},
		&deleteTime,
	)

	err := image.CanUpload()

	if !errors.Is(err, ErrImageDeleted) {
		t.Fatalf("expected ErrImageDeleted, got %v", err)
	}

	if !image.IsPending() {
		t.Error("expected image to remain pending")
	}

	if events := image.PullEvents(); len(events) != 0 {
		t.Errorf("expected no events, got %d", len(events))
	}
}

func TestImage_Upload(t *testing.T) {
	image := ReconstituteImage(
		vo.ImageID{},
		vo.Tags{},
		vo.ObjectKey{},
		vo.ObjectKey{},
		vo.ImageStatePending,
		time.Time{},
		nil,
	)

	now := time.Now().UTC()

	err := image.Upload(now)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if !image.IsUploaded() {
		t.Error("expected image to be uploaded")
	}

	events := image.PullEvents()

	if len(events) != 1 {
		t.Fatalf("expected 1 event, got %d", len(events))
	}

	if _, ok := events[0].(event.ImageUpdated); !ok {
		t.Errorf(
			"expected ImageUpdated, got %T",
			events[0],
		)
	}
}

func TestImage_Upload_WhenAlreadyUploaded(t *testing.T) {
	image := ReconstituteImage(
		vo.ImageID{},
		vo.Tags{},
		vo.ObjectKey{},
		vo.ObjectKey{},
		vo.ImageStateUploaded,
		time.Time{},
		nil,
	)

	err := image.Upload(time.Now().UTC())

	if !errors.Is(err, ErrImageAlreadyUploaded) {
		t.Fatalf("expected ErrImageAlreadyUploaded, got %v", err)
	}

	if !image.IsUploaded() {
		t.Error("expected image to remain uploaded")
	}

	if events := image.PullEvents(); len(events) != 0 {
		t.Errorf("expected no events, got %d", len(events))
	}
}

func TestImage_Upload_WhenProcessed(t *testing.T) {
	image := ReconstituteImage(
		vo.ImageID{},
		vo.Tags{},
		vo.ObjectKey{},
		vo.ObjectKey{},
		vo.ImageStateProcessed,
		time.Time{},
		nil,
	)

	err := image.Upload(time.Now().UTC())

	if !errors.Is(err, ErrImageProcessed) {
		t.Fatalf("expected ErrImageProcessed, got %v", err)
	}

	if !image.IsProcessed() {
		t.Error("expected image to remain processed")
	}

	if events := image.PullEvents(); len(events) != 0 {
		t.Errorf("expected no events, got %d", len(events))
	}
}

func TestImage_Upload_WhenDeleted(t *testing.T) {
	deleteTime := time.Now().UTC()

	image := ReconstituteImage(
		vo.ImageID{},
		vo.Tags{},
		vo.ObjectKey{},
		vo.ObjectKey{},
		vo.ImageStatePending,
		time.Time{},
		&deleteTime,
	)

	err := image.Upload(time.Now().UTC())

	if !errors.Is(err, ErrImageDeleted) {
		t.Fatalf("expected ErrImageDeleted, got %v", err)
	}

	if !image.IsPending() {
		t.Error("expected image to remain pending")
	}

	if events := image.PullEvents(); len(events) != 0 {
		t.Errorf("expected no events, got %d", len(events))
	}
}

func TestImage_CanProcess(t *testing.T) {
	image := ReconstituteImage(
		vo.ImageID{},
		vo.Tags{},
		vo.ObjectKey{},
		vo.ObjectKey{},
		vo.ImageStateUploaded,
		time.Time{},
		nil,
	)

	err := image.CanProcess()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if !image.IsUploaded() {
		t.Error("expected image to remain uploaded")
	}

	if events := image.PullEvents(); len(events) != 0 {
		t.Errorf("expected no events, got %d", len(events))
	}
}

func TestImage_CanProcess_WhenAlreadyProcessed(t *testing.T) {
	image := ReconstituteImage(
		vo.ImageID{},
		vo.Tags{},
		vo.ObjectKey{},
		vo.ObjectKey{},
		vo.ImageStateProcessed,
		time.Time{},
		nil,
	)

	err := image.CanProcess()

	if !errors.Is(err, ErrImageAlreadyProcessed) {
		t.Fatalf("expected ErrImageAlreadyProcessed, got %v", err)
	}

	if !image.IsProcessed() {
		t.Error("expected image to remain processed")
	}

	if events := image.PullEvents(); len(events) != 0 {
		t.Errorf("expected no events, got %d", len(events))
	}
}

func TestImage_CanProcess_WhenPending(t *testing.T) {
	image := ReconstituteImage(
		vo.ImageID{},
		vo.Tags{},
		vo.ObjectKey{},
		vo.ObjectKey{},
		vo.ImageStatePending,
		time.Time{},
		nil,
	)

	err := image.CanProcess()

	if !errors.Is(err, ErrImageNotUploaded) {
		t.Fatalf("expected ErrImageNotUploaded, got %v", err)
	}

	if !image.IsPending() {
		t.Error("expected image to remain pending")
	}

	if events := image.PullEvents(); len(events) != 0 {
		t.Errorf("expected no events, got %d", len(events))
	}
}

func TestImage_CanProcess_WhenDeleted(t *testing.T) {
	deleteTime := time.Now().UTC()

	image := ReconstituteImage(
		vo.ImageID{},
		vo.Tags{},
		vo.ObjectKey{},
		vo.ObjectKey{},
		vo.ImageStateUploaded,
		time.Time{},
		&deleteTime,
	)

	err := image.CanProcess()

	if !errors.Is(err, ErrImageDeleted) {
		t.Fatalf("expected ErrImageDeleted, got %v", err)
	}

	if !image.IsUploaded() {
		t.Error("expected image to remain uploaded")
	}

	if events := image.PullEvents(); len(events) != 0 {
		t.Errorf("expected no events, got %d", len(events))
	}
}

func TestImage_Process(t *testing.T) {
	image := ReconstituteImage(
		vo.ImageID{},
		vo.Tags{},
		vo.ObjectKey{},
		vo.ObjectKey{},
		vo.ImageStateUploaded,
		time.Time{},
		nil,
	)

	now := time.Now().UTC()

	err := image.Process(now)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if !image.IsProcessed() {
		t.Error("expected image to be processed")
	}

	events := image.PullEvents()

	if len(events) != 1 {
		t.Fatalf("expected 1 event, got %d", len(events))
	}

	if _, ok := events[0].(event.ImageUpdated); !ok {
		t.Errorf(
			"expected ImageUpdated, got %T",
			events[0],
		)
	}
}

func TestImage_Process_WhenAlreadyProcessed(t *testing.T) {
	image := ReconstituteImage(
		vo.ImageID{},
		vo.Tags{},
		vo.ObjectKey{},
		vo.ObjectKey{},
		vo.ImageStateProcessed,
		time.Time{},
		nil,
	)

	err := image.Process(time.Now().UTC())

	if !errors.Is(err, ErrImageAlreadyProcessed) {
		t.Fatalf("expected ErrImageAlreadyProcessed, got %v", err)
	}

	if !image.IsProcessed() {
		t.Error("expected image to remain processed")
	}

	if events := image.PullEvents(); len(events) != 0 {
		t.Errorf("expected no events, got %d", len(events))
	}
}

func TestImage_Process_WhenPending(t *testing.T) {
	image := ReconstituteImage(
		vo.ImageID{},
		vo.Tags{},
		vo.ObjectKey{},
		vo.ObjectKey{},
		vo.ImageStatePending,
		time.Time{},
		nil,
	)

	err := image.Process(time.Now().UTC())

	if !errors.Is(err, ErrImageNotUploaded) {
		t.Fatalf("expected ErrImageNotUploaded, got %v", err)
	}

	if !image.IsPending() {
		t.Error("expected image to remain pending")
	}

	if events := image.PullEvents(); len(events) != 0 {
		t.Errorf("expected no events, got %d", len(events))
	}
}

func TestImage_Process_WhenDeleted(t *testing.T) {
	deleteTime := time.Now().UTC()

	image := ReconstituteImage(
		vo.ImageID{},
		vo.Tags{},
		vo.ObjectKey{},
		vo.ObjectKey{},
		vo.ImageStateUploaded,
		time.Time{},
		&deleteTime,
	)

	err := image.Process(time.Now().UTC())

	if !errors.Is(err, ErrImageDeleted) {
		t.Fatalf("expected ErrImageDeleted, got %v", err)
	}

	if !image.IsUploaded() {
		t.Error("expected image to remain uploaded")
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
		vo.ImageStatePending,
		time.Time{},
		nil,
	)

	err := image.CanDelete()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if image.IsDeleted() {
		t.Error("expected image to remain not deleted")
	}

	if events := image.PullEvents(); len(events) != 0 {
		t.Errorf("expected no events, got %d", len(events))
	}
}

func TestImage_CanDelete_WhenAlreadyDeleted(t *testing.T) {
	deleteTime := time.Now().UTC()

	image := ReconstituteImage(
		vo.ImageID{},
		vo.Tags{},
		vo.ObjectKey{},
		vo.ObjectKey{},
		vo.ImageStatePending,
		time.Time{},
		&deleteTime,
	)

	err := image.CanDelete()

	if !errors.Is(err, ErrImageAlreadyDeleted) {
		t.Fatalf("expected ErrImageAlreadyDeleted, got %v", err)
	}

	actualDeleteTime := image.DeleteTime()
	if actualDeleteTime == nil {
		t.Fatal("expected delete time")
	}

	if !actualDeleteTime.Equal(deleteTime) {
		t.Errorf("expected delete time %v, got %v", deleteTime, *actualDeleteTime)
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
		vo.ImageStatePending,
		time.Time{},
		nil,
	)

	now := time.Now().UTC()

	err := image.Delete(now)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if !image.IsDeleted() {
		t.Error("expected image to be deleted")
	}

	deleteTime := image.DeleteTime()
	if deleteTime == nil {
		t.Fatal("expected delete time")
	}

	if !deleteTime.Equal(now) {
		t.Errorf("expected delete time %v, got %v", now, *deleteTime)
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
	firstDeleteTime := time.Now().UTC()

	image := ReconstituteImage(
		vo.ImageID{},
		vo.Tags{},
		vo.ObjectKey{},
		vo.ObjectKey{},
		vo.ImageStatePending,
		time.Time{},
		&firstDeleteTime,
	)

	err := image.Delete(firstDeleteTime.Add(time.Hour))

	if !errors.Is(err, ErrImageAlreadyDeleted) {
		t.Fatalf("expected ErrImageAlreadyDeleted, got %v", err)
	}

	deleteTime := image.DeleteTime()
	if deleteTime == nil {
		t.Fatal("expected delete time")
	}

	if !deleteTime.Equal(firstDeleteTime) {
		t.Errorf("expected original delete time %v, got %v", firstDeleteTime, *deleteTime)
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
		vo.ImageStatePending,
		time.Time{},
		nil,
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
