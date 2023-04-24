package controllers

import (
	"errors"
	"go-mygram/helpers"
	"go-mygram/models"
	"go-mygram/repository"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
)

type ClientInputPhoto struct {
	Title    string `json:"title"`
	Caption  string `json:"caption"`
	PhotoURL string `json:"photo_url"`
}

func GetAllPhoto(c *gin.Context) {
	userData := c.MustGet("userData").(jwt.MapClaims)
	userId := uint(userData["id"].(float64))
	photo, err := repository.FindAllPhoto(userId)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Error getting photo data",
			"err":     err.Error(),
		})
		return
	}

	for _, photo := range photo {
		photo.User.Password = ""
	}
	c.JSON(http.StatusOK, photo)
}

func GetOnePhoto(c *gin.Context) {

	photoID, _ := strconv.Atoi(c.Param("id"))
	photo, err := repository.FindPhotoByID(uint(photoID))

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{
				"message": "photo not found",
				"err":     "not found",
			})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Error getting photo",
			"err":     err.Error(),
		})
		return
	}

	photo.User.Password = ""
	c.JSON(http.StatusOK, &photo)
}

func CreatePhoto(c *gin.Context) {
	userData := c.MustGet("userData").(jwt.MapClaims)
	contentType := helpers.GetContentType(c)

	Photo := models.Photo{}

	userID := uint(userData["id"].(float64))

	if contentType == appJSON {
		c.ShouldBindJSON(&Photo)
	} else {
		c.ShouldBind(&Photo)
	}

	Photo.UserID = userID
	Photo.Title = strings.TrimSpace(Photo.Title)
	Photo.PhotoUrl = strings.TrimSpace(Photo.PhotoUrl)
	Photo.Caption = strings.TrimSpace(Photo.Caption)

	err := repository.CreatePhoto(&Photo)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"err":     "Bad Request",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, &Photo)
}

func UpdatePhoto(c *gin.Context) {
	userData := c.MustGet("userData").(jwt.MapClaims)
	contentType := helpers.GetContentType(c)
	PhotoInput := models.Photo{}

	photoID, _ := strconv.Atoi(c.Param("id"))
	userID := uint(userData["id"].(float64))

	if contentType == appJSON {
		c.ShouldBindJSON(&PhotoInput)
	} else {
		c.ShouldBind(&PhotoInput)
	}

	PhotoInput.UserID = userID
	PhotoInput.ID = uint(photoID)
	PhotoInput.Title = strings.TrimSpace(PhotoInput.Title)
	PhotoInput.PhotoUrl = strings.TrimSpace(PhotoInput.PhotoUrl)
	PhotoInput.Caption = strings.TrimSpace(PhotoInput.Caption)

	_, err := url.ParseRequestURI(PhotoInput.PhotoUrl)
	if err != nil && PhotoInput.PhotoUrl != "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"err":     "Bad Request",
			"message": "invalid url",
		})
		return
	}

	updatedPhoto, err := repository.UpdatePhoto(&PhotoInput)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"err":     "Bad Request",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, &updatedPhoto)
}

func DeletePhoto(c *gin.Context) {
	photoID, _ := strconv.Atoi(c.Param("id"))

	err := repository.DeletePhoto(uint(photoID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   err.Error(),
			"message": "Can't delete photo",
		})
		return
	}

	c.JSON(http.StatusOK, "Photo successfully deleted")
}
