package settings

import ("errors"; "gorm.io/gorm")

type Repository interface {
	Get() (*Setting, error)
	Create(s *Setting) error
	Update(s *Setting) error
	Upsert(s *Setting) error
}
type repository struct{ db *gorm.DB }
func NewRepository(db *gorm.DB) Repository { return &repository{db: db} }
func (r *repository) Get() (*Setting, error) {
	var s Setting
	err := r.db.Order("id ASC").First(&s).Error
	if err != nil { return nil, err }
	return &s, nil
}
func (r *repository) Create(s *Setting) error { return r.db.Create(s).Error }
func (r *repository) Update(s *Setting) error { return r.db.Save(s).Error }
func (r *repository) Upsert(s *Setting) error {
	ex, err := r.Get()
	if errors.Is(err, gorm.ErrRecordNotFound) { return r.Create(s) }
	if err != nil { return err }
	s.ID = ex.ID; s.CreatedAt = ex.CreatedAt
	return r.db.Save(s).Error
}
