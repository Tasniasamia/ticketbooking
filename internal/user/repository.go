package user

import (
	"errors"
    "gorm.io/gorm"
)

type Repository interface {
	CreateUser(user *User) error
	GetUserByEmail(email string) (*User, error)
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


