package blog

import (
	"errors"
	"ticketBooking/internal/utils/query"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Repository interface {
	Create(blog *Blog) error
	GetAll(p query.Params, lang string) ([]*Blog, int64, error)
	GetByID(blogID uint) (*Blog, error)
	GetByIDAdmin(blogID uint) (*Blog, error)
	Update(blog *Blog) error
	Delete(blog *Blog) error

	// Like — concurrent-safe
	ToggleLike(blogID, userID uint) (liked bool, likeCount int64, err error)
	IsLiked(blogID, userID uint) (bool, error)

	// Comment count — atomic
	IncrementCommentCount(blogID uint) error
	DecrementCommentCount(blogID uint) error
}

type repository struct{ db *gorm.DB }

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) Create(blog *Blog) error {
	if err := r.db.Create(blog).Error; err != nil {
		return err
	}
	return r.db.Preload("Author").First(blog, blog.ID).Error
}

func (r *repository) GetAll(p query.Params, lang string) ([]*Blog, int64, error) {
	var blogs []*Blog
	var total int64

	db := r.db.Model(&Blog{})
	db = query.Apply(db, p, nil, []string{"title", "short_description", "long_description"}, lang, nil, nil)

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	db = query.Paginate(db, p)

	err := db.Preload("Author").Find(&blogs).Error
	if err != nil {
		return nil, 0, err
	}

	return blogs, total, nil
}

func (r *repository) GetByID(blogID uint) (*Blog, error) {
	var blog Blog
	err := r.db.
		Where("status = ?", "approved").
		Preload("Author").
		First(&blog, blogID).Error
	if err != nil {
		return nil, err
	}
	return &blog, nil
}

func (r *repository) GetByIDAdmin(blogID uint) (*Blog, error) {
	var blog Blog
	err := r.db.Preload("Author").First(&blog, blogID).Error
	if err != nil {
		return nil, err
	}
	return &blog, nil
}

func (r *repository) Update(blog *Blog) error {
	return r.db.Save(blog).Error
}

func (r *repository) Delete(blog *Blog) error {
	return r.db.Delete(blog).Error
}

// ToggleLike — race-safe:
// 1) unique (blog_id, user_id) → double-like impossible
// 2) transaction + atomic UPDATE like_count
func (r *repository) ToggleLike(blogID, userID uint) (liked bool, likeCount int64, err error) {
	err = r.db.Transaction(func(tx *gorm.DB) error {
		// blog exists?
		var blog Blog
		if err := tx.Select("id", "like_count").First(&blog, blogID).Error; err != nil {
			return errors.New("blog not found")
		}

		var existing BlogLike
		findErr := tx.Where("blog_id = ? AND user_id = ?", blogID, userID).First(&existing).Error

		if findErr == nil {
			// already liked → unlike
			if err := tx.Delete(&existing).Error; err != nil {
				return err
			}
			// atomic decrement (never go below 0)
			if err := tx.Model(&Blog{}).Where("id = ?", blogID).
				Update("like_count", gorm.Expr("GREATEST(like_count - 1, 0)")).Error; err != nil {
				return err
			}
			liked = false
		} else if errors.Is(findErr, gorm.ErrRecordNotFound) {
			// not liked → like
			// ON CONFLICT DO NOTHING handles concurrent double-insert race
			like := BlogLike{BlogID: blogID, UserID: userID}
			res := tx.Clauses(clause.OnConflict{
				Columns:   []clause.Column{{Name: "blog_id"}, {Name: "user_id"}},
				DoNothing: true,
			}).Create(&like)
			if res.Error != nil {
				return res.Error
			}
			if res.RowsAffected > 0 {
				if err := tx.Model(&Blog{}).Where("id = ?", blogID).
					Update("like_count", gorm.Expr("like_count + 1")).Error; err != nil {
					return err
				}
				liked = true
			} else {
				// concurrent insert already happened — treat as liked
				liked = true
			}
		} else {
			return findErr
		}

		// fresh count
		if err := tx.Model(&Blog{}).Select("like_count").Where("id = ?", blogID).Scan(&likeCount).Error; err != nil {
			return err
		}
		return nil
	})
	return
}

func (r *repository) IsLiked(blogID, userID uint) (bool, error) {
	var count int64
	err := r.db.Model(&BlogLike{}).
		Where("blog_id = ? AND user_id = ?", blogID, userID).
		Count(&count).Error
	return count > 0, err
}

func (r *repository) IncrementCommentCount(blogID uint) error {
	return r.db.Model(&Blog{}).Where("id = ?", blogID).
		Update("comment_count", gorm.Expr("comment_count + 1")).Error
}

func (r *repository) DecrementCommentCount(blogID uint) error {
	return r.db.Model(&Blog{}).Where("id = ?", blogID).
		Update("comment_count", gorm.Expr("GREATEST(comment_count - 1, 0)")).Error
}
