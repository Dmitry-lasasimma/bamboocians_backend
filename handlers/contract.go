package handlers

import (
	"bamboocians/database"
	"bamboocians/models"
	"bamboocians/utils"

	"github.com/gofiber/fiber/v2"
)

func CreateContract(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uint)

	var input struct {
		BookingID uint   `json:"booking_id"`
		Content   string `json:"content"`
	}
	if err := c.BodyParser(&input); err != nil {
		return utils.Fail(c, fiber.StatusBadRequest, "invalid request body")
	}

	// Ensure the booking belongs to this organizer
	var booking models.Booking
	if err := database.DB.Where("id = ? AND booked_by_id = ?", input.BookingID, userID).First(&booking).Error; err != nil {
		return utils.Fail(c, fiber.StatusNotFound, "booking not found")
	}

	if booking.Status != models.BookingStatusConfirmed {
		return utils.Fail(c, fiber.StatusBadRequest, "contracts can only be created for confirmed bookings")
	}

	var existing int64
	database.DB.Model(&models.Contract{}).Where("booking_id = ?", booking.ID).Count(&existing)
	if existing > 0 {
		return utils.Fail(c, fiber.StatusConflict, "a contract already exists for this booking")
	}

	contract := models.Contract{
		BookingID:       input.BookingID,
		Content:         input.Content,
		OrganizerSigned: true, // organizer creates and auto-signs
		Status:          models.ContractStatusPending,
	}
	database.DB.Create(&contract)
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"success": true, "data": contract})
}

func GetContract(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uint)
	id := c.Params("id")

	var contract models.Contract
	if err := database.DB.Preload("Booking.Event").Preload("Booking.BookedBy").Preload("Booking.BookedTo").
		First(&contract, id).Error; err != nil {
		return utils.Fail(c, fiber.StatusNotFound, "contract not found")
	}

	// Only parties involved can view
	b := contract.Booking
	if b.BookedByID != userID && b.BookedToID != userID {
		return utils.Fail(c, fiber.StatusForbidden, "access denied")
	}

	return utils.OK(c, contract)
}

func SignContract(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uint)
	id := c.Params("id")

	var contract models.Contract
	if err := database.DB.Preload("Booking").First(&contract, id).Error; err != nil {
		return utils.Fail(c, fiber.StatusNotFound, "contract not found")
	}

	b := contract.Booking
	if b.BookedByID == userID {
		contract.OrganizerSigned = true
	} else if b.BookedToID == userID {
		contract.TalentVenueSigned = true
	} else {
		return utils.Fail(c, fiber.StatusForbidden, "you are not a party to this contract")
	}

	if contract.OrganizerSigned && contract.TalentVenueSigned {
		contract.Status = models.ContractStatusActive
	}

	database.DB.Save(&contract)
	return utils.OK(c, contract)
}

func GetMyContracts(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uint)
	var contracts []models.Contract
	database.DB.Preload("Booking.Event").Preload("Booking.BookedBy").Preload("Booking.BookedTo").
		Joins("JOIN bookings ON bookings.id = contracts.booking_id").
		Where("bookings.booked_by_id = ? OR bookings.booked_to_id = ?", userID, userID).
		Find(&contracts)
	return utils.OK(c, contracts)
}
