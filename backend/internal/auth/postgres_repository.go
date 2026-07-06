package auth

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresRepository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) *PostgresRepository {
	return &PostgresRepository{db: db}
}

func (r *PostgresRepository) SaveOTP(ctx context.Context, phone, code string, expiresAt time.Time) error {
	_, err := r.db.Exec(ctx, `
		INSERT INTO otp_codes (phone, code, expires_at)
		VALUES ($1, $2, $3)
	`, phone, code, expiresAt)

	return err
}

func (r *PostgresRepository) VerifyOTP(ctx context.Context, phone, code string) (bool, error) {
	var id string

	err := r.db.QueryRow(ctx, `
		SELECT id::text
		FROM otp_codes
		WHERE phone = $1
		  AND code = $2
		  AND used_at IS NULL
		  AND expires_at > NOW()
		ORDER BY created_at DESC
		LIMIT 1
	`, phone, code).Scan(&id)

	if err != nil {
		return false, nil
	}

	_, err = r.db.Exec(ctx, `
		UPDATE otp_codes
		SET used_at = NOW()
		WHERE id = $1
	`, id)

	if err != nil {
		return false, err
	}

	return true, nil
}

func (r *PostgresRepository) FindOrCreateUser(ctx context.Context, phone, name string) (string, error) {
	var id string

	err := r.db.QueryRow(ctx, `
		SELECT id::text
		FROM users
		WHERE phone = $1
		  AND is_deleted = false
	`, phone).Scan(&id)

	if err == nil {
		if name != "" {
			if _, updateErr := r.db.Exec(ctx, `
				UPDATE users
				SET name = $1, updated_at = NOW()
				WHERE id = $2
			`, name, id); updateErr != nil {
				return "", updateErr
			}
		}

		return id, nil
	}

	if !errors.Is(err, pgx.ErrNoRows) {
		return "", err
	}

	err = r.db.QueryRow(ctx, `
		INSERT INTO users (phone, name)
		VALUES ($1, $2)
		RETURNING id::text
	`, phone, name).Scan(&id)

	return id, err
}
