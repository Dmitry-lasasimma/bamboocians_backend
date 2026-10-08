package routes

import (
	"bamboocians/handlers"
	"bamboocians/middleware"
	"bamboocians/models"

	"github.com/gofiber/fiber/v2"
)

func Setup(app *fiber.App) {
	api := app.Group("/api")

	// ── Auth (public) ────────────────────────────────────────────────────────
	auth := api.Group("/auth")
	auth.Post("/register", handlers.Register)
	auth.Post("/login", handlers.Login)
	auth.Get("/me", middleware.Auth(), handlers.Me)

	// ── Marketplace (public, no auth required) ───────────────────────────────
	market := api.Group("/marketplace")
	market.Get("/talents", handlers.ListTalents)
	market.Get("/talents/:id", handlers.GetTalentPublicProfile)
	market.Get("/venues", handlers.ListVenues)
	market.Get("/venues/:id", handlers.GetVenuePublicProfile)

	// ── Organizer routes ─────────────────────────────────────────────────────
	org := api.Group("/organizer", middleware.Auth(), middleware.RequireRole(models.RoleOrganizer))
	org.Get("/dashboard", handlers.OrganizerDashboard)
	org.Get("/profile", handlers.GetOrganizerProfile)
	org.Put("/profile", handlers.UpdateOrganizerProfile)

	// Events
	org.Get("/events", handlers.GetOrganizerEvents)
	org.Post("/events", handlers.CreateEvent)
	org.Get("/events/:id", handlers.GetOrganizerEvent)
	org.Put("/events/:id", handlers.UpdateEvent)
	org.Delete("/events/:id", handlers.DeleteEvent)

	// Guests
	org.Get("/events/:eventId/guests", handlers.GetGuests)
	org.Post("/events/:eventId/guests", handlers.AddGuest)
	org.Put("/guests/:id", handlers.UpdateGuestStatus)

	// Tickets
	org.Get("/events/:eventId/tickets", handlers.GetTickets)
	org.Post("/events/:eventId/tickets", handlers.CreateTicket)

	// Bookings (organizer creates bookings for their events)
	org.Get("/bookings", handlers.GetOrganizerBookings)
	org.Post("/bookings", handlers.CreateBooking)
	org.Put("/bookings/:id/cancel", handlers.CancelBooking)

	// Contracts (organizer creates contracts for confirmed bookings)
	org.Post("/contracts", handlers.CreateContract)

	// ── Talent routes ────────────────────────────────────────────────────────
	talent := api.Group("/talent", middleware.Auth(), middleware.RequireRole(models.RoleTalent))
	talent.Get("/dashboard", handlers.TalentDashboard)
	talent.Get("/profile", handlers.GetTalentProfile)
	talent.Put("/profile", handlers.UpdateTalentProfile)
	talent.Get("/bookings", handlers.GetTalentBookings)
	talent.Put("/bookings/:id/respond", handlers.RespondToBooking)

	// ── Venue routes ─────────────────────────────────────────────────────────
	venue := api.Group("/venue", middleware.Auth(), middleware.RequireRole(models.RoleVenue))
	venue.Get("/dashboard", handlers.VenueDashboard)
	venue.Get("/profile", handlers.GetVenueProfile)
	venue.Put("/profile", handlers.UpdateVenueProfile)
	venue.Get("/bookings", handlers.GetVenueBookings)
	venue.Put("/bookings/:id/respond", handlers.RespondToVenueBooking)
	venue.Get("/calendar", handlers.GetVenueCalendar)

	// ── Shared (any authenticated user) ─────────────────────────────────────
	shared := api.Group("/", middleware.Auth())
	shared.Get("/contracts", handlers.GetMyContracts)
	shared.Get("/contracts/:id", handlers.GetContract)
	shared.Put("/contracts/:id/sign", handlers.SignContract)
	shared.Post("/messages", handlers.SendMessage)
	shared.Get("/messages/inbox", handlers.GetInbox)
	shared.Get("/messages/unread", handlers.GetUnreadCount)
	shared.Get("/messages/:userId", handlers.GetConversation)
}
