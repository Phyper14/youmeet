package repositories

import (
	"context"
	"github.com/google/uuid"
	"youmeet/internal/core/domain/appointment"
)

type AppointmentRepository struct {
	db DBClient
}

func NewAppointmentRepository(db DBClient) *AppointmentRepository {
	return &AppointmentRepository{db: db}
}

func (r *AppointmentRepository) CreateAppointment(ctx context.Context, appt *appointment.Appointment) error {
	return r.db.Create(appt)
}

func (r *AppointmentRepository) GetAppointmentByID(ctx context.Context, id uuid.UUID) (*appointment.Appointment, error) {
	var appt appointment.Appointment
	err := r.db.First(&appt, "id = ?", id)
	return &appt, err
}

func (r *AppointmentRepository) ListAppointments(ctx context.Context, clientID uuid.UUID) ([]*appointment.Appointment, error) {
	var appointments []*appointment.Appointment
	err := r.db.Find(&appointments, "client_id = ?", clientID)
	return appointments, err
}

func (r *AppointmentRepository) ListByProfessional(ctx context.Context, professionalID uuid.UUID) ([]*appointment.Appointment, error) {
	var appointments []*appointment.Appointment
	err := r.db.Find(&appointments, "professional_id = ?", professionalID)
	return appointments, err
}

type AvailabilityRepository struct {
	db DBClient
}

func NewAvailabilityRepository(db DBClient) *AvailabilityRepository {
	return &AvailabilityRepository{db: db}
}

func (r *AvailabilityRepository) CreateAvailability(ctx context.Context, availability *appointment.Availability) error {
	return r.db.Create(availability)
}

func (r *AvailabilityRepository) GetByProfessional(ctx context.Context, professionalID uuid.UUID) ([]*appointment.Availability, error) {
	var availabilities []*appointment.Availability
	err := r.db.Find(&availabilities, "professional_id = ?", professionalID)
	return availabilities, err
}