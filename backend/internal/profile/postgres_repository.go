package profile

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var ErrProfileNotFound = errors.New("profile not found")

type PostgresRepository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) *PostgresRepository {
	return &PostgresRepository{db: db}
}

func (r *PostgresRepository) GetByID(ctx context.Context, userID string) (*UserProfile, error) {
	var user UserProfile

	err := r.db.QueryRow(ctx, `
		SELECT id::text, name, phone, created_at, updated_at
		FROM users
		WHERE id = $1
		  AND is_deleted = false
	`, userID).Scan(
		&user.ID,
		&user.Name,
		&user.Phone,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrProfileNotFound
		}
		return nil, err
	}

	return &user, nil
}

func (r *PostgresRepository) Update(ctx context.Context, userID string, req UpdateProfileRequest) (*UserProfile, error) {
	var user UserProfile

	err := r.db.QueryRow(ctx, `
		UPDATE users
		SET
			name = COALESCE($2, name),
			phone = COALESCE($3, phone),
			updated_at = NOW()
		WHERE id = $1
		  AND is_deleted = false
		RETURNING id::text, name, phone, created_at, updated_at
	`, userID, req.Name, req.Phone).Scan(
		&user.ID,
		&user.Name,
		&user.Phone,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrProfileNotFound
		}
		return nil, err
	}

	return &user, nil
}

func (r *PostgresRepository) Delete(ctx context.Context, userID string) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	rows, err := tx.Query(ctx, `
		SELECT slot_id, seats_count, rental_count
		FROM bookings
		WHERE user_id = $1
		  AND status = 'active'
		FOR UPDATE
	`, userID)
	if err != nil {
		return err
	}
	defer rows.Close()

	type activeBooking struct {
		SlotID      string
		SeatsCount  int
		RentalCount int
	}

	activeBookings := make([]activeBooking, 0)

	for rows.Next() {
		var item activeBooking

		if err := rows.Scan(&item.SlotID, &item.SeatsCount, &item.RentalCount); err != nil {
			return err
		}

		activeBookings = append(activeBookings, item)
	}

	if err := rows.Err(); err != nil {
		return err
	}

	for _, booking := range activeBookings {
		_, err = tx.Exec(ctx, `
			UPDATE slots
			SET booked_seats = booked_seats - $2,
			    rental_equipment_booked = rental_equipment_booked - $3,
			    updated_at = NOW()
			WHERE id = $1
		`, booking.SlotID, booking.SeatsCount, booking.RentalCount)

		if err != nil {
			return err
		}
	}

	_, err = tx.Exec(ctx, `
		UPDATE bookings
		SET status = 'cancelled',
		    cancelled_at = NOW()
		WHERE user_id = $1
		  AND status = 'active'
	`, userID)
	if err != nil {
		return err
	}

	result, err := tx.Exec(ctx, `
		UPDATE users
		SET
			name = NULL,
			phone = 'deleted_' || id::text,
			is_deleted = true,
			updated_at = NOW()
		WHERE id = $1
		  AND is_deleted = false
	`, userID)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return ErrProfileNotFound
	}

	return tx.Commit(ctx)
}
