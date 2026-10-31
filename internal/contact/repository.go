package contact
import ("ticketBooking/internal/utils/query"; "gorm.io/gorm")
type Repository interface {
	Create(c *Contact) error; Update(c *Contact) error; GetByID(id uint) (*Contact, error)
	GetAll(p query.Params) ([]*Contact, int64, error)
}
type repository struct{ db *gorm.DB }
func NewRepository(db *gorm.DB) Repository { return &repository{db: db} }
func (r *repository) Create(c *Contact) error { return r.db.Create(c).Error }
func (r *repository) Update(c *Contact) error { return r.db.Save(c).Error }
func (r *repository) GetByID(id uint) (*Contact, error) {
	var c Contact
	if err := r.db.First(&c, id).Error; err != nil { return nil, err }
	return &c, nil
}
func (r *repository) GetAll(p query.Params) ([]*Contact, int64, error) {
	var items []*Contact; var total int64
	db := query.Apply(r.db.Model(&Contact{}), p, []string{"name","email","subject","message"}, nil,nil,nil,nil)
	if err := db.Count(&total).Error; err != nil { return nil,0,err }
	db = query.Paginate(db, p)
	if err := db.Find(&items).Error; err != nil { return nil,0,err }
	return items, total, nil
}
