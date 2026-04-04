package controllers

import (
	"net/http"

	"github.com/Transcendence/models"
	"github.com/Transcendence/services"
	"github.com/gin-gonic/gin"
)

type UserController struct {
	userService *services.UserService
}

func NewUserController(userService *services.UserService) *UserController {
	return &UserController{userService: userService}
}

func (uc *UserController) GetUsers(c *gin.Context) {
	users, err := uc.userService.GetUsers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, users)
}

func (uc *UserController) GetUser(c *gin.Context) {
	id := c.Param("id")
	user, err := uc.userService.GetUser(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}
	c.JSON(http.StatusOK, user)
}

func (uc *UserController) UpdateUser(c *gin.Context) {
	id := c.Param("id")
	var input models.UpdateUserInput
	err := c.ShouldBindJSON(&input)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error:": err.Error()})
		return
	}
	user, errorU := uc.userService.UpdateUser(id, input)
	if errorU != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error:": errorU.Error()})
		return
	}
	c.JSON(http.StatusOK, user)
}

func (uc *UserController) DeleteUser(c *gin.Context) {
	id := c.Param("id")
	err := uc.userService.DeleteUser(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "User deleted"})
}
