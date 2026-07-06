package bookings

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	ErrSlotNotFound       = errors.New("slot not found")
	ErrSlotUnavailable    = errors.New("slot unavailable")
	ErrSlotFull           = errors.New("slot full")
	ErrRentalUnavailable  = errors.New("rental equipment unavailable")
	ErrInvalidRequest     = errors.New("invalid booking request")
	ErrBookingNotFound    = errors.New("booking not found")
	ErrBookingNotActive   = errors.New("booking is not active")
	ErrSlotAlreadyStarted = errors.New("slot already started")
)

type Repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Create(
	ctx context.Context,
	userID string,
	req CreateBookingRequest,
	idempotencyKey string,
) (*Booking, error) {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	var existing Booking

	err = tx.QueryRow(ctx, `
		SELECT id::text, user_id::text, slot_id::text, seats_count, rental_count,
		       price_total, currency, status, created_at
		FROM bookings
		WHERE user_id = $1
		  AND idempotency_key = $2
	`, userID, idempotencyKey).Scan(
		&existing.ID,
		&existing.UserID,
		&existing.SlotID,
		&existing.SeatsCount,
		&existing.RentalCount,
		&existing.PriceTotal,
		&existing.Currency,
		&existing.Status,
		&existing.CreatedAt,
	)

	if err == nil {
		return &existing, nil
	}

	var slot struct {
		ID                    string
		TotalSeats            int
		BookedSeats           int
		RentalEquipmentTotal  int
		RentalEquipmentBooked int
		Price                 int
		RentalPrice           int
		Currency              string
		Status                string
	}

	err = tx.QueryRow(ctx, `
		SELECT id::text, total_seats, booked_seats,
		       rental_equipment_total, rental_equipment_booked,
		       price, rental_price, currency, status
		FROM slots
		WHERE id = $1
		FOR UPDATE
	`, req.SlotID).Scan(
		&slot.ID,
		&slot.TotalSeats,
		&slot.BookedSeats,
		&slot.RentalEquipmentTotal,
		&slot.RentalEquipmentBooked,
		&slot.Price,
		&slot.RentalPrice,
		&slot.Currency,
		&slot.Status,
	)

	if err != nil {
		return nil, ErrSlotNotFound
	}

	if slot.Status != "active" {
		return nil, ErrSlotUnavailable
	}

	seatsCount := len(req.Equipment)
	if seatsCount < 1 || seatsCount > 3 {
		return nil, ErrInvalidRequest
	}

	rentalCount := 0
	for _, item := range req.Equipment {
		if item != "own" && item != "rental" {
			return nil, ErrInvalidRequest
		}
		if item == "rental" {
			rentalCount++
		}
	}

	freeSeats := slot.TotalSeats - slot.BookedSeats
	freeRental := slot.RentalEquipmentTotal - slot.RentalEquipmentBooked

	if seatsCount > freeSeats {
		return nil, ErrSlotFull
	}

	if rentalCount > freeRental {
		return nil, ErrRentalUnavailable
	}

	priceTotal := slot.Price*seatsCount + slot.RentalPrice*rentalCount

	var booking Booking

	err = tx.QueryRow(ctx, `
		INSERT INTO bookings (
			user_id, slot_id, seats_count, rental_count,
			price_total, currency, status, idempotency_key
		)
		VALUES ($1, $2, $3, $4, $5, $6, 'active', $7)
		RETURNING id::text, user_id::text, slot_id::text, seats_count, rental_count,
		          price_total, currency, status, created_at
	`, userID, req.SlotID, seatsCount, rentalCount, priceTotal, slot.Currency, idempotencyKey).Scan(
		&booking.ID,
		&booking.UserID,
		&booking.SlotID,
		&booking.SeatsCount,
		&booking.RentalCount,
		&booking.PriceTotal,
		&booking.Currency,
		&booking.Status,
		&booking.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	for _, equipmentType := range req.Equipment {
		_, err = tx.Exec(ctx, `
			INSERT INTO booking_seats (booking_id, equipment_type)
			VALUES ($1, $2)
		`, booking.ID, equipmentType)

		if err != nil {
			return nil, err
		}
	}

	_, err = tx.Exec(ctx, `
		UPDATE slots
		SET booked_seats = booked_seats + $2,
		    rental_equipment_booked = rental_equipment_booked + $3,
		    updated_at = NOW()
		WHERE id = $1
	`, req.SlotID, seatsCount, rentalCount)

	if err != nil {
		return nil, err
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}

	return &booking, nil
}

func (r *Repository) ListByUser(
	ctx context.Context,
	userID string,
	limit int,
	offset int,
) ([]BookingListItem, error) {
	rows, err := r.db.Query(ctx, `
		SELECT
			b.id::text,
			b.user_id::text,
			b.slot_id::text,
			b.seats_count,
			b.rental_count,
			b.price_total,
			b.currency,
			b.status,

			s.start_at,
			s.arrival_at,
			tc.name,
			tc.code,
			COALESCE(m.name, ''),
			s.address,
			s.meeting_point_name,
			s.status,
			s.cancel_reason,

			b.created_at
		FROM bookings b
		JOIN slots s ON s.id = b.slot_id
		JOIN track_configs tc ON tc.id = s.track_config_id
		LEFT JOIN marshals m ON m.id = s.marshal_id
		WHERE b.user_id = $1
		ORDER BY s.start_at ASC
		LIMIT $2 OFFSET $3
	`, userID, limit, offset)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := make([]BookingListItem, 0)

	for rows.Next() {
		var item BookingListItem

		err := rows.Scan(
			&item.ID,
			&item.UserID,
			&item.SlotID,
			&item.SeatsCount,
			&item.RentalCount,
			&item.PriceTotal,
			&item.Currency,
			&item.Status,

			&item.SlotStartAt,
			&item.SlotArrivalAt,
			&item.TrackConfigName,
			&item.TrackConfigCode,
			&item.MarshalName,
			&item.Address,
			&item.MeetingPointName,
			&item.SlotStatus,
			&item.SlotCancelReason,

			&item.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

		items = append(items, item)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return items, nil
}

func (r *Repository) Cancel(
	ctx context.Context,
	userID string,
	bookingID string,
) (*Booking, error) {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	var booking Booking
	var slotStartAt time.Time

	err = tx.QueryRow(ctx, `
		SELECT
			b.id::text,
			b.user_id::text,
			b.slot_id::text,
			b.seats_count,
			b.rental_count,
			b.price_total,
			b.currency,
			b.status,
			b.created_at,
			s.start_at
		FROM bookings b
		JOIN slots s ON s.id = b.slot_id
		WHERE b.id = $1
		  AND b.user_id = $2
		FOR UPDATE
	`, bookingID, userID).Scan(
		&booking.ID,
		&booking.UserID,
		&booking.SlotID,
		&booking.SeatsCount,
		&booking.RentalCount,
		&booking.PriceTotal,
		&booking.Currency,
		&booking.Status,
		&booking.CreatedAt,
		&slotStartAt,
	)

	if err != nil {
		return nil, ErrBookingNotFound
	}

	if booking.Status != "active" {
		return nil, ErrBookingNotActive
	}

	var newStatus string

	err = tx.QueryRow(ctx, `
		SELECT
			CASE
				WHEN NOW() >= $1 THEN 'started'
				WHEN NOW() <= $1 - INTERVAL '3 hours' THEN 'cancelled'
				ELSE 'late_cancel'
			END
	`, slotStartAt).Scan(&newStatus)

	if err != nil {
		return nil, err
	}

	if newStatus == "started" {
		return nil, ErrSlotAlreadyStarted
	}

	_, err = tx.Exec(ctx, `
		UPDATE bookings
		SET status = $1,
		    cancelled_at = NOW()
		WHERE id = $2
	`, newStatus, booking.ID)

	if err != nil {
		return nil, err
	}

	if newStatus == "cancelled" {
		_, err = tx.Exec(ctx, `
			UPDATE slots
			SET booked_seats = booked_seats - $2,
			    rental_equipment_booked = rental_equipment_booked - $3,
			    updated_at = NOW()
			WHERE id = $1
		`, booking.SlotID, booking.SeatsCount, booking.RentalCount)

		if err != nil {
			return nil, err
		}
	}

	booking.Status = newStatus

	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}

	return &booking, nil
}

func (r *Repository) GetByID(
	ctx context.Context,
	userID string,
	bookingID string,
) (*BookingDetails, error) {
	var item BookingDetails

	err := r.db.QueryRow(ctx, `
		SELECT
			b.id::text,
			b.user_id::text,
			b.slot_id::text,
			b.seats_count,
			b.rental_count,
			b.price_total,
			b.currency,
			b.status,
			b.created_at,

			s.start_at,
			s.arrival_at,
			tc.name,
			tc.code,
			COALESCE(m.name, ''),
			s.address,
			s.meeting_point_name,
			s.status,
			s.cancel_reason
		FROM bookings b
		JOIN slots s ON s.id = b.slot_id
		JOIN track_configs tc ON tc.id = s.track_config_id
		LEFT JOIN marshals m ON m.id = s.marshal_id
		WHERE b.id = $1
		  AND b.user_id = $2
	`, bookingID, userID).Scan(
		&item.ID,
		&item.UserID,
		&item.SlotID,
		&item.SeatsCount,
		&item.RentalCount,
		&item.PriceTotal,
		&item.Currency,
		&item.Status,
		&item.CreatedAt,

		&item.SlotStartAt,
		&item.SlotArrivalAt,
		&item.TrackConfigName,
		&item.TrackConfigCode,
		&item.MarshalName,
		&item.Address,
		&item.MeetingPointName,
		&item.SlotStatus,
		&item.SlotCancelReason,
	)

	if err != nil {
		return nil, ErrBookingNotFound
	}

	rows, err := r.db.Query(ctx, `
		SELECT equipment_type
		FROM booking_seats
		WHERE booking_id = $1
		ORDER BY created_at ASC
	`, bookingID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	item.Equipment = make([]string, 0)

	for rows.Next() {
		var equipment string

		if err := rows.Scan(&equipment); err != nil {
			return nil, err
		}

		item.Equipment = append(item.Equipment, equipment)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return &item, nil
}
