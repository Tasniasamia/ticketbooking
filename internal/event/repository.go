package event

import (
	"errors"
	"ticketBooking/internal/utils/query"
	"gorm.io/gorm"
)

type Repository interface {
	Create(event *Event) error
	GetAll(p query.Params) ([]*Event, int64, error)
	GetByID(eventId uint) (*Event, error)
	Update(event *Event) error
	DecrementTickets(eventID uint, quantity int) error
	IncrementTickets(eventID uint, quantity int) error
}

type repository struct{ db *gorm.DB }

func NewRepository(db *gorm.DB) Repository { return &repository{db: db} }

func (r *repository) Create(event *Event) error { return r.db.Create(event).Error }

func (r *repository) GetAll(p query.Params) ([]*Event, int64, error) {
	var events []*Event
	var total int64
	db := r.db.Model(&Event{})
	db = query.Apply(db, p, nil, []string{"title", "description", "location"})
	if err := db.Count(&total).Error; err != nil { return nil, 0, err }
	db = query.Paginate(db, p)
	if err := db.Find(&events).Error; err != nil { return nil, 0, err }
	return events, total, nil
}

func (r *repository) GetByID(eventId uint) (*Event, error) {
	var event Event
	if err := r.db.First(&event, eventId).Error; err != nil { return nil, err }
	return &event, nil
}

func (r *repository) Update(event *Event) error { return r.db.Save(event).Error }

func (r *repository) DecrementTickets(eventID uint, quantity int) error {
	res := r.db.Exec(`UPDATE events SET available_tickets = available_tickets - ? WHERE id = ? AND available_tickets >= ?`, quantity, eventID, quantity)
	if res.Error != nil { return res.Error }
	if res.RowsAffected == 0 { return errors.New("Not enough tickets available") }
	return nil
}

func (r *repository) IncrementTickets(eventID uint, quantity int) error {
	return r.db.Exec(`UPDATE events SET available_tickets = available_tickets + ? WHERE id = ?`, quantity, eventID).Error
}
