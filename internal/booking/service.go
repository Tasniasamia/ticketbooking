package booking

import (
	"ticketBooking/internal/booking/dto"
	"ticketBooking/internal/event"

	"github.com/google/uuid"
)

type Service interface {
	CreateBooking(userID uint, req dto.CreateRequest) (*dto.Response, error)
	GetByID(id uint) (*dto.Response, error)
	GetByUserID(userID uint) ([]*dto.Response, error)
}

type service struct {
	bookingRepo Repository
	eventRepo   event.Repository
}

func NewService(bookingRepo Repository, eventRepo event.Repository) Service {
	return &service{bookingRepo: bookingRepo, eventRepo: eventRepo}
}

func generateBookingCode() string {
	return "GT-" + uuid.New().String()
}

func (s *service) CreateBooking(userID uint, req dto.CreateRequest) (*dto.Response, error) {
	if err := s.eventRepo.DecrementTickets(req.EventID, req.Quantity); err != nil {
		return nil, ErrNotEnoughTickets
	}
	eventData, err := s.eventRepo.GetByID(req.EventID)
	if err != nil {
		_ = s.eventRepo.IncrementTickets(req.EventID, req.Quantity)
		return nil, err
	}
	b := &Booking{
		UserID: userID, EventID: req.EventID, Quantity: req.Quantity,
		Status: dto.BookingConfirmed,
		TotalPrice: float64(req.Quantity) * float64(eventData.Price),
		BookingCode: generateBookingCode(),
	}
	if err := s.bookingRepo.Create(b); err != nil {
		_ = s.eventRepo.IncrementTickets(req.EventID, req.Quantity)
		return nil, err
	}
	return b.ToResponse(), nil
}

func (s *service) GetByID(id uint) (*dto.Response, error) {
	b, err := s.bookingRepo.GetByID(id)
	if err != nil {
		return nil, err
	}
	return b.ToResponse(), nil
}

func (s *service) GetByUserID(userID uint) ([]*dto.Response, error) {
	list, err := s.bookingRepo.GetByUserID(userID)
	if err != nil {
		return nil, err
	}
	out := make([]*dto.Response, 0, len(list))
	for i := range list {
		out = append(out, list[i].ToResponse())
	}
	return out, nil
}
