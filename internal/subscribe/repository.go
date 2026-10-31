package subscribe
import ("ticketBooking/internal/utils/query"; "gorm.io/gorm")
type Repository interface {
	Create(s *Subscribe) error; GetByEmail(email string) (*Subscribe, error)
	GetAll(p query.Params) ([]*Subscribe, int64, error)
}
type repository struct{ db *gorm.DB }
func NewRepository(db *gorm.DB) Repository { return &repository{db: db} }
func (r *repository) Create(s *Subscribe) error { return r.db.Create(s).Error }
func (r *repository) GetByEmail(email string) (*Subscribe, error) {
	var s Subscribe
	if err := r.db.Where("email = ?", email).First(&s).Error; err != nil { return nil, err }
	return &s, nil
}
func (r *repository) GetAll(p query.Params) ([]*Subscribe, int64, error) {
	var items []*Subscribe; var total int64
	db := query.Apply(r.db.Model(&Subscribe{}), p, []string{"email"}, nil,nil,nil,nil)
	if err := db.Count(&total).Error; err != nil { return nil,0,err }
	db = query.Paginate(db, p)
	if err := db.Find(&items).Error; err != nil { return nil,0,err }
	return items, total, nil
}
