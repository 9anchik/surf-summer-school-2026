package bookings

import (
	"context"
	"errors"
	"strings"
)

var ErrMissingIdempotencyKey = errors.New("missing idempotency key")

type Repository interface {
	Create(ctx context.Context, userID string, req CreateBookingRequest, idempotencyKey string) (*Booking, error)
	ListByUser(ctx context.Context, userID string, limit int, offset int) ([]BookingListItem, error)
	GetByID(ctx context.Context, userID string, bookingID string) (*BookingDetails, error)
	Cancel(ctx context.Context, userID string, bookingID string) (*Booking, error)
}

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Create(
	ctx context.Context,
	userID string,
	req CreateBookingRequest,
	idempotencyKey string,
) (*Booking, error) {
	idempotencyKey = strings.TrimSpace(idempotencyKey)
	if err := ValidateCreateBookingRequest(req, idempotencyKey); err != nil {
		return nil, err
	}

	return s.repo.Create(ctx, userID, req, idempotencyKey)
}

func (s *Service) ListByUser(
	ctx context.Context,
	userID string,
	req ListBookingsRequest,
) ([]BookingListItem, error) {
	if req.Limit <= 0 {
		req.Limit = 20
	}

	if req.Limit > 100 {
		req.Limit = 100
	}

	if req.Offset < 0 {
		req.Offset = 0
	}

	return s.repo.ListByUser(ctx, userID, req.Limit, req.Offset)
}

func (s *Service) Cancel(
	ctx context.Context,
	userID string,
	bookingID string,
) (*Booking, error) {
	bookingID = strings.TrimSpace(bookingID)
	if bookingID == "" {
		return nil, ErrInvalidRequest
	}

	return s.repo.Cancel(ctx, userID, bookingID)
}

func (s *Service) GetByID(
	ctx context.Context,
	userID string,
	bookingID string,
) (*BookingDetails, error) {
	bookingID = strings.TrimSpace(bookingID)
	if bookingID == "" {
		return nil, ErrInvalidRequest
	}

	return s.repo.GetByID(ctx, userID, bookingID)
}

func ValidateCreateBookingRequest(req CreateBookingRequest, idempotencyKey string) error {
	idempotencyKey = strings.TrimSpace(idempotencyKey)
	if idempotencyKey == "" {
		return ErrMissingIdempotencyKey
	}

	if strings.TrimSpace(req.SlotID) == "" {
		return ErrInvalidRequest
	}

	if len(req.Equipment) < 1 || len(req.Equipment) > 3 {
		return ErrInvalidRequest
	}

	for _, item := range req.Equipment {
		if item != "own" && item != "rental" {
			return ErrInvalidRequest
		}
	}

	return nil
}
