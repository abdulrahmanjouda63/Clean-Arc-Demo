package models

import (
	"time"
)

// User represents a user in the system
type User struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Name      string    `json:"name" gorm:"size:100;not null" binding:"required,min=2,max=100"`
	Email     string    `json:"email" gorm:"uniqueIndex;size:100;not null" binding:"required,email"`
	Password  string    `json:"-" gorm:"not null"`
	CreatedAt time.Time `json:"created_at"`
}
