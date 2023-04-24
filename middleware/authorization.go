package middleware

import (
	"errors"
	"go-mygram/repository"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

var Error error 

func Authorization() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Method == "PUT" {
			Error = errors.New("you are not allowed to edit this data")
		} else if c.Request.Method == "DELETE" {
			Error= errors.New("you are not allowed to delete this data")
		}

		id, errConvert := strconv.Atoi(c.Param("id"))
		if errConvert != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error": "invalid id",
			})
			return
		}

		userData := c.MustGet("userData").(jwt.MapClaims)
		userID := uint(userData["id"].(float64))

		var userIDData uint 
		var err error
	
		if strings.Contains(c.Request.URL.Path, "photo") {
			photo, dberr := repository.FindPhotoUser(uint(id))
			err = dberr
			userIDData = photo.UserID
		} else if strings.Contains(c.Request.URL.Path, "comment") {
			comment, dberr := repository.FindCommentUser(uint(id))
			err = dberr
			userIDData = comment.UserID
		} else if strings.Contains(c.Request.URL.Path, "social-media") {
			sm, dberr := repository.FindSocialMediaUser(uint(id))
			err = dberr
			userIDData = sm.UserID
		}

		if err != nil {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
				"error":   "Data not found",
				"message": err.Error(),
			})
			return
		}

		if userIDData != userID {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error":   "Unauthorized",
				"message": Error.Error(),
			})
			return
		}

		c.Next()
	}
}
