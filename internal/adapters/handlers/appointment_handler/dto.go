package appointment_handler

import "github.com/google/uuid"

type BookAppointmentRequest struct {
	ServiceID uuid.UUID `json:"service_id"`
	ClientID  uuid.UUID `json:"client_id"`
	StartTime string    `json:"start_time"`
}