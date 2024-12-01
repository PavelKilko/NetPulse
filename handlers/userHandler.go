package handlers

import (
	"github.com/PavelKilko/NetPulse/database"
	"github.com/PavelKilko/NetPulse/models"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"log"
	"strconv"
)

func CreateGroup(c *fiber.Ctx) error {
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	userID := uint(claims["user_id"].(float64))

	var group models.Group

	if err := c.BodyParser(&group); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Cannot parse request",
		})
	}

	group.UserID = userID

	// Save Group to DB
	if err := database.DB.Create(&group).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create group",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(group)
}

func GetGroups(c *fiber.Ctx) error {
	// Get the user ID from the JWT token
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)

	// Extract the user ID from the claims
	userID, ok := claims["user_id"].(float64)
	if !ok {
		log.Println("Invalid user ID in JWT claims")
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized: invalid user ID",
		})
	}

	// Find groups that belong to the logged-in user
	var groups []models.Group
	if err := database.DB.Where("user_id = ?", uint(userID)).Find(&groups).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to retrieve groups",
		})
	}

	return c.JSON(groups)
}

func UpdateGroup(c *fiber.Ctx) error {
	id := c.Params("group_id")
	var group models.Group

	// Get the user ID from the JWT token
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	userID, ok := claims["user_id"].(float64)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized: invalid user ID",
		})
	}

	// Find the group by ID and check if it belongs to the logged-in user
	if err := database.DB.First(&group, id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Group not found",
		})
	}

	// Ensure the group belongs to the logged-in user
	if group.UserID != uint(userID) {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "You do not have permission to update this group",
		})
	}

	// Parse the request body for new data
	if err := c.BodyParser(&group); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Cannot parse request",
		})
	}

	// Save the updated group to the database
	if err := database.DB.Save(&group).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update group",
		})
	}

	return c.JSON(group)
}

func DeleteGroup(c *fiber.Ctx) error {
	id := c.Params("group_id")
	var group models.Group

	// Get the user ID from the JWT token
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	userID, ok := claims["user_id"].(float64)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized: invalid user ID",
		})
	}

	// Find the group by ID and check if it belongs to the logged-in user
	if err := database.DB.First(&group, id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Group not found",
		})
	}

	// Ensure the group belongs to the logged-in user
	if group.UserID != uint(userID) {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "You do not have permission to delete this group",
		})
	}

	// Convert group.ID to uint for type safety with SQL queries
	groupID := group.ID

	// Delete all URLs associated with this group
	if err := database.DB.Where("group_id = ?", groupID).Delete(&models.URL{}).Error; err != nil {
		log.Printf("Failed to delete associated URLs for group ID %d: %v", groupID, err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete associated URLs",
		})
	}

	// Delete the group from the database
	if err := database.DB.Delete(&group).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete group",
		})
	}

	return c.SendStatus(fiber.StatusNoContent)
}

func CreateURL(c *fiber.Ctx) error {
	// Extract user information from JWT token claims
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	userID := uint(claims["user_id"].(float64))

	// Extract group ID from URL parameters
	groupIDParam := c.Params("group_id")
	groupID, err := strconv.ParseUint(groupIDParam, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid group ID",
		})
	}

	// Check if the group belongs to the user
	var group models.Group
	if err := database.DB.Where("id = ? AND user_id = ?", groupID, userID).First(&group).Error; err != nil {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "Group not found or access denied",
		})
	}

	var url models.URL

	// Parse the request body
	if err := c.BodyParser(&url); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Cannot parse request",
		})
	}

	// Set the group ID to ensure the URL is linked to the correct group
	url.GroupID = uint(groupID)

	// Save URL to DB
	if err := database.DB.Create(&url).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create URL",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(url)
}

