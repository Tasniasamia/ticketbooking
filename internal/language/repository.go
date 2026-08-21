package language

import (
	"gorm.io/gorm"
	"ticketBooking/internal/utils/query"
)

type Repository interface {
	Create(lang *Language) error
    GetAll(p query.Params) ([]*Language, int64, error) 
	GetByID(id uint) (*Language, error)
	GetByCode(code string) (*Language, error)
	Update(lang *Language) error
	Delete(id uint) error
	ClearDefault() error
	GetDefault() (*Language, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) Create(lang *Language) error {
	return r.db.Create(lang).Error
}

func(r *repository) GetAll(p query.Params) ([]*Language, int64, error) {
	var list []*Language;
	var total int64;

	db := r.db.Model(&Language{})

	// multi-lang fields
	searchFields := []string{"code", "name"}
	db = query.Apply(db, p, searchFields, nil,nil,nil,nil)

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	db = query.Paginate(db, p)
	if err := db.Find(&list).Error; err != nil {
		return nil, 0, err
	}
	return list, total, nil
}



func (r *repository) GetByID(id uint) (*Language, error) {
	var lang Language
	err := r.db.First(&lang, id).Error
	if err != nil {
		return nil, err
	}
	return &lang, nil
}

func (r *repository) GetByCode(code string) (*Language, error) {
	var lang Language
	err := r.db.Where("code = ?", code).First(&lang).Error
	if err != nil {
		return nil, err
	}
	return &lang, nil
}

func (r *repository) Update(lang *Language) error {
	return r.db.Save(lang).Error
}

func (r *repository) Delete(id uint) error {
	return r.db.Delete(&Language{}, id).Error
}

func (r *repository) ClearDefault() error {
	return r.db.Model(&Language{}).Where("is_default = ?", true).Update("is_default", false).Error
}
func (r *repository) GetDefault() (*Language, error) {
	var lang Language
	err := r.db.Where("is_default = ? AND is_active = ?", true, true).First(&lang).Error
	if err != nil {
		return nil, err
	}
	return &lang, nil
}
