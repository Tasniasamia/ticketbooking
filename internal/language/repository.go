package language

import "gorm.io/gorm"

type Repository interface {
	Create(lang *Language) error
	GetAll() ([]*Language, error)
	GetByID(id uint) (*Language, error)
	GetByCode(code string) (*Language, error)
	Update(lang *Language) error
	Delete(id uint) error
	ClearDefault() error
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

func (r *repository) GetAll() ([]*Language, error) {
	var list []*Language
	err := r.db.Order("is_default DESC, name ASC").Find(&list).Error
	return list, err
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