package event

import "gorm.io/gorm"

import (
	"ticketBooking/internal/utils/query"
)

type Repository interface {
	Create(event *Event) error
	 GetAll(p query.Params) ([]*Event, int64, error) 
	GetByID(eventId uint) (*Event, error)
	Update(event *Event) error
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) Create(event *Event) error {
	return r.db.Create(event).Error
}

func (r *repository) GetAll(p query.Params) ([]*Event, int64, error) {
	var events []*Event
	var total int64

	db := r.db.Model(&Event{})

	// multi-lang fields
	jsonbFields := []string{"title", "description", "location"}
	db = query.Apply(db, p, nil, jsonbFields)

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	db = query.Paginate(db, p)
	if err := db.Find(&events).Error; err != nil {
		return nil, 0, err
	}
	return events, total, nil
}


func (r *repository) GetByID(eventId uint) (*Event, error) {
	var event Event

	err := r.db.First(&event, eventId).Error
	if err != nil {

		return nil, err
	}

	return &event, nil
}

func (r *repository) Update(event *Event) error {
	return r.db.Save(event).Error
}