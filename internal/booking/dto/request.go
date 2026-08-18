package dto
type BookingStatus string

const (
	BookingPending   BookingStatus = "pending"
	BookingConfirmed BookingStatus = "confirmed"
	BookingCancelled BookingStatus = "cancelled"
)


type CreateRequest struct {
	EventID  uint `json:"event_id" validate:"required"`
	Quantity int  `json:"quantity" validate:"required,min=1"`
}
