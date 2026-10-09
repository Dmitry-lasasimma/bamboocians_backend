package handlers

import (
	"bamboocians/database"
	"bamboocians/models"
	"bamboocians/utils"
	"time"

	"github.com/gofiber/fiber/v2"
)

// ── Public Venue Marketplace ──────────────────────────────────────────────────

func ListVenues(c *fiber.Ctx) error {
	var profiles []models.VenueProfile
	database.DB.Preload("User").Find(&profiles)
	return utils.OK(c, profiles)
}

func GetVenuePublicProfile(c *fiber.Ctx) error {
	id := c.Params("id")
	var profile models.VenueProfile
	if err := database.DB.Preload("User").Where("user_id = ?", id).First(&profile).Error; err != nil {
		return utils.Fail(c, fiber.StatusNotFound, "venue not found")
	}
	return utils.OK(c, profile)
}

// ── Venue Dashboard ───────────────────────────────────────────────────────────

func VenueDashboard(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uint)

	var totalBookings, pendingBookings, confirmedBookings int64
	database.DB.Model(&models.Booking{}).Where("booked_to_id = ?", userID).Count(&totalBookings)
	database.DB.Model(&models.Booking{}).Where("booked_to_id = ? AND status = ?", userID, models.BookingStatusPending).Count(&pendingBookings)
	database.DB.Model(&models.Booking{}).Where("booked_to_id = ? AND status = ?", userID, models.BookingStatusConfirmed).Count(&confirmedBookings)

	var upcomingBookings []models.Booking
	database.DB.Preload("Event").Preload("BookedBy").
		Where("booked_to_id = ? AND status = ?", userID, models.BookingStatusConfirmed).
		Order("date asc").Limit(5).Find(&upcomingBookings)

	return utils.OK(c, fiber.Map{
		"total_bookings":     totalBookings,
		"pending_bookings":   pendingBookings,
		"confirmed_bookings": confirmedBookings,
		"upcoming_bookings":  upcomingBookings,
	})
}

// ── Venue Profile ─────────────────────────────────────────────────────────────

func GetVenueProfile(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uint)
	var profile models.VenueProfile
	database.DB.Preload("User").Where("user_id = ?", userID).First(&profile)
	return utils.OK(c, profile)
}

func UpdateVenueProfile(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uint)
	var profile models.VenueProfile
	database.DB.Where("user_id = ?", userID).First(&profile)
	profileID := profile.ID

	if err := c.BodyParser(&profile); err != nil {
		return utils.Fail(c, fiber.StatusBadRequest, "invalid request body")
	}
	profile.ID = profileID
	profile.UserID = userID
	profile.User = models.User{}
	database.DB.Save(&profile)
	return utils.OK(c, profile)
}

// ── Venue Bookings ────────────────────────────────────────────────────────────

func GetVenueBookings(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uint)
	var bookings []models.Booking
	database.DB.Preload("Event").Preload("BookedBy").
		Where("booked_to_id = ?", userID).
		Order("created_at desc").Find(&bookings)
	return utils.OK(c, bookings)
}

func RespondToVenueBooking(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uint)
	id := c.Params("id")

	var booking models.Booking
	if err := database.DB.Where("id = ? AND booked_to_id = ?", id, userID).First(&booking).Error; err != nil {
		return utils.Fail(c, fiber.StatusNotFound, "booking not found")
	}

	var input struct {
		Status string `json:"status"` // confirmed | declined
	}
	if err := c.BodyParser(&input); err != nil {
		return utils.Fail(c, fiber.StatusBadRequest, "invalid request body")
	}

	if input.Status != models.BookingStatusConfirmed && input.Status != models.BookingStatusDeclined {
		return utils.Fail(c, fiber.StatusBadRequest, "status must be confirmed or declined")
	}

	booking.Status = input.Status
	database.DB.Save(&booking)
	return utils.OK(c, booking)
}

// ── Venue Calendar (confirmed bookings as calendar events) ────────────────────

func GetVenueCalendar(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uint)
	var bookings []models.Booking
	database.DB.Preload("Event").
		Where("booked_to_id = ? AND status = ?", userID, models.BookingStatusConfirmed).
		Order("date asc").Find(&bookings)

	type CalendarEntry struct {
		BookingID  uint   `json:"booking_id"`
		EventTitle string `json:"event_title"`
		Date       string `json:"date"`
		Location   string `json:"location"`
		Status     string `json:"status"`
	}

	entries := []CalendarEntry{}
	for _, b := range bookings {
		entries = append(entries, CalendarEntry{
			BookingID:  b.ID,
			EventTitle: b.Event.Title,
			Date:       b.Date.Format(time.RFC3339),
			Location:   b.Event.Location,
			Status:     b.Status,
		})
	}
	return utils.OK(c, entries)
}
