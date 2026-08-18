package translation

import (
	"gorm.io/gorm"
	"ticketBooking/internal/utils/query"
)

type Repository interface {
	Create(t *Translation) error
	GetAll(p query.Params) ([]*Translation, int64, error)
	GetByKey(key string) (*Translation, error)
	GetByID(id uint) (*Translation, error)
	Update(t *Translation) error
	Delete(id uint) error
	DeleteByKey(key string) error
	RemoveLangCodeFromAll(code string) error
	GetAllSimple() ([]*Translation, error)   // ← add
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) Create(t *Translation) error {
	return r.db.Create(t).Error
}

func (r *repository) GetAll(p query.Params) ([]*Translation, int64, error) {
	var list []*Translation
	var total int64

	db := r.db.Model(&Translation{})

	searchFields := []string{"key"}
	db = query.Apply(db, p, searchFields, nil,nil)

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	db = query.Paginate(db, p)
	if err := db.Find(&list).Error; err != nil {
		return nil, 0, err
	}
	return list, total, nil
}
func (r *repository) GetAllSimple() ([]*Translation, error) {
	var list []*Translation
	if err := r.db.Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}
func (r *repository) GetByKey(key string) (*Translation, error) {
	var t Translation
	err := r.db.Where("key = ?", key).First(&t).Error
	if err != nil {
		return nil, err
	}
	return &t, nil
}

func (r *repository) GetByID(id uint) (*Translation, error) {
	var t Translation
	err := r.db.First(&t, id).Error
	if err != nil {
		return nil, err
	}
	return &t, nil
}

func (r *repository) Update(t *Translation) error {
	return r.db.Save(t).Error
}

func (r *repository) Delete(id uint) error {
	return r.db.Delete(&Translation{}, id).Error
}

func (r *repository) DeleteByKey(key string) error {
	return r.db.Where("key = ?", key).Delete(&Translation{}).Error
}

// Language delete হলে সব translation থেকে ওই code remove হয়
func (r *repository) RemoveLangCodeFromAll(code string) error {
	return r.db.Exec(
		`UPDATE translations SET lang_values = lang_values - ? WHERE lang_values ? ?`,
		code, code,
	).Error
}