package eventCategory;

import (

	"ticketBooking/internal/utils/query"
	"gorm.io/gorm"
)

type Repository interface {
	Create(event *EventCategory) error
	GetAll(p query.Params,lang string) ([]*EventCategory, int64, error)
	GetByID(eventId uint) (*EventCategory, error)
	Update(event *EventCategory) error
	Delete(eventId uint) error
	
}

type repository struct{ db *gorm.DB }

func NewRepository(db *gorm.DB) Repository { return &repository{db: db} }

func (r *repository) Create(event *EventCategory) error { return r.db.Create(event).Error }

func (r *repository) GetAll(p query.Params,lang string) ([]*EventCategory, int64, error) {
	var events []*EventCategory
	var total int64
	db := r.db.Model(&EventCategory{})
	db = query.Apply(db, p, nil, []string{"name", "description", "created_at"},lang)
	if err := db.Count(&total).Error; err != nil { return nil, 0, err }
	db = query.Paginate(db, p)
	if err := db.Find(&events).Error; err != nil { return nil, 0, err }
	return events, total, nil
}

func (r *repository) GetByID(eventId uint) (*EventCategory, error) {
	var event EventCategory
	if err := r.db.First(&event, eventId).Error; err != nil { return nil, err }
	return &event, nil
}

func (r *repository) Update(event *EventCategory) error { return r.db.Save(event).Error }


func (r *repository) Delete(eventId uint) error { return r.db.Delete(&EventCategory{}, eventId).Error }