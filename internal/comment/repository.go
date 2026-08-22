package comment

import (
	"gorm.io/gorm"
)

type Repository interface {
	Create(c *Comment) error
	GetByID(id uint) (*Comment, error)
	// Get all comments of a blog (flat list) — tree is built in service
	GetByBlogID(blogID uint) ([]*Comment, error)
	Update(c *Comment) error
	// soft delete flag
	SoftDelete(id uint) error
	// hard delete (admin)
	Delete(c *Comment) error
}

type repository struct{ db *gorm.DB }

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) Create(c *Comment) error {
	if err := r.db.Create(c).Error; err != nil {
		return err
	}
	return r.db.Preload("User").First(c, c.ID).Error
}

func (r *repository) GetByID(id uint) (*Comment, error) {
	var c Comment
	err := r.db.Preload("User").First(&c, id).Error
	if err != nil {
		return nil, err
	}
	return &c, nil
}

func (r *repository) GetByBlogID(blogID uint) ([]*Comment, error) {
	var comments []*Comment
	err := r.db.
		Where("blog_id = ?", blogID).
		Preload("User").
		Order("created_at ASC").
		Find(&comments).Error
	if err != nil {
		return nil, err
	}
	return comments, nil
}

func (r *repository) Update(c *Comment) error {
	return r.db.Save(c).Error
}

func (r *repository) SoftDelete(id uint) error {
	return r.db.Model(&Comment{}).Where("id = ?", id).Update("is_deleted", true).Error
}

func (r *repository) Delete(c *Comment) error {
	return r.db.Delete(c).Error
}
