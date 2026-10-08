package handlers

import (
	"bamboocians/database"
	"bamboocians/middleware"
	"bamboocians/models"
	"bamboocians/utils"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
)

type RegisterInput struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Role     string `json:"role"` // organizer | talent | venue
}

type LoginInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func Register(c *fiber.Ctx) error {
	var input RegisterInput
	if err := c.BodyParser(&input); err != nil {
		return utils.Fail(c, fiber.StatusBadRequest, "invalid request body")
	}

	if input.Name == "" || input.Email == "" || input.Password == "" {
		return utils.Fail(c, fiber.StatusBadRequest, "name, email and password are required")
	}

	validRoles := map[string]bool{
		models.RoleOrganizer: true,
		models.RoleTalent:    true,
		models.RoleVenue:     true,
	}
	if !validRoles[input.Role] {
		return utils.Fail(c, fiber.StatusBadRequest, "role must be organizer, talent, or venue")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return utils.Fail(c, fiber.StatusInternalServerError, "failed to hash password")
	}

	user := models.User{
		Name:     input.Name,
		Email:    input.Email,
		Password: string(hash),
		Role:     input.Role,
	}

	if err := database.DB.Create(&user).Error; err != nil {
		return utils.Fail(c, fiber.StatusConflict, "email already registered")
	}

	// Create an empty profile for the chosen role
	switch input.Role {
	case models.RoleOrganizer:
		database.DB.Create(&models.OrganizerProfile{UserID: user.ID})
	case models.RoleTalent:
		database.DB.Create(&models.TalentProfile{UserID: user.ID})
	case models.RoleVenue:
		database.DB.Create(&models.VenueProfile{UserID: user.ID})
	}

	token, _ := utils.GenerateToken(user.ID, user.Email, user.Role, middleware.JWTSecret)
	return utils.OK(c, fiber.Map{"token": token, "user": user})
}

func Login(c *fiber.Ctx) error {
	var input LoginInput
	if err := c.BodyParser(&input); err != nil {
		return utils.Fail(c, fiber.StatusBadRequest, "invalid request body")
	}

	var user models.User
	if err := database.DB.Where("email = ?", input.Email).First(&user).Error; err != nil {
		return utils.Fail(c, fiber.StatusUnauthorized, "invalid email or password")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
		return utils.Fail(c, fiber.StatusUnauthorized, "invalid email or password")
	}

	token, _ := utils.GenerateToken(user.ID, user.Email, user.Role, middleware.JWTSecret)
	return utils.OK(c, fiber.Map{"token": token, "user": user})
}

func Me(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uint)

	var user models.User
	if err := database.DB.First(&user, userID).Error; err != nil {
		return utils.Fail(c, fiber.StatusNotFound, "user not found")
	}
	return utils.OK(c, user)
}
