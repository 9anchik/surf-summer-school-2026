package slots

import "time"

type Slot struct {
	ID                     string    `json:"id"`
	StartAt                time.Time `json:"start_at"`
	ArrivalAt              time.Time `json:"arrival_at"`
	TrackConfigName        string    `json:"track_config_name"`
	TrackConfigCode        string    `json:"track_config_code"`
	TrackConfigType        string    `json:"track_config_type"`
	TrackConfigDescription string    `json:"track_config_description"`
	MarshalID              *string   `json:"marshal_id,omitempty"`
	MarshalName            string    `json:"marshal_name"`
	TotalSeats             int       `json:"total_seats"`
	BookedSeats            int       `json:"booked_seats"`
	FreeSeats              int       `json:"free_seats"`
	FreeRentalEquipment    int       `json:"free_rental_equipment"`
	Price                  int       `json:"price"`
	RentalPrice            int       `json:"rental_price"`
	Currency               string    `json:"currency"`
	Address                string    `json:"address"`
	MeetingPointName       string    `json:"meeting_point_name"`
	MeetingPointLat        string    `json:"meeting_point_lat"`
	MeetingPointLng        string    `json:"meeting_point_lng"`
	Status                 string    `json:"status"`
	CancelReason           *string   `json:"cancel_reason,omitempty"`
}

type ListFilters struct {
	DateFrom      string
	DateTo        string
	TrackConfig   string
	MarshalID     string
	OnlyAvailable bool
}
