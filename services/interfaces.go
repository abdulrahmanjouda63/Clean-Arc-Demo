package services

import "temp/models"

// UserServiceInterface defines the interface for user service operations
type UserServiceInterface interface {
	Register(name, email, password string) (*models.User, error)
	Authenticate(email, password string) (string, *models.User, error)
}
