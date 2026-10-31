package faq
import (
	"ticketBooking/internal/utils/query"
	"gorm.io/gorm"
)
type Repository interface {
	Create(f *FAQ) error
	Update(f *FAQ) error
	Delete(id uint) error
	GetByID(id uint) (*FAQ, error)
	GetAll(p query.Params) ([]*FAQ, int64, error)
	GetAllActive(p query.Params) ([]*FAQ, int64, error)
}
type repository struct{ db *gorm.DB }
func NewRepository(db *gorm.DB) Repository { return &repository{db: db} }
func (r *repository) Create(f *FAQ) error { return r.db.Create(f).Error }
func (r *repository) Update(f *FAQ) error { return r.db.Save(f).Error }
func (r *repository) Delete(id uint) error { return r.db.Delete(&FAQ{}, id).Error }
func (r *repository) GetByID(id uint) (*FAQ, error) {
	var f FAQ
	if err := r.db.First(&f, id).Error; err != nil { return nil, err }
	return &f, nil
}
func (r *repository) GetAll(p query.Params) ([]*FAQ, int64, error) {
	var items []*FAQ; var total int64
	db := query.Apply(r.db.Model(&FAQ{}), p, []string{"question", "answer"}, nil, nil, nil, nil)
	if err := db.Count(&total).Error; err != nil { return nil, 0, err }
	db = query.Paginate(db, p)
	if err := db.Find(&items).Error; err != nil { return nil, 0, err }
	return items, total, nil
}
func (r *repository) GetAllActive(p query.Params) ([]*FAQ, int64, error) {
	var items []*FAQ; var total int64
	db := r.db.Model(&FAQ{}).Where("is_active = ?", true)
	db = query.Apply(db, p, []string{"question", "answer"}, nil, nil, nil, nil)
	if err := db.Count(&total).Error; err != nil { return nil, 0, err }
	if p.SortBy == "created_at" { db = db.Order("sort_order ASC, created_at DESC") }
	db = query.Paginate(db, p)
	if err := db.Find(&items).Error; err != nil { return nil, 0, err }
	return items, total, nil
}
