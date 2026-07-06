package slots

import (
	"context"
	"errors"
	"testing"
	"time"
)

type fakeRepository struct {
	listFn    func(ctx context.Context, filters ListFilters) ([]Slot, error)
	getByIDFn func(ctx context.Context, id string) (*Slot, error)
}

func (f *fakeRepository) List(ctx context.Context, filters ListFilters) ([]Slot, error) {
	return f.listFn(ctx, filters)
}

func (f *fakeRepository) GetByID(ctx context.Context, id string) (*Slot, error) {
	return f.getByIDFn(ctx, id)
}

func TestServiceListSuccess(t *testing.T) {
	expectedSlots := []Slot{
		{
			ID:                  "slot-1",
			StartAt:             time.Now().Add(24 * time.Hour),
			ArrivalAt:           time.Now().Add(24*time.Hour - 15*time.Minute),
			TrackConfigName:     "Короткая трасса",
			TrackConfigCode:     "short",
			MarshalName:         "Алексей",
			TotalSeats:          8,
			BookedSeats:         2,
			FreeSeats:           6,
			FreeRentalEquipment: 3,
			Price:               2500,
			RentalPrice:         500,
			Currency:            "RUB",
			Status:              "active",
		},
	}

	repo := &fakeRepository{
		listFn: func(ctx context.Context, filters ListFilters) ([]Slot, error) {
			if filters.OnlyAvailable != true {
				t.Fatalf("expected only_available=true")
			}

			if filters.TrackConfig != "short" {
				t.Fatalf("expected track_config=short, got %s", filters.TrackConfig)
			}

			return expectedSlots, nil
		},
	}

	service := NewService(repo)

	result, err := service.List(context.Background(), ListFilters{
		TrackConfig:   "short",
		OnlyAvailable: true,
	})

	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	if len(result) != 1 {
		t.Fatalf("expected 1 slot, got %d", len(result))
	}

	if result[0].ID != "slot-1" {
		t.Fatalf("expected slot id slot-1, got %s", result[0].ID)
	}
}

func TestServiceListRepositoryError(t *testing.T) {
	expectedErr := errors.New("db error")

	repo := &fakeRepository{
		listFn: func(ctx context.Context, filters ListFilters) ([]Slot, error) {
			return nil, expectedErr
		},
	}

	service := NewService(repo)

	_, err := service.List(context.Background(), ListFilters{})
	if !errors.Is(err, expectedErr) {
		t.Fatalf("expected db error, got %v", err)
	}
}

func TestServiceGetByIDSuccess(t *testing.T) {
	expectedSlot := &Slot{
		ID:              "slot-1",
		TrackConfigName: "Длинная трасса",
		TrackConfigCode: "long",
		TotalSeats:      14,
		FreeSeats:       10,
		Price:           3000,
		Currency:        "RUB",
		Status:          "active",
	}

	repo := &fakeRepository{
		getByIDFn: func(ctx context.Context, id string) (*Slot, error) {
			if id != "slot-1" {
				t.Fatalf("expected id slot-1, got %s", id)
			}

			return expectedSlot, nil
		},
	}

	service := NewService(repo)

	result, err := service.GetByID(context.Background(), "slot-1")
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	if result.ID != "slot-1" {
		t.Fatalf("expected slot id slot-1, got %s", result.ID)
	}

	if result.TrackConfigCode != "long" {
		t.Fatalf("expected track_config_code long, got %s", result.TrackConfigCode)
	}
}

func TestServiceGetByIDNotFound(t *testing.T) {
	repo := &fakeRepository{
		getByIDFn: func(ctx context.Context, id string) (*Slot, error) {
			return nil, ErrSlotNotFound
		},
	}

	service := NewService(repo)

	_, err := service.GetByID(context.Background(), "missing-slot")
	if !errors.Is(err, ErrSlotNotFound) {
		t.Fatalf("expected ErrSlotNotFound, got %v", err)
	}
}
