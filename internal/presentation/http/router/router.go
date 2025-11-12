package router

import (
	"github.com/anshjamwal15/hsb_backend/internal/application/services"
	"github.com/anshjamwal15/hsb_backend/internal/infrastructure/database"
	"github.com/anshjamwal15/hsb_backend/internal/infrastructure/payment"
	repo "github.com/anshjamwal15/hsb_backend/internal/infrastructure/repositories"
	"github.com/anshjamwal15/hsb_backend/internal/interfaces/http/handlers"
	"github.com/anshjamwal15/hsb_backend/internal/interfaces/http/middleware"
	"github.com/gin-gonic/gin"
)

// SetupRoutes configures all the routes for the application
func SetupRoutes(r *gin.Engine, db *database.MongoDB, jwtSecret, razorpayKey, razorpaySecret string) {
	// Initialize repositories
	userRepo := repo.NewUserRepository(db.Database)
	doctorRepo := repo.NewDoctorRepository(db.Database)
	bookingRepo := repo.NewBookingRepository(db.Database)

	// Initialize payment client
	razorpayClient := payment.NewRazorpayClient(razorpayKey, razorpaySecret)

	// Initialize services
	authService := services.NewAuthService(userRepo, jwtSecret)
	userService := services.NewUserService(userRepo)
	doctorService := services.NewDoctorService(doctorRepo)
	bookingService := services.NewBookingService(bookingRepo, doctorRepo, razorpayClient)

	// Initialize handlers
	authHandler := handlers.NewAuthHandler(authService)
	userProfileHandler := handlers.NewUserProfileHandler(userService)
	doctorHandler := handlers.NewDoctorHandler(doctorService)
	bookingHandler := handlers.NewBookingHandler(bookingService)
	timeSlotHandler := handlers.NewTimeSlotHandler(doctorService, bookingService)

	// Health check endpoint
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
		})
	})

	// Authentication routes (public)
	userGroup := r.Group("/user")
	{
		userGroup.POST("/register", authHandler.Register)
		userGroup.POST("/login", authHandler.Login)
		userGroup.POST("/forgot-password", authHandler.ForgotPassword)
		userGroup.POST("/verify-otp", authHandler.VerifyOTP)
		userGroup.POST("/reset-password", authHandler.ResetPassword)
	}

	// API routes
	api := r.Group("/api")

	// Protected routes (require authentication)
	protected := api.Group("")
	protected.Use(middleware.AuthMiddleware(jwtSecret))
	{
		// User profile routes
		protected.GET("/user/profile", userProfileHandler.GetProfile)
		protected.PUT("/user/profile", userProfileHandler.UpdateProfile)
		protected.POST("/user/change-password", authHandler.ChangePassword)

		// Doctor routes
		protected.GET("/doctors", doctorHandler.GetDoctors)
		protected.GET("/doctors/:doctorId", doctorHandler.GetDoctorByID)

		// Time slots route
		protected.GET("/time-slots", timeSlotHandler.GetAvailableTimeSlots)

		// Booking routes
		protected.POST("/bookings", bookingHandler.CreateBooking)
		protected.POST("/bookings/verify", bookingHandler.VerifyPayment)
		protected.GET("/bookings/my-with-doctors", bookingHandler.GetUserBookings)
		protected.GET("/sessions/active", bookingHandler.GetActiveBookings)
	}
}
