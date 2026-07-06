package bookings

import (
	"context"
	"errors"
	"strings"
)

var ErrMissingIdempotencyKey = errors.New("missing idempotency key")

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Create(
	ctx context.Context,
	userID string,
	req CreateBookingRequest,
	idempotencyKey string,
) (*Booking, error) {
	idempotencyKey = strings.TrimSpace(idempotencyKey)
	if idempotencyKey == "" {
		return nil, ErrMissingIdempotencyKey
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
