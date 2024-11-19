package controllers

import (
	"net/http"

	"github.com/ecabigting/letsgo-brrr/usermanager-api/models"
	"github.com/ecabigting/letsgo-brrr/usermanager-api/services"
	"github.com/ecabigting/letsgo-brrr/usermanager-api/utils"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type UserController struct {
	service *services.UserService
}

func NewUserController(service *services.UserService) *UserController {
	return &UserController{service: service}
}

func (uc *UserController) CreateUser(c *gin.Context) {
	var user models.User
	// Validate password complexity
	validate := validator.New()
	err := validate.RegisterValidation("password_complexity", utils.PasswordComplexity)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Bind request body to the the User Model
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check for role and auto set it
	if user.Role == "" {
		user.Role = "User"
	}

	if err := validate.Struct(user); err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			if err.Tag() == "password_complexity" {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Password must be at least 8 characters long, contain at least one uppercase letter, one lowercase letter, one number, and one special character."})
				return
			}
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if email exist
	exist := uc.service.CheckIfEmailExist(user.Email)
	if exist {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email already exist!"})
		return
	}

	// Create user
	err = uc.service.CreateUser(&user)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to create user", "desc": err.Error()})
		return
	}

	// Log the transaction
	utils.Logger.Info("User created: ", user.Email)

	c.JSON(http.StatusCreated, gin.H{"message": "User created successfully"})
}

// VerifyUser handles user verification
func (uc *UserController) VerifyUser(c *gin.Context) {
	userID := c.Param("userId")
	var request struct {
		Token string `json:"token" binding:"required"`
	}
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}
	err := uc.service.VerifyUser(userID, request.Token)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User verified successfully"})
}

// GetUser handles retrieving a user by ID
func (uc *UserController) GetUser(c *gin.Context) {
	userID := c.Param("userId")
	user, err := uc.service.GetUser(userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, user)
}

// UpdateUser handles updating user information
func (uc *UserController) UpdateUser(c *gin.Context) {
	userID := c.Param("userId")
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	err := uc.service.UpdateUser(userID, &user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User updated successfully"})
}

// DeleteUser handles deleting a user by ID
func (uc *UserController) DeleteUser(c *gin.Context) {
	userID := c.Param("userId")
	err := uc.service.DeleteUser(userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}
