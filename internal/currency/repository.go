package currency

import (
	"errors"
	"strings"
	"ticketBooking/internal/utils/query"
	"gorm.io/gorm"
)

type Repository interface {
	GetByCode(code string) (*Currency, error)
	GetByCodeAnyStatus(code string) (*Currency, error)
	GetByID(id uint) (*Currency, error)
	GetDefault() (*Currency, error)
	GetAll(p query.Params) ([]*Currency, int64, error)
	 GetAllEnabled(p query.Params) ([]*Currency, int64, error)
	Create(c *Currency) error
	Update(c *Currency) error
	Delete(id uint) error
	ClearDefault() error
	Upsert(c *Currency) error
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) GetByCode(code string) (*Currency, error) {
	var c Currency
	err := r.db.Where("UPPER(code) = ? AND status = ?", strings.ToUpper(code), "enable").First(&c).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("currency not found or disabled")
		}
		return nil, err
	}
	return &c, nil
}

func (r *repository) GetByCodeAnyStatus(code string) (*Currency, error) {
	var c Currency
	err := r.db.Where("UPPER(code) = ?", strings.ToUpper(code)).First(&c).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("currency not found")
		}
		return nil, err
	}
	return &c, nil
}

func (r *repository) GetByID(id uint) (*Currency, error) {
	var c Currency
	err := r.db.First(&c, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("currency not found")
		}
		return nil, err
	}
	return &c, nil
}

func (r *repository) GetDefault() (*Currency, error) {
	var c Currency
	err := r.db.Where("is_default = ? AND status = ?", true, "enable").First(&c).Error
	if err != nil {
		return r.GetByCode("BDT")
	}
	return &c, nil
}




func (r *repository) GetAll(p query.Params) ([]*Currency, int64, error) {
	var lists []*Currency
	var total int64

	db := r.db.Model(&Currency{})

	
	// multi-lang fields
	stringFields := []string{"name", "code", "status"}
	db = query.Apply(db, p, stringFields, nil)

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	db = query.Paginate(db, p)
	if err := db.Find(&lists).Error; err != nil {
		return nil, 0, err
	}
	return lists, total, nil
}



func (r *repository) GetAllEnabled(p query.Params) ([]*Currency, int64, error) {
	var lists []*Currency
	var total int64

	db := r.db.Where("status = ?", "enable").Order("is_default DESC, code ASC").Find(&lists)

	
	// multi-lang fields
	stringFields := []string{"name", "code", "status"}
	db = query.Apply(db, p, stringFields, nil)

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	db = query.Paginate(db, p)
	if err := db.Find(&lists).Error; err != nil {
		return nil, 0, err
	}
	return lists, total, nil
}














// func (r *repository) GetAllEnabled() ([]Currency, error) {
// 	var list []Currency
// 	err := r.db.Where("status = ?", "enable").Order("is_default DESC, code ASC").Find(&list).Error
// 	return list, err
// }

func (r *repository) Create(c *Currency) error {
	return r.db.Create(c).Error
}

func (r *repository) Update(c *Currency) error {
	return r.db.Save(c).Error
}

func (r *repository) Delete(id uint) error {
	return r.db.Delete(&Currency{}, id).Error
}

func (r *repository) ClearDefault() error {
	return r.db.Model(&Currency{}).Where("is_default = ?", true).Update("is_default", false).Error
}

func (r *repository) Upsert(c *Currency) error {
	var existing Currency
	err := r.db.Where("UPPER(code) = ?", strings.ToUpper(c.Code)).First(&existing).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return r.db.Create(c).Error
	}
	if err != nil {
		return err
	}
	// seed: only fill missing fields, don't overwrite admin changes to rate
	c.ID = existing.ID
	if existing.RateToBDT > 0 && c.RateToBDT == existing.RateToBDT {
		// keep existing
	}
	c.IsDefault = existing.IsDefault // don't reset default on seed
	if existing.Status != "" {
		c.Status = existing.Status
	}
	return r.db.Model(&existing).Updates(map[string]interface{}{
		"name":   c.Name,
		"symbol": c.Symbol,
	}).Error
}