func GetURLs(c *fiber.Ctx) error {
	// Extract user information from JWT token claims
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	userID := uint(claims["user_id"].(float64))

	// Extract group ID from URL parameters
	groupIDParam := c.Params("group_id")
	groupID, err := strconv.ParseUint(groupIDParam, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid group ID",
		})
	}

	// Check if the group belongs to the user
	var group models.Group
	if err := database.DB.Where("id = ? AND user_id = ?", groupID, userID).First(&group).Error; err != nil {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "Group not found or access denied",
		})
	}

	var urls []models.URL

	// Fetch URLs for the given group
	if err := database.DB.Where("group_id = ?", groupID).Find(&urls).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch URLs",
		})
	}

	return c.JSON(urls)
}

func UpdateURL(c *fiber.Ctx) error {
	// Extract user information from JWT token claims
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	userID := uint(claims["user_id"].(float64))

	// Extract group ID and URL ID from URL parameters
	groupIDParam := c.Params("group_id")
	groupID, err := strconv.ParseUint(groupIDParam, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid group ID",
		})
	}

	urlIDParam := c.Params("url_id")
	urlID, err := strconv.ParseUint(urlIDParam, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid URL ID",
		})
	}

	// Check if the group belongs to the user
	var group models.Group
	if err := database.DB.Where("id = ? AND user_id = ?", groupID, userID).First(&group).Error; err != nil {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "Group not found or access denied",
		})
	}

	// Find the URL by ID and ensure it belongs to the specified group
	var url models.URL
	if err := database.DB.Where("id = ? AND group_id = ?", urlID, groupID).First(&url).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "URL not found or access denied",
		})
	}

	// Parse the request body for new data
	if err := c.BodyParser(&url); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Cannot parse request",
		})
	}

	// Save the updated URL to the database
	if err := database.DB.Save(&url).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update URL",
		})
	}

	return c.JSON(url)
}

func DeleteURL(c *fiber.Ctx) error {
	// Extract user information from JWT token claims
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	userID := uint(claims["user_id"].(float64))

	// Extract group ID and URL ID from URL parameters
	groupIDParam := c.Params("group_id")
	groupID, err := strconv.ParseUint(groupIDParam, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid group ID",
		})
	}

	urlIDParam := c.Params("url_id")
	urlID, err := strconv.ParseUint(urlIDParam, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid URL ID",
		})
	}

	// Check if the group belongs to the user
	var group models.Group
	if err := database.DB.Where("id = ? AND user_id = ?", groupID, userID).First(&group).Error; err != nil {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "Group not found or access denied",
		})
	}

	// Find the URL by ID and ensure it belongs to the specified group
	var url models.URL
	if err := database.DB.Where("id = ? AND group_id = ?", urlID, groupID).First(&url).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "URL not found or access denied",
		})
	}

	// Delete the URL from the database
	if err := database.DB.Delete(&url).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete URL",
		})
	}

	return c.SendStatus(fiber.StatusNoContent)
}

func ToggleMonitoring(c *fiber.Ctx) error {
	// Extract user information from JWT token claims
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	userID := uint(claims["user_id"].(float64))

	// Extract group ID and URL ID from URL parameters
	groupIDParam := c.Params("group_id")
	groupID, err := strconv.ParseUint(groupIDParam, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid group ID",
		})
	}

	urlIDParam := c.Params("url_id")
	urlID, err := strconv.ParseUint(urlIDParam, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid URL ID",
		})
	}

	// Check if the group belongs to the user
	var group models.Group
	if err := database.DB.Where("id = ? AND user_id = ?", groupID, userID).First(&group).Error; err != nil {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "Group not found or access denied",
		})
	}

	// Find the URL by ID and ensure it belongs to the specified group
	var url models.URL
	if err := database.DB.Where("id = ? AND group_id = ?", urlID, groupID).First(&url).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "URL not found or access denied",
		})
	}

	// Toggle the monitoring value
	url.Monitoring = !url.Monitoring

	// Update URL in DB
	if err := database.DB.Save(&url).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update monitoring status",
		})
	}

	return c.JSON(url)
}
