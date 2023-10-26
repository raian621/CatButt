package server

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"os"

	"catbook.com/auth"
	"catbook.com/db"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

var EmailFieldError = errors.New("invalid email field")
var UsernameFieldError = errors.New("invalid username field")
var PasswordFieldError = errors.New("invalid password field")
var MultipleFieldError = errors.New("multple invalid fields")

func NewRouter(database *sql.DB) *gin.Engine {
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
		if cookie, err := c.Request.Cookie("sessionid"); err == nil {
			sessionId := cookie.Value
			validSession, err := db.ValidSession(sessionId, database)
			if err != nil {
				fmt.Println(err)
			}
			if validSession {
				c.AbortWithStatusJSON(http.StatusAccepted, gin.H{
					"message": "already signed in",
				})
				return
			}
		}

		var userCreds auth.UserCredentials
		c.BindJSON(&userCreds)

		if valid, err := db.ValidUserCredentials(&userCreds, database); err == nil {
			if valid {
				user, err := db.GetUserByUsername(userCreds.Username, database)
				if err != nil {
					c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
						"message": "could not find user",
					})
					return
				}

				session, err := db.NewSession(user.UserId, c.Request.UserAgent(), c.ClientIP(), database)
				if err != nil {
					fmt.Println(err)
					c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
						"message": "could not create session to login",
					})
					return
				}

				if err := db.CreateSession(session, database); err != nil {
					fmt.Println(err)
					c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
						"message": "could not create session to login",
					})
					return
				}

				c.SetCookie(
					"sessionid",
					session.SessionId,
					60*60*24*7,
					"/",
					"localhost",
					false,
					false,
				)

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
		if valid, err := validateRegistrationFields(&userRegInfo); !valid {
			fmt.Println(err)
			c.AbortWithStatusJSON(400, gin.H{"message": "invalid information"})
			return
		}

		if err := db.CreateUser(&userRegInfo, database); err != nil {
			fmt.Println(err, "User", userRegInfo.Username, "already exists")
			c.JSON(http.StatusConflict, gin.H{
				"message": "username or email already exists",
			})
			return
		} else {
			c.JSON(http.StatusCreated, gin.H{
				"message": "user account successfully created",
			})
			return
		}
	})

	return r
}

func NewDatabase() (*sql.DB, error) {
	dbParams := &db.DatabaseParams{
		Provider: os.Getenv("DB_PROV"),
		Hostname: os.Getenv("DB_HOST"),
		Username: os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASS"),
		Database: os.Getenv("DB_NAME"),
		Port:     os.Getenv("DB_PORT"),
		SSLMode:  os.Getenv("DB_SSL"),
	}

	database, err := db.ConnectToDB(dbParams)
	dbParams = nil
	if err != nil {
		fmt.Println("Database offline!")
		return nil, err
	}

	if err = db.CreateTables(database); err != nil {
		return nil, err
	}

	return database, nil
}

func validateRegistrationFields(userRegInfo *auth.UserRegistrationInfo) (_ bool, err error) {
	if unameLen := len(userRegInfo.Username); unameLen == 0 || unameLen >= 50 {
		err = UsernameFieldError
	}
	if passwdLen := len(userRegInfo.Password); passwdLen < 8 {
		if err != nil {
			err = MultipleFieldError
		} else {
			err = PasswordFieldError
		}
	}
	if emailLen := len(userRegInfo.Email); emailLen == 0 {
		if err != nil {
			err = MultipleFieldError
		} else {
			err = EmailFieldError
		}
	}

	return err == nil, err
}
