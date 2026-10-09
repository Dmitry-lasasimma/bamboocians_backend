package handlers

import (
	"bamboocians/database"
	"bamboocians/models"
	"bamboocians/utils"
	"time"

	"github.com/gofiber/fiber/v2"
)

type BookingInput struct {
	EventID     uint    `json:"event_id"`
	BookedToID  uint    `json:"booked_to_id"`
	BookingType string  `json:"booking_type"` // talent | venue
	Price       float64 `json:"price"`
	Notes       string  `json:"notes"`
	Date        string  `json:"date"` // RFC3339
}

func CreateBooking(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uint)

	var input BookingInput
	if err := c.BodyParser(&input); err != nil {
		return utils.Fail(c, fiber.StatusBadRequest, "invalid request body")
	}

	if input.BookingType != models.BookingTypeTalent && input.BookingType != models.BookingTypeVenue {
		return utils.Fail(c, fiber.StatusBadRequest, "booking_type must be talent or venue")
	}

	// Verify the event belongs to this organizer
	var event models.Event
	if err := database.DB.Where("id = ? AND organizer_id = ?", input.EventID, userID).First(&event).Error; err != nil {
		return utils.Fail(c, fiber.StatusNotFound, "event not found or not yours")
	}

	date, err := time.Parse(time.RFC3339, input.Date)
	if err != nil {
		return utils.Fail(c, fiber.StatusBadRequest, "date must be RFC3339 format")
	}

	// The booked user must exist and have the role matching the booking type
	var bookedTo models.User
	if err := database.DB.Where("id = ? AND role = ?", input.BookedToID, input.BookingType).First(&bookedTo).Error; err != nil {
		return utils.Fail(c, fiber.StatusNotFound, input.BookingType+" not found")
	}

	booking := models.Booking{
		EventID:     input.EventID,
		BookedByID:  userID,
		BookedToID:  input.BookedToID,
		BookingType: input.BookingType,
		Price:       input.Price,
		Notes:       input.Notes,
		Date:        date,
		Status:      models.BookingStatusPending,
	}
	database.DB.Create(&booking)
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"success": true, "data": booking})
}

func GetOrganizerBookings(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uint)
	var bookings []models.Booking
	database.DB.Preload("Event").Preload("BookedTo").
		Where("booked_by_id = ?", userID).
		Order("created_at desc").Find(&bookings)
	return utils.OK(c, bookings)
}

func CancelBooking(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uint)
	id := c.Params("id")

	var booking models.Booking
	if err := database.DB.Where("id = ? AND booked_by_id = ?", id, userID).First(&booking).Error; err != nil {
		return utils.Fail(c, fiber.StatusNotFound, "booking not found")
	}

	booking.Status = models.BookingStatusCancelled
	database.DB.Save(&booking)
	return utils.OK(c, booking)
}
