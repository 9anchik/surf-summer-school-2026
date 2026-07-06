package profile

import (
	"context"
	"errors"
	"testing"
	"time"
)

type fakeRepository struct {
	getByIDFn func(ctx context.Context, userID string) (*UserProfile, error)
	updateFn  func(ctx context.Context, userID string, req UpdateProfileRequest) (*UserProfile, error)
	deleteFn  func(ctx context.Context, userID string) error
}

func (f *fakeRepository) GetByID(ctx context.Context, userID string) (*UserProfile, error) {
	return f.getByIDFn(ctx, userID)
}

func (f *fakeRepository) Update(ctx context.Context, userID string, req UpdateProfileRequest) (*UserProfile, error) {
	return f.updateFn(ctx, userID, req)
}

func (f *fakeRepository) Delete(ctx context.Context, userID string) error {
	return f.deleteFn(ctx, userID)
}

func TestServiceGetSuccess(t *testing.T) {
	name := "Daniil"

	expected := &UserProfile{
		ID:        "user-1",
		Name:      &name,
		Phone:     "+79990000000",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	repo := &fakeRepository{
		getByIDFn: func(ctx context.Context, userID string) (*UserProfile, error) {
			if userID != "user-1" {
				t.Fatalf("expected userID user-1, got %s", userID)
			}

			return expected, nil
		},
	}

	service := NewService(repo)

	result, err := service.Get(context.Background(), "user-1")
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	if result.ID != "user-1" {
		t.Fatalf("expected user id user-1, got %s", result.ID)
	}

	if result.Name == nil || *result.Name != "Daniil" {
		t.Fatalf("expected name Daniil, got %v", result.Name)
	}
}

func TestServiceGetRepositoryError(t *testing.T) {
	expectedErr := errors.New("db error")

	repo := &fakeRepository{
		getByIDFn: func(ctx context.Context, userID string) (*UserProfile, error) {
			return nil, expectedErr
		},
	}

	service := NewService(repo)

	_, err := service.Get(context.Background(), "user-1")
	if !errors.Is(err, expectedErr) {
		t.Fatalf("expected db error, got %v", err)
	}
}

func TestServiceUpdateSuccessTrimNameAndPhone(t *testing.T) {
	repo := &fakeRepository{
		updateFn: func(ctx context.Context, userID string, req UpdateProfileRequest) (*UserProfile, error) {
			if userID != "user-1" {
				t.Fatalf("expected userID user-1, got %s", userID)
			}

			if req.Name == nil || *req.Name != "Daniil" {
				t.Fatalf("expected trimmed name Daniil, got %v", req.Name)
			}

			if req.Phone == nil || *req.Phone != "+79990000000" {
				t.Fatalf("expected trimmed phone +79990000000, got %v", req.Phone)
			}

			return &UserProfile{
				ID:        "user-1",
				Name:      req.Name,
				Phone:     *req.Phone,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			}, nil
		},
	}

	service := NewService(repo)

	name := "  Daniil  "
	phone := "  +79990000000  "

	result, err := service.Update(context.Background(), "user-1", UpdateProfileRequest{
		Name:  &name,
		Phone: &phone,
	})
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	if result.Name == nil || *result.Name != "Daniil" {
		t.Fatalf("expected name Daniil, got %v", result.Name)
	}

	if result.Phone != "+79990000000" {
		t.Fatalf("expected phone +79990000000, got %s", result.Phone)
	}
}

func TestServiceUpdateEmptyName(t *testing.T) {
	repo := &fakeRepository{}

	service := NewService(repo)

	name := "   "

	_, err := service.Update(context.Background(), "user-1", UpdateProfileRequest{
		Name: &name,
	})

	if !errors.Is(err, ErrInvalidProfileData) {
		t.Fatalf("expected ErrInvalidProfileData, got %v", err)
	}
}

func TestServiceUpdateEmptyPhone(t *testing.T) {
	repo := &fakeRepository{}

	service := NewService(repo)

	phone := "   "

	_, err := service.Update(context.Background(), "user-1", UpdateProfileRequest{
		Phone: &phone,
	})

	if !errors.Is(err, ErrInvalidProfileData) {
		t.Fatalf("expected ErrInvalidProfileData, got %v", err)
	}
}

func TestServiceUpdateRepositoryError(t *testing.T) {
	expectedErr := errors.New("db error")

	repo := &fakeRepository{
		updateFn: func(ctx context.Context, userID string, req UpdateProfileRequest) (*UserProfile, error) {
			return nil, expectedErr
		},
	}

	service := NewService(repo)

	name := "Daniil"

	_, err := service.Update(context.Background(), "user-1", UpdateProfileRequest{
		Name: &name,
	})

	if !errors.Is(err, expectedErr) {
		t.Fatalf("expected db error, got %v", err)
	}
}

func TestServiceDeleteSuccess(t *testing.T) {
	repo := &fakeRepository{
		deleteFn: func(ctx context.Context, userID string) error {
			if userID != "user-1" {
				t.Fatalf("expected userID user-1, got %s", userID)
			}

			return nil
		},
	}

	service := NewService(repo)

	err := service.Delete(context.Background(), "user-1")
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
}

func TestServiceDeleteRepositoryError(t *testing.T) {
	expectedErr := errors.New("db error")

	repo := &fakeRepository{
		deleteFn: func(ctx context.Context, userID string) error {
			return expectedErr
		},
	}

	service := NewService(repo)

	err := service.Delete(context.Background(), "user-1")
	if !errors.Is(err, expectedErr) {
		t.Fatalf("expected db error, got %v", err)
	}
}
