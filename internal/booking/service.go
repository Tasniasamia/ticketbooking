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
	// 1. আগে atomic ভাবে টিকেট কমাও
	if err := s.eventRepo.DecrementTickets(req.EventID, req.Quantity); err != nil {
		return nil, err // not enough tickets বা DB error
	}

	// 2. event নিয়ে price বের করো
	eventData, err := s.eventRepo.GetByID(req.EventID)
	if err != nil {
		return nil, err
	}

	// 3. booking তৈরি
	b := &Booking{
		UserID:      userID,
		EventID:     req.EventID,
		Quantity:    req.Quantity,
		Status:      dto.BookingConfirmed,
		TotalPrice:  float64(req.Quantity) * float64(eventData.Price),
		BookingCode: generateBookingCode(),
	}

	if err := s.bookingRepo.Create(b); err != nil {
		// optional: টিকেট ফেরত দিতে পারো (compensate)
		return nil, err
	}

	return b.ToResponse(), nil
}