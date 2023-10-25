package server

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"

	"catbook.com/auth"
	"catbook.com/db"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func Start(database *sql.DB) {
	r := gin.Default()
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = []string{"*"}
	r.Use(cors.New(corsConfig))
	r.Use(authMiddleware())

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	r.POST("/signin", func(c *gin.Context) {
		var userCreds auth.UserCredentials
		c.BindJSON(&userCreds)

		if valid, err := db.ValidUserCredentials(&userCreds, database); err == nil {
			fmt.Println(valid)
			if valid {
				c.JSON(http.StatusOK, gin.H{
					"message": "signin successful",
				})
			} else {
				c.JSON(http.StatusBadRequest, gin.H{
					"message": "invalid username or password",
				})
			}
		} else {
			c.Status(http.StatusInternalServerError)
		}
	})

	r.POST("/register", func(c *gin.Context) {
		var userRegInfo auth.UserRegistrationInfo
		c.BindJSON(&userRegInfo)

		if err := db.CreateUser(&userRegInfo, database); err != nil {
			fmt.Println(err)
			c.JSON(http.StatusConflict, gin.H{
				"message": "username or email already exists",
			})
		} else {
			c.JSON(http.StatusCreated, gin.H{
				"message": "user account successfully created",
			})
		}
	})

	host := os.Getenv("HOST")
	port := os.Getenv("PORT")
	if host == "" || port == "" {
		r.Run()
	} else {
		r.Run(fmt.Sprintf("%s:%s", host, port))
	}
}
