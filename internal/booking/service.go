package booking;

import (
	"ticketBooking/internal/booking/dto"
	"ticketBooking/internal/event"
	"github.com/google/uuid"
)

type service struct {
	bookingRepo Repository
	eventRepo   event.Repository
}

func NewService(bookingRepo Repository, eventRepo event.Repository) *service {
	return &service{
		bookingRepo: bookingRepo,
		eventRepo:   eventRepo,
	}
}

func generateBookingCode() string {
	return "GT-" + uuid.New().String()
}

func (s *service) CreateBooking(userID uint, req dto.CreateRequest) (*dto.Response, error) {
	eventData, err := s.eventRepo.GetByID(req.EventID)
	if err != nil {
		return nil, err
	}

	if eventData.AvailableTickets < req.Quantity {
		return nil, ErrNotEnoughTickets
	}

	booking := &Booking{
		UserID:      userID,
		EventID:     req.EventID,
		Quantity:    req.Quantity,
		Status:     dto.BookingConfirmed,
		TotalPrice:  float64(req.Quantity) * float64(eventData.Price),
		BookingCode: generateBookingCode(),
	}

	if err := s.bookingRepo.Create(booking); err != nil {
		return nil, err
	}

	// টিকেট কমিয়ে দাও
	eventData.AvailableTickets -= req.Quantity
	if err := s.eventRepo.Update(eventData); err != nil {
		return nil, err
	}

	return booking.ToResponse(), nil
}