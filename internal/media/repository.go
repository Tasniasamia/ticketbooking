package media

import (
	"ticketBooking/internal/media/dto"
	"ticketBooking/internal/utils/query"
	"gorm.io/gorm"
)

type Repository interface {
	Create(m *Media) error
	GetByID(id uint) (*Media, error)
	GetByImageID(imageID string) (*Media, error)
	List(p query.Params, filter dto.ListFilter) ([]*Media, int64, error)
	SoftDelete(id uint) error
	HardDelete(id uint) error
	Update(m *Media) error
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) Create(m *Media) error {
	return r.db.Create(m).Error
}

func (r *repository) GetByID(id uint) (*Media, error) {
	var m Media
	err := r.db.First(&m, id).Error
	if err != nil {
		return nil, err
	}
	return &m, nil
}

func (r *repository) GetByImageID(imageID string) (*Media, error) {
	var m Media
	err := r.db.Where("image_id = ?", imageID).First(&m).Error
	if err != nil {
		return nil, err
	}
	return &m, nil
}

func (r *repository) List(p query.Params, filter dto.ListFilter) ([]*Media, int64, error) {
	var list []*Media
	var total int64

	db := r.db.Model(&Media{})

	if filter.ModelName != "" {
		db = db.Where("model_name = ?", filter.ModelName)
	}
	if filter.ModelID != nil {
		db = db.Where("model_id = ?", *filter.ModelID)
	}
	if filter.Type != "" {
		db = db.Where("type = ?", filter.Type)
	}

	// generic search / sort / pagination from shared query helper
	db = query.Apply(db, p, nil, nil,nil,nil,nil)

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	db = query.Paginate(db, p)
	if err := db.Find(&list).Error; err != nil {
		return nil, 0, err
	}
	return list, total, nil
}

func (r *repository) SoftDelete(id uint) error {
	return r.db.Delete(&Media{}, id).Error // GORM soft-delete because of gorm.Model
}

func (r *repository) HardDelete(id uint) error {
	return r.db.Unscoped().Delete(&Media{}, id).Error
}

func (r *repository) Update(m *Media) error {
	return r.db.Save(m).Error
}
