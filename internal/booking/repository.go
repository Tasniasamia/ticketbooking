package booking

import (
	"errors"

	"gorm.io/gorm"
)

var (
	ErrBookingNotFound         = errors.New("booking not found")
	ErrNotEnoughTickets        = errors.New("not enough tickets available")
	ErrBookingAlreadyCancelled = errors.New("booking already cancelled")
	ErrForbiddenBookingAccess  = errors.New("you do not own this booking")
)

type Repository interface {
	Create(booking *Booking) error
	GetByID(bookingID uint) (*Booking, error)
	GetByUserID(userID uint) ([]Booking, error)
	Update(booking *Booking) error
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) Create(booking *Booking) error {
	return r.db.Create(booking).Error
}

func (r *repository) GetByID(bookingID uint) (*Booking, error) {
	var b Booking
	err := r.db.First(&b, bookingID).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrBookingNotFound
	}
	return &b, err
}

func (r *repository) GetByUserID(userID uint) ([]Booking, error) {
	var bookings []Booking
	err := r.db.Where("user_id = ?", userID).Order("created_at DESC").Find(&bookings).Error
	return bookings, err
}

func (r *repository) Update(booking *Booking) error {
	return r.db.Save(booking).Error
}