package handlers

import (
	"bamboocians/database"
	"bamboocians/models"
	"bamboocians/utils"
	"time"

	"github.com/gofiber/fiber/v2"
)

// ── Profile ──────────────────────────────────────────────────────────────────

func GetOrganizerProfile(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uint)
	var profile models.OrganizerProfile
	database.DB.Preload("User").Where("user_id = ?", userID).First(&profile)
	return utils.OK(c, profile)
}

func UpdateOrganizerProfile(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uint)
	var profile models.OrganizerProfile
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

// ── Dashboard ─────────────────────────────────────────────────────────────────

func OrganizerDashboard(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uint)

	var totalEvents, upcomingEvents, pendingBookings int64
	database.DB.Model(&models.Event{}).Where("organizer_id = ?", userID).Count(&totalEvents)
	database.DB.Model(&models.Event{}).Where("organizer_id = ? AND date > ? AND status != ?", userID, time.Now(), models.EventStatusCancelled).Count(&upcomingEvents)
	database.DB.Model(&models.Booking{}).Where("booked_by_id = ? AND status = ?", userID, models.BookingStatusPending).Count(&pendingBookings)

	var recentEvents []models.Event
	database.DB.Where("organizer_id = ?", userID).Order("created_at desc").Limit(5).Find(&recentEvents)

	return utils.OK(c, fiber.Map{
		"total_events":     totalEvents,
		"upcoming_events":  upcomingEvents,
		"pending_bookings": pendingBookings,
		"recent_events":    recentEvents,
	})
}

// ── Events ────────────────────────────────────────────────────────────────────

func GetOrganizerEvents(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uint)
	var events []models.Event
	database.DB.Where("organizer_id = ?", userID).Order("date desc").Find(&events)
	return utils.OK(c, events)
}

func GetOrganizerEvent(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uint)
	id := c.Params("id")
	var event models.Event
	if err := database.DB.Where("id = ? AND organizer_id = ?", id, userID).First(&event).Error; err != nil {
		return utils.Fail(c, fiber.StatusNotFound, "event not found")
	}
	return utils.OK(c, event)
}

type EventInput struct {
	Title       string  `json:"title"`
	Description string  `json:"description"`
	EventType   string  `json:"event_type"`
	Date        string  `json:"date"` // RFC3339
	Location    string  `json:"location"`
	Budget      float64 `json:"budget"`
	Capacity    int     `json:"capacity"`
	Status      string  `json:"status"`
}

func CreateEvent(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uint)
	var input EventInput
	if err := c.BodyParser(&input); err != nil {
		return utils.Fail(c, fiber.StatusBadRequest, "invalid request body")
	}

	date, err := time.Parse(time.RFC3339, input.Date)
	if err != nil {
		return utils.Fail(c, fiber.StatusBadRequest, "date must be RFC3339 format (e.g. 2024-12-31T20:00:00Z)")
	}

	status := input.Status
	if status == "" {
		status = models.EventStatusDraft
	}

	event := models.Event{
		OrganizerID: userID,
		Title:       input.Title,
		Description: input.Description,
		EventType:   input.EventType,
		Date:        date,
		Location:    input.Location,
		Budget:      input.Budget,
		Capacity:    input.Capacity,
		Status:      status,
	}
	database.DB.Create(&event)
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"success": true, "data": event})
}

func UpdateEvent(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uint)
	id := c.Params("id")

	var event models.Event
	if err := database.DB.Where("id = ? AND organizer_id = ?", id, userID).First(&event).Error; err != nil {
		return utils.Fail(c, fiber.StatusNotFound, "event not found")
	}

	var input EventInput
	if err := c.BodyParser(&input); err != nil {
		return utils.Fail(c, fiber.StatusBadRequest, "invalid request body")
	}

	if input.Date != "" {
		date, err := time.Parse(time.RFC3339, input.Date)
		if err != nil {
			return utils.Fail(c, fiber.StatusBadRequest, "date must be RFC3339 format")
		}
		event.Date = date
	}

	event.Title = input.Title
	event.Description = input.Description
	event.EventType = input.EventType
	event.Location = input.Location
	event.Budget = input.Budget
	event.Capacity = input.Capacity
	if input.Status != "" {
		event.Status = input.Status
	}

	database.DB.Save(&event)
	return utils.OK(c, event)
}

