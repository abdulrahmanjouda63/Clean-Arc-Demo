package repositories

import (
	"errors"
	"temp/global"
	"temp/models"

	"gorm.io/gorm"
)

var ErrUserNotFound = errors.New("user not found")

// UserRepo handles DB operations for users
type UserRepo struct{}

// Ensure UserRepo implements UserRepository interface
var _ UserRepository = (*UserRepo)(nil)

func NewUserRepo() *UserRepo {
	return &UserRepo{}
}

func (r *UserRepo) Migrate() error {
	return global.DB.AutoMigrate(&models.User{})
}

func (r *UserRepo) Create(u *models.User) error {
	return global.DB.Create(u).Error
}

func (r *UserRepo) FindByEmail(email string) (*models.User, error) {
	var u models.User
	res := global.DB.Where("email = ?", email).First(&u)
	if res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, res.Error
	}
	return &u, nil
}

func (r *UserRepo) FindByID(id uint) (*models.User, error) {
	var u models.User
	res := global.DB.First(&u, id)
	if res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, res.Error
	}
	return &u, nil
}
