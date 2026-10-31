package testimonial
import ("ticketBooking/internal/utils/query"; "gorm.io/gorm")
type Repository interface {
	Create(t *Testimonial) error; Update(t *Testimonial) error; Delete(id uint) error
	GetByID(id uint) (*Testimonial, error); GetAll(p query.Params) ([]*Testimonial, int64, error)
	GetAllActive(p query.Params) ([]*Testimonial, int64, error)
}
type repository struct{ db *gorm.DB }
func NewRepository(db *gorm.DB) Repository { return &repository{db: db} }
func (r *repository) Create(t *Testimonial) error { return r.db.Create(t).Error }
func (r *repository) Update(t *Testimonial) error { return r.db.Save(t).Error }
func (r *repository) Delete(id uint) error { return r.db.Delete(&Testimonial{}, id).Error }
func (r *repository) GetByID(id uint) (*Testimonial, error) {
	var t Testimonial
	if err := r.db.First(&t, id).Error; err != nil { return nil, err }
	return &t, nil
}
func (r *repository) GetAll(p query.Params) ([]*Testimonial, int64, error) {
	var items []*Testimonial; var total int64
	db := query.Apply(r.db.Model(&Testimonial{}), p, []string{"name","role","text"}, nil,nil,nil,nil)
	if err := db.Count(&total).Error; err != nil { return nil,0,err }
	db = query.Paginate(db, p)
	if err := db.Find(&items).Error; err != nil { return nil,0,err }
	return items, total, nil
}
func (r *repository) GetAllActive(p query.Params) ([]*Testimonial, int64, error) {
	var items []*Testimonial; var total int64
	db := r.db.Model(&Testimonial{}).Where("is_active = ?", true)
	db = query.Apply(db, p, []string{"name","role","text"}, nil,nil,nil,nil)
	if err := db.Count(&total).Error; err != nil { return nil,0,err }
	if p.SortBy == "created_at" { db = db.Order("sort_order ASC, created_at DESC") }
	db = query.Paginate(db, p)
	if err := db.Find(&items).Error; err != nil { return nil,0,err }
	return items, total, nil
}