func DeleteEvent(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uint)
	id := c.Params("id")

	var event models.Event
	if err := database.DB.Where("id = ? AND organizer_id = ?", id, userID).First(&event).Error; err != nil {
		return utils.Fail(c, fiber.StatusNotFound, "event not found")
	}
	database.DB.Delete(&event)
	return utils.OK(c, fiber.Map{"message": "event deleted"})
}

// ── Guest List ────────────────────────────────────────────────────────────────

func GetGuests(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uint)
	eventID := c.Params("eventId")

	var event models.Event
	if err := database.DB.Where("id = ? AND organizer_id = ?", eventID, userID).First(&event).Error; err != nil {
		return utils.Fail(c, fiber.StatusNotFound, "event not found")
	}

	var guests []models.Guest
	database.DB.Where("event_id = ?", eventID).Find(&guests)
	return utils.OK(c, guests)
}

func AddGuest(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uint)
	eventID := c.Params("eventId")

	var event models.Event
	if err := database.DB.Where("id = ? AND organizer_id = ?", eventID, userID).First(&event).Error; err != nil {
		return utils.Fail(c, fiber.StatusNotFound, "event not found")
	}

	var guest models.Guest
	if err := c.BodyParser(&guest); err != nil {
		return utils.Fail(c, fiber.StatusBadRequest, "invalid request body")
	}
	guest.EventID = event.ID
	guest.Status = models.GuestStatusInvited

	database.DB.Create(&guest)
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"success": true, "data": guest})
}

func UpdateGuestStatus(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uint)
	id := c.Params("id")
	var input struct {
		Status string `json:"status"`
	}
	if err := c.BodyParser(&input); err != nil {
		return utils.Fail(c, fiber.StatusBadRequest, "invalid request body")
	}

	switch input.Status {
	case models.GuestStatusInvited, models.GuestStatusConfirmed, models.GuestStatusDeclined:
	default:
		return utils.Fail(c, fiber.StatusBadRequest, "status must be invited, confirmed or declined")
	}

	// Only the organizer who owns the guest's event may update it
	var guest models.Guest
	if err := database.DB.Joins("JOIN events ON events.id = guests.event_id").
		Where("guests.id = ? AND events.organizer_id = ?", id, userID).
		First(&guest).Error; err != nil {
		return utils.Fail(c, fiber.StatusNotFound, "guest not found")
	}
	database.DB.Model(&guest).Update("status", input.Status)
	return utils.OK(c, fiber.Map{"message": "guest status updated"})
}

// ── Tickets ───────────────────────────────────────────────────────────────────

func GetTickets(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uint)
	eventID := c.Params("eventId")

	var event models.Event
	if err := database.DB.Where("id = ? AND organizer_id = ?", eventID, userID).First(&event).Error; err != nil {
		return utils.Fail(c, fiber.StatusNotFound, "event not found")
	}

	var tickets []models.Ticket
	database.DB.Where("event_id = ?", eventID).Find(&tickets)
	return utils.OK(c, tickets)
}

func CreateTicket(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uint)
	eventID := c.Params("eventId")

	var event models.Event
	if err := database.DB.Where("id = ? AND organizer_id = ?", eventID, userID).First(&event).Error; err != nil {
		return utils.Fail(c, fiber.StatusNotFound, "event not found")
	}

	var ticket models.Ticket
	if err := c.BodyParser(&ticket); err != nil {
		return utils.Fail(c, fiber.StatusBadRequest, "invalid request body")
	}
	ticket.EventID = event.ID
	database.DB.Create(&ticket)
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"success": true, "data": ticket})
}
