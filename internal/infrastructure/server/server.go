package server

import (
	"context"
	"log"
	"net/http"

	"github.com/anshjamwal15/hsb_backend/internal/application/services"
	"github.com/anshjamwal15/hsb_backend/internal/config"
	"github.com/anshjamwal15/hsb_backend/internal/infrastructure/database"
	"github.com/anshjamwal15/hsb_backend/internal/infrastructure/payment"
	"github.com/anshjamwal15/hsb_backend/internal/infrastructure/repositories"
	"github.com/anshjamwal15/hsb_backend/internal/interfaces/http/handlers"
	"github.com/anshjamwal15/hsb_backend/internal/interfaces/http/middleware"
	"github.com/gin-gonic/gin"
)

type Server struct {
	httpServer *http.Server
	router     *gin.Engine
	db         *database.MongoDB
}

func NewServer(cfg *config.Config, db *database.MongoDB) *Server {
	// Set Gin mode
	if cfg.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	// Initialize Gin router
	router := gin.Default()

	// Setup routes
	setupRoutes(router, db, cfg)

	// Create HTTP server
	httpServer := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: router,
	}

	return &Server{
		httpServer: httpServer,
		router:     router,
		db:         db,
	}
}

func (s *Server) Start() error {
	log.Printf("Server starting on %s", s.httpServer.Addr)
	return s.httpServer.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}

