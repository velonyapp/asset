package entity

import (
	"errors"
	"reflect"
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
		vo.ImageWidth{},
		vo.ImageHeight{},
		vo.ImageFormatPNG,
		time.Now().UTC(),
	)

	if !errors.Is(err, ErrImageObjectKeysEqual) {
		t.Fatalf("expected ErrImageObjectKeysEqual, got %v", err)
	}

	if image != nil {
		t.Error("expected no image")
	}
}

func TestImage_Transitions(t *testing.T) {
	type testCase struct {
		name    string
		state   vo.ImageState
		wantErr error
	}

	tests := []struct {
		name      string
		can       func(*Image) error
		apply     func(*Image, time.Time) error
		nextState vo.ImageState
		wantEvent reflect.Type
		cases     []testCase
	}{
		{
			name:      "ConfirmUpload",
			can:       (*Image).CanConfirmUpload,
			apply:     (*Image).ConfirmUpload,
			nextState: vo.ImageStatePending,
			wantEvent: reflect.TypeOf(event.ImageConfirmed{}),
			cases: []testCase{
				{"AwaitingUpload", vo.ImageStateAwaitingUpload, nil},
				{"Pending", vo.ImageStatePending, ErrImageAlreadyPending},
				{"Processing", vo.ImageStateProcessing, ErrImageProcessing},
				{"Active", vo.ImageStateActive, ErrImageActive},
				{"Deleted", vo.ImageStateDeleted, ErrImageDeleted},
			},
		},
		{
			name:      "StartProcessing",
			can:       (*Image).CanStartProcessing,
			apply:     (*Image).StartProcessing,
			nextState: vo.ImageStateProcessing,
			cases: []testCase{
				{"AwaitingUpload", vo.ImageStateAwaitingUpload, ErrImageNotPending},
				{"Pending", vo.ImageStatePending, nil},
				{"Processing", vo.ImageStateProcessing, ErrImageAlreadyProcessing},
				{"Active", vo.ImageStateActive, ErrImageActive},
				{"Deleted", vo.ImageStateDeleted, ErrImageDeleted},
			},
		},
		{
			name:      "CancelProcessing",
			can:       (*Image).CanCancelProcessing,
			apply:     (*Image).CancelProcessing,
			nextState: vo.ImageStatePending,
			cases: []testCase{
				{"AwaitingUpload", vo.ImageStateAwaitingUpload, ErrImageNotProcessing},
				{"Pending", vo.ImageStatePending, ErrImageNotProcessing},
				{"Processing", vo.ImageStateProcessing, nil},
				{"Active", vo.ImageStateActive, ErrImageNotProcessing},
				{"Deleted", vo.ImageStateDeleted, ErrImageDeleted},
			},
		},
		{
			name:      "CompleteProcessing",
			can:       (*Image).CanCompleteProcessing,
			apply:     (*Image).CompleteProcessing,
			nextState: vo.ImageStateActive,
			wantEvent: reflect.TypeOf(event.ImageActive{}),
			cases: []testCase{
				{"AwaitingUpload", vo.ImageStateAwaitingUpload, ErrImageNotProcessing},
				{"Pending", vo.ImageStatePending, ErrImageNotProcessing},
				{"Processing", vo.ImageStateProcessing, nil},
				{"Active", vo.ImageStateActive, ErrImageAlreadyActive},
				{"Deleted", vo.ImageStateDeleted, ErrImageDeleted},
			},
		},
		{
			name:      "Delete",
			can:       (*Image).CanDelete,
			apply:     (*Image).Delete,
			nextState: vo.ImageStateDeleted,
			wantEvent: reflect.TypeOf(event.ImageDeleted{}),
			cases: []testCase{
				{"AwaitingUpload", vo.ImageStateAwaitingUpload, nil},
				{"Pending", vo.ImageStatePending, nil},
				{"Processing", vo.ImageStateProcessing, nil},
				{"Active", vo.ImageStateActive, nil},
				{"Deleted", vo.ImageStateDeleted, ErrImageAlreadyDeleted},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for _, tc := range tt.cases {
				t.Run(tc.name, func(t *testing.T) {
					image := ReconstituteImage(
						vo.ImageID{},
						vo.Tags{},
						vo.ObjectKey{},
						vo.ObjectKey{},
						tc.state,
						vo.ImageWidth{},
						vo.ImageHeight{},
						vo.ImageFormatPNG,
						time.Time{},
						time.Time{},
					)

					err := tt.can(image)

					if !errors.Is(err, tc.wantErr) {
						t.Fatalf("expected %v, got %v", tc.wantErr, err)
					}

					if !image.State().Equal(tc.state) {
						t.Error("expected image state to remain unchanged")
					}

					if !image.UpdateTime().IsZero() {
						t.Errorf("expected update time to remain zero, got %v", image.UpdateTime())
					}

					if events := image.PullEvents(); len(events) != 0 {
						t.Errorf("expected no events, got %d", len(events))
					}

					now := time.Now().UTC()

					err = tt.apply(image, now)

					if !errors.Is(err, tc.wantErr) {
						t.Fatalf("expected %v, got %v", tc.wantErr, err)
					}

					if tc.wantErr != nil {
						if !image.State().Equal(tc.state) {
							t.Error("expected image state to remain unchanged")
						}

						if !image.UpdateTime().IsZero() {
							t.Errorf("expected update time to remain zero, got %v", image.UpdateTime())
						}

						if events := image.PullEvents(); len(events) != 0 {
							t.Errorf("expected no events, got %d", len(events))
						}

						return
					}

					if !image.State().Equal(tt.nextState) {
						t.Errorf("expected image state %v, got %v", tt.nextState, image.State())
					}

					if !image.UpdateTime().Equal(now) {
						t.Errorf("expected update time %v, got %v", now, image.UpdateTime())
					}

					events := image.PullEvents()

					if tt.wantEvent == nil {
						if len(events) != 0 {
							t.Errorf("expected no events, got %d", len(events))
						}

						return
					}

					if len(events) != 1 {
						t.Fatalf("expected 1 event, got %d", len(events))
					}

					if reflect.TypeOf(events[0]) != tt.wantEvent {
						t.Errorf("expected %v, got %T", tt.wantEvent, events[0])
					}
				})
			}
		})
	}
}

func TestImage_PullEvents(t *testing.T) {
	image := ReconstituteImage(
		vo.ImageID{},
		vo.Tags{},
		vo.ObjectKey{},
		vo.ObjectKey{},
		vo.ImageStatePending,
		vo.ImageWidth{},
		vo.ImageHeight{},
		vo.ImageFormatPNG,
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
		vo.ImageStateAwaitingUpload,
		vo.ImageWidth{},
		vo.ImageHeight{},
		vo.ImageFormatPNG,
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

	if _, ok := events[0].(event.ImageConfirmed); !ok {
		t.Errorf("expected ImageConfirmed, got %T", events[0])
	}

	if _, ok := events[1].(event.ImageActive); !ok {
		t.Errorf("expected ImageActive, got %T", events[1])
	}

	if _, ok := events[2].(event.ImageDeleted); !ok {
		t.Errorf("expected ImageDeleted, got %T", events[2])
	}

	if events := image.PullEvents(); len(events) != 0 {
		t.Errorf("expected events to be consumed, got %d", len(events))
	}
}
