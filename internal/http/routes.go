package http

import (
	"net/http"

	"github.com/anshjamwal15/hsb_backend/internal/application/services"
	"github.com/anshjamwal15/hsb_backend/internal/config"
	"github.com/anshjamwal15/hsb_backend/internal/http/handlers"
	"github.com/anshjamwal15/hsb_backend/internal/http/middleware"
	"github.com/anshjamwal15/hsb_backend/internal/infrastructure/database"
	"github.com/anshjamwal15/hsb_backend/internal/infrastructure/payment"
	"github.com/anshjamwal15/hsb_backend/internal/infrastructure/repositories"
	"github.com/gin-gonic/gin"
)

// SetupRouter initializes and configures all routes for the application
func SetupRouter(db *database.MongoDB, cfg *config.Config) *gin.Engine {
	// Set Gin mode
	if cfg.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.Default()

	// Initialize dependencies
	deps := initializeDependencies(db, cfg)

	// Setup all routes
	setupDocumentation(router)
	setupHealthCheck(router)
	setupAuthRoutes(router, deps)
	setupProtectedRoutes(router, cfg.JWTSecret, deps)

	return router
}

// Dependencies holds all handlers
type Dependencies struct {
	auth         *handlers.AuthHandler
	user         *handlers.UserHandler
	doctor       *handlers.DoctorHandler
	booking      *handlers.BookingHandler
	timeSlot     *handlers.TimeSlotHandler
	clinic       *handlers.ClinicHandler
	diagnostic   *handlers.DiagnosticHandler
	period       *handlers.PeriodHandler
	pregnancy    *handlers.PregnancyHandler
	fsfi         *handlers.FSFIHandler
	mentalHealth *handlers.MentalHealthHandler
	pcos         *handlers.PCOSHandler
	symptoms     *handlers.SymptomsHandler
	weight       *handlers.WeightHandler
	journal      *handlers.JournalHandler
}

func initializeDependencies(db *database.MongoDB, cfg *config.Config) *Dependencies {
	// Repositories
	userRepo := repositories.NewUserRepository(db.Database)
	doctorRepo := repositories.NewDoctorRepository(db.Database)
	bookingRepo := repositories.NewBookingRepository(db.Database)
	clinicRepo := repositories.NewClinicRepository(db.Database)
	diagnosticRepo := repositories.NewDiagnosticRepository(db.Database)
	periodRepo := repositories.NewPeriodRepository(db.Database)
	pregnancyRepo := repositories.NewPregnancyRepository(db.Database)
	mentalHealthRepo := repositories.NewMentalHealthRepository(db.Database)
	pcosRepo := repositories.NewPCOSRepository(db.Database)
	symptomsRepo := repositories.NewSymptomsRepository(db.Database)
	weightRepo := repositories.NewWeightRepository(db.Database)
	journalRepo := repositories.NewJournalRepository(db.Database)

	// Payment client
	razorpayClient := payment.NewRazorpayClient(cfg.RazorpayKey, cfg.RazorpaySecret)

	// Services
	authService := services.NewAuthService(userRepo, cfg.JWTSecret)
	userService := services.NewUserService(userRepo)
	doctorService := services.NewDoctorService(doctorRepo)
	bookingService := services.NewBookingService(bookingRepo, doctorRepo, razorpayClient)
	clinicService := services.NewClinicService(clinicRepo)
	diagnosticService := services.NewDiagnosticService(diagnosticRepo)
	periodService := services.NewPeriodService(periodRepo)
	pregnancyService := services.NewPregnancyService(pregnancyRepo)
	fsfiService := services.NewFSFIService(mentalHealthRepo)
	mentalHealthService := services.NewMentalHealthService(mentalHealthRepo)
	pcosService := services.NewPCOSService(pcosRepo)
	symptomsService := services.NewSymptomsService(symptomsRepo)
	weightService := services.NewWeightService(weightRepo)
	journalService := services.NewJournalService(journalRepo)

	// Handlers
	return &Dependencies{
		auth:         handlers.NewAuthHandler(authService),
		user:         handlers.NewUserHandler(userService),
		doctor:       handlers.NewDoctorHandler(doctorService),
		booking:      handlers.NewBookingHandler(bookingService),
		timeSlot:     handlers.NewTimeSlotHandler(doctorService, bookingService),
		clinic:       handlers.NewClinicHandler(clinicService),
		diagnostic:   handlers.NewDiagnosticHandler(diagnosticService),
		period:       handlers.NewPeriodHandler(periodService),
		pregnancy:    handlers.NewPregnancyHandler(pregnancyService),
		fsfi:         handlers.NewFSFIHandler(fsfiService),
		mentalHealth: handlers.NewMentalHealthHandler(mentalHealthService),
		pcos:         handlers.NewPCOSHandler(pcosService),
		symptoms:     handlers.NewSymptomsHandler(symptomsService),
		weight:       handlers.NewWeightHandler(weightService),
		journal:      handlers.NewJournalHandler(journalService),
	}
}

