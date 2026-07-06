package bookings

import (
	"errors"
	"testing"
)

func TestValidateCreateBookingRequestSuccess(t *testing.T) {
	req := CreateBookingRequest{
		SlotID:    "slot-1",
		Equipment: []string{"own", "rental"},
	}

	err := ValidateCreateBookingRequest(req, "key-1")
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
}

func TestValidateCreateBookingRequestMissingIdempotencyKey(t *testing.T) {
	req := CreateBookingRequest{
		SlotID:    "slot-1",
		Equipment: []string{"own"},
	}

	err := ValidateCreateBookingRequest(req, "")
	if !errors.Is(err, ErrMissingIdempotencyKey) {
		t.Fatalf("expected ErrMissingIdempotencyKey, got %v", err)
	}
}

func TestValidateCreateBookingRequestTooManySeats(t *testing.T) {
	req := CreateBookingRequest{
		SlotID:    "slot-1",
		Equipment: []string{"own", "own", "rental", "rental"},
	}

	err := ValidateCreateBookingRequest(req, "key-1")
	if !errors.Is(err, ErrInvalidRequest) {
		t.Fatalf("expected ErrInvalidRequest, got %v", err)
	}
}

func TestValidateCreateBookingRequestInvalidEquipment(t *testing.T) {
	req := CreateBookingRequest{
		SlotID:    "slot-1",
		Equipment: []string{"own", "helmet"},
	}

	err := ValidateCreateBookingRequest(req, "key-1")
	if !errors.Is(err, ErrInvalidRequest) {
		t.Fatalf("expected ErrInvalidRequest, got %v", err)
	}
}

func TestValidateCreateBookingRequestEmptySlotID(t *testing.T) {
	req := CreateBookingRequest{
		SlotID:    "",
		Equipment: []string{"own"},
	}

	err := ValidateCreateBookingRequest(req, "key-1")
	if !errors.Is(err, ErrInvalidRequest) {
		t.Fatalf("expected ErrInvalidRequest, got %v", err)
	}
}
