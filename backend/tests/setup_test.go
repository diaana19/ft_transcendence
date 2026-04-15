package tests

import (
	"os"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/Transcendence/config"
	"github.com/Transcendence/models"
	"github.com/Transcendence/routes"
)

func SetupTestEnv() (*gin.Engine, *gorm.DB) {
	gin.SetMode(gin.TestMode)

	os.Setenv("DB_HOST", "localhost")
	os.Setenv("DB_NAME", "app_db")

	db, err := config.ConnectDB()
	if err != nil {
		panic(err)
	}

	rdb, err := config.InitRedis()
	if err != nil {
		panic(err)
	}

	db.Exec("DROP TABLE IF EXISTS users CASCADE")
	db.AutoMigrate(&models.User{})

	router := gin.Default()
	routes.SetupRoutes(router, db, rdb)

	return router, db
}
