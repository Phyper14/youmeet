package main

import (
	"log"
	"youmeet/internal/adapters/handlers/appointment_handler"
	"youmeet/internal/adapters/handlers/auth_handler"
	"youmeet/internal/adapters/repositories"
	"youmeet/internal/core/domain/appointment"
	"youmeet/internal/core/domain/service"
	"youmeet/internal/core/domain/user"
	"youmeet/internal/core/services"
	"youmeet/internal/infra/database"

	"github.com/gin-gonic/gin"
)

func main() {
	// Database connection (PostgreSQL or SQLite based on env)
	db, err := database.NewDBClient()
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Auto-migrate tables
	err = db.AutoMigrate(
		&user.User{},
		&user.Company{},
		&user.Professional{},
		&appointment.Appointment{},
		&appointment.Availability{},
		&service.Service{},
	)
	if err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	// Repositórios
	userRepo := repositories.NewUserRepository(db)
	companyRepo := repositories.NewCompanyRepository(db)
	profRepo := repositories.NewProfessionalRepository(db)
	appointmentRepo := repositories.NewAppointmentRepository(db)
	serviceRepo := repositories.NewServiceRepository(db)

	// Serviços
	authService := services.NewAuthService(userRepo, companyRepo, profRepo)
	bookingService := services.NewBookingService(appointmentRepo, serviceRepo)

	// Handlers
	authHandler := auth_handler.NewHandler(authService)
	appointmentHandler := appointment_handler.NewHandler(bookingService)

	r := gin.Default()

	// Rotas de autenticação
	auth := r.Group("/auth")
	{
		auth.POST("/register", authHandler.Register)
		auth.POST("/login", authHandler.Login)
	}

	// Rotas de agendamentos
	r.POST("/appointments", appointmentHandler.BookAppointment)
	r.GET("/appointments/:clientId", appointmentHandler.GetAppointments)

	log.Println("Server starting on :8080")
	r.Run(":8080")
}
