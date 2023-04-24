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

func GetAllSocialMedia(c *gin.Context) {
	userData := c.MustGet("userData").(jwt.MapClaims)
	userId := uint(userData["id"].(float64))
	social_media, err := repository.FindAllSocialMedia(userId)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Error getting social media data",
			"err":     err.Error(),
		})
		return
	}

	for _, social_media := range social_media {
		social_media.User.Password = ""
	}
	c.JSON(http.StatusOK, &social_media)
}

func GetOneSocialMedia(c *gin.Context) {
	socialmediaID, _ := strconv.Atoi(c.Param("id"))
	social_media, err := repository.FindSocialMediaById(uint(socialmediaID))

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{
				"message": "social media not found",
				"err":     "not found",
			})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Error getting social media",
			"err":     err.Error(),
		})
		return
	}

	social_media.User.Password = ""
	c.JSON(http.StatusOK, &social_media)
}

type SocialMediaInput struct {
	Name           string `json:"name" form:"name"`
	SocialMediaURL string `json:"social_media_url" form:"social_media_url"`
}

func CreateSocialMedia(c *gin.Context) {
	userData := c.MustGet("userData").(jwt.MapClaims)
	contentType := helpers.GetContentType(c)

	SocialMedia := models.SocialMedia{}

	userID := uint(userData["id"].(float64))

	if contentType == appJSON {
		c.ShouldBindJSON(&SocialMedia)
	} else {
		c.ShouldBind(&SocialMedia)
	}

	SocialMedia.UserID = userID
	SocialMedia.Name = strings.TrimSpace(SocialMedia.Name)
	SocialMedia.SocialUrl = strings.TrimSpace(SocialMedia.SocialUrl)

	err := repository.CreateSocialMedia(&SocialMedia)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"err":     "Bad Request",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "your social media successfully added",
		"data":    &SocialMedia,
	})
}
func UpdateSocialMedia(c *gin.Context) {
	userData := c.MustGet("userData").(jwt.MapClaims)
	contentType := helpers.GetContentType(c)
	SocialMediainput := models.SocialMedia{}

	socialMediaID, _ := strconv.Atoi(c.Param("id"))
	userID := uint(userData["id"].(float64))

	if contentType == appJSON {
		c.ShouldBindJSON(&SocialMediainput)
	} else {
		c.ShouldBind(&SocialMediainput)
	}

	SocialMediainput.UserID = userID
	SocialMediainput.ID = uint(socialMediaID)
	SocialMediainput.Name = strings.TrimSpace(SocialMediainput.Name)
	SocialMediainput.SocialUrl = strings.TrimSpace(SocialMediainput.SocialUrl)

	_, err := url.ParseRequestURI(SocialMediainput.SocialUrl)
	if err != nil && SocialMediainput.SocialUrl != "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"err":     "Bad Request",
			"message": "invalid url",
		})
		return
	}

	updatedSocialMedia, err := repository.UpdateSocialMedia(&SocialMediainput)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"err":     "Bad Request",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "your social media is successfully updated",
		"data":    &updatedSocialMedia,
	})
}
func DeleteSocialMedia(c *gin.Context) {
	socialmediaID, _ := strconv.Atoi(c.Param("id"))

	err := repository.DeleteSocialMedia(uint(socialmediaID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   err.Error(),
			"message": "Can't delete social media",
		})
		return
	}

	c.JSON(http.StatusOK, "Social media successfully deleted")
}