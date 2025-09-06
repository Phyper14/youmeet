package appointment_handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"youmeet/internal/core/services"
)

type Handler struct {
	bookingService *services.BookingService
}

func NewHandler(bookingService *services.BookingService) *Handler {
	return &Handler{
		bookingService: bookingService,
	}
}

func (h *Handler) BookAppointment(c *gin.Context) {
	var req BookAppointmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	appointment, err := h.bookingService.BookAppointment(c.Request.Context(), req.ServiceID, req.ClientID, req.StartTime)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, appointment)
}

func (h *Handler) GetAppointments(c *gin.Context) {
	clientIDStr := c.Param("clientId")
	clientID, err := uuid.Parse(clientIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid client ID"})
		return
	}

	appointments, err := h.bookingService.GetAppointments(c.Request.Context(), clientID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, appointments)
}