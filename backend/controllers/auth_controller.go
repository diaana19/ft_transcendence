package controllers

import (
	"fmt"
	"log"
	"strings"

	"github.com/Transcendence/models"
	"github.com/gin-gonic/gin"
)

func checkEmail(email string) bool {
	if email == "" {
		return false
	}
	if !strings.Contains(email, "@") {
		return false
	}
	after_arobase := email[strings.Index(email, "@"):]
	if strings.Contains(after_arobase, "@") || !strings.Contains(after_arobase, ".") {
		return false
	}
	return true
}

// auth route
func RegisterUser(c *gin.Context) {
	var user models.User 
	err := c.BindJSON(&user);
	if err != nil {
		log.Fatal(err);
	}
	fmt.Println(user)
	// parse email
	checkEmail(user.Email)
	c.JSON(200, gin.H{
		"msg": "mon msg",
	})
}