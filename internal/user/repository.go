package user

import (
	"errors"
    "gorm.io/gorm"
)

type Repository interface {
	CreateUser(user *User) error
	GetUserByEmail(email string) (*User, error)
	GetUserById(userId uint) (*User, error)
	UpdateUser(user *User) (*User, error)
	DeleteUser(userId uint) error
	MarkAsVerified(email string) error
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


func  (r *repository) UpdateUser(user *User) (*User, error){
  var UpdateUser *User;
  err:=r.db.Save(user)
  if err.Error != nil {
    return nil, err.Error
  }
	result := r.db.Where(&User{Email: user.Email}).First(&UpdateUser)

  if(result.Error != nil){
    return nil, result.Error
  }
  return UpdateUser, nil
}

func (r *repository) DeleteUser(userId uint) error {
	return r.db.Delete(&User{}, userId).Error
}

func (r *repository) MarkAsVerified(email string) error {
	return r.db.Model(&User{}).Where("email = ?", email).Update("is_verified", true).Error
}