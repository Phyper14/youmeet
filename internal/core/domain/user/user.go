package user

import (
	"time"
	"github.com/google/uuid"
)

type User struct {
	ID           uuid.UUID `json:"id" gorm:"primaryKey;type:uuid"`
	Name         string    `json:"name" gorm:"not null"`
	Email        string    `json:"email" gorm:"uniqueIndex;not null"`
	PasswordHash string    `json:"-" gorm:"not null"`
	Role         string    `json:"role" gorm:"not null"`
	CreatedAt    time.Time `json:"created_at" gorm:"autoCreateTime"`
}

type Company struct {
	ID        uuid.UUID `json:"id" gorm:"primaryKey;type:uuid"`
	UserID    uuid.UUID `json:"user_id" gorm:"type:uuid;not null"`
	Name      string    `json:"name" gorm:"not null"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	User      User      `gorm:"foreignKey:UserID"`
}

type Professional struct {
	ID        uuid.UUID  `json:"id" gorm:"primaryKey;type:uuid"`
	UserID    uuid.UUID  `json:"user_id" gorm:"type:uuid;not null"`
	Name      string     `json:"name" gorm:"not null"`
	CompanyID *uuid.UUID `json:"company_id,omitempty" gorm:"type:uuid"`
	CreatedAt time.Time  `json:"created_at" gorm:"autoCreateTime"`
	User      User       `gorm:"foreignKey:UserID"`
	Company   *Company   `gorm:"foreignKey:CompanyID"`
}