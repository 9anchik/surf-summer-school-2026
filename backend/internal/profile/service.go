package profile

import (
	"context"
	"errors"
	"strings"
)

var ErrInvalidProfileData = errors.New("invalid profile data")

type Repository interface {
	GetByID(ctx context.Context, userID string) (*UserProfile, error)
	Update(ctx context.Context, userID string, req UpdateProfileRequest) (*UserProfile, error)
	Delete(ctx context.Context, userID string) error
}

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Get(ctx context.Context, userID string) (*UserProfile, error) {
	return s.repo.GetByID(ctx, userID)
}

func (s *Service) Update(ctx context.Context, userID string, req UpdateProfileRequest) (*UserProfile, error) {
	if req.Name != nil {
		name := strings.TrimSpace(*req.Name)
		if name == "" {
			return nil, ErrInvalidProfileData
		}
		req.Name = &name
	}

	if req.Phone != nil {
		phone := strings.TrimSpace(*req.Phone)
		if phone == "" {
			return nil, ErrInvalidProfileData
		}
		req.Phone = &phone
	}

	return s.repo.Update(ctx, userID, req)
}

func (s *Service) Delete(ctx context.Context, userID string) error {
	return s.repo.Delete(ctx, userID)
}
