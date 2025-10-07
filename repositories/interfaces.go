package repositories

import "temp/models"

// UserRepository defines the interface for user repository operations
type UserRepository interface {
	Migrate() error
	Create(u *models.User) error
	FindByEmail(email string) (*models.User, error)
	FindByID(id uint) (*models.User, error)
}
