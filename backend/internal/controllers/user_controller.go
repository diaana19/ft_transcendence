package controllers

import (
	"errors"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"ft_transcendence/backend/internal/models"
	"ft_transcendence/backend/internal/services"
)

// UserController handles the user profile endpoints.
type UserController struct {
	userService   *services.UserService
	friendService *services.FriendService
	mailService   *services.MailService
}

// DeleteAccountInput holds the password used to confirm account deletion.
type DeleteAccountInput struct {
	Password string `json:"password" binding:"required"`
}

// NewUserController creates a new UserController with its services.
func NewUserController(
	userService *services.UserService,
	friendService *services.FriendService,
	mailService *services.MailService,
) *UserController {
	return &UserController{userService: userService, friendService: friendService, mailService: mailService}
}

// GetUsers lists all users, looks one up by exact username, or searches them by text.
// @Summary   List users, look one up by exact username, or search by text (q)
// @Tags      users
// @Security  BearerAuth
// @Produce   json
// @Param     username query string false "exact username; returns the single matching user instead of the list"
// @Param     q        query string false "search text matched against username and display name (paginated)"
// @Param     limit    query int    false "max users to return when searching (default 10, max 50)"
// @Param     offset   query int    false "users to skip when searching (default 0)"
// @Success   200 {array}  models.UserResponse
// @Failure   404 {object} map[string]string
// @Failure   500 {object} map[string]string
// @Router    /users [get]
func (uc *UserController) GetUsers(c *gin.Context) {
	if q, ok := c.GetQuery("q"); ok {
		uc.searchUsers(c, q)
		return
	}

	if username := c.Query("username"); username != "" {
		user, err := uc.userService.GetUserByUsername(username)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				c.JSON(http.StatusNotFound, gin.H{"error": msgUserNotFound})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		response := user.ToResponse()
		response.FollowersCount, _ = uc.friendService.CountFollowers(user.ID)
		response.FollowingCount, _ = uc.friendService.CountFollowing(user.ID)
		c.JSON(http.StatusOK, response)
		return
	}

	users, err := uc.userService.GetUsers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	responses := make([]models.UserResponse, len(users))
	for i, u := range users {
		responses[i] = u.ToResponse()
	}
	c.JSON(http.StatusOK, responses)
}

// searchUsers returns a page of users matching the text query, with the total count.
func (uc *UserController) searchUsers(c *gin.Context, query string) {
	limit, offset := parseLimitOffset(c)

	users, total, err := uc.userService.SearchUsers(query, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	responses := make([]models.UserResponse, len(users))
	for i, u := range users {
		responses[i] = u.ToResponse()
	}
	c.JSON(http.StatusOK, gin.H{"data": responses, "limit": limit, "offset": offset, "total": total})
}

// GetUser returns one user by its id with the follower counts.
// @Summary   Get a user by id
// @Tags      users
// @Security  BearerAuth
// @Produce   json
// @Param     id path string true "user id"
// @Success   200 {object} models.UserResponse
// @Failure   404 {object} map[string]string
// @Failure   500 {object} map[string]string
// @Router    /users/{id} [get]
func (uc *UserController) GetUser(c *gin.Context) {
	id := c.Param("id")
	user, err := uc.userService.GetUser(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": msgUserNotFound})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response := user.ToResponse()
	response.FollowersCount, _ = uc.friendService.CountFollowers(user.ID)
	response.FollowingCount, _ = uc.friendService.CountFollowing(user.ID)

	c.JSON(http.StatusOK, response)
}

// UpdateUser updates the profile. The user can only update his own profile.
// @Summary   Update a user's profile
// @Tags      users
// @Security  BearerAuth
// @Accept    json
// @Produce   json
// @Param     id   path string                 true "user id"
// @Param     body body models.UpdateUserInput true "fields to update"
// @Success   200 {object} models.UserResponse
// @Failure   400 {object} map[string]string
// @Failure   401 {object} map[string]string
// @Failure   403 {object} map[string]string
// @Failure   404 {object} map[string]string
// @Failure   500 {object} map[string]string
// @Router    /users/{id} [put]
func (uc *UserController) UpdateUser(c *gin.Context) {
	targetID := c.Param("id")
	currentUserID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	if currentUserID.(string) != targetID {
		c.JSON(http.StatusForbidden, gin.H{"error": "you can only update your own profile"})
		return
	}

	var input models.UpdateUserInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := uc.userService.UpdateUser(targetID, input)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": msgUserNotFound})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, user.ToResponse())
}

// DeleteUser deletes the account after the password is confirmed. The user can only delete his own.
// @Summary   Delete a user's account
// @Tags      users
// @Security  BearerAuth
// @Accept    json
// @Produce   json
// @Param     id   path string             true "user id"
// @Param     body body DeleteAccountInput true "password confirmation"
// @Success   200 {object} map[string]string
// @Failure   400 {object} map[string]string
// @Failure   401 {object} map[string]string
// @Failure   403 {object} map[string]string
// @Failure   404 {object} map[string]string
// @Failure   500 {object} map[string]string
// @Router    /users/{id} [delete]
func (uc *UserController) DeleteUser(c *gin.Context) {
	targetID := c.Param("id")
	currentUserID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	if currentUserID.(string) != targetID {
		c.JSON(http.StatusForbidden, gin.H{"error": "you can only delete your own profile"})
		return
	}

	var input DeleteAccountInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "password confirmation required"})
		return
	}

	if err := uc.userService.VerifyPassword(targetID, input.Password); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": msgUserNotFound})
			return
		}
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid password"})
		return
	}

	deletedUser, _ := uc.userService.GetUser(targetID)

	if err := uc.userService.DeleteUser(targetID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": msgUserNotFound})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if deletedUser != nil && deletedUser.Email != "" {
		go func(email, username string) {
			subject := "Your Synk account has been deleted"
			body := "Hi " + username + ",\n\nThis confirms that your Synk account and all associated data " +
				"have been permanently deleted, as you requested.\n\n" +
				"If you did not request this, please contact us immediately.\n"
			if err := uc.mailService.SendMail([]string{email}, subject, body); err != nil {
				log.Printf("gdpr: deletion confirmation email failed for %s: %v", email, err)
			}
		}(deletedUser.Email, deletedUser.Username)
	}

	c.JSON(http.StatusOK, gin.H{"message": "User deleted"})
}
