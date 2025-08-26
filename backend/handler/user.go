package handler

import (
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"reonify/clocked/database"
	"reonify/clocked/model"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func validToken(t *jwt.Token, id string) bool {
	n, err := strconv.Atoi(id)
	if err != nil {
		return false
	}

	claims := t.Claims.(jwt.MapClaims)
	uid := int(claims["user_id"].(float64))

	return uid == n
}

func getUidFromToken(t *jwt.Token) int {
	claims := t.Claims.(jwt.MapClaims)
	uid := int(claims["user_id"].(float64))

	return uid
}

func validUser(id string, p string) bool {
	db := database.DB
	var user model.User
	db.First(&user, id)
	if user.Username == "" {
		return false
	}
	if !CheckPasswordHash(p, user.Password) {
		return false
	}
	return true
}

// GetUser get a user
func GetUser(c *fiber.Ctx) error {
	token := c.Locals("user").(*jwt.Token)

	db := database.DB
	var user model.User
	db.Find(&user, getUidFromToken(token))
	if user.Username == "" {
		return fiber.NewError(fiber.StatusNotFound, "No user found with ID")
	}
	return c.JSON(fiber.Map{"status": "success", "message": "User found", "data": user})
}

// CreateUser new user
func CreateUser(c *fiber.Ctx) error {
	type NewUser struct {
		Username string `json:"username"`
		Email    string `json:"email"`
	}

	db := database.DB
	user := new(model.User)
	if err := c.BodyParser(user); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Review your input")
	}

	hash, err := hashPassword(user.Password)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Couldn't hash password")
	}

	user.Password = hash
	if err := db.Create(&user).Error; err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Couldn't create user")
	}

	newUser := NewUser{
		Email:    user.Email,
		Username: user.Username,
	}

	return c.JSON(fiber.Map{"status": "success", "message": "Created user", "data": newUser})
}

// UpdateUser update user
func UpdateUser(c *fiber.Ctx) error {
	type UpdateUserInput struct {
		Names string `json:"names"`
	}
	var uui UpdateUserInput
	if err := c.BodyParser(&uui); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Review your input")
	}
	id := c.Params("id")
	token := c.Locals("user").(*jwt.Token)

	if !validToken(token, id) {
		return fiber.NewError(fiber.StatusInternalServerError, "Invalid token id")
	}

	db := database.DB
	var user model.User

	db.First(&user, id)
	user.Names = uui.Names
	db.Save(&user)

	return c.JSON(fiber.Map{"status": "success", "message": "User successfully updated", "data": user})
}

// DeleteUser delete user
func DeleteUser(c *fiber.Ctx) error {
	type PasswordInput struct {
		Password string `json:"password"`
	}
	var pi PasswordInput
	if err := c.BodyParser(&pi); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Review your input")
	}
	id := c.Params("id")
	token := c.Locals("user").(*jwt.Token)

	if !validToken(token, id) {
		return fiber.NewError(fiber.StatusInternalServerError, "Invalid token id")
	}

	if !validUser(id, pi.Password) {
		return fiber.NewError(fiber.StatusInternalServerError, "Not valid user")
	}

	db := database.DB
	var user model.User

	db.First(&user, id)

	db.Delete(&user)
	return c.JSON(fiber.Map{"status": "success", "message": "User successfully deleted", "data": nil})
}

// UploadAvatar upload avatar
func UploadAvatar(c *fiber.Ctx) error {
	id := c.Params("id")
	token := c.Locals("user").(*jwt.Token)

	if !validToken(token, id) {
		return fiber.NewError(fiber.StatusInternalServerError, "Invalid token id")
	}

	db := database.DB
	var user model.User
	db.First(&user, id)

	file, err := c.FormFile("avatar")
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Couldn't get avatar")
	}

	// Ensure upload directory exists (project-relative for cross-platform)
	baseDir := filepath.Join(".", "uploads", "avatars")
	if mkErr := os.MkdirAll(baseDir, 0o755); mkErr != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Couldn't prepare upload directory")
	}

	// Preserve original extension when saving
	ext := filepath.Ext(file.Filename)
	filename := id + ext
	destPath := filepath.Join(baseDir, filename)

	if err := c.SaveFile(file, destPath); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Couldn't save avatar")
	}

	user.Avatar = destPath
	db.Save(&user)

	return c.JSON(fiber.Map{"status": "success", "message": "Avatar successfully uploaded", "data": user})
}

// GetAvatar get avatar
func GetAvatar(c *fiber.Ctx) error {
	id := c.Params("id")
	db := database.DB
	var user model.User
	db.First(&user, id)

	// Validate avatar path
	if strings.TrimSpace(user.Avatar) == "" {
		return fiber.NewError(fiber.StatusNotFound, "Avatar not set")
	}

	// Clean and normalize path; keep project-relative base allowed
	avatarPath := filepath.Clean(user.Avatar)
	if !filepath.IsAbs(avatarPath) {
		avatarPath = filepath.Join(".", avatarPath)
	}

	// Ensure file exists
	if _, err := os.Stat(avatarPath); err != nil {
		if os.IsNotExist(err) {
			return fiber.NewError(fiber.StatusNotFound, "Avatar not found")
		}
		return fiber.NewError(fiber.StatusInternalServerError, "Unable to read avatar")
	}

	// Set content type by extension for better rendering in browsers
	switch strings.ToLower(filepath.Ext(avatarPath)) {
	case ".png":
		c.Type("png")
	case ".jpg", ".jpeg":
		c.Type("jpeg")
	case ".gif":
		c.Type("gif")
	case ".webp":
		c.Type("webp")
	}

	return c.SendFile(avatarPath)
}
