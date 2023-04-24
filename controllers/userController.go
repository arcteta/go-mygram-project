package controllers

import (
	"fmt"
	"go-mygram/database"
	"go-mygram/helpers"
	"go-mygram/models"
	"go-mygram/repository"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

var (
	appJSON = "application/json"
)

type Creds struct {
	Username string `json:"username" form:"username" valid:"required~Your username is required"`
	Password string `json:"password" form:"password" valid:"required~Your password is required,minstringlength(6)~Password minimum lengths is 6 characters"`
}

type ClientInput struct {
	Username string `json:"username" form:"username"`
	FullName string `json:"full_name" form:"full_name"`
	Email    string `json:"email" form:"email"`
	Password string `json:"password,omitempty" form:"password"`
	Age      uint   `json:"age" form:"age"`
}

func UserRegister(c *gin.Context) {

	db := database.GetDB()
	contentType := helpers.GetContentType(c)
	_, _ = db, contentType
	User := models.User{}

	UserExist := repository.FindUser(User.Username, User.Email)

	if UserExist.Username != "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Bad Request",
			"message": "username already registered",
		})
		return
	}

	if contentType == appJSON {
		c.ShouldBindJSON(&User)
	} else {
		c.ShouldBind(&User)
	}

	err := repository.CreateUser(&User)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Bad request",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"id":       User.ID,
		"name":     User.Name,
		"username": User.Username,
		"email":    User.Email,
		"age":      User.Age,
	})
}

func UserLogin(c *gin.Context) {
	contentType := helpers.GetContentType(c)

	Credentials := Creds{}

	if contentType == appJSON {
		c.ShouldBindJSON(&Credentials)
	} else {
		c.ShouldBind(&Credentials)
	}

	user, err := repository.FindUserByUsername(Credentials.Username)

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":    "unauthorized",
			"messaage": fmt.Sprintf("user %s doesnt exist", Credentials.Username),
		})
		return
	}
	log.Println(user)

	comparePass := helpers.ComparePass([]byte(user.Password), []byte(Credentials.Password))

	if !comparePass {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   "Unauthorized",
			"message": "wrong password, please recheck your password",
		})
		return
	}

	token := helpers.GenerateToken(user.ID, user.Email)

	c.JSON(http.StatusOK, gin.H{
		"token": token,
	})
}
