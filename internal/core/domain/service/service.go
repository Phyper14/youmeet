package service

import (
	"time"
	"github.com/google/uuid"
)

type Service struct {
	ID              uuid.UUID   `json:"id" gorm:"primaryKey;type:uuid"`
	Name            string      `json:"name" gorm:"not null"`
	Description     string      `json:"description"`
	Duration        int         `json:"duration" gorm:"not null"`
	Price           float64     `json:"price" gorm:"not null"`
	ProfessionalIDs []uuid.UUID `json:"professional_ids" gorm:"type:uuid[]"`
	CreatedAt       time.Time   `json:"created_at" gorm:"autoCreateTime"`
}