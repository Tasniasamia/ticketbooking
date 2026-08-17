package user

import (
	"errors"

	"gorm.io/gorm"
	"ticketBooking/internal/utils/query"
)

type Repository interface {
	CreateUser(user *User) error
	GetUserByEmail(email string) (*User, error)
	GetUserById(userId uint) (*User, error)
   UpdateUserFields(userId uint, updates map[string]interface{}) (*User, error)
	DeleteUser(userId uint) error
	MarkAsVerified(email string) error
	 GetAll(p query.Params) ([]*User, int64, error)
}

type repository struct {
	db *gorm.DB;
}


func NewUserRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) CreateUser(user *User) error {
	 result:=r.db.Create(user);

  if result.Error != nil {
		 if errors.Is(result.Error, gorm.ErrRecordNotFound) {
    return errors.New("User with this email already exists");
  }
    return result.Error
  }
  return nil
}


func (r *repository) GetUserByEmail(email string) (*User, error) {
	var user User
	result := r.db.Where(&User{Email: email}).First(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	
	return &user, nil
}

func (r *repository) GetUserById(userId uint) (*User, error) {
	var user User

	result := r.db.First(&user, userId)

	if result.Error != nil {
		return nil, result.Error
	}

	return &user, nil
}


func (r *repository) UpdateUserFields(userId uint, updates map[string]interface{}) (*User, error) {
	var user User

	result := r.db.Model(&User{}).Where("id = ?", userId).Updates(updates)
	if result.Error != nil {
		return nil, result.Error
	}

	// আপডেটের পর fresh data আনছি
	if err := r.db.First(&user, userId).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *repository) DeleteUser(userId uint) error {
	return r.db.Delete(&User{}, userId).Error
}

func (r *repository) MarkAsVerified(email string) error {
	return r.db.Model(&User{}).Where("email = ?", email).Update("is_verified", true).Error
}

func (r *repository) GetAll(p query.Params) ([]*User, int64, error) {
	var users []*User
	var total int64
	db := r.db.Model(&User{})
	db = query.Apply(db, p, nil, []string{"name", "email"})
	if err := db.Count(&total).Error; err != nil { return nil, 0, err }
	db = query.Paginate(db, p)
	if err := db.Find(&users).Error; err != nil { return nil, 0, err }
	return users, total, nil
}