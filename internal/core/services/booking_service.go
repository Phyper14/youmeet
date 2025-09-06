package services

import (
	"context"
	"time"
	"github.com/google/uuid"
	"youmeet/internal/core/domain/appointment"
	"youmeet/internal/core/domain/service"
)

type BookingService struct {
	appointmentRepo appointment.Repository
	serviceRepo     service.Repository
}

func NewBookingService(appointmentRepo appointment.Repository, serviceRepo service.Repository) *BookingService {
	return &BookingService{
		appointmentRepo: appointmentRepo,
		serviceRepo:     serviceRepo,
	}
}

func (s *BookingService) BookAppointment(ctx context.Context, serviceID, clientID uuid.UUID, startTimeStr string) (*appointment.Appointment, error) {
	startTime, err := time.Parse(time.RFC3339, startTimeStr)
	if err != nil {
		return nil, err
	}

	appt := &appointment.Appointment{
		ID:        uuid.New(),
		ServiceID: serviceID,
		ClientID:  clientID,
		StartTime: startTime,
		Status:    "scheduled",
		CreatedAt: time.Now(),
	}

	err = s.appointmentRepo.CreateAppointment(ctx, appt)
	if err != nil {
		return nil, err
	}

	return appt, nil
}

func (s *BookingService) GetAppointments(ctx context.Context, clientID uuid.UUID) ([]*appointment.Appointment, error) {
	return s.appointmentRepo.ListAppointments(ctx, clientID)
}

func (s *BookingService) GetAppointmentsByProfessional(ctx context.Context, professionalID uuid.UUID) ([]*appointment.Appointment, error) {
	return s.appointmentRepo.ListByProfessional(ctx, professionalID)
}