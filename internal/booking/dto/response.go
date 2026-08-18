package dto



type Response struct {
	ID          uint          `json:"id"`
	UserID      uint          `json:"user_id"`
	EventID     uint          `json:"event_id"`
	Quantity    int           `json:"quantity"`
	TotalPrice  float64       `json:"total_price"`
	Status      BookingStatus `json:"status"`
	BookingCode string        `json:"booking_code"`
	CreatedAt   string        `json:"created_at"`
}