func setupDocumentation(r *gin.Engine) {
	r.StaticFile("/swagger.yaml", "./swagger.yaml")
	r.GET("/swagger", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "/swagger-ui")
	})
	r.GET("/swagger-ui", func(c *gin.Context) {
		html := `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Women's Health API Documentation</title>
    <link rel="stylesheet" type="text/css" href="https://unpkg.com/swagger-ui-dist@5.10.0/swagger-ui.css">
    <style>body { margin: 0; padding: 0; }</style>
</head>
<body>
    <div id="swagger-ui"></div>
    <script src="https://unpkg.com/swagger-ui-dist@5.10.0/swagger-ui-bundle.js"></script>
    <script src="https://unpkg.com/swagger-ui-dist@5.10.0/swagger-ui-standalone-preset.js"></script>
    <script>
        window.onload = function() {
            window.ui = SwaggerUIBundle({
                url: "/swagger.yaml",
                dom_id: '#swagger-ui',
                deepLinking: true,
                presets: [SwaggerUIBundle.presets.apis, SwaggerUIStandalonePreset],
                plugins: [SwaggerUIBundle.plugins.DownloadUrl],
                layout: "StandaloneLayout"
            });
        };
    </script>
</body>
</html>`
		c.Header("Content-Type", "text/html; charset=utf-8")
		c.String(http.StatusOK, html)
	})
}

func setupHealthCheck(r *gin.Engine) {
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "healthy",
			"service": "Women's Health API",
		})
	})
}

func setupAuthRoutes(r *gin.Engine, deps *Dependencies) {
	r.POST("/user/register", deps.auth.Register)
	r.POST("/user/login", deps.auth.Login)
	r.POST("/user/forgot-password", deps.auth.ForgotPassword)
	r.POST("/user/verify-otp", deps.auth.VerifyOTP)
	r.POST("/user/reset-password", deps.auth.ResetPassword)
}

func setupProtectedRoutes(r *gin.Engine, jwtSecret string, deps *Dependencies) {
	api := r.Group("/api")
	api.Use(middleware.AuthMiddleware(jwtSecret))

	// User Profile
	users := api.Group("/users")
	{
		users.GET("/me", deps.user.GetProfile)
		users.PUT("/me", deps.user.UpdateProfile)
	}
	api.POST("/user/change-password", deps.auth.ChangePassword)

	// Doctors
	api.GET("/doctors", deps.doctor.GetDoctors)
	api.GET("/doctors/:doctorId", deps.doctor.GetDoctorByID)

	// Bookings & Sessions
	api.GET("/time-slots", deps.timeSlot.GetAvailableTimeSlots)
	api.POST("/bookings", deps.booking.CreateBooking)
	api.POST("/bookings/verify", deps.booking.VerifyPayment)
	api.GET("/bookings/my-with-doctors", deps.booking.GetUserBookings)
	api.GET("/sessions/active", deps.booking.GetActiveBookings)

	// Clinics
	api.GET("/clinics", deps.clinic.GetClinics)
	api.POST("/clinic-bookings", deps.clinic.CreateBooking)
	api.GET("/clinic-bookings/my-bookings", deps.clinic.GetMyBookings)
	api.POST("/clinic-bookings/verify-payment", deps.clinic.VerifyPayment)

	// Diagnostics
	api.GET("/public/diagnostics", deps.diagnostic.GetDiagnostics)
	api.GET("/diagnosticsUsers", deps.diagnostic.GetDiagnosticsUsers)
	api.POST("/diagnostics-bookings", deps.diagnostic.CreateBooking)
	api.POST("/diagnostics-bookings/verify-payment", deps.diagnostic.VerifyPayment)

	// Health Tracking
	setupHealthTracking(api, deps)

	// Journals
	journals := api.Group("/journals")
	{
		journals.GET("", deps.journal.GetJournals)
		journals.POST("", deps.journal.CreateJournal)
		journals.GET("/:journalId", deps.journal.GetJournalByID)
		journals.PUT("/:journalId", deps.journal.UpdateJournal)
		journals.DELETE("/:journalId", deps.journal.DeleteJournal)
		journals.GET("/user/:userId", deps.journal.GetJournalsByUserID)
	}
}

