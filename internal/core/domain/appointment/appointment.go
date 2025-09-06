package appointment

import (
	"time"
	"github.com/google/uuid"
)

type Appointment struct {
	ID             uuid.UUID `json:"id" gorm:"primaryKey;type:uuid"`
	ClientID       uuid.UUID `json:"client_id" gorm:"type:uuid;not null"`
	ProfessionalID uuid.UUID `json:"professional_id" gorm:"type:uuid;not null"`
	ServiceID      uuid.UUID `json:"service_id" gorm:"type:uuid;not null"`
	StartTime      time.Time `json:"start_time" gorm:"not null"`
	EndTime        time.Time `json:"end_time" gorm:"not null"`
	Status         string    `json:"status" gorm:"not null;default:'scheduled'"`
	CreatedAt      time.Time `json:"created_at" gorm:"autoCreateTime"`
}

type Availability struct {
	ID             uuid.UUID `json:"id" gorm:"primaryKey;type:uuid"`
	ProfessionalID uuid.UUID `json:"professional_id" gorm:"type:uuid;not null"`
	DayOfWeek      string    `json:"day_of_week" gorm:"not null"`
	StartTime      string    `json:"start_time" gorm:"not null"`
	EndTime        string    `json:"end_time" gorm:"not null"`
	CreatedAt      time.Time `json:"created_at" gorm:"autoCreateTime"`
}