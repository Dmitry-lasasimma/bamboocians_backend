package handlers

import (
	"bamboocians/database"
	"bamboocians/models"
	"bamboocians/utils"

	"github.com/gofiber/fiber/v2"
)

func SendMessage(c *fiber.Ctx) error {
	senderID := c.Locals("userID").(uint)

	var input struct {
		ReceiverID uint   `json:"receiver_id"`
		Content    string `json:"content"`
	}
	if err := c.BodyParser(&input); err != nil {
		return utils.Fail(c, fiber.StatusBadRequest, "invalid request body")
	}

	if input.Content == "" {
		return utils.Fail(c, fiber.StatusBadRequest, "content cannot be empty")
	}

	msg := models.Message{
		SenderID:   senderID,
		ReceiverID: input.ReceiverID,
		Content:    input.Content,
	}
	database.DB.Create(&msg)
	database.DB.Preload("Sender").Preload("Receiver").First(&msg, msg.ID)
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"success": true, "data": msg})
}

func GetInbox(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uint)
	var messages []models.Message
	database.DB.Preload("Sender").Preload("Receiver").
		Where("receiver_id = ?", userID).
		Order("created_at desc").Find(&messages)
	return utils.OK(c, messages)
}

func GetConversation(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uint)
	otherID := c.Params("userId")

	var messages []models.Message
	database.DB.Preload("Sender").Preload("Receiver").
		Where("(sender_id = ? AND receiver_id = ?) OR (sender_id = ? AND receiver_id = ?)",
			userID, otherID, otherID, userID).
		Order("created_at asc").Find(&messages)

	// Mark received messages as read
	database.DB.Model(&models.Message{}).
		Where("sender_id = ? AND receiver_id = ? AND read = false", otherID, userID).
		Update("read", true)

	return utils.OK(c, messages)
}

func GetUnreadCount(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uint)
	var count int64
	database.DB.Model(&models.Message{}).Where("receiver_id = ? AND read = false", userID).Count(&count)
	return utils.OK(c, fiber.Map{"unread_count": count})
}
