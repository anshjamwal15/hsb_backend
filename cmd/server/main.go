package main

import (
	_ "embed"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/anshjamwal15/hsb_backend/internal/infrastructure/database"
	"github.com/anshjamwal15/hsb_backend/internal/interfaces/http/router"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

var swaggerYAML []byte

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	// Get configuration from environment
	port := getEnv("PORT", "8080")
	mongoURI := getEnv("MONGODB_URI", "mongodb://localhost:27017")
	mongoDatabase := getEnv("MONGODB_DATABASE", "hsb_backend")
	jwtSecret := getEnv("JWT_SECRET", "your-secret-key-change-in-production")

	// Connect to MongoDB
	db, err := database.NewMongoDB(mongoURI, mongoDatabase)
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	defer db.Disconnect()

	log.Println("✅ Connected to MongoDB successfully")

	// Initialize Gin router
	r := gin.Default()

	// CORS middleware
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE, PATCH")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	// Setup routes
	router.SetupRoutes(r, db, jwtSecret)

	// Serve swagger.yaml file (embedded in binary)
	r.GET("/swagger.yaml", func(c *gin.Context) {
		c.Data(http.StatusOK, "application/x-yaml", swaggerYAML)
	})

	// Swagger UI endpoint - serves a simple HTML page with Swagger UI
	r.GET("/swagger", func(c *gin.Context) {
		html := `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>HSB Backend API Documentation</title>
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
            SwaggerUIBundle({
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
</html>
`
		c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(html))
	})

	// Print startup information
	printStartupInfo(port)

	// Start server
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func printStartupInfo(port string) {
	fmt.Println("\n" + "═══════════════════════════════════════════════════════════════")
	fmt.Println("🚀 HSB Backend Server Starting...")
	fmt.Println("═══════════════════════════════════════════════════════════════")
	fmt.Println()
	fmt.Printf("📡 Server URL:          http://localhost:%s\n", port)
	fmt.Printf("📚 Swagger UI:          http://localhost:%s/swagger\n", port)
	fmt.Printf("📄 Swagger YAML:        http://localhost:%s/swagger.yaml\n", port)
	fmt.Printf("❤️  Health Check:        http://localhost:%s/health\n", port)
	fmt.Println()
	fmt.Println("═══════════════════════════════════════════════════════════════")
	fmt.Println("📋 Available Endpoints:")
	fmt.Println("═══════════════════════════════════════════════════════════════")
	fmt.Println()
	fmt.Println("🔐 Authentication:")
	fmt.Println("   • POST   /user/register")
	fmt.Println("   • POST   /user/login")
	fmt.Println("   • POST   /user/forgot-password")
	fmt.Println("   • POST   /user/verify-otp")
	fmt.Println("   • POST   /user/reset-password")
	fmt.Println()
	fmt.Println("👤 User Profile:")
	fmt.Println("   • GET    /api/user/profile")
	fmt.Println("   • PUT    /api/user/profile")
	fmt.Println("   • POST   /api/user/change-password")
	fmt.Println()
	fmt.Println("👨‍⚕️ Doctors:")
	fmt.Println("   • GET    /api/doctors")
	fmt.Println("   • GET    /api/doctors/:doctorId")
	fmt.Println()
	fmt.Println("📅 Bookings:")
	fmt.Println("   • GET    /api/time-slots")
	fmt.Println("   • POST   /api/bookings")
	fmt.Println("   • POST   /api/bookings/verify")
	fmt.Println("   • GET    /api/bookings/my-with-doctors")
	fmt.Println("   • GET    /api/sessions/active")
	fmt.Println()
	fmt.Println("═══════════════════════════════════════════════════════════════")
	fmt.Printf("✨ Server is running on port %s\n", port)
	fmt.Println("═══════════════════════════════════════════════════════════════")
	fmt.Println()
	fmt.Println("💡 Quick Start:")
	fmt.Printf("   1. Open Swagger UI: http://localhost:%s/swagger\n", port)
	fmt.Println("   2. Test APIs directly from the browser")
	fmt.Println("   3. Check health: curl http://localhost:" + port + "/health")
	fmt.Println()
	fmt.Println("🎯 Ready to accept requests!")
	fmt.Println()
}
