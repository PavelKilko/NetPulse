package routes

import (
	"github.com/PavelKilko/NetPulse/handlers"
	"github.com/PavelKilko/NetPulse/middleware"
	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App) {
	// Public routes
	app.Post("/login", handlers.Login)
	app.Post("/signup", handlers.Signup)

	// User routes (protected)
	user := app.Group("/user")
	user.Use(middleware.JWTMiddleware())

	user.Post("/logout", handlers.Logout)

	user.Get("/groups", handlers.GetGroups)                // Fetch all groups
	user.Post("/groups", handlers.CreateGroup)             // Create a new group
	user.Put("/groups/:group_id", handlers.UpdateGroup)    // Update a group
	user.Delete("/groups/:group_id", handlers.DeleteGroup) // Delete a group

	user.Get("/groups/:group_id/urls", handlers.GetURLs)              // Fetch all URLs of group
	user.Post("/groups/:group_id/urls", handlers.CreateURL)           // Create a new URL of group
	user.Put("/groups/:group_id/urls/:url_id", handlers.UpdateURL)    // Update a URL of group
	user.Delete("/groups/:group_id/urls/:url_id", handlers.DeleteURL) // Delete a URL of group

	user.Put("/groups/:group_id/urls/:url_id/monitoring", handlers.ToggleMonitoring)
	user.Get("/groups/:group_id/urls/:url_id/metrics", handlers.GetMonitoringStatistics)
}