func setupRoutes(router *gin.Engine, db *database.MongoDB, cfg *config.Config) {
	// Initialize repositories
	userRepo := repositories.NewUserRepository(db.Database)
	doctorRepo := repositories.NewDoctorRepository(db.Database)
	bookingRepo := repositories.NewBookingRepository(db.Database)
	journalRepo := repositories.NewJournalRepository(db.Database)
	periodRepo := repositories.NewPeriodRepository(db.Database)
	pregnancyRepo := repositories.NewPregnancyRepository(db.Database)
	symptomsRepo := repositories.NewSymptomsRepository(db.Database)
	weightRepo := repositories.NewWeightRepository(db.Database)
	pcosRepo := repositories.NewPCOSRepository(db.Database)
	mentalHealthRepo := repositories.NewMentalHealthRepository(db.Database)
	clinicRepo := repositories.NewClinicRepository(db.Database)
	diagnosticRepo := repositories.NewDiagnosticRepository(db.Database)

	// Initialize payment client
	razorpayClient := payment.NewRazorpayClient(cfg.RazorpayKey, cfg.RazorpaySecret)

	// Initialize services
	authService := services.NewAuthService(userRepo, cfg.JWTSecret)
	userService := services.NewUserService(userRepo)
	doctorService := services.NewDoctorService(doctorRepo)
	bookingService := services.NewBookingService(bookingRepo, doctorRepo, razorpayClient)
	journalService := services.NewJournalService(journalRepo)
	periodService := services.NewPeriodService(periodRepo)
	pregnancyService := services.NewPregnancyService(pregnancyRepo)
	symptomsService := services.NewSymptomsService(symptomsRepo)
	weightService := services.NewWeightService(weightRepo)
	pcosService := services.NewPCOSService(pcosRepo)
	mentalHealthService := services.NewMentalHealthService(mentalHealthRepo)
	fsfiService := services.NewFSFIService(mentalHealthRepo)
	clinicService := services.NewClinicService(clinicRepo)
	diagnosticService := services.NewDiagnosticService(diagnosticRepo)

	// Initialize handlers
	authHandler := handlers.NewAuthHandler(authService)
	userHandler := handlers.NewUserHandler(userService)
	doctorHandler := handlers.NewDoctorHandler(doctorService)
	bookingHandler := handlers.NewBookingHandler(bookingService)
	timeSlotHandler := handlers.NewTimeSlotHandler(doctorService, bookingService)
	journalHandler := handlers.NewJournalHandler(journalService)
	periodHandler := handlers.NewPeriodHandler(periodService)
	pregnancyHandler := handlers.NewPregnancyHandler(pregnancyService)
	symptomsHandler := handlers.NewSymptomsHandler(symptomsService)
	weightHandler := handlers.NewWeightHandler(weightService)
	pcosHandler := handlers.NewPCOSHandler(pcosService)
	mentalHealthHandler := handlers.NewMentalHealthHandler(mentalHealthService)
	fsfiHandler := handlers.NewFSFIHandler(fsfiService)
	clinicHandler := handlers.NewClinicHandler(clinicService)
	diagnosticsHandler := handlers.NewDiagnosticHandler(diagnosticService)

	// Swagger documentation
	router.StaticFile("/swagger.yaml", "./swagger.yaml")
	router.GET("/swagger", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "/swagger-ui")
	})
	router.GET("/swagger-ui", func(c *gin.Context) {
		html := `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Women's Health API Documentation</title>
    <link rel="stylesheet" type="text/css" href="https://unpkg.com/swagger-ui-dist@5.10.0/swagger-ui.css">
    <style>
        body { margin: 0; padding: 0; }
    </style>
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
                presets: [
                    SwaggerUIBundle.presets.apis,
                    SwaggerUIStandalonePreset
                ],
                plugins: [
                    SwaggerUIBundle.plugins.DownloadUrl
                ],
                layout: "StandaloneLayout"
            });
        };
    </script>
</body>
</html>`
		c.Header("Content-Type", "text/html; charset=utf-8")
		c.String(http.StatusOK, html)
	})

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "healthy",
			"service": "Women's Health API",
		})
	})

	// Public routes - Authentication
	router.POST("/user/register", authHandler.Register)
	router.POST("/user/login", authHandler.Login)
	router.POST("/user/forgot-password", authHandler.ForgotPassword)
	router.POST("/user/verify-otp", authHandler.VerifyOTP)
	router.POST("/user/reset-password", authHandler.ResetPassword)

	// Protected routes - User Profile
	userRoutes := router.Group("/user")
	userRoutes.Use(middleware.AuthMiddleware(cfg.JWTSecret))
	{
		userRoutes.POST("/change-password", authHandler.ChangePassword)
		userRoutes.GET("/profile", userHandler.GetProfile)
		userRoutes.PUT("/profile", userHandler.UpdateProfile)
	}

	// Protected API routes
	api := router.Group("/api")
	api.Use(middleware.AuthMiddleware(cfg.JWTSecret))
	{
		// Health Data
		// api.GET("/health-data", healthDataHandler.GetHealthData) // TODO: Implement

		// Doctors
		api.GET("/doctors", doctorHandler.GetDoctors)
		api.GET("/doctors/:doctorId", doctorHandler.GetDoctorByID)

		// Sessions & Bookings
		api.GET("/time-slots", timeSlotHandler.GetAvailableTimeSlots)
		api.POST("/bookings", bookingHandler.CreateBooking)
		api.POST("/bookings/verify", bookingHandler.VerifyPayment)
		api.GET("/bookings/my-with-doctors", bookingHandler.GetUserBookings)
		api.GET("/sessions/active", bookingHandler.GetActiveBookings)
		// api.POST("/agora/session-token", agoraHandler.GenerateToken) // TODO: Implement

		// Clinics
		api.GET("/clinics", clinicHandler.GetClinics)
		api.POST("/clinic-bookings", clinicHandler.CreateBooking)
		api.GET("/clinic-bookings/my-bookings", clinicHandler.GetMyBookings)
		api.POST("/clinic-bookings/verify-payment", clinicHandler.VerifyPayment)

		// Diagnostics
		api.GET("/public/diagnostics", diagnosticsHandler.GetDiagnostics)
		api.GET("/diagnosticsUsers", diagnosticsHandler.GetDiagnosticsUsers)
		api.POST("/diagnostics-bookings", diagnosticsHandler.CreateBooking)
		api.POST("/diagnostics-bookings/verify-payment", diagnosticsHandler.VerifyPayment)

		// Period Tracker
		api.GET("/period-cycle", periodHandler.GetPeriodCycle)
		api.POST("/period-cycle", periodHandler.AddPeriodCycle)
		api.DELETE("/period-cycle", periodHandler.ResetPeriodTracker)

		// Pregnancy Tracker
		api.GET("/pregnancy-tracker", pregnancyHandler.GetPregnancyData)
		api.POST("/pregnancy-tracker", pregnancyHandler.AddPregnancyEntry)

		// Sexual Wellness (FSFI)
		api.GET("/fsfi/test", fsfiHandler.GetTest)
		api.POST("/fsfi/submit", fsfiHandler.SubmitTest)
		api.GET("/fsfi/my-results", fsfiHandler.GetMyResults)

		// Mental Health Tests
		api.GET("/tests", mentalHealthHandler.GetTests)
		api.GET("/tests/:testName", mentalHealthHandler.GetTestByName)
		api.POST("/test-results", mentalHealthHandler.SubmitTestResults)
		api.GET("/test-results", mentalHealthHandler.GetTestResults)

		// PCOS Assessment
		api.GET("/pcos-assessment/questions", pcosHandler.GetQuestions)
		api.POST("/pcos-assessment/submit", pcosHandler.SubmitAssessment)
		api.GET("/pcos-assessment/history", pcosHandler.GetHistory)
		api.GET("/pcos-assessment", pcosHandler.GetLatestAssessment)

		// Symptoms Tracking
		api.POST("/symptoms-tracking", symptomsHandler.SubmitTracking)
		api.GET("/symptoms-tracking", symptomsHandler.GetTrackingHistory)

		// Weight & Metabolic Wellness
		api.GET("/weight-metabolic-wellness", weightHandler.GetData)
		api.POST("/weight-metabolic-wellness", weightHandler.AddEntry)

		// Journals
		api.GET("/journals", journalHandler.GetJournals)
		api.POST("/journals", journalHandler.CreateJournal)
		api.GET("/journals/:journalId", journalHandler.GetJournalByID)
		api.PUT("/journals/:journalId", journalHandler.UpdateJournal)
		api.DELETE("/journals/:journalId", journalHandler.DeleteJournal)
		api.GET("/journals/user/:userId", journalHandler.GetJournalsByUserID)

		// Groups & Community
		// api.GET("/groups", groupHandler.GetGroups) // TODO: Implement
		// api.GET("/groups/joined", groupHandler.GetJoinedGroups) // TODO: Implement
		// api.POST("/groups/:groupId/join", groupHandler.JoinGroup) // TODO: Implement
		// api.POST("/groups/:groupId/leave", groupHandler.LeaveGroup) // TODO: Implement
		// api.GET("/groups/:groupId/posts", groupHandler.GetGroupPosts) // TODO: Implement
		// api.POST("/groups/:groupId/posts", groupHandler.CreateGroupPost) // TODO: Implement
		// api.GET("/groups/posts", groupHandler.GetAllPosts) // TODO: Implement
		// api.GET("/groups/comments/my", groupHandler.GetMyComments) // TODO: Implement

		// Blogs
		// api.GET("/blogs", blogHandler.GetBlogs) // TODO: Implement
		// api.GET("/blogs/:blogId", blogHandler.GetBlogByID) // TODO: Implement

		// Chat
		// api.POST("/chat/send", chatHandler.SendMessage) // TODO: Implement
		// api.GET("/chat/history/:userId", chatHandler.GetChatHistory) // TODO: Implement
		// api.GET("/chat/rooms", chatHandler.GetChatRooms) // TODO: Implement

		// Media Gallery
		// api.GET("/media-gallery", mediaHandler.GetGallery) // TODO: Implement
		// api.POST("/media-gallery", mediaHandler.UploadMedia) // TODO: Implement

		// Surveys
		// api.GET("/surveys", surveyHandler.GetSurveys) // TODO: Implement
		// api.GET("/surveys/:surveyId/questions", surveyHandler.GetSurveyQuestions) // TODO: Implement
	}
}

// Helper function to get database from context
func GetDBFromContext(c *gin.Context) *database.MongoDB {
	db, exists := c.Get("db")
	if !exists {
		return nil
	}
	return db.(*database.MongoDB)
}