func setupHealthTracking(api *gin.RouterGroup, deps *Dependencies) {
	// Period Tracker
	period := api.Group("/period")
	{
		period.GET("", deps.period.GetPeriodCycle)
		period.POST("", deps.period.AddPeriodCycle)
		period.DELETE("", deps.period.ResetPeriodTracker)
	}
	// Legacy routes
	api.GET("/period-cycle", deps.period.GetPeriodCycle)
	api.POST("/period-cycle", deps.period.AddPeriodCycle)
	api.DELETE("/period-cycle/reset", deps.period.ResetPeriodTracker)

	// Pregnancy Tracker
	pregnancy := api.Group("/pregnancy")
	{
		pregnancy.GET("", deps.pregnancy.GetPregnancyData)
		pregnancy.POST("", deps.pregnancy.AddPregnancyEntry)
	}
	// Legacy routes
	api.GET("/pregnancy-tracker", deps.pregnancy.GetPregnancyData)
	api.POST("/pregnancy-tracker", deps.pregnancy.AddPregnancyEntry)

	// FSFI (Sexual Wellness)
	fsfi := api.Group("/fsfi")
	{
		fsfi.GET("", deps.fsfi.GetTest)
		fsfi.POST("", deps.fsfi.SubmitTest)
		fsfi.GET("/results", deps.fsfi.GetMyResults)
		// Legacy routes
		fsfi.GET("/test", deps.fsfi.GetTest)
		fsfi.POST("/submit", deps.fsfi.SubmitTest)
		fsfi.GET("/my-results", deps.fsfi.GetMyResults)
	}

	// Mental Health
	mentalHealth := api.Group("/mental-health")
	{
		mentalHealth.GET("", deps.mentalHealth.GetTests)
		mentalHealth.POST("/submit", deps.mentalHealth.SubmitTestResults)
		mentalHealth.GET("/results", deps.mentalHealth.GetTestResults)
	}
	// Legacy routes
	api.GET("/tests", deps.mentalHealth.GetTests)
	api.GET("/tests/:testName", deps.mentalHealth.GetTestByName)
	api.POST("/test-results", deps.mentalHealth.SubmitTestResults)
	api.GET("/test-results", deps.mentalHealth.GetTestResults)

	// PCOS Assessment
	pcos := api.Group("/pcos")
	{
		pcos.GET("", deps.pcos.GetQuestions)
		pcos.POST("", deps.pcos.SubmitAssessment)
		pcos.GET("/history", deps.pcos.GetHistory)
		pcos.GET("/latest", deps.pcos.GetLatestAssessment)
	}
	// Legacy routes
	api.GET("/pcos-assessment/questions", deps.pcos.GetQuestions)
	api.POST("/pcos-assessment/submit", deps.pcos.SubmitAssessment)
	api.GET("/pcos-assessment/history", deps.pcos.GetHistory)
	api.GET("/pcos-assessment", deps.pcos.GetLatestAssessment)

	// Symptoms Tracking
	symptoms := api.Group("/symptoms")
	{
		symptoms.POST("", deps.symptoms.SubmitTracking)
		symptoms.GET("", deps.symptoms.GetTrackingHistory)
	}
	// Legacy routes
	api.POST("/symptoms-tracking", deps.symptoms.SubmitTracking)
	api.GET("/symptoms-tracking", deps.symptoms.GetTrackingHistory)

	// Weight & Metabolic Wellness
	api.GET("/weight-metabolic-wellness", deps.weight.GetData)
	api.POST("/weight-metabolic-wellness", deps.weight.AddEntry)
}
