package bookings

import "time"

type CreateBookingRequest struct {
	SlotID    string   `json:"slot_id"`
	Equipment []string `json:"equipment"`
}

type Booking struct {
	ID          string    `json:"id"`
	UserID      string    `json:"user_id"`
	SlotID      string    `json:"slot_id"`
	SeatsCount  int       `json:"seats_count"`
	RentalCount int       `json:"rental_count"`
	PriceTotal  int       `json:"price_total"`
	Currency    string    `json:"currency"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
}

type BookingListItem struct {
	ID          string `json:"id"`
	UserID      string `json:"user_id"`
	SlotID      string `json:"slot_id"`
	SeatsCount  int    `json:"seats_count"`
	RentalCount int    `json:"rental_count"`
	PriceTotal  int    `json:"price_total"`
	Currency    string `json:"currency"`
	Status      string `json:"status"`

	SlotStartAt      time.Time `json:"slot_start_at"`
	SlotArrivalAt    time.Time `json:"slot_arrival_at"`
	TrackConfigName  string    `json:"track_config_name"`
	TrackConfigCode  string    `json:"track_config_code"`
	MarshalName      string    `json:"marshal_name"`
	Address          string    `json:"address"`
	MeetingPointName string    `json:"meeting_point_name"`
	SlotStatus       string    `json:"slot_status"`
	SlotCancelReason *string   `json:"slot_cancel_reason,omitempty"`

	CreatedAt time.Time `json:"created_at"`
}

type ListBookingsRequest struct {
	Limit  int
	Offset int
}

type BookingDetails struct {
	ID          string    `json:"id"`
	UserID      string    `json:"user_id"`
	SlotID      string    `json:"slot_id"`
	SeatsCount  int       `json:"seats_count"`
	RentalCount int       `json:"rental_count"`
	PriceTotal  int       `json:"price_total"`
	Currency    string    `json:"currency"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`

	Equipment []string `json:"equipment"`

	SlotStartAt      time.Time `json:"slot_start_at"`
	SlotArrivalAt    time.Time `json:"slot_arrival_at"`
	TrackConfigName  string    `json:"track_config_name"`
	TrackConfigCode  string    `json:"track_config_code"`
	MarshalName      string    `json:"marshal_name"`
	Address          string    `json:"address"`
	MeetingPointName string    `json:"meeting_point_name"`

	SlotStatus       string  `json:"slot_status"`
	SlotCancelReason *string `json:"slot_cancel_reason,omitempty"`
}
