package services

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"temp/global"
	"temp/models"
	"temp/repositories"
	"temp/utils"

	"github.com/golang-jwt/jwt/v5"
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
)

type UserService struct {
	repo        repositories.UserRepository
	jwtSecret   string
	expDuration time.Duration
}

// Ensure UserService implements UserServiceInterface interface
var _ UserServiceInterface = (*UserService)(nil)

func NewUserService(repo repositories.UserRepository, jwtSecret string, expHours int) *UserService {
	return &UserService{
		repo:        repo,
		jwtSecret:   jwtSecret,
		expDuration: time.Duration(expHours) * time.Hour,
	}
}

func (s *UserService) Register(name, email, password string) (*models.User, error) {
	// check if user already exists
	if _, err := s.repo.FindByEmail(email); err == nil {
		return nil, errors.New("user with this email already exists")
	}

	hash, err := utils.GenerateHash(password)
	if err != nil {
		return nil, err
	}

	u := &models.User{
		Name:     name,
		Email:    email,
		Password: hash,
	}

	if err := s.repo.Create(u); err != nil {
		return nil, err
	}

		// Store user registration in Redis for demonstration
		if global.Redis != nil {
			ctx := context.Background()
			redisKey := fmt.Sprintf("user:registered:%s", u.Email)
			if err := global.Redis.Set(ctx, redisKey, u.ID, 10*time.Minute).Err(); err != nil {
				// Log the error but don't fail registration if Redis fails
				log.Printf("Failed to store user registration in Redis: %v", err)
			}
		}

	return u, nil
}

func (s *UserService) Authenticate(email, password string) (string, *models.User, error) {
	u, err := s.repo.FindByEmail(email)
	if err != nil {
		return "", nil, ErrInvalidCredentials
	}

	if !utils.CompareHash(password, u.Password) {
		return "", nil, ErrInvalidCredentials
	}

	// create token
	claims := jwt.MapClaims{
		"sub": u.ID,
		"exp": time.Now().Add(s.expDuration).Unix(),
	}
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := t.SignedString([]byte(s.jwtSecret))
	if err != nil {
		return "", nil, err
	}

	// Store JWT token in Redis after successful authentication
	if global.Redis != nil {
		ctx := context.Background()
		redisKey := fmt.Sprintf("user:token:%d", u.ID)
		if err := global.Redis.Set(ctx, redisKey, signed, s.expDuration).Err(); err != nil {
			// Log the error but don't fail login if Redis fails
			log.Printf("Failed to store user token in Redis: %v", err)
		}
	}

	return signed, u, nil
}
