package auth

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type fakeRepository struct {
	saveOTPFn          func(ctx context.Context, phone, code string, expiresAt time.Time) error
	verifyOTPFn        func(ctx context.Context, phone, code string) (bool, error)
	findOrCreateUserFn func(ctx context.Context, phone, name string) (string, error)
}

func (f *fakeRepository) SaveOTP(ctx context.Context, phone, code string, expiresAt time.Time) error {
	return f.saveOTPFn(ctx, phone, code, expiresAt)
}

func (f *fakeRepository) VerifyOTP(ctx context.Context, phone, code string) (bool, error) {
	return f.verifyOTPFn(ctx, phone, code)
}

func (f *fakeRepository) FindOrCreateUser(ctx context.Context, phone, name string) (string, error) {
	return f.findOrCreateUserFn(ctx, phone, name)
}

func TestServiceSendOTPSuccess(t *testing.T) {
	repo := &fakeRepository{
		saveOTPFn: func(ctx context.Context, phone, code string, expiresAt time.Time) error {
			if phone != "+79990000000" {
				t.Fatalf("expected phone +79990000000, got %s", phone)
			}

			if len(code) != 6 {
				t.Fatalf("expected 6 digit code, got %s", code)
			}

			if time.Until(expiresAt) <= 0 {
				t.Fatalf("expected expires_at in future")
			}

			return nil
		},
	}

	service := NewService(repo, "test-secret")

	code, err := service.SendOTP(context.Background(), "+79990000000")
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	if len(code) != 6 {
		t.Fatalf("expected 6 digit code, got %s", code)
	}
}

func TestServiceSendOTPRepositoryError(t *testing.T) {
	expectedErr := errors.New("db error")

	repo := &fakeRepository{
		saveOTPFn: func(ctx context.Context, phone, code string, expiresAt time.Time) error {
			return expectedErr
		},
	}

	service := NewService(repo, "test-secret")

	_, err := service.SendOTP(context.Background(), "+79990000000")
	if !errors.Is(err, expectedErr) {
		t.Fatalf("expected db error, got %v", err)
	}
}

func TestServiceVerifyOTPSuccess(t *testing.T) {
	repo := &fakeRepository{
		verifyOTPFn: func(ctx context.Context, phone, code string) (bool, error) {
			if phone != "+79990000000" {
				t.Fatalf("expected phone +79990000000, got %s", phone)
			}

			if code != "123456" {
				t.Fatalf("expected code 123456, got %s", code)
			}

			return true, nil
		},
		findOrCreateUserFn: func(ctx context.Context, phone, name string) (string, error) {
			if phone != "+79990000000" {
				t.Fatalf("expected phone +79990000000, got %s", phone)
			}

			if name != "Daniil" {
				t.Fatalf("expected name Daniil, got %s", name)
			}

			return "user-1", nil
		},
	}

	service := NewService(repo, "test-secret")

	resp, err := service.VerifyOTP(context.Background(), "+79990000000", "123456", "Daniil")
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	if resp.UserID != "user-1" {
		t.Fatalf("expected user_id user-1, got %s", resp.UserID)
	}

	if resp.AccessToken == "" {
		t.Fatalf("expected access token")
	}

	claims := jwt.MapClaims{}

	token, err := jwt.ParseWithClaims(resp.AccessToken, claims, func(token *jwt.Token) (any, error) {
		return []byte("test-secret"), nil
	})
	if err != nil {
		t.Fatalf("failed to parse jwt: %v", err)
	}

	if !token.Valid {
		t.Fatalf("expected valid token")
	}

	if claims["user_id"] != "user-1" {
		t.Fatalf("expected user_id claim user-1, got %v", claims["user_id"])
	}
}

func TestServiceVerifyOTPInvalidCode(t *testing.T) {
	repo := &fakeRepository{
		verifyOTPFn: func(ctx context.Context, phone, code string) (bool, error) {
			return false, nil
		},
	}

	service := NewService(repo, "test-secret")

	_, err := service.VerifyOTP(context.Background(), "+79990000000", "000000", "Daniil")
	if !errors.Is(err, ErrInvalidOTP) {
		t.Fatalf("expected ErrInvalidOTP, got %v", err)
	}
}

func TestServiceVerifyOTPRepositoryError(t *testing.T) {
	expectedErr := errors.New("db error")

	repo := &fakeRepository{
		verifyOTPFn: func(ctx context.Context, phone, code string) (bool, error) {
			return false, expectedErr
		},
	}

	service := NewService(repo, "test-secret")

	_, err := service.VerifyOTP(context.Background(), "+79990000000", "123456", "Daniil")
	if !errors.Is(err, expectedErr) {
		t.Fatalf("expected db error, got %v", err)
	}
}

func TestServiceVerifyOTPFindOrCreateUserError(t *testing.T) {
	expectedErr := errors.New("create user error")

	repo := &fakeRepository{
		verifyOTPFn: func(ctx context.Context, phone, code string) (bool, error) {
			return true, nil
		},
		findOrCreateUserFn: func(ctx context.Context, phone, name string) (string, error) {
			return "", expectedErr
		},
	}

	service := NewService(repo, "test-secret")

	_, err := service.VerifyOTP(context.Background(), "+79990000000", "123456", "Daniil")
	if !errors.Is(err, expectedErr) {
		t.Fatalf("expected create user error, got %v", err)
	}
}
