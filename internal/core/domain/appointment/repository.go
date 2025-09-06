package appointment

import (
	"context"
	"github.com/google/uuid"
)

type Repository interface {
	CreateAppointment(ctx context.Context, appointment *Appointment) error
	GetAppointmentByID(ctx context.Context, id uuid.UUID) (*Appointment, error)
	ListAppointments(ctx context.Context, clientID uuid.UUID) ([]*Appointment, error)
	ListByProfessional(ctx context.Context, professionalID uuid.UUID) ([]*Appointment, error)
}

type AvailabilityRepository interface {
	CreateAvailability(ctx context.Context, availability *Availability) error
	GetByProfessional(ctx context.Context, professionalID uuid.UUID) ([]*Availability, error)
}