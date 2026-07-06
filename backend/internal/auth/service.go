package auth

import (
	"context"
	"errors"
	"math/rand"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var ErrInvalidOTP = errors.New("invalid otp")

type Repository interface {
	SaveOTP(ctx context.Context, phone, code string, expiresAt time.Time) error
	VerifyOTP(ctx context.Context, phone, code string) (bool, error)
	FindOrCreateUser(ctx context.Context, phone, name string) (string, error)
}

type Service struct {
	repo      Repository
	jwtSecret string
}

func NewService(repo Repository, jwtSecret string) *Service {
	return &Service{
		repo:      repo,
		jwtSecret: jwtSecret,
	}
}

func (s *Service) SendOTP(ctx context.Context, phone string) (string, error) {
	code := strconv.Itoa(100000 + rand.Intn(900000))
	expiresAt := time.Now().Add(5 * time.Minute)

	err := s.repo.SaveOTP(ctx, phone, code, expiresAt)
	if err != nil {
		return "", err
	}

	return code, nil
}

func (s *Service) VerifyOTP(ctx context.Context, phone, code, name string) (*VerifyOTPResponse, error) {
	ok, err := s.repo.VerifyOTP(ctx, phone, code)
	if err != nil {
		return nil, err
	}

	if !ok {
		return nil, ErrInvalidOTP
	}

	userID, err := s.repo.FindOrCreateUser(ctx, phone, name)
	if err != nil {
		return nil, err
	}

	token, err := s.generateAccessToken(userID)
	if err != nil {
		return nil, err
	}

	return &VerifyOTPResponse{
		AccessToken: token,
		UserID:      userID,
	}, nil
}

func (s *Service) generateAccessToken(userID string) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(15 * time.Minute).Unix(),
		"iat":     time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(s.jwtSecret))
}
