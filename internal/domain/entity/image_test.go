package entity

import (
	"errors"
	"testing"
	"time"

	"github.com/velonyapp/asset/internal/domain/event"
	"github.com/velonyapp/asset/internal/domain/vo"
)

func TestImage_Reconcile(t *testing.T) {
	var (
		imageID   vo.ImageID
		tags      vo.Tags
		objectKey vo.ObjectKey
	)

	image := ReconstituteImage(
		imageID,
		tags,
		objectKey,
		false,
		time.Time{},
		nil,
	)

	now := time.Now().UTC()

	err := image.Reconcile(true, now)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if !image.ObjectExists() {
		t.Error("expected object to exist")
	}

	events := image.PullEvents()

	if len(events) != 1 {
		t.Fatalf("expected 1 event, got %d", len(events))
	}

	if _, ok := events[0].(event.ImageReconciled); !ok {
		t.Errorf(
			"expected ImageReconciled, got %T",
			events[0],
		)
	}
}

func TestImage_Reconcile_WhenValueDoesNotChange(t *testing.T) {
	image := ReconstituteImage(
		vo.ImageID{},
		vo.Tags{},
		vo.ObjectKey{},
		true,
		time.Time{},
		nil,
	)

	err := image.Reconcile(true, time.Now().UTC())
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	events := image.PullEvents()

	if len(events) != 0 {
		t.Errorf("expected no events, got %d", len(events))
	}
}

func TestImage_Reconcile_WhenDeleted(t *testing.T) {
	deleteTime := time.Now().UTC()

	image := ReconstituteImage(
		vo.ImageID{},
		vo.Tags{},
		vo.ObjectKey{},
		false,
		time.Time{},
		&deleteTime,
	)

	err := image.Reconcile(true, time.Now().UTC())

	if !errors.Is(err, ErrImageDeleted) {
		t.Fatalf("expected ErrImageDeleted, got %v", err)
	}

	if image.ObjectExists() {
		t.Error("expected object existence to remain unchanged")
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
		false,
		time.Time{},
		nil,
	)

	now := time.Now().UTC()

	image.Delete(now)

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
		false,
		time.Time{},
		&firstDeleteTime,
	)

	image.Delete(firstDeleteTime.Add(time.Hour))

	deleteTime := image.DeleteTime()
	if deleteTime == nil {
		t.Fatal("expected delete time")
	}

	if !deleteTime.Equal(firstDeleteTime) {
		t.Errorf(
			"expected original delete time %v, got %v",
			firstDeleteTime,
			*deleteTime,
		)
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
		false,
		time.Time{},
		nil,
	)

	image.Delete(time.Now().UTC())

	first := image.PullEvents()

	if len(first) != 1 {
		t.Fatalf("expected 1 event, got %d", len(first))
	}

	second := image.PullEvents()

	if len(second) != 0 {
		t.Errorf("expected events to be consumed, got %d", len(second))
	}
}
