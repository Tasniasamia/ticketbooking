package booking

import (
	"errors"
	"ticketBooking/internal/booking/dto"
	"ticketBooking/internal/utils/query"

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
    GetByManagerID(managerID uint) ([]Booking, error) 
	 GetAll(p query.Params) ([]*Booking, int64, error)
	Update(booking *Booking) error
	UpdateStatus(bookingID uint, status dto.BookingStatus) error
}

type repository struct{ db *gorm.DB }

func NewRepository(db *gorm.DB) Repository { return &repository{db: db} }

func (r *repository) Create(booking *Booking) error {
		if err := r.db.Create(booking).Error; err != nil {
		return err

	}
    return r.db.
		Preload("UserInfo").
		Preload("EventInfo").                
		Preload("EventInfo.Manager").       
		Preload("EventInfo.Category").
		First(booking, booking.ID).Error
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
	err := r.db.Where("user_id = ?", userID).Preload("EventInfo").
		Preload("EventInfo.Manager").
		Preload("EventInfo.Category").
		Preload("UserInfo").Order("created_at DESC").Find(&bookings).Error
	return bookings, err
}
func (r *repository) GetByManagerID(managerID uint) ([]Booking, error) {
	var bookings []Booking
err := r.db.
		Joins("EventInfo").                                    // EventInfo টেবিলের সাথে join
		Where("\"EventInfo\".manager_id = ?", managerID).      // অথবা event_infos.manager_id (টেবিলের নাম অনুযায়ী)
		Preload("EventInfo").
		Preload("EventInfo.Manager").
		Preload("EventInfo.Category").
		Preload("UserInfo").
		Find(&bookings).Error

	return bookings, err
}

func (r *repository) Update(booking *Booking) error {
	return r.db.Save(booking).Error
}

func (r *repository) UpdateStatus(bookingID uint, status dto.BookingStatus) error {
	return r.db.Model(&Booking{}).Where("id = ?", bookingID).Update("status", status).Error
}


func (r *repository) GetAll(p query.Params) ([]*Booking, int64, error) {
	var bookings []*Booking
	var total int64

	db := r.db.Model(&Booking{})
	db = query.Apply(db, p, []string{"booking_code", "created_at", "status","total_price"}, nil, nil,nil,nil)
if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	db = query.Paginate(db, p)

	err := db.
		Preload("EventInfo").
		Preload("EventInfo.Manager").
		Preload("EventInfo.Category").
		Preload("UserInfo").
		Find(&bookings).Error

	if err != nil {
		return nil, 0, err
	}

	return bookings, total, nil
  
}
// func (r *repository) GetAll(p query.Params, lang string) ([]*Event, int64, error) {
// 	var events []*Event
// 	var total int64

// 	db := r.db.Model(&Event{})
// 	db = query.Apply(db, p, nil, []string{"title", "description", "location"}, lang)

// 	if err := db.Count(&total).Error; err != nil {
// 		return nil, 0, err
// 	}

// 	db = query.Paginate(db, p)

// 	err := db.
// 		Preload("Manager").
// 		Preload("Category").
// 		Find(&events).Error

// 	if err != nil {
// 		return nil, 0, err
// 	}

// 	return events, total, nil
// }