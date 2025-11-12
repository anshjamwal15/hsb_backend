package router

import (
	"github.com/anshjamwal15/hsb_backend/internal/application/services"
	"github.com/anshjamwal15/hsb_backend/internal/infrastructure/database"
	repo "github.com/anshjamwal15/hsb_backend/internal/infrastructure/repositories"
	"github.com/anshjamwal15/hsb_backend/internal/interfaces/http/handlers"
	"github.com/anshjamwal15/hsb_backend/internal/interfaces/http/middleware"
	"github.com/gin-gonic/gin"
)

// SetupRoutes configures all the routes for the application
func SetupRoutes(r *gin.Engine, db *database.MongoDB, jwtSecret string) {
	// Get the MongoDB database client
	mongoDB := db.Database

	// Initialize repositories
	userRepo := repo.NewUserRepository(mongoDB)
	// FSFI uses the same repository as mental health since they share the same collection
	mentalHealthRepo := repo.NewMentalHealthRepository(mongoDB)
	pcosRepo := repo.NewPCOSRepository(mongoDB)
	pregnancyRepo := repo.NewPregnancyRepository(mongoDB)
	periodRepo := repo.NewPeriodRepository(mongoDB)
	symptomsRepo := repo.NewSymptomsRepository(mongoDB)

	// Initialize services
	userService := services.NewUserService(userRepo)
	// FSFI uses the same service as mental health since they share the same repository
	fsfiService := services.NewFSFIService(mentalHealthRepo)
	mentalHealthService := services.NewMentalHealthService(mentalHealthRepo)
	pcosService := services.NewPCOSService(pcosRepo)
	pregnancyService := services.NewPregnancyService(pregnancyRepo)
	periodService := services.NewPeriodService(periodRepo)
	symptomsService := services.NewSymptomsService(symptomsRepo)

	// Initialize handlers
	userHandler := handlers.NewUserHandler(userService)
	fsfiHandler := handlers.NewFSFIHandler(fsfiService)
	mentalHealthHandler := handlers.NewMentalHealthHandler(mentalHealthService)
	pcosHandler := handlers.NewPCOSHandler(pcosService)
	pregnancyHandler := handlers.NewPregnancyHandler(pregnancyService)
	periodHandler := handlers.NewPeriodHandler(periodService)
	symptomsHandler := handlers.NewSymptomsHandler(symptomsService)

	// Health check endpoint
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
		})
	})

	// API routes
	api := r.Group("/api")
	{
		// Public routes
		// TODO: Add any public routes here

		// Protected routes (require authentication)
		protected := api.Group("")
		protected.Use(middleware.AuthMiddleware(jwtSecret))
		{
			// User routes
			user := protected.Group("/users")
			{
				user.GET("/me", userHandler.GetProfile)
				user.PUT("/me", userHandler.UpdateProfile)
			}

			// FSFI routes
			fsfi := protected.Group("/fsfi")
			{
				fsfi.GET("", fsfiHandler.GetTest)           // Get FSFI test questions
				fsfi.POST("", fsfiHandler.SubmitTest)        // Submit FSFI test answers
				fsfi.GET("/results", fsfiHandler.GetMyResults) // Get user's FSFI results
			}

			// Mental Health routes
			mentalHealth := protected.Group("/mental-health")
			{
				mentalHealth.GET("", mentalHealthHandler.GetTests)           // Get available mental health tests
				mentalHealth.POST("/submit", mentalHealthHandler.SubmitTestResults) // Submit mental health test results
				mentalHealth.GET("/results", mentalHealthHandler.GetTestResults)   // Get user's mental health test results
			}

			// PCOS routes
			pcos := protected.Group("/pcos")
			{
				pcos.GET("", pcosHandler.GetQuestions)          // Get PCOS assessment questions
				pcos.POST("", pcosHandler.SubmitAssessment)     // Submit PCOS assessment
				pcos.GET("/history", pcosHandler.GetHistory)     // Get user's PCOS assessment history
				pcos.GET("/latest", pcosHandler.GetLatestAssessment) // Get user's latest PCOS assessment
			}

			// Pregnancy routes
			pregnancy := protected.Group("/pregnancy")
			{
				pregnancy.GET("", pregnancyHandler.GetPregnancyData)  // Get user's pregnancy data
				pregnancy.POST("", pregnancyHandler.AddPregnancyEntry) // Add new pregnancy entry
			}

			// Period routes
			period := protected.Group("/period")
			{
				period.GET("", periodHandler.GetPeriodCycle)     // Get user's period cycle data
				period.POST("", periodHandler.AddPeriodCycle)    // Add new period cycle
				period.DELETE("", periodHandler.ResetPeriodTracker) // Reset period tracker
			}

			// Symptoms routes
			symptoms := protected.Group("/symptoms")
			{
				symptoms.POST("", symptomsHandler.SubmitTracking)    // Submit symptom tracking data
				symptoms.GET("", symptomsHandler.GetTrackingHistory) // Get user's symptom tracking history
			}
		}
	}
}
