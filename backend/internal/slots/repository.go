package slots

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var ErrSlotNotFound = errors.New("slot not found")

type Repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{db: db}
}

func (r *Repository) List(ctx context.Context, filters ListFilters) ([]Slot, error) {
	query := `
		SELECT
			s.id::text,
			s.start_at,
			s.arrival_at,
			tc.name,
			tc.code,
			tc.type,
			tc.description,
			s.marshal_id::text,
			COALESCE(m.name, ''),
			s.total_seats,
			s.booked_seats,
			(s.total_seats - s.booked_seats) AS free_seats,
			(s.rental_equipment_total - s.rental_equipment_booked) AS free_rental_equipment,
			s.price,
			s.rental_price,
			s.currency,
			s.address,
			s.meeting_point_name,
			s.meeting_point_lat::text,
			s.meeting_point_lng::text,
			s.status,
			s.cancel_reason
		FROM slots s
		JOIN track_configs tc ON tc.id = s.track_config_id
		LEFT JOIN marshals m ON m.id = s.marshal_id
		WHERE s.start_at >= COALESCE(NULLIF($1, '')::timestamptz, NOW())
		  AND s.start_at <= COALESCE(NULLIF($2, '')::timestamptz, NOW() + INTERVAL '7 days')
		  AND ($3 = '' OR tc.code = $3)
		  AND ($4 = '' OR s.marshal_id::text = $4)
		  AND ($5 = false OR (s.total_seats - s.booked_seats) > 0)
		ORDER BY s.start_at ASC
	`

	rows, err := r.db.Query(
		ctx,
		query,
		filters.DateFrom,
		filters.DateTo,
		filters.TrackConfig,
		filters.MarshalID,
		filters.OnlyAvailable,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make([]Slot, 0)

	for rows.Next() {
		slot, err := scanSlot(rows)
		if err != nil {
			return nil, err
		}

		result = append(result, slot)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return result, nil
}

func (r *Repository) GetByID(ctx context.Context, id string) (*Slot, error) {
	query := `
		SELECT
			s.id::text,
			s.start_at,
			s.arrival_at,
			tc.name,
			tc.code,
			tc.type,
			tc.description,
			s.marshal_id::text,
			COALESCE(m.name, ''),
			s.total_seats,
			s.booked_seats,
			(s.total_seats - s.booked_seats) AS free_seats,
			(s.rental_equipment_total - s.rental_equipment_booked) AS free_rental_equipment,
			s.price,
			s.rental_price,
			s.currency,
			s.address,
			s.meeting_point_name,
			s.meeting_point_lat::text,
			s.meeting_point_lng::text,
			s.status,
			s.cancel_reason
		FROM slots s
		JOIN track_configs tc ON tc.id = s.track_config_id
		LEFT JOIN marshals m ON m.id = s.marshal_id
		WHERE s.id = $1
	`

	row := r.db.QueryRow(ctx, query, id)

	slot, err := scanSlot(row)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrSlotNotFound
		}
		return nil, err
	}

	return &slot, nil
}

type slotScanner interface {
	Scan(dest ...any) error
}

func scanSlot(scanner slotScanner) (Slot, error) {
	var slot Slot

	err := scanner.Scan(
		&slot.ID,
		&slot.StartAt,
		&slot.ArrivalAt,
		&slot.TrackConfigName,
		&slot.TrackConfigCode,
		&slot.TrackConfigType,
		&slot.TrackConfigDescription,
		&slot.MarshalID,
		&slot.MarshalName,
		&slot.TotalSeats,
		&slot.BookedSeats,
		&slot.FreeSeats,
		&slot.FreeRentalEquipment,
		&slot.Price,
		&slot.RentalPrice,
		&slot.Currency,
		&slot.Address,
		&slot.MeetingPointName,
		&slot.MeetingPointLat,
		&slot.MeetingPointLng,
		&slot.Status,
		&slot.CancelReason,
	)

	return slot, err
}
